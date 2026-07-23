package docker

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type ContainerConfig struct {
	Image   string            `json:"image"`
	Cmd     []string          `json:"cmd"`
	Env     []string          `json:"env"`
	Ports   map[string]string `json:"ports"`
	Name    string            `json:"name"`
	Network string            `json:"network"`
}

type ContainerInfo struct {
	ID     string `json:"id"`
	Image  string `json:"image"`
	Status string `json:"image"`
	State  string `json:"state"`
}

type DockerCLIClient struct {
	timeout time.Duration
}

func NewDockerCLIClient() (*DockerCLIClient, error) {
	// Verify docker CLI is available
	cmd := exec.Command("docker", "version", "--format", "{{.Server.Version}}")
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("docker CLI not available: %w", err)
	}
	return &DockerCLIClient{
		timeout: 30 * time.Second,
	}, nil
}

func (d *DockerCLIClient) Ping() error {
	cmd := exec.Command("docker", "info", "--format", "{{.ServerVersion}}")
	return cmd.Run()
}

func (d *DockerCLIClient) ListContainers(all bool) ([]ContainerInfo, error) {
	args := []string{"ps", "--no-trunc", "--format", "{{.ID}}|{{.Image}}|{{.Status}}"}
	if all {
		args = append([]string{"ps", "-a"}, args[1:]...)
	}
	cmd := exec.Command("docker", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var containers []ContainerInfo
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 3)
		if len(parts) == 3 {
			containers = append(containers, ContainerInfo{
				ID:     parts[0],
				Image:  parts[1],
				Status: parts[2],
			})
		}
	}
	return containers, nil
}

func (d *DockerCLIClient) PullImage(image string) error {
	cmd := exec.Command("docker", "pull", image)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func (d *DockerCLIClient) CreateContainer(config *ContainerConfig) (string, error) {
	args := []string{"run", "-d", "--rm"}
	
	for _, env := range config.Env {
		args = append(args, "-e", env)
	}
	
	for hostPort, containerPort := range config.Ports {
		args = append(args, "-p", hostPort+":"+containerPort)
	}
	
	if config.Name != "" {
		args = append(args, "--name", config.Name)
	}
	
	if config.Network != "" {
		args = append(args, "--network", config.Network)
	}
	
	args = append(args, config.Image)
	args = append(args, config.Cmd...)
	
	cmd := exec.Command("docker", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("docker run failed: %w", err)
	}
	
	return strings.TrimSpace(string(out)), nil
}

func (d *DockerCLIClient) StartContainer(id string) error {
	cmd := exec.Command("docker", "start", id)
	return cmd.Run()
}

func (d *DockerCLIClient) StopContainer(id string) error {
	cmd := exec.Command("docker", "stop", id)
	return cmd.Run()
}

func (d *DockerCLIClient) RemoveContainer(id string, force bool) error {
	args := []string{"rm"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, id)
	cmd := exec.Command("docker", args...)
	return cmd.Run()
}

func (d *DockerCLIClient) InspectContainer(id string) (*ContainerInfo, error) {
	cmd := exec.Command("docker", "inspect", "--format", "{{.Id}}|{{.Image}}|{{.State.Status}}", id)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	parts := strings.SplitN(strings.TrimSpace(string(out)), "|", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("unexpected inspect output: %s", out)
	}
	
	return &ContainerInfo{
		ID:     parts[0],
		Image:  parts[1],
		State:  parts[2],
	}, nil
}

func (d *DockerCLIClient) ContainerExists(id string) (bool, error) {
	_, err := d.InspectContainer(id)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != 0 {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}
