package database

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/przebro/couchdb/connection"
	"github.com/przebro/couchdb/context"
	"github.com/przebro/couchdb/response"
)

const host = "127.0.0.1"
const port = 5300
const username string = "admin"
const password string = "notsecure"
const database string = "couchtest"

type SampleDoc struct {
	Name  string `json:"Name"`
	Group string `json:"Group"`
	Age   int    `json:"Age"`
	ID    string `json:"_id,omitempty"`
	Rev   string `json:"_rev,omitempty"`
}

type testfunc func()

var builder connection.ConnectionBuilder
var conn *connection.Connection

func init() {
	var err error
	builder = connection.NewBuilder()
	conn, err = builder.WithAuthentication(context.Basic, username, password).WithAddress(host, port).Build()
	if err != nil {
		fmt.Println("error")
	}

}

func TestCreateDatabase(t *testing.T) {

	_, db, err := CreateDatabase(database, conn.GetContext())
	if err != nil {
		t.Error("unexpected result")
	}
	_, err = db.Stat()
	if err != nil {
		t.Error("unexpected result")
	}

}

func TestInsertDocument(t *testing.T) {

	resmsg := struct {
		Ok  bool   `json:"ok"`
		ID  string `json:"id"`
		Rev string `json:"rev"`
	}{}

	result, database, err := GetDatabsase(database, conn.GetContext())
	if err != nil {
		t.Log(result)
	}

	sample := SampleDoc{ID: "test_document_id_01", Name: "User Name", Age: 21, Group: "group_1"}
	result, err = database.Insert(&sample)
	if err != nil {
		t.Error(err)
	}

	sample = SampleDoc{ID: "test_document_id_02", Name: "User Name", Age: 21, Group: "group_1"}
	database.Insert(&sample)

	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(result[response.ResultMessage].([]byte), &resmsg)

	if err != nil {
		t.Error(err)
	}

	if resmsg.Ok != true {
		t.Error("unexpected result")
	}

}

func TestSelectDocument(t *testing.T) {

	result, database, err := GetDatabsase(database, conn.GetContext())
	if err != nil {
		t.Log(result)
	}

	docs := []SampleDoc{}

	expr := "{}"

	result, err = database.Select(&docs, expr, nil, nil)
	if err != nil {
		t.Error(err)
	}

	if len(docs) != 2 {

		t.Error("unexpected result")
	}

	docs = []SampleDoc{}

	expr = `{ "_id" : "test_document_id_01" }`

	result, err = database.Select(&docs, expr, nil, nil)
	if err != nil {
		t.Error(err)
	}

	if len(docs) != 1 {

		t.Error("unexpected result")
	}

}

func TestGetSingleDocument(t *testing.T) {

	_, db, err := GetDatabsase(database, conn.GetContext())
	if err != nil {
		t.Error(err)

	}

	db.Get("", nil)
}

func TestCheckDocumentType(t *testing.T) {

	intfslc := make([]interface{}, 0)
	funcslc := make([]testfunc, 0)
	byteslc := make([]byte, 0)
	docs := make([]SampleDoc, 0)

	isValidSlice(docs)
	isValidSlice(&intfslc)
	isValidSlice(&funcslc)
	isValidSlice(&byteslc)
	isValidSlice(&docs)

	document := SampleDoc{ID: "sdasdas", Rev: "f323"}

	id, rev, _ := requiredFields(&document)
	fmt.Println(id, rev)

}
