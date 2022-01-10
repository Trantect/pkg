package objectstorage

import (
	"context"
	"io"
)

// ObjectHandler is an interface for storing and retrieving objects.
type ObjectHandler interface {
	// PutObject stores an object.
	PutObject(ctx context.Context, key string, r io.Reader) error
	// GetObject retrieves an object.
	GetObject(ctx context.Context, key string) (io.ReadCloser, error)
	// DeleteObject deletes an object.
	DeleteObject(ctx context.Context, key string) error
}
