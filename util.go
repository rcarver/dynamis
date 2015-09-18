package dynamo

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Str returns a string from a DynamoDB attribute value. If anything goes wrong
// reading the value, an empty string is returned.
func Str(item map[string]*dynamodb.AttributeValue, key string) string {
	if item == nil {
		return ""
	}
	if val, ok := item[key]; ok {
		if val.S == nil {
			return ""
		}
		return *val.S
	}
	return ""
}

// Str returns an int from a DynamoDB attribute value. If anything goes wrong
// reading or parsing the value, 0 is returned.
func Int(item map[string]*dynamodb.AttributeValue, key string) int {
	if item == nil {
		return 0
	}
	if val, ok := item[key]; ok {
		if val.N == nil {
			return 0
		}
		if i, err := strconv.Atoi(*val.N); err == nil {
			return i
		}
	}
	return 0
}

// SetStr stores a string attribute. If the string is empty, it is not stored.
func SetStr(item map[string]*dynamodb.AttributeValue, key string, val string) {
	if key != "" && val != "" {
		item[key] = &dynamodb.AttributeValue{
			S: aws.String(val),
		}
	}
}

// SetInt stores an int attribute.
func SetInt(item map[string]*dynamodb.AttributeValue, key string, val int) {
	if key != "" {
		item[key] = &dynamodb.AttributeValue{
			N: aws.String(strconv.FormatInt(int64(val), 10)),
		}
	}
}

// RowCount returns the number of records in a table.
func RowCount(db *dynamodb.DynamoDB, tableName string) int {
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

// Rows returns a Row accessor for every row in a table. The order of rows is
// unspecified. If an error occurs, the rows will be nil. The returned
// ValueDefiner can be used to define access to custom types across all rows.
func Rows(db *dynamodb.DynamoDB, tableName string) ([]Row, ValueDefiner) {
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
