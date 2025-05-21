package testing

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/creack/pty"
)

const (
	// Configuration constants
	elasticsearchPort = "9200"
	kibanaPort        = "5601"
	elasticsearchURL  = "https://localhost:" + elasticsearchPort
	kibanaURL         = "https://localhost:" + kibanaPort

	// Timeouts
	composeUpTimeout     = 5 * time.Minute
	elasticsearchTimeout = 3 * time.Minute
	kibanaTimeout        = 4 * time.Minute
	composeDownTimeout   = 2 * time.Minute

	// Environment variables that will be set for tests
	envKibanaURL        = "TEST_KIBANA_URL"
	envElasticsearchURL = "TEST_ELASTICSEARCH_URL"
	envTestMode         = "TEST_MODE"

	// Default memory limit for containers. Set MEM_LIMIT to change
	memLimit = "2147483648"
)

var (
	keepContainersRunning = true
)

func TestMain(m *testing.M) {
	if os.Getenv("STACK_VERSION") == "" {
		log.Fatalln("Required environment variable [STACK_VERSION] not set")
	}

	if os.Getenv(envTestMode) == "ci" {
		log.Println("Running in CI mode, skipping Docker setup")
		os.Exit(m.Run())
	}

	ctx, cancel := context.WithTimeout(context.Background(), composeUpTimeout)
	defer cancel()

	if err := startDockerEnvironment(ctx); err != nil {
		log.Fatalf("Failed to start Docker environment: %v", err)
	}

	// Set environment variables for tests
	os.Setenv(envKibanaURL, kibanaURL)
	os.Setenv(envElasticsearchURL, elasticsearchURL)

	// Run tests
	exitCode := m.Run()

	// Cleanup
	if !keepContainersRunning {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), composeDownTimeout)
		defer cleanupCancel()

		if err := stopDockerEnvironment(cleanupCtx); err != nil {
			log.Printf("Warning: Failed to stop Docker environment: %v", err)
		}
	}

	os.Exit(exitCode)
}

// startDockerEnvironment starts the Docker Compose environment
func startDockerEnvironment(ctx context.Context) error {
	stackVersion := os.Getenv("STACK_VERSION")

	// Determine the docker-compose file path
	composeFilePath := filepath.Join("docker", "docker-compose.yml")

	if os.Getenv("MEM_LIMIT") == "" {
		os.Setenv("MEM_LIMIT", memLimit)
	}

	// check if compose is already running and get the services
	existingServices, err := getExistingServices(ctx, composeFilePath)
	if err != nil {
		return fmt.Errorf("failed to check existing services: %w", err)
	}

	if len(existingServices) > 0 {
		versionMatch := checkVersionMatch(existingServices, stackVersion)
		if versionMatch {
			err := waitForServicesHealth(ctx, composeFilePath)
			if err != nil {
				log.Println("Existing Docker services not healthy, recreating...")
				// Services exist but aren't healthy, shut them down before restarting
				if err := stopDockerEnvironment(ctx); err != nil {
					return fmt.Errorf("failed to stop unhealthy services: %w", err)
				}
			}
			// Services healthy and versions match
			log.Println("Services healthy, running tests...")
			return nil
		} else {
			log.Printf("Version mismatch: expected %s, recreating services...", stackVersion)
			// Version mismatch, shut down before restarting
			if err := stopDockerEnvironment(ctx); err != nil {
				return fmt.Errorf("failed to stop services with version mismatch: %w", err)
			}
		}
	}

	// Build the docker-compose command
	dockerArgs := []string{"compose", "-f", composeFilePath, "up", "-d"}

	cmd := exec.CommandContext(ctx, "docker", dockerArgs...)

	env := os.Environ()
	env = append(env, fmt.Sprintf("MEM_LIMIT=%s", memLimit))
	cmd.Env = env

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

	log.Println("Docker environment is ready")
	return nil
}
