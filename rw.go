package dynamo

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// ValueReader is the interface for anything that can read values from a
// dynamodb row.
type ValueReader interface {
	StdValueReader
	CustomValueReader
}

// StdValueReader defines access for builtin types.
type StdValueReader interface {
	Str(key string) string
	Int(key string) int
}

// CustomValueReader lets you define and get custom types.
type CustomValueReader interface {
	Def(key string, f func() interface{})
	Get(key string) interface{}
}

type customValueReader struct {
	gen map[string]func() interface{}
}

func newCustomReader() customValueReader {
	return customValueReader{make(map[string]func() interface{})}
}

func (r customValueReader) Def(key string, f func() interface{}) {
	r.gen[key] = f
}
func (r customValueReader) Get(key string) interface{} {
	if f, ok := r.gen[key]; ok {
		return f()
	}
	return fmt.Errorf("Missing def for: %s", key)
}

type stdValueReader struct {
	item map[string]*dynamodb.AttributeValue
}

func (r stdValueReader) Str(key string) string {
	return Str(r.item, key)
}
func (r stdValueReader) Int(key string) int {
	return Int(r.item, key)
}

type valueReader struct {
	stdValueReader
	customValueReader
}

// NewValueReader initializes a ValueReader over an item.
func NewValueReader(item map[string]*dynamodb.AttributeValue) ValueReader {
	return valueReader{stdValueReader{item}, customValueReader{}}
}

// ValueWriter lets you easily set values in an DynamoDB AttributeValue map.
type ValueWriter struct {
	item map[string]*dynamodb.AttributeValue
}

// NewValueWriter initializes a ValueWriter over an item.
func NewValueWriter(item map[string]*dynamodb.AttributeValue) ValueWriter {
	return ValueWriter{item}
}

func (w ValueWriter) Str(key string, val string) {
	SetStr(w.item, key, val)
}
func (w ValueWriter) Int(key string, val int) {
	SetInt(w.item, key, val)
}
