package dynamis

import (
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
)

type fakeSchema struct {
	fail       bool
	wantCalled bool
	called     bool
}

func newFakeSchema(fail bool, wantCall bool) *fakeSchema {
	return &fakeSchema{fail, wantCall, false}
}

func (s *fakeSchema) Create(cfg *aws.Config) error {
	s.called = true
	if s.fail {
		return errors.New("failed")
	}
	return nil
}
func (s *fakeSchema) Delete(cfg *aws.Config) error {
	s.called = true
	if s.fail {
		return errors.New("failed")
	}
	return nil
}

func TestCreate(t *testing.T) {
	tests := []struct {
		schemes []Schema
		abort   bool
		want    error
	}{
		{
			// All succeed.
			schemes: []Schema{
				newFakeSchema(false, true),
				newFakeSchema(false, true),
			},
			abort: true,
			want:  nil,
		},
		{
			// First fails and abort on error
			schemes: []Schema{
				newFakeSchema(true, true),
				newFakeSchema(false, false),
			},
			abort: true,
			want:  errors.New("failed"),
		},
		{
			// First fails, do not abort.
			schemes: []Schema{
				newFakeSchema(true, true),
				newFakeSchema(false, true),
			},
			abort: false,
			want:  nil,
		},
	}
	for i, test := range tests {
		got := Create(nil, test.schemes, test.abort)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%d Create() got %#v, want %#v", i, got, test.want)
		}
		for j, s := range test.schemes {
			if fs, ok := s.(*fakeSchema); ok {
				if fs.wantCalled != fs.called {
					t.Errorf("%d/%d Create called got %#v, want %#v", i, j, fs.called, fs.wantCalled)
				}
			}
		}
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		schemes []Schema
		abort   bool
		want    error
	}{
		{
			// All succeed.
			schemes: []Schema{
				newFakeSchema(false, true),
				newFakeSchema(false, true),
			},
			abort: true,
			want:  nil,
		},
		{
			// First fails and abort on error
			schemes: []Schema{
				newFakeSchema(true, true),
				newFakeSchema(false, false),
			},
			abort: true,
			want:  errors.New("failed"),
		},
		{
			// First fails, do not abort.
			schemes: []Schema{
				newFakeSchema(true, true),
				newFakeSchema(false, true),
			},
			abort: false,
			want:  nil,
		},
	}
	for i, test := range tests {
		got := Delete(nil, test.schemes, test.abort)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%d Delete() got %#v, want %#v", i, got, test.want)
		}
		for j, s := range test.schemes {
			if fs, ok := s.(*fakeSchema); ok {
				if fs.wantCalled != fs.called {
					t.Errorf("%d/%d Delete called got %#v, want %#v", i, j, fs.called, fs.wantCalled)
				}
			}
		}
	}
}
