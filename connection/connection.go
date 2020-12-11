package connection

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/przebro/couchdb/client"
	"github.com/przebro/couchdb/request"
	"github.com/przebro/couchdb/response"
)

const (
	endPointSession = "_session"
	endPointUp      = "_up"
	endPointAllDbs  = "_all_dbs"
	endPointDbsInfo = "_dbs_info"
	endPointUuids   = "_uuids"
)

//Connection - Represents server connection
type Connection struct {
	cli *client.CouchClient
}

//GetClient - Returns context
func (c *Connection) GetClient() *client.CouchClient {
	return c.cli
}

//GetSession - Gets a session information
func (c *Connection) GetSession(ctx context.Context) (*response.CouchResult, error) {

	b := request.NewRequestBuilder()
	rq, err := b.WithEndpoint(endPointSession).WithMethod(request.MethodGet).Build(c.cli)
	if err != nil {
		return nil, err
	}
	rs, err := rq.Execute(ctx)

	return response.NewResult(rs.CouchStatus, rs.Rdr), err

}

/*Session - Establishes a new session. If successful then the returned cookie will be atached to context making it possible to call
db specific endpoints.
*/
func (c *Connection) Session(ctx context.Context, user, password string) (*response.CouchResult, error) {

	b := request.NewRequestBuilder()

	body := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}{
		Name:     user,
		Password: password,
	}

	doc, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}

	request, err := b.WithEndpoint(endPointSession).WithMethod(request.MethodPost).WithBody(doc).Build(c.cli)

	if err != nil {
		return nil, err
	}

	rs, err := request.Execute(ctx)
	c.cli.AuthData = rs.Cookie.Value

	return response.NewResult(rs.CouchStatus, rs.Rdr), err
}

//Up - Checks database connection
func (c *Connection) Up(ctx context.Context) (*response.CouchResult, error) {

	b := request.NewRequestBuilder()
	rq, err := b.WithEndpoint(endPointUp).WithMethod(request.MethodGet).Build(c.cli)

	if err != nil {
		return nil, err
	}

	rs, err := rq.Execute(ctx)

	return response.NewResult(rs.CouchStatus, rs.Rdr), err
}

//Uuid - Generates uuid
func (c *Connection) Uuid(ctx context.Context, count int) (*response.CouchResult, error) {

	if count < 0 || count > 1000 {
		return nil, errors.New("invalid count parameter")
	}

	b := request.NewRequestBuilder()
	rq, err := b.WithEndpoint(endPointUuids).WithMethod(request.MethodGet).
		WithParameters(map[string]string{"count": fmt.Sprintf("%d", count)}).
		Build(c.cli)

	if err != nil {
		return nil, err
	}

	rs, err := rq.Execute(ctx)

	return response.NewResult(rs.CouchStatus, rs.Rdr), err
}

//AllDbs - Returns information about all databases
func (c *Connection) AllDbs(ctx context.Context) (*response.CouchResult, error) {

	b := request.NewRequestBuilder()
	rq, err := b.WithEndpoint(endPointAllDbs).WithMethod(request.MethodGet).Build(c.cli)

	if err != nil {
		return nil, err
	}

	rs, err := rq.Execute(ctx)

	return response.NewResult(rs.CouchStatus, rs.Rdr), err

}

//DbsInfo - Returns information about specific database
func (c *Connection) DbsInfo(ctx context.Context, name string) (*response.CouchResult, error) {

	if name == "" {
		return nil, errors.New("db name required")
	}

	body := map[string][]string{
		"keys": []string{name},
	}

	doc, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}

	b := request.NewRequestBuilder()
	request, err := b.WithEndpoint(endPointDbsInfo).WithMethod(request.MethodPost).
		WithBody(doc).
		Build(c.cli)

	if err != nil {
		return nil, err
	}

	rs, err := request.Execute(ctx)

	return response.NewResult(rs.CouchStatus, rs.Rdr), err

}
