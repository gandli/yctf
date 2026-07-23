package db

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestNewPostgresPool(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}

	pool, err := NewPostgresPool(dsn)
	if err != nil {
		t.Fatalf("NewPostgresPool failed: %v", err)
	}
	defer pool.Close()

	if pool == nil {
		t.Error("pool should not be nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		t.Errorf("pool.Ping failed: %v", err)
	}
}

func TestNewPostgresPoolInvalidDSN(t *testing.T) {
	_, err := NewPostgresPool("://malformed-url")
	if err == nil {
		t.Error("expected error for malformed DSN")
	}
}

func TestNewPostgresPoolContext(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := NewPostgresPoolContext(ctx, dsn)
	if err != nil {
		t.Fatalf("NewPostgresPoolContext failed: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		t.Errorf("pool.Ping failed: %v", err)
	}
}
