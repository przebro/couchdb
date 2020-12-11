package connection

import (
	"context"
	"fmt"
	"testing"

	"github.com/przebro/couchdb/client"
)

const host = "127.0.0.1"
const port = 5300
const sport = 6300
const username string = "admin"
const password string = "notsecure"

type UpResponse struct {
	Status string
	Seeds  interface{}
}

type SessionResponse struct {
	Ok    bool     `json:"ok"`
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}

type UuidResponse struct {
	Uuid []string `json:"uuids"`
}

func TestAuth(t *testing.T) {

	builder := NewBuilder()
	//Build should fails if authentication is not set
	conn, err := builder.WithAddress(host, port).Build(true)
	if err == nil {
		t.Error(err)
	}

	//
	conn, err = builder.WithAuthentication(client.None, username, password).WithAddress(host, port).Build(true)

	if err != nil {
		t.Error(err)
	}
	result, err := conn.Up(context.Background())
	if err != nil {
		t.Error(err)
	}
	if result.Code != 200 {
		t.Error("Unexpected result")
	}

	conn, err = builder.WithAuthentication(client.Basic, username, password).WithAddress(host, port).Build(true)

	t.Log(conn.cli.BaseAddr)

	if err != nil {
		t.Error(err)
	}
	result, err = conn.Up(context.Background())
	if err != nil {
		t.Error(err)
	}
	if result.Code != 200 {
		t.Error("Unexpected result")
	}

	conn, err = builder.WithAuthentication(client.Cookie, username, password).WithAddress(host, port).Build(true)

	if err != nil {
		t.Error(err)
	}
	result, err = conn.Up(context.Background())
	if err != nil {
		t.Error(err)
	}
	if result.Code != 200 {
		t.Error("Unexpected result")
	}

	ctx := conn.GetClient()
	if ctx == nil {
		t.Error("Unexpected result")
	}
}

func TestUp(t *testing.T) {
	builder := NewBuilder()
	conn, _ := builder.WithAuthentication(client.Cookie, username, password).WithAddress(host, port).Build(true)

	result, err := conn.Up(context.Background())
	if err != nil {
		t.Error(err)
	}

	up := UpResponse{}
	err = result.Decode(&up)
	if err != nil {
		t.Error(err)
	}

}

func TestGetSession(t *testing.T) {

	builder := NewBuilder()
	conn, err := builder.WithAuthentication(client.Cookie, username, password).WithAddress(host, port).Build(true)
	if err != nil {
		t.Error(err)
	}
	result, err := conn.GetSession(context.Background())
	if err != nil {
		t.Error(err)
	}

	if result.Code != 200 {
		t.Error("Unexpected status code", result.Code)
	}

	ses := SessionResponse{}

	err = result.Decode(&ses)
	if err != nil {
		t.Error(err)
	}

	if ses.Ok != true && ses.Name != username {
		t.Error("unexpected result")
	}

}

func TestUuid(t *testing.T) {

	builder := NewBuilder()
	conn, err := builder.WithAuthentication(client.Cookie, username, password).WithAddress(host, port).Build(true)
	if err != nil {
		t.Error(err)
	}
	rs, err := conn.Uuid(context.Background(), 5)
	if err != nil {
		t.Error(err)
	}

	if rs.Code != 200 {
		t.Error("Unexpected message code")
	}

	uuids := UuidResponse{}
	err = rs.Decode(&uuids)

	if len(uuids.Uuid) != 5 {
		t.Error("Unexpected result")
	}

	rs, err = conn.Uuid(context.Background(), 1235)
	if err == nil {
		t.Error("Unexpected result")
	}
}

func TestAllDbs(t *testing.T) {

	builder := NewBuilder()
	conn, err := builder.WithAuthentication(client.Cookie, username, password).WithAddress(host, port).Build(true)
	if err != nil {
		t.Error(err)
	}

	result, err := conn.AllDbs(context.Background())
	if err != nil {
		t.Error(err)
	}

	if result.Code != 200 {
		t.Error("Unexpected result")
	}

	res := []string{}
	err = result.Decode(&res)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
	if len(res) < 2 {
		t.Error("Unexpected result")
	}

}

func TestDbInfo(t *testing.T) {

	builder := NewBuilder()
	conn, err := builder.WithAuthentication(client.Cookie, username, password).WithAddress(host, port).Build(true)
	if err != nil {
		t.Error(err)
	}

	_, err = conn.DbsInfo(context.Background(), "")

	if err == nil {
		t.Error(err)
	}

	r, err := conn.DbsInfo(context.Background(), "sd")
	if err != nil {
		t.Error(err)
	}

	if r.Code != 200 {
		t.Error("Unexpected status code")
	}

	r, err = conn.DbsInfo(context.Background(), "_users")
	if err != nil {
		t.Error(err)
	}

	if r.Code != 200 {
		t.Error("Unexpected status code")
	}

}
