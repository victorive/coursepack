package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var EmbedMigrations embed.FS

func open(connString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("migrations: open db: %w", err)
	}
	return db, nil
}

func provider(db *sql.DB) (*goose.Provider, error) {
	return goose.NewProvider(goose.DialectPostgres, db, EmbedMigrations)
}

func Up(ctx context.Context, connString string) error {
	db, err := open(connString)
	if err != nil {
		return err
	}
	defer db.Close()

	p, err := provider(db)
	if err != nil {
		return fmt.Errorf("migrations: create provider: %w", err)
	}

	results, err := p.Up(ctx)
	if err != nil {
		return fmt.Errorf("migrations: up: %w", err)
	}
	for _, r := range results {
		fmt.Printf("  migration %s OK (%s)\n", r.Source.Path, r.Duration)
	}
	return nil
}

func UpTo(ctx context.Context, connString string, version int64) error {
	db, err := open(connString)
	if err != nil {
		return err
	}
	defer db.Close()

	p, err := provider(db)
	if err != nil {
		return err
	}

	results, err := p.UpTo(ctx, version)
	if err != nil {
		return fmt.Errorf("migrations: up-to %d: %w", version, err)
	}
	for _, r := range results {
		fmt.Printf("  migration %s OK (%s)\n", r.Source.Path, r.Duration)
	}
	return nil
}

func Down(ctx context.Context, connString string) error {
	db, err := open(connString)
	if err != nil {
		return err
	}
	defer db.Close()

	p, err := provider(db)
	if err != nil {
		return err
	}

	result, err := p.Down(ctx)
	if err != nil {
		return fmt.Errorf("migrations: down: %w", err)
	}
	fmt.Printf("  rolled back %s (%s)\n", result.Source.Path, result.Duration)
	return nil
}

func DownTo(ctx context.Context, connString string, version int64) error {
	db, err := open(connString)
	if err != nil {
		return err
	}
	defer db.Close()

	p, err := provider(db)
	if err != nil {
		return err
	}

	results, err := p.DownTo(ctx, version)
	if err != nil {
		return fmt.Errorf("migrations: down-to %d: %w", version, err)
	}
	for _, r := range results {
		fmt.Printf("  rolled back %s (%s)\n", r.Source.Path, r.Duration)
	}
	return nil
}

func Reset(ctx context.Context, connString string) error {
	return DownTo(ctx, connString, 0)
}

func Status(ctx context.Context, connString string) error {
	db, err := open(connString)
	if err != nil {
		return err
	}
	defer db.Close()

	p, err := provider(db)
	if err != nil {
		return err
	}

	sources, err := p.Status(ctx)
	if err != nil {
		return fmt.Errorf("migrations: status: %w", err)
	}
	for _, s := range sources {
		fmt.Printf("  %-5s  %s\n", s.State, s.Source.Path)
	}
	return nil
}
