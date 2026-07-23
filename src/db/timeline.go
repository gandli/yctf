package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type ScoreEvent struct {
	TeamID        string `json:"team_id"`
	TeamName      string `json:"team_name"`
	ChallengeName string `json:"challenge_name"`
	Points        int64  `json:"points"`
	Timestamp     string `json:"timestamp"`
}

type ScoreTimeline struct {
	rdb *redis.Client
}

func NewScoreTimeline(rdb *redis.Client) *ScoreTimeline {
	return &ScoreTimeline{rdb: rdb}
}

func (st *ScoreTimeline) Add(ctx context.Context, key, teamID, teamName, challengeName string, points int64) error {
	event := fmt.Sprintf("%s|%s|%s|%d|%d", teamID, teamName, challengeName, points, time.Now().Unix())
	return st.rdb.LPush(ctx, key, event).Err()
}

func (st *ScoreTimeline) GetRecent(ctx context.Context, key string, count int64) ([]string, error) {
	return st.rdb.LRange(ctx, key, 0, count-1).Result()
}

func (st *ScoreTimeline) GetByTeam(ctx context.Context, key string, teamID string) ([]string, error) {
	all, err := st.rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	var result []string
	for _, event := range all {
		if strings.HasPrefix(event, teamID+"|") {
			result = append(result, event)
		}
	}
	return result, nil
}

func (st *ScoreTimeline) Count(ctx context.Context, key string) (int64, error) {
	return st.rdb.LLen(ctx, key).Result()
}
