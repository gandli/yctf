package docker

import (
	"strings"
	"testing"
	"time"
)

func TestNewContainerManager(t *testing.T) {
	cm, err := NewContainerManager()
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	if cm == nil {
		t.Error("manager should not be nil")
	}
}

func TestContainerManagerStartInstance(t *testing.T) {
	cm, err := NewContainerManager()
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}

	config := &InstanceConfig{
		ChallengeID: "chal-123",
		TeamID:      "team-456",
		Image:       "alpine:latest",
		Cmd:         []string{"sleep", "30"},
		Env:         []string{"FLAG=flag{test-flag}"},
		ExpiresIn:   5 * time.Minute,
	}

	instance, err := cm.StartInstance(config)
	if err != nil {
		t.Fatalf("StartInstance failed: %v", err)
	}

	if instance.ID == "" {
		t.Error("instance ID should not be empty")
	}
	if instance.ChallengeID != config.ChallengeID {
		t.Errorf("expected challenge ID %s, got %s", config.ChallengeID, instance.ChallengeID)
	}
	if instance.TeamID != config.TeamID {
		t.Errorf("expected team ID %s, got %s", config.TeamID, instance.TeamID)
	}
	if instance.Status != "running" {
		t.Errorf("expected status running, got %s", instance.Status)
	}

	// Cleanup
	_ = cm.StopInstance(instance.ID)
}

func TestContainerManagerStopInstance(t *testing.T) {
	cm, err := NewContainerManager()
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}

	config := &InstanceConfig{
		ChallengeID: "chal-456",
		TeamID:      "team-789",
		Image:       "alpine:latest",
		Cmd:         []string{"sleep", "30"},
		Env:         []string{"FLAG=flag{stop-test}"},
		ExpiresIn:   5 * time.Minute,
	}

	instance, err := cm.StartInstance(config)
	if err != nil {
		t.Fatalf("StartInstance failed: %v", err)
	}

	err = cm.StopInstance(instance.ID)
	if err != nil {
		t.Errorf("StopInstance failed: %v", err)
	}

	// Verify stopped
	info, err := cm.GetInstance(instance.ID)
	if err != nil {
		t.Errorf("GetInstance failed: %v", err)
	}
	if info.Status != "stopped" && info.Status != "exited" {
		t.Errorf("expected status stopped/exited, got %s", info.Status)
	}

	// Cleanup
	_ = cm.RemoveInstance(instance.ID)
}

func TestContainerManagerGetInstance(t *testing.T) {
	cm, err := NewContainerManager()
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}

	config := &InstanceConfig{
		ChallengeID: "chal-get",
		TeamID:      "team-get",
		Image:       "alpine:latest",
		Cmd:         []string{"sleep", "30"},
		Env:         []string{"FLAG=flag{get-test}"},
		ExpiresIn:   5 * time.Minute,
	}

	instance, err := cm.StartInstance(config)
	if err != nil {
		t.Fatalf("StartInstance failed: %v", err)
	}
	defer cm.RemoveInstance(instance.ID)

	info, err := cm.GetInstance(instance.ID)
	if err != nil {
		t.Errorf("GetInstance failed: %v", err)
	}

	if info.ID != instance.ID {
		t.Errorf("expected ID %s, got %s", instance.ID, info.ID)
	}
	if info.ChallengeID != config.ChallengeID {
		t.Errorf("expected challenge ID %s, got %s", config.ChallengeID, info.ChallengeID)
	}
}

func TestContainerManagerListInstances(t *testing.T) {
	cm, err := NewContainerManager()
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}

	// Start 2 instances
	config1 := &InstanceConfig{
		ChallengeID: "chal-list-1",
		TeamID:      "team-list-1",
		Image:       "alpine:latest",
		Cmd:         []string{"sleep", "30"},
		Env:         []string{"FLAG=flag{list-1}"},
		ExpiresIn:   5 * time.Minute,
	}
	config2 := &InstanceConfig{
		ChallengeID: "chal-list-2",
		TeamID:      "team-list-2",
		Image:       "alpine:latest",
		Cmd:         []string{"sleep", "30"},
		Env:         []string{"FLAG=flag{list-2}"},
		ExpiresIn:   5 * time.Minute,
	}

	instance1, err := cm.StartInstance(config1)
	if err != nil {
		t.Fatalf("StartInstance 1 failed: %v", err)
	}
	instance2, err := cm.StartInstance(config2)
	if err != nil {
		t.Fatalf("StartInstance 2 failed: %v", err)
	}

	// List all
	instances, err := cm.ListInstances()
	if err != nil {
		t.Errorf("ListInstances failed: %v", err)
	}

	if len(instances) < 2 {
		t.Errorf("expected at least 2 instances, got %d", len(instances))
	}

	// Cleanup
	_ = cm.RemoveInstance(instance1.ID)
	_ = cm.RemoveInstance(instance2.ID)
}

