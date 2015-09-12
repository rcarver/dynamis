package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Table is a convient wrapper over any DynamoDB table.
type Table struct {
	db        *dynamodb.DynamoDB
	tableName string
}

// NewTable initializes a new wrapper over the table.
func NewTable(cfg *aws.Config, tableName string) Table {
	return Table{dynamodb.New(cfg), tableName}
}

// RowCount returns the number of rows in the table.
func (t Table) RowCount() int {
	return RowCount(t.db, t.tableName)
}

// Rows returns a simple accessor for each row in the table. The rows are
// returned in no particular order. The returned ValueDefiner can be used to
// initialize access to complex values.
func (t Table) Rows() ([]Row, ValueDefiner) {
	return Rows(t.db, t.tableName)
}
