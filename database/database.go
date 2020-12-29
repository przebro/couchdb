package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/przebro/couchdb/client"
	"github.com/przebro/couchdb/request"
	"github.com/przebro/couchdb/response"
)

//CreateDatabase - Creates new database
func CreateDatabase(ctx context.Context, name string, cli *client.CouchClient) (*response.CouchResult, CouchDatabase, error) {

	var database CouchDatabase
	builder := request.NewRequestBuilder()

	request, err := builder.WithMethod(request.MethodPut).WithEndpoint(name).Build(cli)
	if err != nil {
		return nil, database, err
	}
	rs, err := request.Execute(ctx)

	if err == nil {
		database.cli = cli
		database.Name = name

	}

	if rs.Code >= response.StatusCode400BadRequest {
		err = errors.New(rs.Status)
	}

	return response.NewResult(rs.CouchStatus, rs.Rdr), database, err
}

//GetDatabsase - Gets database
func GetDatabsase(ctx context.Context, name string, cli *client.CouchClient) (*response.CouchResult, CouchDatabase, error) {

	var database CouchDatabase
	builder := request.NewRequestBuilder()

	rq, err := builder.WithMethod(request.MethodGet).WithEndpoint(name).Build(cli)
	if err != nil {
		return nil, database, err
	}
	rs, err := rq.Execute(ctx)

	if err == nil {
		database.cli = cli
		database.Name = name
	}

	if rs.Code == response.StatusCode404NotFound {
		err = errors.New(rs.Status)
	}

	return response.NewResult(rs.CouchStatus, rs.Rdr), database, err

}

//DropDatabase - Removes database
func DropDatabase(ctx context.Context, name string, cli *client.CouchClient) (*response.CouchResult, error) {

	builder := request.NewRequestBuilder()

	rq, err := builder.WithMethod(request.MethodDelete).WithEndpoint(name).Build(cli)
	if err != nil {
		return nil, err
	}
	rs, err := rq.Execute(ctx)

	if rs.Code >= response.StatusCode400BadRequest {
		err = errors.New(rs.Status)
	}

	return response.NewResult(rs.CouchStatus, rs.Rdr), err
}

//Get - Gets a single document with given id
func (db *CouchDatabase) Get(ctx context.Context, id string) (*response.CouchResult, error) {

	if id == "" {
		return nil, errEmptyDocumentID
	}

	endpoint := fmt.Sprintf("%s/%s", db.Name, id)
	rqb := request.NewRequestBuilder()

	request, err := rqb.WithEndpoint(endpoint).WithMethod(request.MethodGet).Build(db.cli)
	if err != nil {
		return nil, err
	}

	rs, err := request.Execute(ctx)

	if rs.Code >= response.StatusCode400BadRequest {
		err = errors.New(rs.Status)
	}

	return response.NewResult(rs.CouchStatus, rs.Rdr), err
}

//Stat -returns details about current database
func (db *CouchDatabase) Stat(ctx context.Context) (*response.CouchResult, error) {

	builder := request.NewRequestBuilder()

	request, err := builder.WithMethod(request.MethodGet).WithEndpoint(db.Name).Build(db.cli)
	if err != nil {
		return nil, err
	}
	rs, err := request.Execute(ctx)

	if rs.Code >= response.StatusCode400BadRequest {
		err = errors.New(rs.Status)
	}

	return response.NewResult(rs.CouchStatus, rs.Rdr), err
}

//Select - Selects documents from the database.
func (db *CouchDatabase) Select(ctx context.Context, sel string, fld []string, opt map[FindOption]interface{}) (*response.CouchMultiResult, error) {

	if sel == "" {
		return nil, errNilSelector
	}
	query := DataSelector{
		Selector: []byte(sel),
		Fields:   fld,
	}

	db.setFindOptions(&query, opt)

	body, err := json.Marshal(query)

	endpoint := fmt.Sprintf("%s/%s", db.Name, endPointFind)
	rqb := request.NewRequestBuilder()

	rq, err := rqb.WithEndpoint(endpoint).WithMethod(request.MethodPost).WithBody(body).Build(db.cli)

	if err != nil {
		return nil, err
	}

	rs, err := rq.Execute(ctx)

	if rs.Code >= response.StatusCode400BadRequest {
		err = errors.New(rs.Status)
	}

	return response.NewMultiResult(rs.CouchStatus, newBufferedCursor(rs.Rdr, endpoint, query, db.cli)), err
}

//Insert - Inserts a new document to databsase
func (db *CouchDatabase) Insert(ctx context.Context, doc interface{}) (*response.CouchResult, error) {

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
	rq, err := rqb.WithEndpoint(endpoint).WithMethod(method).WithBody(data).Build(db.cli)
	if err != nil {
		return nil, err
	}
	rs, err := rq.Execute(ctx)

	if rs.Code >= response.StatusCode400BadRequest {
		err = errors.New(rs.Status)
	}

	return response.NewResult(rs.CouchStatus, rs.Rdr), err

}