func TestContainerManagerListInstancesByTeam(t *testing.T) {
	cm, err := NewContainerManager()
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}

	config := &InstanceConfig{
		ChallengeID: "chal-team-list",
		TeamID:      "team-filter-test",
		Image:       "alpine:latest",
		Cmd:         []string{"sleep", "30"},
		Env:         []string{"FLAG=flag{team-list}"},
		ExpiresIn:   5 * time.Minute,
	}

	instance, err := cm.StartInstance(config)
	if err != nil {
		t.Fatalf("StartInstance failed: %v", err)
	}

	// Filter by team
	instances, err := cm.ListInstancesByTeam("team-filter-test")
	if err != nil {
		t.Errorf("ListInstancesByTeam failed: %v", err)
	}

	if len(instances) != 1 {
		t.Errorf("expected 1 instance, got %d", len(instances))
	}

	if instances[0].TeamID != "team-filter-test" {
		t.Errorf("expected team ID team-filter-test, got %s", instances[0].TeamID)
	}

	_ = cm.RemoveInstance(instance.ID)
}

func TestContainerManagerRemoveInstance(t *testing.T) {
	cm, err := NewContainerManager()
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}

	config := &InstanceConfig{
		ChallengeID: "chal-remove",
		TeamID:      "team-remove",
		Image:       "alpine:latest",
		Cmd:         []string{"sleep", "30"},
		Env:         []string{"FLAG=flag{remove-test}"},
		ExpiresIn:   5 * time.Minute,
	}

	instance, err := cm.StartInstance(config)
	if err != nil {
		t.Fatalf("StartInstance failed: %v", err)
	}

	err = cm.RemoveInstance(instance.ID)
	if err != nil {
		t.Errorf("RemoveInstance failed: %v", err)
	}

	// Verify removed
	_, err = cm.GetInstance(instance.ID)
	if err == nil {
		t.Error("GetInstance should fail for removed instance")
	}
}

func TestContainerManagerInstanceExists(t *testing.T) {
	cm, err := NewContainerManager()
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}

	config := &InstanceConfig{
		ChallengeID: "chal-exists",
		TeamID:      "team-exists",
		Image:       "alpine:latest",
		Cmd:         []string{"sleep", "30"},
		Env:         []string{"FLAG=flag{exists-test}"},
		ExpiresIn:   5 * time.Minute,
	}

	instance, err := cm.StartInstance(config)
	if err != nil {
		t.Fatalf("StartInstance failed: %v", err)
	}

	exists, err := cm.InstanceExists(instance.ID)
	if err != nil {
		t.Errorf("InstanceExists failed: %v", err)
	}
	if !exists {
		t.Error("instance should exist")
	}

	_ = cm.RemoveInstance(instance.ID)

	exists, err = cm.InstanceExists(instance.ID)
	if err != nil {
		t.Errorf("InstanceExists failed: %v", err)
	}
	if exists {
		t.Error("instance should not exist after removal")
	}
}

func TestContainerManagerExpiredInstances(t *testing.T) {
	cm, err := NewContainerManager()
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}

	config := &InstanceConfig{
		ChallengeID: "chal-expired",
		TeamID:      "team-expired",
		Image:       "alpine:latest",
		Cmd:         []string{"sleep", "30"},
		Env:         []string{"FLAG=flag{expired-test}"},
		ExpiresIn:   1 * time.Second, // Very short expiry
	}

	instance, err := cm.StartInstance(config)
	if err != nil {
		t.Fatalf("StartInstance failed: %v", err)
	}

	// Wait for expiry
	time.Sleep(2 * time.Second)

	// Check expired
	expired, err := cm.GetExpiredInstances()
	if err != nil {
		t.Errorf("GetExpiredInstances failed: %v", err)
	}

	found := false
	for _, exp := range expired {
		if exp.ID == instance.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("instance should be expired")
	}

	_ = cm.RemoveInstance(instance.ID)
}

func TestContainerManagerPortMapping(t *testing.T) {
	cm, err := NewContainerManager()
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}

	config := &InstanceConfig{
		ChallengeID: "chal-port",
		TeamID:      "team-port",
		Image:       "alpine:latest",
		Cmd:         []string{"sleep", "30"},
		Env:         []string{"FLAG=flag{port-test}"},
		Ports:       map[string]string{"8080": "80"},
		ExpiresIn:   5 * time.Minute,
	}

	instance, err := cm.StartInstance(config)
	if err != nil {
		t.Fatalf("StartInstance failed: %v", err)
	}

	if instance.Port == 0 {
		t.Error("expected port mapping")
	}

	_ = cm.RemoveInstance(instance.ID)
}

func TestContainerManagerFlagInjection(t *testing.T) {
	cm, err := NewContainerManager()
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}

	config := &InstanceConfig{
		ChallengeID: "chal-flag",
		TeamID:      "team-flag",
		Image:       "alpine:latest",
		Cmd:         []string{"sleep", "30"},
		Env:         []string{"FLAG=flag{injected-test}", "OTHER_VAR=value"},
		ExpiresIn:   5 * time.Minute,
	}

	instance, err := cm.StartInstance(config)
	if err != nil {
		t.Fatalf("StartInstance failed: %v", err)
	}

	// Verify flag is in env
	if !strings.Contains(instance.Flag, "flag{") {
		t.Errorf("expected flag in instance, got %s", instance.Flag)
	}

	_ = cm.RemoveInstance(instance.ID)
}
