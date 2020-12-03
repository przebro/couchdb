package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/przebro/couchdb/context"
	"github.com/przebro/couchdb/query"
	"github.com/przebro/couchdb/request"
	"github.com/przebro/couchdb/response"
)

const (
	endPointFind = "_find"
	endPointBulk = "_bulk_docs"
)

var (
	errNilSelector     = errors.New("nil selector specified")
	errInvalidDocKind  = errors.New("invalid kind of document, not a ptr to slice, or not a ptr to struct")
	errEmptyDocumentID = errors.New("document id cannot be empty")
)

type responseDocument struct {
	Documents map[string]*json.RawMessage `json:"docs"`
}

type arrrayDocument struct {
	Documents []interface{} `json:"docs"`
}

//CouchDatabase - Represents a CouchDB database
type CouchDatabase struct {
	Name string
	Ctx  *context.CouchContext
}

//CreateDatabase - Creates new database
func CreateDatabase(name string, ctx *context.CouchContext) (response.CouchResult, CouchDatabase, error) {

	var database CouchDatabase
	builder := request.NewRequestBuilder()

	request, err := builder.WithMethod(request.MethodPut).WithEndpoint(name).Build(ctx)
	if err != nil {
		return response.CouchResult{}, database, err
	}
	response, err := request.Execute()

	if err == nil {
		database.Ctx = ctx
		database.Name = name

	}
	return response.CouchResult, database, err

}

//GetDatabsase - Gets database
func GetDatabsase(name string, ctx *context.CouchContext) (response.CouchResult, CouchDatabase, error) {

	var database CouchDatabase
	builder := request.NewRequestBuilder()

	request, err := builder.WithMethod(request.MethodGet).WithEndpoint(name).Build(ctx)
	if err != nil {
		return response.CouchResult{}, database, err
	}
	response, err := request.Execute()

	if err == nil {
		database.Ctx = ctx
		database.Name = name

	}
	return response.CouchResult, database, err

}

//Get - Gets a single document with given id
func (db *CouchDatabase) Get(id string, doc interface{}) (response.CouchResult, error) {

	if id == "" {
		return nil, errEmptyDocumentID
	}

	if !isValidStruct(doc) {
		return nil, errInvalidDocKind
	}

	endpoint := fmt.Sprintf("%s/%s", db.Name, id)
	rqb := request.NewRequestBuilder()

	request, err := rqb.WithEndpoint(endpoint).WithMethod(request.MethodGet).Build(db.Ctx)
	if err != nil {
		return nil, err
	}

	result, err := request.Execute()
	if err != nil {
		return result.CouchResult, err
	}

	err = json.Unmarshal(result.Body, doc)

	return result.CouchResult, err
}

func (db *CouchDatabase) Stat() (response.CouchResult, error) {

	builder := request.NewRequestBuilder()

	request, err := builder.WithMethod(request.MethodGet).WithEndpoint(db.Name).Build(db.Ctx)
	if err != nil {
		return response.CouchResult{}, err
	}
	rs, err := request.Execute()

	if err == nil {
		rs.CouchResult[response.ResultMessage] = rs.Body
	}

	return rs.CouchResult, err
}

//Select - Selects documents from database
func (db *CouchDatabase) Select(doc interface{}, sel string, fld []string, opt map[string]string) (response.CouchResult, error) {

	if sel == "" {
		return nil, errNilSelector
	}

	if !isValidSlice(doc) {
		return nil, errInvalidDocKind
	}

	query := query.Data{
		Selector: []byte(sel),
		Fields:   fld,
	}

	selct, err := json.Marshal(query)

	endpoint := fmt.Sprintf("%s/%s", db.Name, endPointFind)
	rqb := request.NewRequestBuilder()

	request, err := rqb.WithEndpoint(endpoint).WithMethod(request.MethodPost).WithBody(selct).Build(db.Ctx)
	if err != nil {
		return nil, err
	}

	result, err := request.Execute()
	if err != nil {
		return result.CouchResult, err
	}
	rawdoc, err := extractDocFromBody(result.Body)
	if err != nil {
		return result.CouchResult, err
	}

	err = json.Unmarshal(*rawdoc.Documents["docs"], doc)

	return result.CouchResult, err

}

//Insert - Inserts a new document to databsase
func (db *CouchDatabase) Insert(doc interface{}) (response.CouchResult, error) {

	//Insert with existing id and rev makes and modified fileds makes an new revision
	//insert existing  without rev returns error
	//Insert without id creates new document even if rev is specified

	id, _, err := requiredFields(doc)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	endpoint, method := func() (string, request.CouchMethod) {
		if id != "" {
			return fmt.Sprintf(`%s/%s`, db.Name, id), request.MethodPut
		}
		return db.Name, request.MethodPost
	}()

	rqb := request.NewRequestBuilder()
	rq, err := rqb.WithEndpoint(endpoint).WithMethod(method).WithBody(data).Build(db.Ctx)
	if err != nil {
		return nil, err
	}

	rs, err := rq.Execute()
	if err == nil {
		rs.CouchResult[response.ResultMessage] = rs.Body
	}

	return rs.CouchResult, err

}

