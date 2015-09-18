package dynamo

import (
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
			item: nil,
			key:  "k",
			want: "",
		},
		{
			item: map[string]*dynamodb.AttributeValue{},
			key:  "k",
			want: "",
		},
		{
			item: map[string]*dynamodb.AttributeValue{"k": {}},
			key:  "k",
			want: "",
		},
		{
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
			item: nil,
			key:  "k",
			want: 0,
		},
		{
			item: map[string]*dynamodb.AttributeValue{},
			key:  "k",
			want: 0,
		},
		{
			item: map[string]*dynamodb.AttributeValue{"k": {}},
			key:  "k",
			want: 0,
		},
		{
			item: map[string]*dynamodb.AttributeValue{"k": {N: aws.String("0")}},
			key:  "k",
			want: 0,
		},
		{
			item: map[string]*dynamodb.AttributeValue{"k": {N: aws.String("NaN")}},
			key:  "k",
			want: 0,
		},
		{
			item: map[string]*dynamodb.AttributeValue{"k": {N: aws.String("33")}},
			key:  "k",
			want: 33,
		},
	}
	for i, test := range tests {
		got := Int(test.item, test.key)
		if got != test.want {
			t.Errorf("%d Int() got %#v, want %#v", i, got, test.want)
		}
	}
}
