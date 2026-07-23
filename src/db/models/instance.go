package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type Instance struct {
	ID            string    `json:"id"`
	ChallengeID   string    `json:"challenge_id"`
	TeamID        string    `json:"team_id"`
	ContainerID   string    `json:"container_id"`
	Image         string    `json:"image"`
	Status        string    `json:"status"`
	Port          int       `json:"port"`
	Flag          string    `json:"flag"`
	ExpiresAt     time.Time `json:"expires_at"`
	CreatedAt     time.Time `json:"created_at"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
}

func NewInstance(challengeID, teamID, containerID, image string) *Instance {
	return &Instance{
		ID:          generateID(),
		ChallengeID: challengeID,
		TeamID:      teamID,
		ContainerID: containerID,
		Image:       image,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}
}

func (i *Instance) SetRunning() {
	i.Status = "running"
	i.LastHeartbeat = time.Now()
}

func (i *Instance) SetStopped() {
	i.Status = "stopped"
}

func (i *Instance) SetExpired() {
	i.Status = "expired"
}

func (i *Instance) SetError() {
	i.Status = "error"
}

func (i *Instance) IsExpired() bool {
	if i.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(i.ExpiresAt)
}

func (i *Instance) IsRunning() bool {
	return i.Status == "running"
}

func (i *Instance) RefreshHeartbeat() {
	i.LastHeartbeat = time.Now()
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
