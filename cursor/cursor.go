package cursor

import (
	"context"
)

//QueryMeta - Contains bookmark and number of documents
type QueryMeta struct {
	Bookmark  string
	Documents int
	Warning   string
}

//ResultCursor - Helps iterate over returned result
type ResultCursor interface {
	//All - returns all documents from cursor
	All(ctx context.Context, v interface{}) error
	//Next - check if there is a next document in current resultset, if not fetch next resultset
	Next(ctx context.Context) bool
	//Decode - decodes current document
	Decode(v interface{}) error
	//Meta - returns meta information for current resultset. At least a bookmark is guaranteed to be present in meta
	Meta() QueryMeta
	//Close closes cursor
	Close(ctx context.Context) error
}
