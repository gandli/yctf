package db

import (
	"context"
	"testing"
	"time"
)

func TestScoreTimelineAdd(t *testing.T) {
	rdb, err := NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer rdb.Close()

	ctx := context.Background()
	key := "test:timeline:" + time.Now().Format("20060102150405")
	defer rdb.Del(ctx, key)

	st := NewScoreTimeline(rdb)

	err = st.Add(ctx, key, "team-1", "Team A", "Web 101", 100)
	if err != nil {
		t.Errorf("Add failed: %v", err)
	}

	events, err := st.GetRecent(ctx, key, 10)
	if err != nil {
		t.Errorf("GetRecent failed: %v", err)
	}

	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
}

func TestScoreTimelineGetRecent(t *testing.T) {
	rdb, err := NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer rdb.Close()

	ctx := context.Background()
	key := "test:timeline:recent:" + time.Now().Format("20060102150405")
	defer rdb.Del(ctx, key)

	st := NewScoreTimeline(rdb)

	for i := 1; i <= 5; i++ {
		st.Add(ctx, key, "team-"+string(rune('0'+i)), "Team "+string(rune('A'+i-1)), "Challenge "+string(rune('0'+i)), int64(i*100))
	}

	events, err := st.GetRecent(ctx, key, 3)
	if err != nil {
		t.Errorf("GetRecent failed: %v", err)
	}

	if len(events) != 3 {
		t.Errorf("expected 3 events, got %d", len(events))
	}
}

func TestScoreTimelineGetByTeam(t *testing.T) {
	rdb, err := NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer rdb.Close()

	ctx := context.Background()
	key := "test:timeline:team:" + time.Now().Format("20060102150405")
	defer rdb.Del(ctx, key)

	st := NewScoreTimeline(rdb)

	st.Add(ctx, key, "team-1", "Team A", "Web 101", 100)
	st.Add(ctx, key, "team-2", "Team B", "PWN 101", 200)
	st.Add(ctx, key, "team-1", "Team A", "Crypto 101", 150)

	events, err := st.GetByTeam(ctx, key, "team-1")
	if err != nil {
		t.Errorf("GetByTeam failed: %v", err)
	}

	if len(events) != 2 {
		t.Errorf("expected 2 events for team-1, got %d", len(events))
	}
}

func TestScoreTimelineCount(t *testing.T) {
	rdb, err := NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer rdb.Close()

	ctx := context.Background()
	key := "test:timeline:count:" + time.Now().Format("20060102150405")
	defer rdb.Del(ctx, key)

	st := NewScoreTimeline(rdb)

	st.Add(ctx, key, "team-1", "Team A", "Web 101", 100)
	st.Add(ctx, key, "team-2", "Team B", "PWN 101", 200)
	st.Add(ctx, key, "team-3", "Team C", "RE 101", 300)

	count, err := st.Count(ctx, key)
	if err != nil {
		t.Errorf("Count failed: %v", err)
	}

	if count != 3 {
		t.Errorf("expected count 3, got %d", count)
	}
}
