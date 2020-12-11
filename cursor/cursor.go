package cursor

//ResultCursor - Helps iterate over returned result
type ResultCursor interface {
	All(v interface{}) error
	Next(v interface{}) error
	Meta() map[string]interface{}
}
