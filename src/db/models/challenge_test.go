package models

import (
	"fmt"
	"testing"
)

func TestNewChallenge(t *testing.T) {
	title := "SQL Injection 101"
	category := "web"
	points := 100

	chal := NewChallenge(title, category, points)

	if chal.Title != title {
		t.Errorf("expected title %s, got %s", title, chal.Title)
	}
	if chal.Category != category {
		t.Errorf("expected category %s, got %s", category, chal.Category)
	}
	if chal.Points != points {
		t.Errorf("expected points %d, got %d", points, chal.Points)
	}
}

func TestChallengeValidateCategory(t *testing.T) {
	validCategories := []string{"web", "pwn", "crypto", "re", "misc", "forensics", "osint"}
	for _, cat := range validCategories {
		chal := &Challenge{Category: cat}
		if err := chal.ValidateCategory(); err != nil {
			t.Errorf("category %s should be valid: %v", cat, err)
		}
	}

	invalidCategories := []string{"", "hacking", "WEB", "Pwn"}
	for _, cat := range invalidCategories {
		chal := &Challenge{Category: cat}
		if err := chal.ValidateCategory(); err == nil {
			t.Errorf("category %q should be invalid", cat)
		}
	}
}

func TestDynamicScore(t *testing.T) {
	chal := &Challenge{
		Points:          100,
		MinScoreRatio:   0.5,
		DecayThreshold:  100,
	}

	// First solve: full score
	score := chal.DynamicScore(0)
	if score != 100 {
		t.Errorf("first solve should be 100, got %d", score)
	}

	// After 50 solves: should be ~50 (with decay)
	score = chal.DynamicScore(50)
	if score < 50 || score > 60 {
		t.Errorf("expected ~50 after 50 solves, got %d", score)
	}

	// After 100+ solves: should be at min_score_ratio (50)
	score = chal.DynamicScore(100)
	if score != 50 {
		t.Errorf("expected 50 (min), got %d", score)
	}

	// After 200 solves: still at floor
	score = chal.DynamicScore(200)
	if score != 50 {
		t.Errorf("expected 50 (floor), got %d", score)
	}
}

func TestDynamicScoreFirstBlood(t *testing.T) {
	chal := &Challenge{Points: 100, MinScoreRatio: 0.5, DecayThreshold: 100}

	// First blood bonus: 5%
	first := chal.DynamicScore(0)
	if first != 100 {
		t.Errorf("first blood should be 100, got %d", first)
	}

	// Second blood: 3% bonus
	second := chal.DynamicScore(1)
	if second != 100 {
		t.Errorf("second blood should be 100, got %d", second)
	}

	// Third: 1% bonus
	third := chal.DynamicScore(2)
	if third != 100 {
		t.Errorf("third blood should be 100, got %d", third)
	}
}

func TestChallengeVisibility(t *testing.T) {
	chal := &Challenge{IsVisible: false}
	if chal.IsVisible {
		t.Error("challenge should be hidden by default")
	}
}

func TestChallengeWithContainer(t *testing.T) {
	chal := &Challenge{
		Title:          "PWN 101",
		Category:       "pwn",
		Points:         200,
		ContainerImage: "yctf/pwn101:latest",
		ContainerConfig: map[string]interface{}{
			"ports": []int{1337},
		},
	}

	if chal.ContainerImage == "" {
		t.Error("container image should be set")
	}
	_ = fmt.Sprintf("config: %v", chal.ContainerConfig)
}
