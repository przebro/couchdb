package database

import (
	"context"
	"fmt"
	"testing"

	"github.com/przebro/couchdb/client"
	"github.com/przebro/couchdb/connection"
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

type TestDocument struct {
	ID     string  `json:"_id,omitempty"`
	REV    string  `json:"_rev,omitempty"`
	Title  string  `json:"title"`
	Score  float32 `json:"score"`
	Year   int     `json:"year"`
	Oscars bool    `json:"oscars"`
}

type InsertResult struct {
	ID  string `json:"id,omitempty"`
	OK  bool   `json:"ok,omitempty"`
	REV string `json:"rev,omitempty"`
}

var collectionInterface []interface{} = []interface{}{
	TestDocument{ID: "movie_1", Title: "The Godfather", Score: 9.2, Year: 1972, Oscars: true},
	TestDocument{ID: "movie_2", Title: "The Big Lebowski", Score: 8.1, Year: 1998, Oscars: false},
	TestDocument{ID: "movie_3", Title: "Terminator 2:Judgment Day", Score: 8.5, Year: 1991, Oscars: true},
	TestDocument{ID: "movie_4", Title: "The Shining", Score: 8.4, Year: 1980, Oscars: false},
	TestDocument{ID: "movie_5", Title: "Star Wars: Episode V - The Empire Strikes Back", Score: 8.7, Year: 1980, Oscars: true},
	TestDocument{ID: "movie_6", Title: "The Thing", Score: 8.1, Year: 1982, Oscars: false},
	TestDocument{ID: "movie_7", Title: "Platoon", Score: 8.1, Year: 1986, Oscars: true},
	TestDocument{ID: "movie_8", Title: "Der Name der Rose", Score: 7.7, Year: 1986, Oscars: false},
	TestDocument{ID: "another_unique_id", Title: "Predator", Score: 7.8, Year: 1987, Oscars: false},
}

type testfunc func()

var builder connection.ConnectionBuilder
var conn *connection.Connection

func init() {
	var err error
	builder = connection.NewBuilder()
	conn, err = builder.WithAuthentication(client.Basic, username, password).WithAddress(host, port).Build(true)
	if err != nil {
		fmt.Println("error")
	}

}

func TestCreateDatabase(t *testing.T) {

	result, db, err := CreateDatabase(context.Background(), database, conn.GetClient())
	if err != nil {
		t.Error("unexpected result")
		t.FailNow()
	}

	if result.Code != 201 {
		t.Error("unexpected result")
	}
	if db.Name != database {
		t.Error("Unexpected result")
	}
	result, db, err = CreateDatabase(context.Background(), database, conn.GetClient())
	if err == nil {
		t.Error("unexpected result")
	}

	if result.Code != 412 {
		t.Error("unexpected result")
	}

}

func TestGetDatabase(t *testing.T) {

	_, _, err := GetDatabsase(context.Background(), "_users", conn.GetClient())
	if err != nil {
		t.Error("unexpected result")
	}

	r, _, err := GetDatabsase(context.Background(), "db_that_does_not_exists", conn.GetClient())

	if err == nil {
		t.Error("Unexpected result")
	}

	if r.Code != 404 {
		t.Error("Unexpected result")
	}

}

func TestInsertDocument(t *testing.T) {

	resmsg := struct {
		Ok  bool   `json:"ok"`
		ID  string `json:"id"`
		Rev string `json:"rev"`
	}{}

	result, database, err := GetDatabsase(context.Background(), database, conn.GetClient())
	if err != nil {
		t.Error(result)
	}

	withoutID := SampleDoc{Name: "User Name", Age: 21, Group: "group_1"}
	rd, err := database.Insert(context.Background(), &withoutID)

	if err != nil {
		t.Error(err)
	}

	if rd.Code != 201 {
		t.Error("unexpected result")
	}

	rd, err = database.Insert(context.Background(), withoutID)

	if err != errInvalidDocKind {
		t.Error("Unexpected result:", err)
	}

	sample := SampleDoc{ID: "test_document_id_01", Name: "User Name", Age: 21, Group: "group_1"}
	result, err = database.Insert(context.Background(), &sample)
	if err != nil {
		t.Error(err)
	}

	sample = SampleDoc{ID: "test_document_id_02", Name: "User Name", Age: 21, Group: "group_1"}
	rd, err = database.Insert(context.Background(), &sample)

	if err != nil {
		t.Error(err)
	}
	rd.Decode(&resmsg)

	if resmsg.Ok != true {
		t.Error("unexpected result")
	}

}

func TestRevision(t *testing.T) {

	_, database, err := GetDatabsase(context.Background(), database, conn.GetClient())
	if err != nil {
		t.Error(err)
	}

	result, err := database.Revision(context.Background(), "test_document_id_01")
	if err != nil {
		t.Error(err)
	}
	if result.Code != 200 {
		t.Error("unexpected stastus code")
	}

	rr := map[string]interface{}{}

	errx := result.Decode(&rr)
	fmt.Println("RRL:", rr, errx)

	result, err = database.Revision(context.Background(), "test_document_id_03")

	if err != nil {
		t.Error("unexpected result")
	}

	if result.Code != 404 {
		t.Error("unexpected stastus code")
	}

}

func TestSelectDocument(t *testing.T) {

	r, database, err := GetDatabsase(context.Background(), database, conn.GetClient())
	if err != nil {
		t.Log(r)
	}

	expr := "{}"

	result, err := database.Select(context.Background(), "", nil, nil)

	if err == nil {
		t.Error("Unexpected result")
	}

	result, err = database.Select(context.Background(), expr, nil, nil)
	if err != nil {
		t.Error(err)
	}

	if result.Code != 200 {
		t.Error("Unexpected status")
	}

	docs := []TestDocument{}

	result.All(context.Background(), &docs)

	if len(docs) == 1 {
		t.Error("unexpected result")
	}

	expr = `{"_id": "test_document_id_01"}`
	docs = []TestDocument{}

	result, err = database.Select(context.Background(), expr, nil, map[FindOption]interface{}{OptionStat: true})
	if err != nil {
		t.Error(err)
	}

	result.All(context.Background(), &docs)

	if len(docs) != 1 {
		t.Error("unexpected result, actual:", len(docs))
	}
}

func TestGetSingleDocument(t *testing.T) {

	_, db, err := GetDatabsase(context.Background(), database, conn.GetClient())
	if err != nil {
		t.Error(err)

	}

	rd, err := db.Get(context.Background(), "")

	if err != errEmptyDocumentID {
		t.Error("unexpected result")
	}

	doc := SampleDoc{}

	rd, err = db.Get(context.Background(), "test_document_id_01")
	if err != nil {
		t.Error(err)
	}

	err = rd.Decode(&doc)
	if err != nil {
		t.Error("unexpected result")
	}

}

func TestUpdateDocument(t *testing.T) {

	_, db, err := GetDatabsase(context.Background(), database, conn.GetClient())
	if err != nil {
		t.Error(err)

	}
	doc := SampleDoc{}

	rd, err := db.Get(context.Background(), "test_document_id_01")
	if err != nil {
		t.Error(err)
	}

	rd.Decode(&doc)
	doc.Group = "Abcdef"

	r, err := db.Update(context.Background(), &doc)
	if err != nil {
		t.Error(err)
	}

	if r.Code != 201 {
		t.Error("unexpected result")
	}

	astruct := struct {
		Name string `json:"_id"`
		Rev  string `json:"_rev"`
	}{Name: "name", Rev: ""}

	notstruct := "test string"

	r, err = db.Update(context.Background(), &notstruct)
	if err != errInvalidDocKind {
		t.Error("Unexpected result", err)
	}

	r, err = db.Update(context.Background(), &astruct)
	if err != errIDandRevRequired {
		t.Error("unexpected result")
	}

}

func TestDeleteDocument(t *testing.T) {

	_, db, err := GetDatabsase(context.Background(), database, conn.GetClient())
	if err != nil {
		t.Error(err)
	}

	resmsg := struct {
		Ok  bool   `json:"ok"`
		ID  string `json:"id"`
		Rev string `json:"rev"`
	}{}

	sample := SampleDoc{ID: "test_document_id_05", Name: "User Name", Age: 21, Group: "group_1"}
	result, err := db.Insert(context.Background(), &sample)
	if err != nil {
		t.Error(err)
	}

	result.Decode(&resmsg)
	result, err = db.Delete(context.Background(), "test_document_id_05", "")

	if err != errIDandRevRequired {
		t.Error("Unexpected result", err)
	}

	result, err = db.Delete(context.Background(), "test_document_id_05", resmsg.Rev)

	if err != nil {
		t.Error("Unexpected result", err)
	}

}

func TestPurgeDocument(t *testing.T) {

	documentID := "test_document_id_13"
	_, db, err := GetDatabsase(context.Background(), database, conn.GetClient())
	if err != nil {
		t.Error(err)
	}

	resmsg := struct {
		Ok  bool   `json:"ok"`
		ID  string `json:"id"`
		Rev string `json:"rev"`
	}{}

	sample := SampleDoc{ID: documentID, Name: "User Name", Age: 21, Group: "group_1"}
	result, err := db.Insert(context.Background(), &sample)
	if err != nil {
		t.Error(err)
	}

	result.Decode(&resmsg)

	res, err := db.Purge(context.Background(), "", nil)
	if err != errEmptyDocumentID {
		t.Error("unexpected result:", err)
	}

	res, err = db.Purge(context.Background(), documentID, nil)

	if err != errRevListRequired {
		t.Error("unexpected result:", err)
	}

	res, err = db.Purge(context.Background(), documentID, []string{})

	if err != errRevListRequired {
		t.Error("unexpected result:", err)
	}

	res, err = db.Purge(context.Background(), documentID, []string{resmsg.Rev})

	if err != nil {
		t.Error("unexpected result")
	}

	fmt.Println(res.Code)

	resp := map[string]interface{}{}
	res.Decode(&resp)
	fmt.Println(resp)

}

func TestCheckDocumentType(t *testing.T) {

	intfslc := make([]interface{}, 0)
	funcslc := make([]testfunc, 0)
	byteslc := make([]byte, 0)
	docs := make([]SampleDoc, 0)
	sdocs := SampleDoc{}

	result := isValidSlice(intfslc)
	if result == true {
		t.Error("Unexpected result")
	}

	result = isValidSlice(&intfslc)
	if result != true {
		t.Error("Unexpected result")
	}

	result = isValidSlice(&funcslc)

	if result == true {
		t.Error("Unexpected result")
	}

	result = isValidSlice(&byteslc)

	if result == true {
		t.Error("Unexpected result")
	}

	result = isValidSlice(&docs)

	if result != true {
		t.Error("Unexpected result")
	}

	result = isValidSlice(docs)

	if result == true {
		t.Error("Unexpected result")
	}

	result = isValidSlice(sdocs)

	if result == true {
		t.Error("Unexpected result")
	}

	result = isValidSlice(&sdocs)

	if result == true {
		t.Error("Unexpected result")
	}

}

func TestInsertMany(t *testing.T) {
	_, db, err := GetDatabsase(context.Background(), database, conn.GetClient())
	if err != nil {
		t.Error(err)
	}

	r, err := db.InsertMany(context.Background(), []interface{}{[]string{"abc", "def"}})

	if err != errInvalidDocKind {
		t.Error("unexpected result:", err)
	}

	r, err = db.InsertMany(context.Background(), collectionInterface)
	if err != nil {
		t.Error(err)
	}

	ir := []InsertResult{}
	r.Decode(&ir)

	fmt.Println(ir)

}
func TestStat(t *testing.T) {

	_, db, err := GetDatabsase(context.Background(), database, conn.GetClient())
	if err != nil {
		t.Error(err)
	}

	r, err := db.Stat(context.Background())
	if err != nil {
		t.Error(err)
	}

	val := map[string]interface{}{}

	err = r.Decode(&val)
	if err != nil {
		t.Error(err)
	}

}

func TestCopy(t *testing.T) {

	_, db, err := GetDatabsase(context.Background(), database, conn.GetClient())
	if err != nil {
		t.Error(err)

	}

	r, err := db.Copy(context.Background(), "", "", "")
	if err != errRequiredDocumentID {
		t.Error("Unexpected result:", err)
	}

	r, err = db.Copy(context.Background(), "test_document_id_01", "", "")
	if err != errRequiredestinationID {
		t.Error("Unexpected result:", err)
	}

	r, err = db.Copy(context.Background(), "test_document_id_01", "test_document_id_11", "")
	fmt.Println(err)
	if r != nil {
		mm := map[string]interface{}{}
		r.Decode(&mm)
		fmt.Println(mm)
	}

}

func TestSelectIterator(t *testing.T) {
	_, db, err := GetDatabsase(context.Background(), database, conn.GetClient())
	if err != nil {
		t.Error(err)
	}

	r, err := db.Select(context.Background(), "{}", nil, map[FindOption]interface{}{OptionStat: true})
	if err != nil {
		t.Error(err)
	}

	cnt := 0
	result := r.Meta()

	for r.Next(context.Background()) {
		td := TestDocument{}
		err := r.Decode(&td)
		if err != nil {
			break
		}
		cnt++
	}

	if result.Documents != cnt {
		t.Error("unexpected result, expected:", result.Documents, "actual:", cnt)
	}
}
