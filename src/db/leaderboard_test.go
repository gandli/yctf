package db

import (
	"context"
	"strconv"
	"testing"
	"time"
)

func TestNewLeaderboard(t *testing.T) {
	rdb, err := NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer rdb.Close()

	lb := NewLeaderboard(rdb)
	if lb == nil {
		t.Error("leaderboard should not be nil")
	}
}

func TestLeaderboardAddScore(t *testing.T) {
	rdb, err := NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer rdb.Close()

	ctx := context.Background()
	key := "test:lb:" + time.Now().Format("20060102150405")
	defer rdb.Del(ctx, key)

	lb := NewLeaderboard(rdb)

	err = lb.AddScore(ctx, key, "team-1", 100)
	if err != nil {
		t.Errorf("AddScore failed: %v", err)
	}

	err = lb.AddScore(ctx, key, "team-2", 200)
	if err != nil {
		t.Errorf("AddScore failed: %v", err)
	}

	top, err := lb.GetTop(ctx, key, 2)
	if err != nil {
		t.Errorf("GetTop failed: %v", err)
	}

	if len(top) != 2 {
		t.Errorf("expected 2 teams, got %d", len(top))
	}

	if top[0].Member != "team-2" || top[0].Score != 200 {
		t.Errorf("expected team-2 with 200, got %v", top[0])
	}
}

func TestLeaderboardGetTop(t *testing.T) {
	rdb, err := NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer rdb.Close()

	ctx := context.Background()
	key := "test:lb:top:" + time.Now().Format("20060102150405")
	defer rdb.Del(ctx, key)

	lb := NewLeaderboard(rdb)

	for i := 1; i <= 5; i++ {
		lb.AddScore(ctx, key, "team-"+strconv.Itoa(i), int64(i*100))
	}

	top, err := lb.GetTop(ctx, key, 3)
	if err != nil {
		t.Errorf("GetTop failed: %v", err)
	}

	if len(top) != 3 {
		t.Errorf("expected 3, got %d", len(top))
	}

	if top[0].Member != "team-5" || top[0].Score != 500 {
		t.Errorf("expected team-5 with 500, got %v", top[0])
	}
}

func TestLeaderboardGetRank(t *testing.T) {
	rdb, err := NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer rdb.Close()

	ctx := context.Background()
	key := "test:lb:rank:" + time.Now().Format("20060102150405")
	defer rdb.Del(ctx, key)

	lb := NewLeaderboard(rdb)

	lb.AddScore(ctx, key, "team-1", 100)
	lb.AddScore(ctx, key, "team-2", 200)
	lb.AddScore(ctx, key, "team-3", 300)

	rank2, err := lb.GetRank(ctx, key, "team-2")
	if err != nil {
		t.Errorf("GetRank failed: %v", err)
	}

	if rank2 != 2 {
		t.Errorf("expected rank 2, got %d", rank2)
	}
}

func TestLeaderboardGetScore(t *testing.T) {
	rdb, err := NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer rdb.Close()

	ctx := context.Background()
	key := "test:lb:score:" + time.Now().Format("20060102150405")
	defer rdb.Del(ctx, key)

	lb := NewLeaderboard(rdb)

	lb.AddScore(ctx, key, "team-1", 150)

	score, err := lb.GetScore(ctx, key, "team-1")
	if err != nil {
		t.Errorf("GetScore failed: %v", err)
	}

	if score != 150 {
		t.Errorf("expected score 150, got %f", score)
	}
}

func TestLeaderboardUpdateScore(t *testing.T) {
	rdb, err := NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer rdb.Close()

	ctx := context.Background()
	key := "test:lb:update:" + time.Now().Format("20060102150405")
	defer rdb.Del(ctx, key)

	lb := NewLeaderboard(rdb)

	lb.AddScore(ctx, key, "team-1", 100)
	lb.AddScore(ctx, key, "team-1", 50)

	score, _ := lb.GetScore(ctx, key, "team-1")
	if score != 150 {
		t.Errorf("expected score 150 after update, got %f", score)
	}
}

func TestLeaderboardRemove(t *testing.T) {
	rdb, err := NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer rdb.Close()

	ctx := context.Background()
	key := "test:lb:remove:" + time.Now().Format("20060102150405")
	defer rdb.Del(ctx, key)

	lb := NewLeaderboard(rdb)

	lb.AddScore(ctx, key, "team-1", 100)
	lb.AddScore(ctx, key, "team-2", 200)

	err = lb.Remove(ctx, key, "team-1")
	if err != nil {
		t.Errorf("Remove failed: %v", err)
	}

	score, _ := lb.GetScore(ctx, key, "team-1")
	if score != 0 {
		t.Errorf("expected score 0 for removed team, got %f", score)
	}

	score, _ = lb.GetScore(ctx, key, "team-2")
	if score != 200 {
		t.Errorf("expected score 200 for team-2, got %f", score)
	}
}

func TestLeaderboardCount(t *testing.T) {
	rdb, err := NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer rdb.Close()

	ctx := context.Background()
	key := "test:lb:count:" + time.Now().Format("20060102150405")
	defer rdb.Del(ctx, key)

	lb := NewLeaderboard(rdb)

	lb.AddScore(ctx, key, "team-1", 100)
	lb.AddScore(ctx, key, "team-2", 200)
	lb.AddScore(ctx, key, "team-3", 300)

	count, err := lb.Count(ctx, key)
	if err != nil {
		t.Errorf("Count failed: %v", err)
	}

	if count != 3 {
		t.Errorf("expected count 3, got %d", count)
	}
}
