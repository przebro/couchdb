package database

import (
	"reflect"
	"strings"
)

//isValidSlice - Checks if dcoument is a pointer to slice of a struct
func isValidSlice(doc interface{}) bool {

	v := reflect.ValueOf(doc)

	if v.Kind() == reflect.Ptr {
		if v.Elem().Kind() != reflect.Slice {
			return false
		}

		if reflect.TypeOf(v.Elem().Interface()).Elem().Kind() != reflect.Struct && reflect.TypeOf(v.Elem().Interface()).Elem().Kind() != reflect.Interface {
			return false
		}

		return true
	}

	return false
}

//isSliceOfStructs - checks if given document is a slice of structs
func isSliceOfStructs(docs []interface{}) bool {
	for x := range docs {
		v := reflect.ValueOf(docs[x])
		if v.Kind() != reflect.Struct {
			return false
		}
	}
	return true
}

//isValidStruct - checks if given document is a valid ptr to struct
func isValidStruct(doc interface{}) bool {

	v := reflect.ValueOf(doc)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return false
	}
	return true
}

//requiredFields - extracts required fields from doc
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
