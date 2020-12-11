package database

import (
	"encoding/json"
	"fmt"
	"io"
)

type selectCursor struct {
	dec  *json.Decoder
	rdr  io.ReadCloser
	docs bool
}

func (s *selectCursor) Next(v interface{}) error {

	if s.dec == nil {
		s.dec = json.NewDecoder(s.rdr)
		s.dec.Token()
		s.dec.Token()
		s.dec.Token()
	}

	if !s.dec.More() {
		s.docs = true
		return io.EOF
	}

	return s.dec.Decode(v)
}

func (s *selectCursor) All(v interface{}) error {

	if !isValidSlice(v) {
		return errInvalidDocKind
	}

	if s.dec == nil {
		s.dec = json.NewDecoder(s.rdr)
		s.dec.Token()
		s.dec.Token()
	}

	return s.dec.Decode(v)
}
func (s *selectCursor) Meta() map[string]interface{} {

	fmt.Println(s.dec.Token())
	fmt.Println(s.dec.Token())
	if s.docs == true {
		meta := map[string]interface{}{}
		s.dec.Decode(&meta)
		return meta
	}
	return nil
}

func (s *selectCursor) postionDecoder() {

}
