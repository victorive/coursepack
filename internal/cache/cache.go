package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)

	Set(ctx context.Context, key string, value any, ttl time.Duration) error

	SetNX(ctx context.Context, key string, value any, ttl time.Duration) (bool, error)

	Delete(ctx context.Context, keys ...string) error

	Exists(ctx context.Context, key string) (bool, error)

	GetJSON(ctx context.Context, key string, dest any) error

	SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error

	Close() error
}
