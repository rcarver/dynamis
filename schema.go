package dynamis

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
)

// Schema is the interface for managing table schemas.
type Schema interface {

	// Create ensures that the table(s) exist. It should return an error if
	// the table already exists or there is any problem creating it.
	Create(*aws.Config) error

	// Delete ensures that the table(s) do not exist. It should return an
	// error if the table cannot be deleted or if it does not exist.
	Delete(*aws.Config) error
}

// Create makes sure that all of the schemas exist. If abortOnErr is false, it
// iterates through all schemas even if one returns an error. This is generally
// what you want to do since each table schema is independent, and it makes
// calling this function idempotent.
func Create(cfg *aws.Config, schema []Schema, abortOnErr bool) error {
	for _, s := range schema {
		log.Printf("%T Creating...", s)
		if err := s.Create(cfg); err != nil {
			log.Printf("%T Error: %s", s, err)
			if abortOnErr {
				return err
			}
		}
	}
	return nil
}

// Delete makes sure that none of the schemas exist. If abortOnErr is false, it
// continues deleting even if one causes an error. This is generally what you
// want to do, since each table is independent. It also makes this function
// idempotent.
func Delete(cfg *aws.Config, schema []Schema, abortOnErr bool) error {
	for _, s := range schema {
		log.Printf("%T Deleting...", s)
		if err := s.Delete(cfg); err != nil {
			log.Printf("%T Error: %s", s, err)
			if abortOnErr {
				return err
			}
		}
	}
	return nil
}
