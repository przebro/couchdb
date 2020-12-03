package query

import (
	"encoding/json"
)

//Data -
type Data struct {
	Selector json.RawMessage `json:"selector"`
	Fields   []string        `json:"fields,omitempty"`
	Limit    int             `json:"limit,omitempty"`
	Bookmark string          `json:"bookmark,omitempty"`
	Index    string          `json:"use_index,omitempty"`
}
