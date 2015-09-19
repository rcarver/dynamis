package dynamis

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TestStr(t *testing.T) {
	tests := []struct {
		item map[string]*dynamodb.AttributeValue
		key  string
		want string
	}{
		{
			// Item is nil, returns the zero value.
			item: nil,
			key:  "k",
			want: "",
		},
		{
			// Key does not exist, returns the zero value.
			item: map[string]*dynamodb.AttributeValue{},
			key:  "k",
			want: "",
		},
		{
			// Key exists, but not a value. Returns the zero value.
			item: map[string]*dynamodb.AttributeValue{"k": {}},
			key:  "k",
			want: "",
		},
		{
			// Key and value exist, returns the value.
			item: map[string]*dynamodb.AttributeValue{"k": {S: aws.String("v")}},
			key:  "k",
			want: "v",
		},
	}
	for i, test := range tests {
		got := Str(test.item, test.key)
		if got != test.want {
			t.Errorf("%d Str() got %#v, want %#v", i, got, test.want)
		}
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		item map[string]*dynamodb.AttributeValue
		key  string
		want int
	}{
		{
			// Item is nil, returns the zero value.
			item: nil,
			key:  "k",
			want: 0,
		},
		{
			// Key does not exist, returns the zero value.
			item: map[string]*dynamodb.AttributeValue{},
			key:  "k",
			want: 0,
		},
		{
			// Key exists with no value, returns the zero value.
			item: map[string]*dynamodb.AttributeValue{"k": {}},
			key:  "k",
			want: 0,
		},
		{
			// Key exists, fails to parse. Returns the zero value.
			item: map[string]*dynamodb.AttributeValue{"k": {N: aws.String("NaN")}},
			key:  "k",
			want: 0,
		},
		{
			// Key exists, is number.
			item: map[string]*dynamodb.AttributeValue{"k": {N: aws.String("33")}},
			key:  "k",
			want: 33,
		},
		{
			// Key exists, is negative number.
			item: map[string]*dynamodb.AttributeValue{"k": {N: aws.String("-1")}},
			key:  "k",
			want: -1,
		},
	}
	for i, test := range tests {
		got := Int(test.item, test.key)
		if got != test.want {
			t.Errorf("%d Int() got %#v, want %#v", i, got, test.want)
		}
	}
}

func TestSetStr(t *testing.T) {
	tests := []struct {
		item map[string]*dynamodb.AttributeValue
		key  string
		val  string
		want map[string]*dynamodb.AttributeValue
	}{
		{
			item: map[string]*dynamodb.AttributeValue{},
			key:  "k",
			val:  "",
			want: map[string]*dynamodb.AttributeValue{},
		},
		{
			item: map[string]*dynamodb.AttributeValue{},
			key:  "",
			val:  "v",
			want: map[string]*dynamodb.AttributeValue{},
		},
		{
			item: map[string]*dynamodb.AttributeValue{},
			key:  "k",
			val:  "v",
			want: map[string]*dynamodb.AttributeValue{"k": {S: aws.String("v")}},
		},
	}
	for i, test := range tests {
		SetStr(test.item, test.key, test.val)
		if !reflect.DeepEqual(test.item, test.want) {
			t.Errorf("%d SetStr() got %#v, want %#v", i, test.item, test.want)
		}
	}
}

func TestSetInt(t *testing.T) {
	tests := []struct {
		item map[string]*dynamodb.AttributeValue
		key  string
		val  int
		want map[string]*dynamodb.AttributeValue
	}{
		{
			item: map[string]*dynamodb.AttributeValue{},
			key:  "k",
			val:  0,
			want: map[string]*dynamodb.AttributeValue{"k": {N: aws.String("0")}},
		},
		{
			item: map[string]*dynamodb.AttributeValue{},
			key:  "",
			val:  33,
			want: map[string]*dynamodb.AttributeValue{},
		},
		{
			item: map[string]*dynamodb.AttributeValue{},
			key:  "k",
			val:  33,
			want: map[string]*dynamodb.AttributeValue{"k": {N: aws.String("33")}},
		},
		{
			item: map[string]*dynamodb.AttributeValue{},
			key:  "k",
			val:  -1,
			want: map[string]*dynamodb.AttributeValue{"k": {N: aws.String("-1")}},
		},
	}
	for i, test := range tests {
		SetInt(test.item, test.key, test.val)
		if !reflect.DeepEqual(test.item, test.want) {
			t.Errorf("%d SetInt() got %#v, want %#v", i, test.item, test.want)
		}
	}
}
