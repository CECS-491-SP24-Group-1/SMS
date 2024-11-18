package qpage

import (
	"context"
	"fmt"
	"reflect"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/pkg/util"
)

// Represents a qmgo paginator.
type QPage struct {
	collection *qmgo.Collection //The backing qmgo collection.
	sortFields []sorter         //A list of the sorting keys/orders to apply.
}

// Creates a new `QPage` object.
func NewQPage[T any](col *qmgo.Collection) (*QPage, error) {
	//Get a pointer to the incoming object
	//dt := new(T)

	//Ensure the incoming object has an ID field

	return &QPage{
		collection: col,
		//datatype:   *dt,
		sortFields: make([]sorter, 0),
	}, nil
}

// Adds a sort field to the paginator.
func (q *QPage) Sort(key string, order int) *QPage {
	//Ensure the order is valid
	if order != -1 && order != 1 {
		panic(fmt.Sprintf("%d is not a valid sort order; expecting 1 or -1", order))
	}
	q.sortFields = append(q.sortFields, sorter{Name: key, Order: order})
	return q
}

// Runs an aggregation pipeline on the target collection and outputs pagination data.
func (q QPage) Aggregate(dest any, ctx context.Context, pipeline bson.A, params Params) (*Pagination, error) {
	//If SkipToID is provided, adjust the pipeline to include it as the first item
	if params.SkipToID != nil {
		pipeline = append(pipeline,
			bson.M{"$match": bson.M{"_id": bson.M{"$gte": params.SkipToID}}}, //Match for documents greater than or equal to SkipToID
			bson.M{"$sort": bson.M{"_id": 1}},                                //Ensure we sort by ID to bring SkipToID as the first item
		)
		//params.Page = 1 // Reset page to 1 when skipping to a specific ID
	}

	//Add the key sorts to the pipeline if there is at least one
	if len(q.sortFields) > 0 {
		sorts := bson.M{}
		for _, key := range q.sortFields {
			sorts[key.Name] = key.Order
		}
		pipeline = append(pipeline, bson.M{"$sort": sorts})
	}

	//Calculate total count and paginated results in a single pipeline using $facet
	pipeline = append(pipeline,
		//Set up the facet query
		bson.M{"$facet": bson.M{
			//Query 1: count the number of documents in the aggregation
			"metadata": bson.A{bson.M{"$count": "total"}},
			//Query 2: get only a single page of documents
			"data": bson.A{
				bson.M{"$skip": (params.Page - 1) * params.ItemsPerPage}, //Skips x documents in the database
				bson.M{"$limit": params.ItemsPerPage},                    //Limits the number of documents
			},
		}},
		//Move the total down to the root of the document
		bson.M{"$project": bson.M{
			"total": bson.M{"$arrayElemAt": bson.A{"$metadata.total", 0}},
			"data":  1,
		}},
	)

	//Perform the aggregation
	var result aresult
	err := q.collection.Aggregate(ctx, pipeline).One(&result)
	if err != nil {
		return nil, err
	}

	//Paginate the results
	return paginate(dest, result.Data, result.Total, params)
}

// Runs an aggregation pipeline on the target collection and outputs pagination data.
func (q QPage) Find(dest any, ctx context.Context, query interface{}, params Params) (*Pagination, error) {
	//Perform the find to get the query cursor
	res := q.collection.Find(ctx, query)

	//Get the number of documents in the result
	count, err := res.Count()
	if err != nil {
		return nil, err
	}

	//Add the key sorts to the cursor if there is at least one
	if len(q.sortFields) > 0 {
		sorts := make([]string, len(q.sortFields))
		for i, key := range q.sortFields {
			sorts[i] = util.If(key.Order == 1, key.Name, "-"+key.Name)
		}
		res.Sort(sorts...)
	}

	//Add the skip amount to the cursor based on the page number
	res.Skip(int64((params.Page - 1) * params.ItemsPerPage))

	//Get the documents from the database, limiting the count to the max per page
	var docs []bson.D
	res.Limit(int64(params.ItemsPerPage)).All(&docs)

	//Paginate the results
	return paginate(dest, docs, count, params)
}

// Contains the common backend logic for pagination queries.
func paginate(dest any, docs []bson.D, count int64, params Params) (*Pagination, error) {
	//Ensure the destination is a pointer to a slice
	destVal, err := assertCorrectOutputType(dest)
	if err != nil {
		return nil, err
	}

	//Calculate the attributes of the entire pagination run
	paginationInfo := Pagination{
		PerPage:    params.ItemsPerPage,
		TotalPages: (count + int64(params.ItemsPerPage) - 1) / int64(params.ItemsPerPage),
		TotalItems: count,
	}

	//Get the total number of documents in the current page
	psize := len(docs)

	//Calculate the attributes of the singular page that was just pulled out
	firstIdx := (params.Page-1)*params.ItemsPerPage + 1
	pageInfo := Page{
		Num:  params.Page,
		Size: psize,
		//IsLast:   params.ItemsPerPage > psize,
		IsLast:   params.Page == int(paginationInfo.TotalPages),
		IsEmpty:  psize == 0,
		FirstIdx: firstIdx,
		LastIdx:  firstIdx + psize - 1,
	}

	if !pageInfo.IsEmpty {
		//Get the first and last document IDs
		firstID, ok := getValueFromBsonD(docs[0], "_id")
		if ok {
			pageInfo.FirstID = idString(firstID)
		}
		lastID, ok := getValueFromBsonD(docs[psize-1], "_id")
		if ok {
			pageInfo.LastID = idString(lastID)
		}
	}

	//Add the current page info to the pagination info
	paginationInfo.CurrentPage = pageInfo

	//Create a new slice of the same type as dest and allocate space
	destElemType := destVal.Elem().Type().Elem()
	out := reflect.MakeSlice(reflect.SliceOf(destElemType), len(docs), len(docs))

	//Unmarshal each document in the results data array
	for i, doc := range docs {
		//Unmarshal the current doc
		data, err := unmarshalBsonD(doc, destElemType)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling doc #%d: %w", i, err)
		}

		//Set the doc in the output slice at position i
		out.Index(i).Set(reflect.ValueOf(data))
	}

	//Replace the original dest slice with the new one
	destVal.Elem().Set(out)

	//Return the pagination data
	return &paginationInfo, nil
}
