package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/przebro/couchdb/context"
	"github.com/przebro/couchdb/response"
)

//CouchMethod - restricts possible methods
type CouchMethod string

// Valid request methods
const (
	MethodHead   CouchMethod = http.MethodHead
	MethodGet    CouchMethod = http.MethodGet
	MethodPut    CouchMethod = http.MethodPut
	MethodPost   CouchMethod = http.MethodPost
	MethodDelete CouchMethod = http.MethodDelete

	StatusCode200OK           = 200
	StatusCode304NotModified  = 304
	StatusCode400BadRequest   = 400
	StatusCode400Unauthorized = 401
	StatusCode400NotFound     = 404
)

var (
	errRequest = errors.New("request error")
)

//CouchRequest - Wraps http request
type CouchRequest struct {
	cli     *http.Client
	request *http.Request
}

//Execute - executes request
func (req *CouchRequest) Execute() (response.CouchResponse, error) {

	var err error
	couchResponse := response.CouchResponse{CouchResult: make(map[response.ResultKey]interface{})}
	req.request.Body.Close()

	resp, err := req.cli.Do(req.request)

	//Something really bad happen
	if err != nil {
		return couchResponse, err
	}

	couchResponse.CouchResult[response.ResponseStatusCode] = resp.StatusCode
	couchResponse.CouchResult[response.ResponseStatus] = resp.Status
	couchResponse.CouchResult[response.ResponseServer] = resp.Header.Get("Server")

	couchResponse.Body = readBody(resp)
	if resp.StatusCode >= StatusCode400BadRequest {

		return couchResponse, fmt.Errorf("%s:%v;", errRequest, strings.Trim(string(couchResponse.Body), "\r\n"))
	}

	ck := resp.Cookies()

	if len(ck) > 0 {
		couchResponse.Cookie = *ck[0]
	}

	return couchResponse, nil
}

func readBody(r *http.Response) []byte {
	var err error = nil
	var num int
	buffer := make([]byte, 1024)
	data := make([]byte, 0)
	for err != io.EOF {
		num, err = r.Body.Read(buffer)
		data = append(data, buffer[0:num]...)
	}

	return data
}

//Builder - Helps build a new CouchDB request
type Builder interface {
	WithBody(doc []byte) Builder
	WithMethod(method CouchMethod) Builder
	WithParameters(params map[string]string) Builder
	WithEndpoint(endpoint string) Builder
	Build(conn *context.CouchContext) (*CouchRequest, error)
}

type requestBuilder struct {
	endpoint string
	method   CouchMethod
	params   map[string]string
	body     []byte
}

//NewRequestBuilder - Creates a new instance of RequestBuilder
func NewRequestBuilder() Builder {
	builder := &requestBuilder{}

	return builder
}

func (rb *requestBuilder) WithBody(doc []byte) Builder {
	rb.body = doc
	return rb
}
func (rb *requestBuilder) WithMethod(method CouchMethod) Builder {
	rb.method = method
	return rb

}
func (rb *requestBuilder) WithParameters(params map[string]string) Builder {

	rb.params = params
	return rb
}
func (rb *requestBuilder) WithEndpoint(endpoint string) Builder {

	rb.endpoint = endpoint
	return rb
}
func (rb *requestBuilder) Build(ctx *context.CouchContext) (*CouchRequest, error) {

	var err error
	var method string

	switch rb.method {
	case MethodDelete, MethodGet, MethodPost, MethodPut, MethodHead:
		{
			method = string(rb.method)
		}
	default:
		{
			return nil, errors.New("invalid http method")
		}
	}

	endp := fmt.Sprintf(`%s/%s`, ctx.BaseAddr, rb.endpoint)

	if rb.params != nil {
		params := []string{}
		for k, v := range rb.params {
			params = append(params, fmt.Sprintf("%s=%s", k, v))
		}
		qstring := strings.Join(params, "&")
		endp = fmt.Sprintf("%s?%s", endp, qstring)
	}

	r := &CouchRequest{}
	r.cli = ctx.Client

	bodyRdr := bytes.NewReader(rb.body)
	r.request, err = http.NewRequest(method, endp, bodyRdr)

	if err != nil {
		return nil, err
	}

	r.request.Header.Add("Content-Type", "application/json")

	if ctx.Authentication == context.JwtToken {
		data := fmt.Sprintf("Bearer %s", ctx.AuthData)
		r.request.Header.Add("Authorization", data)
	}
	if ctx.Authentication == context.Basic {
		data := fmt.Sprintf("Basic %s", ctx.AuthData)
		r.request.Header.Add("Authorization", data)
	}

	if ctx.Authentication == context.Cookie {
		ck := &http.Cookie{Name: "AuthSession", Value: ctx.AuthData}
		r.request.AddCookie(ck)
	}

	return r, nil

}
