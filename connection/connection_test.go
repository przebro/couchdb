package connection

import (
	"fmt"
	"testing"

	"github.com/przebro/couchdb/context"
	"github.com/przebro/couchdb/response"
)

const host = "127.0.0.1"
const port = 5300
const username string = "admin"
const password string = "notsecure"

func TestAuth(t *testing.T) {

	builder := NewBuilder()
	//Build should fails if authentication is not set
	conn, err := builder.WithAddress(host, port).Build()
	if err == nil {
		t.Error(err)
	}
	//
	conn, err = builder.WithAuthentication(context.None, username, password).WithAddress(host, port).Build()

	t.Log(conn.ctx.BaseAddr)

	if err != nil {
		t.Error(err)
	}
	result, err := conn.Up()
	if err != nil {
		t.Error(err)
	}
	if result[response.ResponseStatusCode] != 200 {
		t.Error("Unexpected result")
	}

	conn, err = builder.WithAuthentication(context.Basic, username, password).WithAddress(host, port).Build()

	t.Log(conn.ctx.BaseAddr)

	if err != nil {
		t.Error(err)
	}
	result, err = conn.Up()
	if err != nil {
		t.Error(err)
	}
	if result[response.ResponseStatusCode] != 200 {
		t.Error("Unexpected result")
	}

	conn, err = builder.WithAuthentication(context.Cookie, username, password).WithAddress(host, port).Build()

	if err != nil {
		t.Error(err)
	}
	result, err = conn.Up()
	if err != nil {
		t.Error(err)
	}
	if result[response.ResponseStatusCode] != 200 {
		t.Error("Unexpected result")
	}

	ctx := conn.GetContext()
	if ctx == nil {
		t.Error("Unexpected result")
	}
}

func TestGetSession(t *testing.T) {

	builder := NewBuilder()
	conn, err := builder.WithAuthentication(context.Cookie, username, password).WithAddress(host, port).Build()
	if err != nil {
		t.Error(err)
	}
	_, err = conn.GetSession()
	if err != nil {
		t.Error(err)
	}
}

func TestUuid(t *testing.T) {

	builder := NewBuilder()
	conn, err := builder.WithAuthentication(context.Cookie, username, password).WithAddress(host, port).Build()
	if err != nil {
		t.Error(err)
	}
	rs, err := conn.Uuid(5)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(rs)

	rs, err = conn.Uuid(1235)
	if err == nil {
		t.Error("Unexpected result")
	}
}
