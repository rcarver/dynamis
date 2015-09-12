package dynamo

import "github.com/aws/aws-sdk-go/aws"

// Schema is the interface for managing table schemas.
type Schema interface {

	// Create ensures that the table(s) exists, to its current definition.
	Create(*aws.Config) error

	// Delete ensures that the table(s) do not exist.
	Delete(*aws.Config) error
}
