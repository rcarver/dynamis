package dynamis

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
