package dynamo

import (
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TestValueReader(t *testing.T) {
	tests := []struct {
		item map[string]*dynamodb.AttributeValue
		strK string
		strV string
		intK string
		intV int
		defK string
		defF DefFunc
		defV interface{}
	}{
		{
			// Zero values.
			item: map[string]*dynamodb.AttributeValue{},
			strK: "s",
			strV: "",
			intK: "i",
			intV: 0,
			defK: "d",
			defF: func(vr ValueReader) interface{} { return "ok" },
			defV: "ok",
		},
		{
			// Good values.
			item: map[string]*dynamodb.AttributeValue{
				"s": {S: aws.String("hello")},
				"i": {N: aws.String("33")},
				"d": {S: aws.String("2015-09-16")},
			},
			strK: "s",
			strV: "hello",
			intK: "i",
			intV: 33,
			defK: "d",
			defF: func(vr ValueReader) interface{} {
				t, _ := time.Parse("2006-01-02", vr.Str("d"))
				return t
			},
			defV: time.Date(2015, 9, 16, 0, 0, 0, 0, time.UTC),
		},
	}
	for i, test := range tests {
		r := NewValueReader(test.item)
		if got, want := r.Str(test.strK), test.strV; got != want {
			t.Errorf("%d ValueReader#Str() got %#v, want %#v", i, got, want)
		}
		if got, want := r.Int(test.intK), test.intV; got != want {
			t.Errorf("%d ValueReader#Int() got %#v, want %#v", i, got, want)
		}
		r.Def(test.defK, test.defF)
		if got, want := r.Get(test.defK), test.defV; got != want {
			t.Errorf("%d ValueReader#Int() got %#v, want %#v", i, got, want)
		}
	}
	var didPanic string
	r := NewValueReader(map[string]*dynamodb.AttributeValue{})
	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = r.(string)
			}
		}()
		r.Get("nope")

	}()
	if didPanic != "Missing def for: nope" {
		t.Errorf("Expected Get() to panic, got %s", didPanic)
	}
}

func TestValueWriter(t *testing.T) {
	tests := []struct {
		strK string
		strV string
		intK string
		intV int
		want map[string]*dynamodb.AttributeValue
	}{
		{
			// Zero values.
			strK: "s",
			strV: "",
			intK: "i",
			intV: 0,
			want: map[string]*dynamodb.AttributeValue{
				"i": {N: aws.String("0")},
			},
		},
		{
			// Good values.
			strK: "s",
			strV: "hello",
			intK: "i",
			intV: 33,
			want: map[string]*dynamodb.AttributeValue{
				"s": {S: aws.String("hello")},
				"i": {N: aws.String("33")},
			},
		},
	}
	for i, test := range tests {
		item := make(map[string]*dynamodb.AttributeValue)
		w := NewValueWriter(item)
		w.Str(test.strK, test.strV)
		w.Int(test.intK, test.intV)
		if !reflect.DeepEqual(item, test.want) {
			t.Errorf("%d ValueWriter got %#v, want %#v", i, item, test.want)
		}
	}
}
