package qpage

import (
	"net/http"
	"strconv"

	"github.com/creasty/defaults"
)

// Represents the pagination params given when performing pagination operations on a collection.
type Params struct {
	//The page to fetch.
	Page int `json:"page" default:"1"`
	//The maximum number of items to include on each page.
	ItemsPerPage int `json:"items_per_page" default:"50"`
	//The ID of the document to skip to.
	SkipToID interface{} `json:"skip_to_id"`
}

// Gets the default options for the pagination params.
func DefaultParams() Params {
	params := Params{}
	defaults.Set(&params)
	return params
}

// Parses the query parameters of a URL to get the paging parameters. Errors are silently ignored.
func ParseQuery(r *http.Request) Params {
	//Derive a default paging params object
	pagingParms := DefaultParams()

	//Get the query params and attempt to get the page and the items per page
	queryParams := r.URL.Query()
	pageNum := queryParams.Get("page")
	perPage := queryParams.Get("items_per_page")

	//Attempt to derive uint values from both the page and per page
	//`err == nil` conditions indicate no errors during conversion
	page, err := strconv.ParseUint(pageNum, 10, 0)
	if err == nil {
		pagingParms.Page = int(page)
	}
	ppage, err := strconv.ParseUint(perPage, 10, 0)
	if err == nil {
		pagingParms.ItemsPerPage = int(ppage)
	}

	return pagingParms
}
