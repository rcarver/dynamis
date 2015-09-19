package dynamis

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func newDynamoTestConfig() *aws.Config {
	hostport := os.Getenv("DYNAMODB_HOSTPORT")
	if hostport == "" {
		println("DYNAMODB_HOSTPORT is undefined")
		if testing.Short() {
			println("This test should not run when testing.short")
		}
		os.Exit(1)
	}
	endpoint := "http://" + hostport
	cfg := aws.NewConfig().
		WithCredentials(credentials.NewStaticCredentials("aws_id", "aws_secret", "")).
		WithEndpoint(endpoint).
		WithRegion("us-east-1").
		WithLogger(aws.NewDefaultLogger())
	if testing.Verbose() {
		cfg = cfg.WithLogLevel(aws.LogDebugWithHTTPBody)
	}
	return cfg
}

type table struct {
	db   *dynamodb.DynamoDB
	name string
}

func newTable() table {
	var (
		cfg       = newDynamoTestConfig()
		db        = dynamodb.New(cfg)
		tableName = fmt.Sprintf("users-%d", time.Now().UnixNano())
	)
	return table{db, tableName}
}