/*InsertMany - Inserts document in bulk - this method does not validate if every document contains
_rev and _id, it only checks if elements are structs so make sure to add these fields.
*/
func (db *CouchDatabase) InsertMany(ctx context.Context, docs []interface{}) (*response.CouchResult, error) {

	endpoint := fmt.Sprintf("%s/%s", db.Name, endPointBulk)
	rqb := request.NewRequestBuilder()
	if valid := isSliceOfStructs(docs); !valid {
		return nil, errInvalidDocKind

	}

	arr := arrrayDocument{Documents: docs}

	data, err := json.Marshal(arr)
	if err != nil {
		return nil, err
	}

	rq, err := rqb.WithEndpoint(endpoint).WithMethod(request.MethodPost).WithBody(data).Build(db.cli)
	if err != nil {
		return nil, err
	}

	rs, err := rq.Execute(ctx)

	if rs.Code >= response.StatusCode400BadRequest {
		err = errors.New(rs.Status)
	}

	return response.NewResult(rs.CouchStatus, rs.Rdr), err
}

//Revision - Gets all revisions of the document
func (db *CouchDatabase) Revision(ctx context.Context, id string) (*response.CouchResult, error) {

	endpoint := fmt.Sprintf("%s/%s", db.Name, id)
	rqb := request.NewRequestBuilder()

	rq, err := rqb.WithEndpoint(endpoint).WithParameters(map[string]string{"revs": "true"}).
		WithMethod(request.MethodGet).Build(db.cli)
	if err != nil {
		return nil, err
	}

	rs, err := rq.Execute(ctx)

	return response.NewResult(rs.CouchStatus, rs.Rdr), err
}

//Update - Updates document
func (db *CouchDatabase) Update(ctx context.Context, doc interface{}) (*response.CouchResult, error) {

	id, rev, err := requiredFields(doc)
	if err != nil {
		return nil, err
	}

	if id == "" || rev == "" {
		return nil, errIDandRevRequired
	}

	data, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/%s", db.Name, id)
	rqb := request.NewRequestBuilder()

	rq, err := rqb.WithEndpoint(endpoint).WithBody(data).
		WithMethod(request.MethodPut).Build(db.cli)
	if err != nil {
		return nil, err
	}

	rs, err := rq.Execute(ctx)

	if rs.Code >= response.StatusCode400BadRequest {
		err = errors.New(rs.Status)
	}

	return response.NewResult(rs.CouchStatus, rs.Rdr), err
}

/*Delete - Deletes a document from the database. Note that calling this method will
mark the document as _delete and leaves a tombstone. To completely remove the document a purge method should be called.
See the CouchDB documentation for more details.
*/
func (db *CouchDatabase) Delete(ctx context.Context, id, rev string) (*response.CouchResult, error) {

	if id == "" || rev == "" {
		return nil, errIDandRevRequired
	}

	endpoint := fmt.Sprintf("%s/%s", db.Name, id)
	rqb := request.NewRequestBuilder()

	rq, err := rqb.WithEndpoint(endpoint).WithMethod(request.MethodDelete).WithParameters(map[string]string{"rev": rev}).Build(db.cli)
	if err != nil {
		return nil, err
	}

	rs, err := rq.Execute(ctx)

	if rs.Code >= response.StatusCode400BadRequest {
		err = errors.New(rs.Status)
	}

	return response.NewResult(rs.CouchStatus, rs.Rdr), err

}

//Copy - Creates a new copy of a document
func (db *CouchDatabase) Copy(ctx context.Context, id, destID, destREV string) (*response.CouchResult, error) {

	if id == "" {
		return nil, errRequiredDocumentID
	}

	if destID == "" {
		return nil, errRequiredestinationID
	}

	headers := map[string]string{"Destination": destID}

	qParams := map[string]string{}

	if destREV != "" {
		qParams["rev"] = destREV
	}

	endpoint := fmt.Sprintf("%s/%s", db.Name, id)
	rqb := request.NewRequestBuilder()
	rq, err := rqb.WithEndpoint(endpoint).WithParameters(qParams).WithMethod(request.MethodCopy).
		WithHeaders(headers).
		Build(db.cli)
	if err != nil {
		return nil, err
	}

	rs, err := rq.Execute(ctx)

	if rs.Code >= response.StatusCode400BadRequest {
		err = errors.New(rs.Status)
	}

	return response.NewResult(rs.CouchStatus, rs.Rdr), err

}

//Purge - Purges revisions from document with a given id
func (db *CouchDatabase) Purge(ctx context.Context, id string, rev []string) (*response.CouchResult, error) {

	if id == "" {
		return nil, errEmptyDocumentID
	}

	if rev == nil || len(rev) == 0 {
		return nil, errRevListRequired
	}

	pdocs := map[string][]string{id: rev}
	data, err := json.Marshal(&pdocs)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/%s", db.Name, endPointPurge)
	rqb := request.NewRequestBuilder()
	rq, err := rqb.WithEndpoint(endpoint).WithMethod(request.MethodPost).
		WithBody(data).
		Build(db.cli)

	if err != nil {
		return nil, err
	}

	rs, err := rq.Execute(ctx)

	if rs.Code >= response.StatusCode400BadRequest {
		err = errors.New(rs.Status)
	}

	return response.NewResult(rs.CouchStatus, rs.Rdr), err
}

func (db *CouchDatabase) setFindOptions(s *DataSelector, opt map[FindOption]interface{}) {

	for k, v := range opt {
		switch k {
		case OptionBookmark:
			{
				if val, ok := v.(string); ok {
					s.Bookmark = val
				}
			}
		case OptionLimit:
			{
				if val, ok := v.(int); ok {
					s.Limit = val
				}
			}
		case OptionStat:
			{
				if val, ok := v.(bool); ok {
					s.Stats = val
				}
			}
		case OptionIndex:
			{
				if val, ok := v.(string); ok {
					s.Index = val
				}
			}
		}
	}
}
