package dynamo

import "github.com/aws/aws-sdk-go/aws"

type Schema interface {
	Create(*aws.Config) error
	Delete(*aws.Config) error
}
