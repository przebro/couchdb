package response

import (
	"net/http"
)

/*ResultKey - A key for a value stored in a CouchResult. It is not guaranteed that every result will contain
the same set of a key-value pair, for instance, requests that end up successfully will not contain an error key.
*/
type ResultKey string

//Keys for data returned by CoucDB
const (
	ResponseStatusCode ResultKey = "code"
	ResponseStatus     ResultKey = "status"
	ResponseServer     ResultKey = "server"
	ResultMessage      ResultKey = "result"
)

//CouchResult - Contains data returned in response
type CouchResult map[ResultKey]interface{}

//CouchResponse - Wraps CouchResult and returned cookie
type CouchResponse struct {
	CouchResult
	Cookie http.Cookie
	Body   []byte
}
