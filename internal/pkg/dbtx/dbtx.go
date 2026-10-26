// Package dbtx carries a database transaction on the context so repositories
// take part in it without threading *gorm.DB through every interface.
package dbtx

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

// WithTx runs fn in a transaction and stores it on the context. If ctx already
// carries a transaction, fn joins it instead of starting a new one.
func WithTx(ctx context.Context, db *gorm.DB, fn func(context.Context) error) error {
	if From(ctx) != nil {
		return fn(ctx)
	}
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, txKey{}, tx))
	})
}

// From returns the transaction stored on ctx, or nil when there is none.
func From(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	if !ok {
		return nil
	}
	return tx
}

// Conn returns the transaction from ctx, falling back to db bound to ctx.
func Conn(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx := From(ctx); tx != nil {
		return tx
	}
	return db.WithContext(ctx)
}
