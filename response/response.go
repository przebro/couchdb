package response

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/przebro/couchdb/cursor"
)

const (
	StatusCode200OK                 = 200
	StatusCode304NotModified        = 304
	StatusCode400BadRequest         = 400
	StatusCode401Unauthorized       = 401
	StatusCode404NotFound           = 404
	StatusCode409Conflict           = 409
	StatusCode412PreconditionFailed = 412
)

//CouchStatus  - Contains http response status and additional info
type CouchStatus struct {
	Code   int
	Status string
	Server string
}

//CouchResult - Contains data returned in response
type CouchResult struct {
	*CouchStatus
	rdr io.ReadCloser
}

//NewResult - creates a new CouchResult, an object that wraps Requests status and body readcloser.
func NewResult(status *CouchStatus, rdr io.ReadCloser) *CouchResult {

	return &CouchResult{CouchStatus: status, rdr: rdr}
}

//Decode - Reads from response body and unmarshal datac into v
func (r *CouchResult) Decode(v interface{}) error {

	data, err := ioutil.ReadAll(r.rdr)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

//CouchMultiResult - Conttains data returned in response. this struct is useful when a body contains
//multiple objects like a result of a _find
type CouchMultiResult struct {
	*CouchStatus
	cursor.ResultCursor
}

//NewMultiResult - Creates a new CouchMultiResult
func NewMultiResult(status *CouchStatus, crsr cursor.ResultCursor) *CouchMultiResult {

	return &CouchMultiResult{CouchStatus: status, ResultCursor: crsr}

}

//CouchResponse - Wraps CouchStatus and returned cookie
type CouchResponse struct {
	*CouchStatus
	Rdr    io.ReadCloser
	Cookie http.Cookie
}
