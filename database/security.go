package database

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/przebro/couchdb/request"
	"github.com/przebro/couchdb/response"
)

//SecurityData - Holds security informations
type SecurityData struct {
	Names []string `json:"names"`
	Roles []string `json:"roles"`
}

//SetMemberSecurity - Sets the security object for the given database.
func (db *CouchDatabase) SetMemberSecurity(ctx context.Context, names, roles []string) (*response.CouchResult, error) {

	if names == nil || roles == nil {
		return nil, errSecurityDataEmpty
	}

	sdata := SecurityData{
		Names: names,
		Roles: roles,
	}

	sobj := map[string]interface{}{"members": sdata}
	data, err := json.Marshal(&sobj)
	if err != nil {
		return nil, err
	}
	return db.setSecurity(ctx, data)

}

//SetAdminSecurity - Sets the security object for the given database.
func (db *CouchDatabase) SetAdminSecurity(ctx context.Context, names, roles []string) (*response.CouchResult, error) {

	if names == nil || roles == nil {
		return nil, errSecurityDataEmpty
	}

	sdata := SecurityData{
		Names: names,
		Roles: roles,
	}

	sobj := map[string]interface{}{"admins": sdata}
	data, err := json.Marshal(&sobj)
	if err != nil {
		return nil, err
	}
	return db.setSecurity(ctx, data)
}

func (db *CouchDatabase) setSecurity(ctx context.Context, data []byte) (*response.CouchResult, error) {

	endpoint := fmt.Sprintf("%s/%s", db.Name, endPointSecurity)

	rqb := request.NewRequestBuilder()

	rq, err := rqb.WithEndpoint(endpoint).WithMethod(request.MethodPut).WithBody(data).Build(db.cli)
	if err != nil {
		return nil, err
	}

	rs, err := rq.Execute(ctx)
	return response.NewResult(rs.CouchStatus, rs.Rdr), err

}

//Security - Returns the current security object from the specified database.
func (db *CouchDatabase) Security(ctx context.Context) (*response.CouchResult, error) {

	endpoint := fmt.Sprintf("%s/%s", db.Name, endPointSecurity)
	rqb := request.NewRequestBuilder()

	rq, err := rqb.WithEndpoint(endpoint).WithMethod(request.MethodGet).Build(db.cli)
	if err != nil {
		return nil, err
	}

	rs, err := rq.Execute(ctx)
	return response.NewResult(rs.CouchStatus, rs.Rdr), err
}
