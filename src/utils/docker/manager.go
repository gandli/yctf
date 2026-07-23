package docker

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type InstanceConfig struct {
	ChallengeID string
	TeamID      string
	Image       string
	Cmd         []string
	Env         []string
	Ports       map[string]string
	ExpiresIn   time.Duration
}

type Instance struct {
	ID           string
	ChallengeID  string
	TeamID       string
	ContainerID  string
	Image        string
	Status       string
	Port         int
	Flag         string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

type ContainerManager struct {
	client   *DockerCLIClient
	instances map[string]*Instance
	mu       sync.RWMutex
	nextPort int
}

func NewContainerManager() (*ContainerManager, error) {
	client, err := NewDockerCLIClient()
	if err != nil {
		return nil, err
	}
	return &ContainerManager{
		client:    client,
		instances: make(map[string]*Instance),
		nextPort:  30000,
	}, nil
}

func (cm *ContainerManager) StartInstance(config *InstanceConfig) (*Instance, error) {
	// Generate unique flag
	flag := generateFlag()

	// Prepare env vars with flag
	env := append(config.Env, "FLAG="+flag)

	// Prepare port mapping
	ports := make(map[string]string)
	if len(config.Ports) > 0 {
		for hostPort, containerPort := range config.Ports {
			ports[hostPort] = containerPort
		}
	} else {
		// Auto-assign port
		port := cm.nextPort
		cm.nextPort++
		ports[strconv.Itoa(port)] = "80"
	}

	// Create container
	containerConfig := &ContainerConfig{
		Image: config.Image,
		Cmd:   config.Cmd,
		Env:   env,
		Ports: ports,
		Name:  fmt.Sprintf("yctf-%s-%s", config.TeamID, config.ChallengeID),
	}

	containerID, err := cm.client.CreateContainer(containerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	// Extract port
	port := 0
	for hostPort := range ports {
		p, _ := strconv.Atoi(hostPort)
		if p > 0 {
			port = p
			break
		}
	}

	// Build instance
	instance := &Instance{
		ID:          generateID(),
		ChallengeID: config.ChallengeID,
		TeamID:      config.TeamID,
		ContainerID: containerID,
		Image:       config.Image,
		Status:      "running",
		Port:        port,
		Flag:        flag,
		ExpiresAt:   time.Now().Add(config.ExpiresIn),
		CreatedAt:   time.Now(),
	}

	// Start container
	if err := cm.client.StartContainer(containerID); err != nil {
		cm.client.RemoveContainer(containerID, true)
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	cm.mu.Lock()
	cm.instances[instance.ID] = instance
	cm.mu.Unlock()

	return instance, nil
}

func (cm *ContainerManager) StopInstance(id string) error {
	cm.mu.RLock()
	instance, ok := cm.instances[id]
	cm.mu.RUnlock()

	if !ok {
		return fmt.Errorf("instance not found: %s", id)
	}

	if err := cm.client.StopContainer(instance.ContainerID); err != nil {
		return err
	}

	instance.Status = "stopped"
	return nil
}

func (cm *ContainerManager) RemoveInstance(id string) error {
	cm.mu.RLock()
	instance, ok := cm.instances[id]
	cm.mu.RUnlock()

	if !ok {
		return fmt.Errorf("instance not found: %s", id)
	}

	force := instance.Status == "running"
	_ = cm.client.RemoveContainer(instance.ContainerID, force)

	cm.mu.Lock()
	delete(cm.instances, id)
	cm.mu.Unlock()

	return nil
}

func (cm *ContainerManager) GetInstance(id string) (*Instance, error) {
	cm.mu.RLock()
	instance, ok := cm.instances[id]
	cm.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("instance not found: %s", id)
	}
	return instance, nil
}

func (cm *ContainerManager) ListInstances() ([]*Instance, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	list := make([]*Instance, 0, len(cm.instances))
	for _, inst := range cm.instances {
		list = append(list, inst)
	}
	return list, nil
}

func (cm *ContainerManager) ListInstancesByTeam(teamID string) ([]*Instance, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var list []*Instance
	for _, inst := range cm.instances {
		if inst.TeamID == teamID {
			list = append(list, inst)
		}
	}
	return list, nil
}

func (cm *ContainerManager) InstanceExists(id string) (bool, error) {
	cm.mu.RLock()
	_, ok := cm.instances[id]
	cm.mu.RUnlock()
	return ok, nil
}

func (cm *ContainerManager) GetExpiredInstances() ([]*Instance, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var expired []*Instance
	now := time.Now()
	for _, inst := range cm.instances {
		if now.After(inst.ExpiresAt) {
			expired = append(expired, inst)
		}
	}
	return expired, nil
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generateFlag() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("flag{%s}", hex.EncodeToString(b)[:16])
}
