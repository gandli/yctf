package models

import (
	"testing"
	"time"
)

func TestNewInstance(t *testing.T) {
	challengeID := "chal-123"
	teamID := "team-456"
	containerID := "abc123def456"
	image := "alpine:latest"

	inst := NewInstance(challengeID, teamID, containerID, image)

	if inst.ChallengeID != challengeID {
		t.Errorf("expected challenge ID %s, got %s", challengeID, inst.ChallengeID)
	}
	if inst.TeamID != teamID {
		t.Errorf("expected team ID %s, got %s", teamID, inst.TeamID)
	}
	if inst.ContainerID != containerID {
		t.Errorf("expected container ID %s, got %s", containerID, inst.ContainerID)
	}
	if inst.Image != image {
		t.Errorf("expected image %s, got %s", image, inst.Image)
	}
	if inst.Status != "pending" {
		t.Errorf("expected status pending, got %s", inst.Status)
	}
}

func TestInstanceSetRunning(t *testing.T) {
	inst := NewInstance("chal-1", "team-1", "cont-1", "alpine:latest")
	inst.SetRunning()

	if inst.Status != "running" {
		t.Errorf("expected status running, got %s", inst.Status)
	}
}

func TestInstanceSetStopped(t *testing.T) {
	inst := NewInstance("chal-1", "team-1", "cont-1", "alpine:latest")
	inst.SetRunning()
	inst.SetStopped()

	if inst.Status != "stopped" {
		t.Errorf("expected status stopped, got %s", inst.Status)
	}
}

func TestInstanceSetExpired(t *testing.T) {
	inst := NewInstance("chal-1", "team-1", "cont-1", "alpine:latest")
	inst.SetExpired()

	if inst.Status != "expired" {
		t.Errorf("expected status expired, got %s", inst.Status)
	}
}

func TestInstanceIsExpired(t *testing.T) {
	inst := NewInstance("chal-1", "team-1", "cont-1", "alpine:latest")
	inst.ExpiresAt = time.Now().Add(-1 * time.Minute)

	if !inst.IsExpired() {
		t.Error("instance should be expired")
	}

	inst.ExpiresAt = time.Now().Add(1 * time.Hour)
	if inst.IsExpired() {
		t.Error("instance should not be expired")
	}
}

func TestInstanceIsRunning(t *testing.T) {
	inst := NewInstance("chal-1", "team-1", "cont-1", "alpine:latest")
	
	if inst.IsRunning() {
		t.Error("new instance should not be running")
	}

	inst.SetRunning()
	if !inst.IsRunning() {
		t.Error("instance should be running after SetRunning")
	}
}

func TestInstancePort(t *testing.T) {
	inst := NewInstance("chal-1", "team-1", "cont-1", "alpine:latest")
	inst.Port = 30001

	if inst.Port != 30001 {
		t.Errorf("expected port 30001, got %d", inst.Port)
	}
}

func TestInstanceFlag(t *testing.T) {
	inst := NewInstance("chal-1", "team-1", "cont-1", "alpine:latest")
	inst.Flag = "flag{test-flag}"

	if inst.Flag != "flag{test-flag}" {
		t.Errorf("expected flag{test-flag}, got %s", inst.Flag)
	}
}

func TestInstanceExpiresAt(t *testing.T) {
	inst := NewInstance("chal-1", "team-1", "cont-1", "alpine:latest")
	expiresAt := time.Now().Add(1 * time.Hour)
	inst.ExpiresAt = expiresAt

	if !inst.ExpiresAt.Equal(expiresAt) {
		t.Errorf("expected expires at %v, got %v", expiresAt, inst.ExpiresAt)
	}
}

func TestInstanceCreatedAt(t *testing.T) {
	inst := NewInstance("chal-1", "team-1", "cont-1", "alpine:latest")
	
	if inst.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
}