func (db *CouchDatabase) InsertMany(docs []interface{}) (response.CouchResult, error) {

	endpoint := fmt.Sprintf("%s/%s", db.Name, endPointBulk)
	rqb := request.NewRequestBuilder()

	arr := arrrayDocument{Documents: docs}

	data, err := json.Marshal(arr)
	if err != nil {
		return nil, err
	}

	rq, err := rqb.WithEndpoint(endpoint).WithMethod(request.MethodPost).WithBody(data).Build(db.Ctx)
	if err != nil {
		return nil, err
	}

	rs, err := rq.Execute()

	if err == nil {
		rs.CouchResult[response.ResultMessage] = rs.Body
	}

	return rs.CouchResult, err
}

//Revision - Gets all revisions of the document
func (db *CouchDatabase) Revision(id string) (response.CouchResult, error) {

	endpoint := fmt.Sprintf("%s/%s", db.Name, id)
	rqb := request.NewRequestBuilder()

	request, err := rqb.WithEndpoint(endpoint).WithParameters(map[string]string{"revs": "true"}).
		WithMethod(request.MethodGet).Build(db.Ctx)
	if err != nil {
		return nil, err
	}

	result, err := request.Execute()
	if err != nil {
		return nil, err
	}
	fmt.Println(result)

	return result.CouchResult, nil
}

//Update - Updates document
func (db *CouchDatabase) Update(doc interface{}) (response.CouchResult, error) {

	id, rev, err := requiredFields(doc)
	if err != nil {
		return nil, err
	}

	if id == "" || rev == "" {
		return nil, errors.New("id and rev fields are required")
	}

	data, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/%s", db.Name, id)
	rqb := request.NewRequestBuilder()

	request, err := rqb.WithEndpoint(endpoint).WithBody(data).
		WithMethod(request.MethodPut).Build(db.Ctx)
	if err != nil {
		return nil, err
	}

	result, err := request.Execute()
	return result.CouchResult, err
}

//Copy - Creates a new copy of a document
func (db *CouchDatabase) Copy(doc []byte) {
	//COPY /recipes/SpaghettiWithMeatballs

}

/*Delete - Deletes a document from the database. Note that calling this method will
mark the document as _delete and leaves a tombstone. To completely remove the document a purge method should be called.
See the CouchDB documentation for more details.
*/
func (db *CouchDatabase) Delete(doc interface{}) (response.CouchResult, error) {

	id, rev, err := requiredFields(doc)
	if err != nil {
		return nil, err
	}

	if id == "" || rev == "" {
		return nil, errors.New("id and rev fields are required")
	}

	endpoint := fmt.Sprintf("%s/%s", db.Name, id)
	rqb := request.NewRequestBuilder()

	request, err := rqb.WithEndpoint(endpoint).WithMethod(request.MethodDelete).WithParameters(map[string]string{"rev": rev}).Build(db.Ctx)
	if err != nil {
		return nil, err
	}

	result, err := request.Execute()
	return result.CouchResult, err

}

func (db *CouchDatabase) Purge(doc interface{}) {

}

//isValidSlice - Checks if dcoument is a pointer to slice of a struct
func isValidSlice(doc interface{}) bool {

	v := reflect.ValueOf(doc)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return false
	}
	tp := v.Elem().Type()

	if tp.Elem().Kind() != reflect.Struct {
		return false
	}
	return true
}

func isValidStruct(doc interface{}) bool {

	v := reflect.ValueOf(doc)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return false
	}
	return true
}

func requiredFields(doc interface{}) (string, string, error) {

	var id string
	var rev string
	if !isValidStruct(doc) {
		return id, rev, errInvalidDocKind
	}
	v := reflect.ValueOf(doc).Elem()

	for i := 0; i < v.NumField(); i++ {

		str, exists := v.Type().Field(i).Tag.Lookup("json")
		if exists {

			if strings.HasPrefix(str, "_id") && v.Field(i).Kind() == reflect.String {
				id = v.Field(i).Interface().(string)
			}

			if strings.HasPrefix(str, "_rev") && v.Field(i).Kind() == reflect.String {
				rev = v.Field(i).Interface().(string)
			}
		}
	}

	return id, rev, nil
}

func extractDocFromBody(raw []byte) (responseDocument, error) {

	rdoc := responseDocument{}
	err := json.Unmarshal(raw, &rdoc.Documents)
	return rdoc, err

}
