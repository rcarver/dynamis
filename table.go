package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Table struct {
	db        *dynamodb.DynamoDB
	tableName string
}

func NewTable(cfg *aws.Config, tableName string) Table {
	return Table{dynamodb.New(cfg), tableName}
}

func (t Table) RowCount() int {
	return RowCount(t.db, t.tableName)
}

func (t Table) Rows() ([]Row, CustomValueReader) {
	return Rows(t.db, t.tableName)
}
