package qpage

import "github.com/creasty/defaults"

// Represents the pagination params given when performing pagination operations on a collection.
type Params struct {
	//The page to fetch.
	Page int `json:"page" default:"1"`
	//The maximum number of items to include on each page.
	ItemsPerPage int `json:"items_per_page" default:"24"`
	//The ID of the document to skip to.
	SkipToID interface{} `json:"skip_to_id"`
}

// Gets the default options for the pagination params.
func DefaultParams() Params {
	params := Params{}
	defaults.Set(&params)
	return params
}
