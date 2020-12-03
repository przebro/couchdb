package connection

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/przebro/couchdb/context"
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
	ctx *context.CouchContext
}

//GetContext - Returns context
func (c *Connection) GetContext() *context.CouchContext {
	return c.ctx
}

//GetSession - Gets a session information
func (c *Connection) GetSession() (response.CouchResult, error) {

	b := request.NewRequestBuilder()
	rq, err := b.WithEndpoint(endPointSession).WithMethod(request.MethodGet).Build(c.ctx)
	if err != nil {
		return response.CouchResult{}, err
	}
	rs, err := rq.Execute()

	if err == nil {
		rs.CouchResult[response.ResultMessage] = strings.Trim(string(rs.Body), "\r\n")

	}

	return rs.CouchResult, err

}

/*Session - Establishes a new session. If successful then the returned cookie will be atached to context making it possible to call
db specific endpoints.
*/
func (c *Connection) Session(user, password string) (response.CouchResult, error) {

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
		return response.CouchResult{}, err
	}

	request, err := b.WithEndpoint(endPointSession).WithMethod(request.MethodPost).WithBody(doc).Build(c.ctx)

	if err != nil {
		return response.CouchResult{}, err
	}

	resp, err := request.Execute()

	if err != nil {
		return response.CouchResult{}, err
	}
	c.ctx.AuthData = resp.Cookie.Value

	return resp.CouchResult, err
}

func (c *Connection) Up() (response.CouchResult, error) {

	b := request.NewRequestBuilder()
	rq, err := b.WithEndpoint(endPointUp).WithMethod(request.MethodGet).Build(c.ctx)

	if err != nil {
		return response.CouchResult{}, err
	}

	rs, err := rq.Execute()
	if err == nil {
		rs.CouchResult[response.ResultMessage] = strings.Trim(string(rs.Body), "\r\n")

	}

	return rs.CouchResult, err
}

//Uuid - Generates uuid
func (c *Connection) Uuid(count int) (response.CouchResult, error) {

	if count < 0 || count > 1000 {
		return response.CouchResult{}, errors.New("invallid count parameter")
	}

	b := request.NewRequestBuilder()
	rq, err := b.WithEndpoint(endPointUuids).WithMethod(request.MethodGet).
		WithParameters(map[string]string{"count": fmt.Sprintf("%d", count)}).
		Build(c.ctx)

	if err != nil {
		return response.CouchResult{}, err
	}

	rs, err := rq.Execute()

	if err == nil {
		rs.CouchResult[response.ResultMessage] = strings.Trim(string(rs.Body), "\r\n")

	}

	return rs.CouchResult, err
}

func (c *Connection) AllDbs() (response.CouchResult, error) {

	b := request.NewRequestBuilder()
	rq, err := b.WithEndpoint(endPointAllDbs).WithMethod(request.MethodGet).Build(c.ctx)

	if err != nil {
		return response.CouchResult{}, err
	}

	rs, err := rq.Execute()

	return rs.CouchResult, err

}

func (c *Connection) DbsInfo() (response.CouchResult, error) {

	b := request.NewRequestBuilder()
	request, err := b.WithEndpoint(endPointAllDbs).WithMethod(request.MethodGet).Build(c.ctx)

	if err != nil {
		return response.CouchResult{}, err
	}

	response, err := request.Execute()

	return response.CouchResult, err

}
