package dynamis

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// ValueReader defines access to different types of value.
type ValueReader interface {
	Str(key string) string
	Int(key string) int
	Get(key string) interface{}
	ValueDefiner
}

// ValueDefiner lets you define accessors for custom types.
type ValueDefiner interface {
	Def(key string, f DefFunc)
}

// DefFunc is the handler for custom types.
type DefFunc func(ValueReader) interface{}

type valueDefiner struct {
	defs map[string]DefFunc
}

func newValueDefiner() valueDefiner {
	return valueDefiner{make(map[string]DefFunc)}
}

func (r valueDefiner) Def(key string, f DefFunc) {
	r.defs[key] = f
}
func (r valueDefiner) call(key string, vr ValueReader) interface{} {
	if f, ok := r.defs[key]; ok {
		return f(vr)
	}
	panic(fmt.Sprintf("Missing def for: %s", key))
}

type valueReader struct {
	item map[string]*dynamodb.AttributeValue
	def  valueDefiner
}

func (r valueReader) Str(key string) string {
	return Str(r.item, key)
}
func (r valueReader) Int(key string) int {
	return Int(r.item, key)
}
func (r valueReader) Def(key string, f DefFunc) {
	r.def.Def(key, f)
}
func (r valueReader) Get(key string) interface{} {
	return r.def.call(key, r)
}

// NewValueReader initializes a ValueReader over an item.
func NewValueReader(item map[string]*dynamodb.AttributeValue) ValueReader {
	return valueReader{item, newValueDefiner()}
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
