package dynamo

import (
	"reflect"
	"testing"
	"time"

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

func TestRowCount(t *testing.T) {
	tests := []struct {
		init func(table) error
		want int
	}{
		{
			init: func(table) error {
				return nil
			},
			want: -1,
		},
		{
			init: func(t table) error {
				_, err := t.db.CreateTable(&dynamodb.CreateTableInput{
					TableName: aws.String(t.name),
					AttributeDefinitions: []*dynamodb.AttributeDefinition{
						{
							AttributeName: aws.String("str"),
							AttributeType: aws.String("S"),
						},
					},
					KeySchema: []*dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("str"),
							KeyType:       aws.String("HASH"),
						},
					},
					ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(1),
						WriteCapacityUnits: aws.Int64(1),
					},
				})
				return err
			},
			want: 0,
		},
		{
			init: func(t table) error {
				var err error
				_, err = t.db.CreateTable(&dynamodb.CreateTableInput{
					TableName: aws.String(t.name),
					AttributeDefinitions: []*dynamodb.AttributeDefinition{
						{
							AttributeName: aws.String("str"),
							AttributeType: aws.String("S"),
						},
					},
					KeySchema: []*dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("str"),
							KeyType:       aws.String("HASH"),
						},
					},
					ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(1),
						WriteCapacityUnits: aws.Int64(1),
					},
				})
				if err != nil {
					return err
				}
				_, err = t.db.PutItem(&dynamodb.PutItemInput{
					TableName: aws.String(t.name),
					Item: map[string]*dynamodb.AttributeValue{
						"str": {
							S: aws.String("one"),
						},
					},
				})
				return err

			},
			want: 1,
		},
	}
	for i, test := range tests {
		tbl := newTable()
		if err := test.init(tbl); err != nil {
			t.Error("%d failed init: %s", err)
			continue
		}
		got := RowCount(tbl.db, tbl.name)
		if got != test.want {
			t.Errorf("%d RowCount() got %#v, want %#v", i, got, test.want)
		}
	}
}

func TestRows(t *testing.T) {
	tests := []struct {
		init func(table) error
		defs func(ValueDefiner)
		len  int
		err  bool
		key  string
		val  interface{}
	}{
		{
			init: func(table) error {
				return nil
			},
			len: 0,
			err: true,
		},
		{
			init: func(t table) error {
				_, err := t.db.CreateTable(&dynamodb.CreateTableInput{
					TableName: aws.String(t.name),
					AttributeDefinitions: []*dynamodb.AttributeDefinition{
						{
							AttributeName: aws.String("str"),
							AttributeType: aws.String("S"),
						},
					},
					KeySchema: []*dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("str"),
							KeyType:       aws.String("HASH"),
						},
					},
					ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(1),
						WriteCapacityUnits: aws.Int64(1),
					},
				})
				return err
			},
			len: 0,
		},
		{
			init: func(t table) error {
				var err error
				_, err = t.db.CreateTable(&dynamodb.CreateTableInput{
					TableName: aws.String(t.name),
					AttributeDefinitions: []*dynamodb.AttributeDefinition{
						{
							AttributeName: aws.String("str"),
							AttributeType: aws.String("S"),
						},
					},
					KeySchema: []*dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("str"),
							KeyType:       aws.String("HASH"),
						},
					},
					ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(1),
						WriteCapacityUnits: aws.Int64(1),
					},
				})
				if err != nil {
					return err
				}
				_, err = t.db.PutItem(&dynamodb.PutItemInput{
					TableName: aws.String(t.name),
					Item: map[string]*dynamodb.AttributeValue{
						"str": {
							S: aws.String("one"),
						},
					},
				})
				return err

			},
			len: 1,
			key: "str",
			val: "one",
		},
		{
			init: func(t table) error {
				var err error
				_, err = t.db.CreateTable(&dynamodb.CreateTableInput{
					TableName: aws.String(t.name),
					AttributeDefinitions: []*dynamodb.AttributeDefinition{
						{
							AttributeName: aws.String("date"),
							AttributeType: aws.String("S"),
						},
					},
					KeySchema: []*dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("date"),
							KeyType:       aws.String("HASH"),
						},
					},
					ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(1),
						WriteCapacityUnits: aws.Int64(1),
					},
				})
				if err != nil {
					return err
				}
				_, err = t.db.PutItem(&dynamodb.PutItemInput{
					TableName: aws.String(t.name),
					Item: map[string]*dynamodb.AttributeValue{
						"date": {
							S: aws.String("2015-09-16"),
						},
					},
				})
				return err

			},
			defs: func(vd ValueDefiner) {
				vd.Def("parsedate", func(r ValueReader) interface{} {
					t, _ := time.Parse("2006-01-02", r.Str("date"))
					return t
				})
			},
			len: 1,
			key: "parsedate",
			val: time.Date(2015, 9, 16, 0, 0, 0, 0, time.UTC),
		},
	}
	for i, test := range tests {
		tbl := newTable()
		if err := test.init(tbl); err != nil {
			t.Error("%d failed init: %s", err)
			continue
		}
		rows, vd := Rows(tbl.db, tbl.name)
		if got, want := len(rows), test.len; got != want {
			t.Errorf("%d Rows() len got %#v, want %#v", i, got, want)
		}
		if test.err {
			if rows != nil {
				t.Errorf("%d want rows to be nil, got %#v", i, rows)
			}
			if vd != nil {
				t.Errorf("%d want vd to be nil, got %#v", i, vd)
			}
		}
		if len(rows) > 0 {
			if test.defs != nil {
				test.defs(vd)
				if got, want := rows[0].Get(test.key), test.val; !reflect.DeepEqual(got, want) {
					t.Errorf("%d RowCount() val got %#v, want %#v", i, got, want)
				}
			} else {
				if got, want := rows[0].Str(test.key), test.val; !reflect.DeepEqual(got, want) {
					t.Errorf("%d RowCount() val got %#v, want %#v", i, got, want)
				}
			}
		}
	}
}
