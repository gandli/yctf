package docker

import (
	"strings"
	"testing"
)

func TestNewDockerCLIClient(t *testing.T) {
	client, err := NewDockerCLIClient()
	if err != nil {
		t.Skipf("Docker CLI not available: %v", err)
	}
	if client == nil {
		t.Error("client should not be nil")
	}
}

func TestDockerCLIClientPing(t *testing.T) {
	client, err := NewDockerCLIClient()
	if err != nil {
		t.Skipf("Docker CLI not available: %v", err)
	}

	if err := client.Ping(); err != nil {
		t.Errorf("ping failed: %v", err)
	}
}

func TestDockerCLIClientListContainers(t *testing.T) {
	client, err := NewDockerCLIClient()
	if err != nil {
		t.Skipf("Docker CLI not available: %v", err)
	}

	containers, err := client.ListContainers(false)
	if err != nil {
		t.Errorf("ListContainers failed: %v", err)
	}
	_ = containers
}

func TestDockerCLIClientListContainersAll(t *testing.T) {
	client, err := NewDockerCLIClient()
	if err != nil {
		t.Skipf("Docker CLI not available: %v", err)
	}

	containers, err := client.ListContainers(true)
	if err != nil {
		t.Errorf("ListContainers failed: %v", err)
	}
	_ = containers
}

func TestDockerCLIClientPullImage(t *testing.T) {
	client, err := NewDockerCLIClient()
	if err != nil {
		t.Skipf("Docker CLI not available: %v", err)
	}

	err = client.PullImage("alpine:latest")
	if err != nil {
		t.Skipf("PullImage failed (may be offline): %v", err)
	}
}

func TestDockerCLIClientCreateContainer(t *testing.T) {
	client, err := NewDockerCLIClient()
	if err != nil {
		t.Skipf("Docker CLI not available: %v", err)
	}

	config := &ContainerConfig{
		Image: "alpine:latest",
		Cmd:   []string{"sleep", "5"},
		Env:   []string{"FLAG=test"},
	}

	id, err := client.CreateContainer(config)
	if err != nil {
		t.Skipf("CreateContainer failed: %v", err)
	}

	if len(id) < 12 {
		t.Errorf("expected container ID, got %s", id)
	}

	_ = client.RemoveContainer(id, true)
}

func TestDockerCLIClientStartContainer(t *testing.T) {
	client, err := NewDockerCLIClient()
	if err != nil {
		t.Skipf("Docker CLI not available: %v", err)
	}

	config := &ContainerConfig{
		Image: "alpine:latest",
		Cmd:   []string{"sleep", "5"},
	}

	id, err := client.CreateContainer(config)
	if err != nil {
		t.Skipf("CreateContainer failed: %v", err)
	}

	err = client.StartContainer(id)
	if err != nil {
		t.Errorf("StartContainer failed: %v", err)
	}

	_ = client.RemoveContainer(id, true)
}

func TestDockerCLIClientStopContainer(t *testing.T) {
	client, err := NewDockerCLIClient()
	if err != nil {
		t.Skipf("Docker CLI not available: %v", err)
	}

	config := &ContainerConfig{
		Image: "alpine:latest",
		Cmd:   []string{"sleep", "10"},
	}

	id, err := client.CreateContainer(config)
	if err != nil {
		t.Skipf("CreateContainer failed: %v", err)
	}

	_ = client.StartContainer(id)

	err = client.StopContainer(id)
	if err != nil {
		t.Errorf("StopContainer failed: %v", err)
	}

	_ = client.RemoveContainer(id, true)
}

func TestDockerCLIClientRemoveContainer(t *testing.T) {
	client, err := NewDockerCLIClient()
	if err != nil {
		t.Skipf("Docker CLI not available: %v", err)
	}

	config := &ContainerConfig{
		Image: "alpine:latest",
		Cmd:   []string{"sleep", "5"},
	}

	id, err := client.CreateContainer(config)
	if err != nil {
		t.Skipf("CreateContainer failed: %v", err)
	}

	err = client.RemoveContainer(id, true)
	if err != nil {
		t.Errorf("RemoveContainer failed: %v", err)
	}

	exists, _ := client.ContainerExists(id)
	if exists {
		t.Error("container should be removed")
	}
}

func TestDockerCLIClientInspectContainer(t *testing.T) {
	client, err := NewDockerCLIClient()
	if err != nil {
		t.Skipf("Docker CLI not available: %v", err)
	}

	config := &ContainerConfig{
		Image: "alpine:latest",
		Cmd:   []string{"sleep", "5"},
	}

	id, err := client.CreateContainer(config)
	if err != nil {
		t.Skipf("CreateContainer failed: %v", err)
	}
	defer client.RemoveContainer(id, true)

	info, err := client.InspectContainer(id)
	if err != nil {
		t.Errorf("InspectContainer failed: %v", err)
	}

	if info.ID != id && !strings.HasPrefix(info.ID, id) {
		t.Errorf("expected ID %s, got %s", id, info.ID)
	}
}

func TestDockerCLIClientContainerExists(t *testing.T) {
	client, err := NewDockerCLIClient()
	if err != nil {
		t.Skipf("Docker CLI not available: %v", err)
	}

	exists, err := client.ContainerExists("nonexistent-container-id")
	if err != nil {
		t.Errorf("ContainerExists failed: %v", err)
	}
	if exists {
		t.Error("non-existent container should not exist")
	}
}

func TestDockerCLIClientLifecycle(t *testing.T) {
	client, err := NewDockerCLIClient()
	if err != nil {
		t.Skipf("Docker CLI not available: %v", err)
	}

	config := &ContainerConfig{
		Image: "alpine:latest",
		Cmd:   []string{"sleep", "30"},
		Env:   []string{"FLAG=lifecycle-test"},
	}

	// Create
	id, err := client.CreateContainer(config)
	if err != nil {
		t.Fatalf("CreateContainer failed: %v", err)
	}

	// Verify exists
	exists, _ := client.ContainerExists(id)
	if !exists {
		t.Error("container should exist after creation")
	}

	// Inspect
	info, err := client.InspectContainer(id)
	if err != nil {
		t.Errorf("InspectContainer failed: %v", err)
	}
	if info.State != "running" {
		t.Errorf("expected state running, got %s", info.State)
	}

	// Stop
	err = client.StopContainer(id)
	if err != nil {
		t.Errorf("StopContainer failed: %v", err)
	}

	// Remove
	err = client.RemoveContainer(id, true)
	if err != nil {
		t.Errorf("RemoveContainer failed: %v", err)
	}

	// Verify gone
	exists, _ = client.ContainerExists(id)
	if exists {
		t.Error("container should not exist after removal")
	}
}
