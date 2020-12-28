package database

import (
	"encoding/json"
	"errors"

	"github.com/przebro/couchdb/client"
)

type FindOption string

const (
	endPointFind     = "_find"
	endPointBulk     = "_bulk_docs"
	endPointPurge    = "_purge"
	endPointSecurity = "_security"

	OptionStat     FindOption = "stat"
	OptionBookmark FindOption = "bookmark"
	OptionLimit    FindOption = "limit"
	OptionIndex    FindOption = "index"
)

var (
	errNilSelector          = errors.New("nil selector specified")
	errInvalidDocKind       = errors.New("invalid kind of document, not a ptr to slice, or not a ptr to struct")
	errEmptyDocumentID      = errors.New("document id cannot be empty")
	errRequiredDocumentID   = errors.New("document id required")
	errRequiredestinationID = errors.New("destination document id required")
	errIDandRevRequired     = errors.New("id and rev fields are required")
	errRevListRequired      = errors.New("revision list cannot be empty")
	errSecurityDataEmpty    = errors.New("empty security data")
)

type arrrayDocument struct {
	Documents []interface{} `json:"docs"`
}

//CopyDocument - struct used to copy document, must contain new ID of a document and optionally revision if copied to existing document
type CopyDocument struct {
	ID  string
	Rev string
}

//DataSelector - Contains strucutrued used in _find request
type DataSelector struct {
	Selector json.RawMessage `json:"selector"`
	Fields   []string        `json:"fields,omitempty"`
	Limit    int             `json:"limit,omitempty"`
	Bookmark string          `json:"bookmark,omitempty"`
	Index    string          `json:"use_index,omitempty"`
	Stats    bool            `json:"execution_stats,omitempty"`
}

//CouchDatabase - Represents a CouchDB database
type CouchDatabase struct {
	Name string
	cli  *client.CouchClient
}
