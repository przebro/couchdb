package database

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/przebro/couchdb/client"

	"github.com/przebro/couchdb/request"

	"github.com/przebro/couchdb/cursor"
)

type bufferedCursor struct {
	resultset []byte
	dec       *json.Decoder
	ep        string
	cli       *client.CouchClient
	sel       DataSelector
	meta      cursor.QueryMeta
}

type responseResult struct {
	Data     json.RawMessage        `json:"docs"`
	Bookmark string                 `json:"bookmark"`
	Warning  string                 `json:"warning"`
	Stats    map[string]interface{} `json:"execution_stats"`
}

func newBufferedCursor(rdr io.ReadCloser, ep string, sel DataSelector, cli *client.CouchClient) *bufferedCursor {

	var docs int
	result, err := getResultset(rdr)

	if err != nil {
		return nil
	}

	if result.Stats != nil {
		docs = int(result.Stats["results_returned"].(float64))
	}

	return &bufferedCursor{
		resultset: result.Data,
		ep:        ep,
		sel:       sel,
		cli:       cli,
		meta:      cursor.QueryMeta{Warning: result.Warning, Bookmark: result.Bookmark, Documents: docs},
	}
}

func (s *bufferedCursor) All(ctx context.Context, v interface{}) error {

	rval := reflect.ValueOf(v)
	if rval.Kind() != reflect.Ptr {
		return errInvalidDocKind
	}

	sval := rval.Elem()
	if sval.Kind() == reflect.Interface {
		sval = sval.Elem()
	}

	if sval.Kind() != reflect.Slice {
		return errInvalidDocKind
	}

	etype := sval.Type().Elem()

	for s.Next(ctx) {

		newElem := reflect.New(etype)
		i := newElem.Interface()
		s.Decode(i)
		sval.Set(reflect.Append(sval, newElem.Elem()))

	}

	return nil
}

func (s *bufferedCursor) Next(ctx context.Context) bool {

	if s.dec == nil {

		s.dec = createDecoder(s.resultset)
		if _, err := s.dec.Token(); err != nil {
			return false
		}
	}

	if !s.dec.More() {
		return s.fetchNextResultset(ctx)
	}

	return true
}

func (s *bufferedCursor) Decode(v interface{}) error {

	return s.dec.Decode(v)
}
func (s *bufferedCursor) Meta() cursor.QueryMeta {
	return s.meta
}
func (s *bufferedCursor) Close(ctx context.Context) error {
	s.resultset = nil
	s.cli = nil
	s.dec = nil

	return nil
}

func (s *bufferedCursor) fetchNextResultset(ctx context.Context) bool {

	if ctx == nil {
		ctx = context.Background()
	}

	s.sel.Bookmark = s.meta.Bookmark
	selct, err := json.Marshal(s.sel)

	rqb := request.NewRequestBuilder()
	request, err := rqb.WithEndpoint(s.ep).WithMethod(request.MethodPost).WithBody(selct).Build(s.cli)

	r, err := request.Execute(ctx)
	if err != nil {
		return false
	}
	doc, err := getResultset(r.Rdr)

	if err != nil {
		return false
	}

	s.resultset = doc.Data
	s.meta = cursor.QueryMeta{Bookmark: doc.Bookmark, Documents: 0, Warning: doc.Warning}

	s.dec = createDecoder(s.resultset)
	s.dec.Token()

	return checkDocuments(s.resultset)
}

func getResultset(rdr io.ReadCloser) (responseResult, error) {

	data, _ := ioutil.ReadAll(rdr)
	dc := responseResult{}

	err := json.Unmarshal(data, &dc)
	if err != nil {
		return dc, err
	}

	return dc, nil
}
func createDecoder(data []byte) *json.Decoder {
	rdr := bytes.NewReader(data)
	return json.NewDecoder(rdr)
}

func checkDocuments(data []byte) bool {

	dec := createDecoder(data)
	_, err := dec.Token()

	if err != nil {

		return false
	}

	return dec.More()
}
