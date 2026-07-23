package db

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type Leaderboard struct {
	rdb *redis.Client
}

func NewLeaderboard(rdb *redis.Client) *Leaderboard {
	return &Leaderboard{rdb: rdb}
}

func (lb *Leaderboard) AddScore(ctx context.Context, key, member string, score int64) error {
	return lb.rdb.ZAdd(ctx, key, redis.Z{
		Score:  float64(score),
		Member: member,
	}).Err()
}

func (lb *Leaderboard) GetTop(ctx context.Context, key string, count int64) ([]redis.Z, error) {
	return lb.rdb.ZRevRangeWithScores(ctx, key, 0, count-1).Result()
}

func (lb *Leaderboard) GetRank(ctx context.Context, key, member string) (int64, error) {
	rank, err := lb.rdb.ZRevRank(ctx, key, member).Result()
	if err != nil {
		return 0, err
	}
	return rank + 1, nil
}

func (lb *Leaderboard) GetScore(ctx context.Context, key, member string) (float64, error) {
	return lb.rdb.ZScore(ctx, key, member).Result()
}

func (lb *Leaderboard) Remove(ctx context.Context, key, member string) error {
	return lb.rdb.ZRem(ctx, key, member).Err()
}

func (lb *Leaderboard) Count(ctx context.Context, key string) (int64, error) {
	return lb.rdb.ZCard(ctx, key).Result()
}
