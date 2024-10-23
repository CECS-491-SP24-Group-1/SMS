package db

import (
	"context"

	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/pkg/util"
)

// QMgoBase is the base struct for all collections.
type QMgoBase struct {
	*qmgo.Collection
}

// Queries for an object by its ID.
func (qb QMgoBase) FindID(ctx context.Context, ID util.UUID, opts ...opts.FindOptions) qmgo.QueryI {
	query := bson.D{{Key: "_id", Value: ID}}
	return qb.Find(ctx, query, opts...)
}
