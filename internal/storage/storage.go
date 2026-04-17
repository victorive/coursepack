package storage

import (
	"context"
	"io"
	"time"
)

type Storage interface {
	Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (string, error)

	Download(ctx context.Context, key string) (io.ReadCloser, error)

	Delete(ctx context.Context, key string) error

	Exists(ctx context.Context, key string) (bool, error)

	GetSignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)

	GetSize(ctx context.Context, key string) (int64, error)

	GetContentType(ctx context.Context, key string) (string, error)
}
