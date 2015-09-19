package dynamo

import (
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TestCheckRowCount(t *testing.T) {
	tests := []struct {
		init func(table) error
		want int
	}{
		{
			// Non-existent table returns -1.
			init: func(table) error {
				return nil
			},
			want: -1,
		},
		{
			// Table with no records returns 0.
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
			// Table with records returns the number of records.
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
			t.Errorf("%d failed init: %s", err)
			continue
		}
		got := CheckRowCount(tbl.db, tbl.name)
		if got != test.want {
			t.Errorf("%d CheckRowCount() got %#v, want %#v", i, got, test.want)
		}
	}
}

func TestCheckRows(t *testing.T) {
	tests := []struct {
		init func(table) error
		defs func(ValueDefiner)
		len  int
		err  bool
		key  string
		val  interface{}
	}{
		{
			// Non-existent table returns nil stuff.
			init: func(table) error {
				return nil
			},
			len: 0,
			err: true,
		},
		{
			// Table with no records.
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
			// Table with records.
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
			// Table with records and parse a custom value.
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
			t.Errorf("%d failed init: %s", err)
			continue
		}
		rows, vd := CheckRows(tbl.db, tbl.name)
		if got, want := len(rows), test.len; got != want {
			t.Errorf("%d CheckRows() len got %#v, want %#v", i, got, want)
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

func TestCheckTable(t *testing.T) {
	tbl := newTable()
	init := func() error {
		_, err := tbl.db.CreateTable(&dynamodb.CreateTableInput{
			TableName: aws.String(tbl.name),
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
		_, err = tbl.db.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(tbl.name),
			Item: map[string]*dynamodb.AttributeValue{
				"str": {
					S: aws.String("hello"),
				},
			},
		})
		return err
	}
	if err := init(); err != nil {
		t.Fatalf("Failed initializing: %s", err)
	}
	table := CheckTable(tbl.db, tbl.name)
	if got, want := table.RowCount(), 1; got != want {
		t.Errorf("RowCount got %d, want %d", got, want)
	}
	rows, vd := table.Rows()
	if got, want := len(rows), 1; got != want {
		t.Errorf("Rows len got %d, want %d", got, want)
	}
	if vd == nil {
		t.Errorf("ValueDefiner is nil")
	}
}
