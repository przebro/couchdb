package request

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/przebro/couchdb/client"
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
	//Nonstandard http method used by copy endpoint
	MethodCopy CouchMethod = "COPY"
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
func (req *CouchRequest) Execute(ctx context.Context) (response.CouchResponse, error) {

	var err error = nil
	var couchResponse response.CouchResponse

	if ctx == nil {
		ctx = context.Background()
	}

	req.request.Body.Close()
	req.request = req.request.WithContext(ctx)

	rc, e := req.execute(req.request)

	select {
	case couchResponse = <-rc:
		{
		}
	case err = <-e:
		{
		}

	case <-ctx.Done():
		{
			err = ctx.Err()
		}

	}

	return couchResponse, err
}
func (req *CouchRequest) execute(rq *http.Request) (<-chan response.CouchResponse, <-chan error) {

	ch := make(chan response.CouchResponse, 1)
	e := make(chan error, 1)
	go func(<-chan response.CouchResponse, <-chan error) {

		rs, err := req.cli.Do(rq)

		if err != nil {
			e <- err
			return
		}

		couchResponse := response.CouchResponse{CouchStatus: &response.CouchStatus{}}
		couchResponse.Code = rs.StatusCode
		couchResponse.Status = rs.Status
		couchResponse.Server = rs.Header.Get("Server")
		couchResponse.Rdr = rs.Body

		ck := rs.Cookies()

		if len(ck) > 0 {
			couchResponse.Cookie = *ck[0]
		}

		ch <- couchResponse

	}(ch, e)

	return ch, e
}

//Builder - Helps build a new CouchDB request
type Builder interface {
	WithBody(doc []byte) Builder
	WithMethod(method CouchMethod) Builder
	WithParameters(params map[string]string) Builder
	WithHeaders(headers map[string]string) Builder
	WithEndpoint(endpoint string) Builder
	Build(conn *client.CouchClient) (*CouchRequest, error)
}

type requestBuilder struct {
	endpoint string
	method   CouchMethod
	params   map[string]string
	headers  map[string]string
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
func (rb *requestBuilder) WithHeaders(headers map[string]string) Builder {
	if headers == nil {
		return rb
	}
	rb.headers = headers
	return rb
}

func (rb *requestBuilder) WithParameters(params map[string]string) Builder {
	if params == nil {
		return rb
	}
	rb.params = params
	return rb
}
func (rb *requestBuilder) WithEndpoint(endpoint string) Builder {

	rb.endpoint = endpoint
	return rb
}
func (rb *requestBuilder) Build(cli *client.CouchClient) (*CouchRequest, error) {

	var err error
	var method string

	switch rb.method {
	case MethodDelete, MethodGet, MethodPost, MethodPut, MethodHead, MethodCopy:
		{
			method = string(rb.method)
		}
	default:
		{
			return nil, errors.New("invalid http method")
		}
	}

	endp := fmt.Sprintf(`%s/%s`, cli.BaseAddr, rb.endpoint)

	if rb.params != nil {
		params := []string{}
		for k, v := range rb.params {
			params = append(params, fmt.Sprintf("%s=%s", k, v))
		}
		qstring := strings.Join(params, "&")
		endp = fmt.Sprintf("%s?%s", endp, qstring)
	}

	r := &CouchRequest{}
	r.cli = cli.Client

	bodyRdr := bytes.NewReader(rb.body)
	r.request, err = http.NewRequest(method, endp, bodyRdr)

	if err != nil {
		return nil, err
	}

	r.request.Header.Add("Content-Type", "application/json")

	if cli.Authentication == client.JwtToken {
		data := fmt.Sprintf("Bearer %s", cli.AuthData)
		r.request.Header.Add("Authorization", data)
	}
	if cli.Authentication == client.Basic {
		data := fmt.Sprintf("Basic %s", cli.AuthData)
		r.request.Header.Add("Authorization", data)
	}

	for k, v := range rb.headers {
		r.request.Header.Add(k, v)
	}

	if cli.Authentication == client.Cookie {
		ck := &http.Cookie{Name: "AuthSession", Value: cli.AuthData}
		r.request.AddCookie(ck)
	}

	return r, nil

}
