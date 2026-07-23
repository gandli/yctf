package models

import (
	"testing"
)

func TestNewSubmission(t *testing.T) {
	userID := "user-1"
	teamID := "team-1"
	challengeID := "chal-1"
	flag := "flag{test}"

	sub := NewSubmission(userID, teamID, challengeID, flag)

	if sub.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, sub.UserID)
	}
	if sub.TeamID != teamID {
		t.Errorf("expected teamID %s, got %s", teamID, sub.TeamID)
	}
	if sub.ChallengeID != challengeID {
		t.Errorf("expected challengeID %s, got %s", challengeID, sub.ChallengeID)
	}
	if sub.FlagSubmitted != flag {
		t.Errorf("expected flag %s, got %s", flag, sub.FlagSubmitted)
	}
	if sub.IsCorrect {
		t.Error("new submission should not be marked correct")
	}
}

func TestSubmissionMarkCorrect(t *testing.T) {
	sub := NewSubmission("u1", "t1", "c1", "flag{test}")
	sub.MarkCorrect()
	if !sub.IsCorrect {
		t.Error("submission should be marked correct")
	}
}

func TestSubmissionMarkIncorrect(t *testing.T) {
	sub := NewSubmission("u1", "t1", "c1", "flag{test}")
	sub.MarkCorrect()
	sub.MarkIncorrect()
	if sub.IsCorrect {
		t.Error("submission should be marked incorrect")
	}
}

func TestNewScoreEvent(t *testing.T) {
	teamID := "team-1"
	challengeID := "chal-1"
	points := 100

	event := NewScoreEvent(teamID, challengeID, points)

	if event.TeamID != teamID {
		t.Errorf("expected teamID %s, got %s", teamID, event.TeamID)
	}
	if event.ChallengeID != challengeID {
		t.Errorf("expected challengeID %s, got %s", challengeID, event.ChallengeID)
	}
	if event.PointsAwarded != points {
		t.Errorf("expected points %d, got %d", points, event.PointsAwarded)
	}
}

func TestScoreEventWithBonus(t *testing.T) {
	event := NewScoreEvent("t1", "c1", 105) // 5% first blood bonus
	if event.PointsAwarded != 105 {
		t.Errorf("expected 105, got %d", event.PointsAwarded)
	}
}
