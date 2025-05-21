package testing

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/creack/pty"
)

// getExistingServices checks for existing Docker compose services
func getExistingServices(ctx context.Context, composeFilePath string) ([]ComposeService, error) {
	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composeFilePath, "ps", "--format", "json")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	var runningServices []ComposeService

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) == 0 {
			continue // Skip empty lines
		}

		var service ComposeService
		if err := json.Unmarshal([]byte(line), &service); err != nil {
			return nil, fmt.Errorf("error unmarshaling service output: %w", err)
		}
		runningServices = append(runningServices, service)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning output: %w", err)
	}

	return runningServices, nil
}

// pullDockerImages pulls the latest Docker images for the services
func pullDockerImages(ctx context.Context, composeFilePath string) error {
	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composeFilePath, "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	var errBuffer bytes.Buffer

	ptty, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("failed to start command with pseudo-tty: %w", err)
	}
	defer ptty.Close()

	mw := io.MultiWriter(os.Stdout, &errBuffer)

	io.Copy(mw, ptty)

	if err := cmd.Wait(); err != nil {
		errorMsg := strings.TrimSpace(errBuffer.String())
		if len(errorMsg) > 0 {
			return fmt.Errorf("docker compose up failed: %w: %s", err, errorMsg)
		}
		return fmt.Errorf("docker compose up failed: %w", err)
	}

	return nil
}

// checkServicesHealth verifies if existing services are healthy
func checkServicesHealth(services []ComposeService) bool {
	for _, service := range services {
		// ecp-agent-1 does not have a health status
		if service.Name == "ecp-agent-1" {
			continue
		}
		if service.Health != "healthy" {
			log.Printf("Service %s is not healthy: %s", service.Name, service.Health)
			return false
		}
	}
	return true
}

// waitForServicesHealth waits for all services to become healthy
func waitForServicesHealth(ctx context.Context, composeFilePath string) error {
	log.Println("Waiting for services to be healthy...")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			services, err := getExistingServices(ctx, composeFilePath)
			if err != nil {
				return err
			}

			if len(services) == 0 {
				continue
			}

			healthy := checkServicesHealth(services)

			if healthy {
				return nil
			}

			log.Println("Services not yet healthy, continuing to wait...")
		}
	}
}

// stopDockerEnvironment stops the Docker Compose environment
func stopDockerEnvironment(ctx context.Context) error {
	log.Println("Stopping Docker environment...")

	composeFilePath := filepath.Join("docker", "docker-compose.yml")
	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composeFilePath, "down", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	var errBuffer bytes.Buffer

	ptty, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("failed to start command with pseudo-tty: %w", err)
	}
	defer ptty.Close()

	mw := io.MultiWriter(os.Stdout, &errBuffer)

	io.Copy(mw, ptty)

	if err := cmd.Wait(); err != nil {
		errorMsg := strings.TrimSpace(errBuffer.String())
		if len(errorMsg) > 0 {
			return fmt.Errorf("docker compose down failed: %w: %s", err, errorMsg)
		}
		return fmt.Errorf("docker compose down failed: %w", err)
	}

	log.Println("Docker environment stopped")
	return nil
}

func checkVersionMatch(services []ComposeService, expectedVersion string) bool {
	for _, service := range services {
		if !strings.Contains(service.Image, expectedVersion) {
			return false
		}
	}
	return true
}

type ComposeService struct {
	Command      string      `json:"Command"`
	CreatedAt    string      `json:"CreatedAt"`
	ExitCode     int         `json:"ExitCode"`
	Health       string      `json:"Health"`
	ID           string      `json:"ID"`
	Image        string      `json:"Image"`
	Labels       string      `json:"Labels"`
	LocalVolumes string      `json:"LocalVolumes"`
	Mounts       string      `json:"Mounts"`
	Name         string      `json:"Name"`
	Names        string      `json:"Names"`
	Networks     string      `json:"Networks"`
	Ports        string      `json:"Ports"`
	Project      string      `json:"Project"`
	Publishers   []Publisher `json:"Publishers"`
	RunningFor   string      `json:"RunningFor"`
	Service      string      `json:"Service"`
	Size         string      `json:"Size"`
	State        string      `json:"State"`
	Status       string      `json:"Status"`
}

type Publisher struct {
	URL           string `json:"URL"`
	TargetPort    int    `json:"TargetPort"`
	PublishedPort int    `json:"PublishedPort"`
	Protocol      string `json:"Protocol"`
}
