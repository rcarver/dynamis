package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// CheckRowCount returns the number of records in a table.
func CheckRowCount(db *dynamodb.DynamoDB, tableName string) int {
	resp, err := db.Scan(&dynamodb.ScanInput{
		TableName: aws.String(tableName),
		Select:    aws.String(dynamodb.SelectCount),
	})
	if err != nil {
		return -1
	}
	return int(*resp.Count)
}

// Row is a generic accessor for any dynamodb row. It implements ValueReader,
// so all of the convenient methods defined there are available.
type Row struct {
	ValueReader
}

// CheckRows returns a Row accessor for every row in a table. The order of rows
// is unspecified. If an error occurs, the rows will be nil. The returned
// ValueDefiner can be used to define access to custom types across all rows.
func CheckRows(db *dynamodb.DynamoDB, tableName string) ([]Row, ValueDefiner) {
	resp, err := db.Scan(&dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, nil
	}
	vd := newValueDefiner()
	rows := make([]Row, len(resp.Items))
	for i, item := range resp.Items {
		reader := valueReader{item, vd}
		rows[i] = Row{reader}
	}
	return rows, vd
}

// Table is a convient wrapper over any DynamoDB table.
type Table struct {
	db        *dynamodb.DynamoDB
	tableName string
}

// CheckTable initializes a new wrapper over the table.
func CheckTable(cfg *aws.Config, tableName string) Table {
	return Table{dynamodb.New(cfg), tableName}
}

// RowCount returns the number of rows in the table.
func (t Table) RowCount() int {
	return CheckRowCount(t.db, t.tableName)
}

// Rows returns a simple accessor for each row in the table. The rows are
// returned in no particular order. The returned ValueDefiner can be used to
// initialize access to complex values.
func (t Table) Rows() ([]Row, ValueDefiner) {
	return CheckRows(t.db, t.tableName)
}
