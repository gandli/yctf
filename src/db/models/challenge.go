package models

import "fmt"

type Challenge struct {
	ID              string                 `json:"id"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Category        string                 `json:"category"`
	Points          int                    `json:"points"`
	FlagTemplate    string                 `json:"flag_template"`
	ContainerImage  string                 `json:"container_image"`
	ContainerConfig map[string]interface{} `json:"container_config"`
	IsVisible       bool                   `json:"is_visible"`
	MinScoreRatio   float64                `json:"min_score_ratio"`
	DecayThreshold  int                    `json:"decay_threshold"`
	Solves          int                    `json:"solves"`
	CreatedBy       string                 `json:"created_by"`
}

func NewChallenge(title, category string, points int) *Challenge {
	return &Challenge{
		Title:          title,
		Category:       category,
		Points:         points,
		MinScoreRatio:  0.5,
		DecayThreshold: 100,
	}
}

func (c *Challenge) ValidateCategory() error {
	switch c.Category {
	case "web", "pwn", "crypto", "re", "misc", "forensics", "osint":
		return nil
	default:
		return fmt.Errorf("invalid category: %s", c.Category)
	}
}

func (c *Challenge) DynamicScore(solves int) int {
	if solves == 0 {
		return c.Points
	}
	if solves == 1 {
		return c.Points
	}
	if solves == 2 {
		return c.Points
	}

	ratio := 1.0 - float64(solves)/float64(c.DecayThreshold)
	if ratio < c.MinScoreRatio {
		ratio = c.MinScoreRatio
	}

	score := int(float64(c.Points) * ratio)
	if score < int(float64(c.Points)*c.MinScoreRatio) {
		score = int(float64(c.Points) * c.MinScoreRatio)
	}
	return score
}
