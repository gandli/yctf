package models

type Submission struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	TeamID         string `json:"team_id"`
	ChallengeID    string `json:"challenge_id"`
	FlagSubmitted  string `json:"flag_submitted"`
	IsCorrect      bool   `json:"is_correct"`
	IPAddress      string `json:"ip_address"`
}

func NewSubmission(userID, teamID, challengeID, flag string) *Submission {
	return &Submission{
		UserID:        userID,
		TeamID:        teamID,
		ChallengeID:   challengeID,
		FlagSubmitted: flag,
		IsCorrect:     false,
	}
}

func (s *Submission) MarkCorrect() {
	s.IsCorrect = true
}

func (s *Submission) MarkIncorrect() {
	s.IsCorrect = false
}

type ScoreEvent struct {
	ID            string `json:"id"`
	TeamID        string `json:"team_id"`
	ChallengeID   string `json:"challenge_id"`
	PointsAwarded int    `json:"points_awarded"`
}

func NewScoreEvent(teamID, challengeID string, points int) *ScoreEvent {
	return &ScoreEvent{
		TeamID:        teamID,
		ChallengeID:   challengeID,
		PointsAwarded: points,
	}
}
