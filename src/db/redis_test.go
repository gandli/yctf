package db

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestNewRedisClient(t *testing.T) {
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		t.Skip("TEST_REDIS_ADDR not set")
	}

	rdb, err := NewRedisClient(addr, "", 0)
	if err != nil {
		t.Fatalf("NewRedisClient failed: %v", err)
	}
	defer rdb.Close()

	if rdb == nil {
		t.Error("redis client should not be nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Errorf("redis ping failed: %v", err)
	}
}

func TestNewRedisClientWithPassword(t *testing.T) {
	addr := os.Getenv("TEST_REDIS_ADDR")
	password := os.Getenv("TEST_REDIS_PASSWORD")
	if addr == "" {
		t.Skip("TEST_REDIS_ADDR not set")
	}

	rdb, err := NewRedisClient(addr, password, 0)
	if err != nil {
		t.Fatalf("NewRedisClient failed: %v", err)
	}
	defer rdb.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Errorf("redis ping failed: %v", err)
	}
}

func TestNewRedisClientInvalidAddr(t *testing.T) {
	_, err := NewRedisClient("invalid-host:9999", "", 0)
	if err == nil {
		t.Error("expected error for invalid redis address")
	}
}

func TestRedisLeaderboard(t *testing.T) {
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		t.Skip("TEST_REDIS_ADDR not set")
	}

	rdb, err := NewRedisClient(addr, "", 0)
	if err != nil {
		t.Fatalf("NewRedisClient failed: %v", err)
	}
	defer rdb.Close()

	ctx := context.Background()
	key := "test:lb:" + time.Now().Format("20060102150405")

	rdb.ZAdd(ctx, key, redis.Z{Score: 100, Member: "team-1"})
	rdb.ZAdd(ctx, key, redis.Z{Score: 200, Member: "team-2"})
	rdb.ZAdd(ctx, key, redis.Z{Score: 150, Member: "team-3"})

	result, err := rdb.ZRevRangeWithScores(ctx, key, 0, 1).Result()
	if err != nil {
		t.Fatalf("ZRevRangeWithScores failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 results, got %d", len(result))
	}

	if result[0].Member != "team-2" || result[0].Score != 200 {
		t.Errorf("expected team-2 with 200, got %v", result[0])
	}

	rdb.Del(ctx, key)
}

func TestRedisRateLimit(t *testing.T) {
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		t.Skip("TEST_REDIS_ADDR not set")
	}

	rdb, err := NewRedisClient(addr, "", 0)
	if err != nil {
		t.Fatalf("NewRedisClient failed: %v", err)
	}
	defer rdb.Close()

	ctx := context.Background()
	key := "test:rl:" + time.Now().Format("20060102150405")

	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		t.Fatalf("Incr failed: %v", err)
	}

	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}

	rdb.Expire(ctx, key, 1*time.Minute)

	ttl, err := rdb.TTL(ctx, key).Result()
	if err != nil {
		t.Fatalf("TTL failed: %v", err)
	}

	if ttl <= 0 {
		t.Error("expected positive TTL")
	}

	rdb.Del(ctx, key)
}
