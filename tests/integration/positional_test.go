package integration_test

import (
	"bytes"
	"context"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/briheet/nozarashi/internal/containers"
)

// TestPositionalServiceLifecycle checks build, up and down with a service argument.
func TestPositionalServiceLifecycle(t *testing.T) {
	// Skip when the Apple Container CLI is not available.
	if _, err := exec.LookPath(containers.ContainerCliName); err != nil {
		t.Skipf("%s is not available: %v", containers.ContainerCliName, err)
	}

	// Resolve the config fixture independently from the test working directory.
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve integration test path")
	}
	configDirectory := filepath.Join(filepath.Dir(testFile), "..", "data")

	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Minute)
	defer cancel()

	options := containers.ContainerOptions{
		FilePath: configDirectory,
		Args:     []string{"redis"},
	}

	// Build only the Redis service image.
	if err := containers.BuildContainers(ctx, options); err != nil {
		t.Fatalf("build positional Redis service: %v", err)
	}

	imageInspectCmd := exec.CommandContext(
		ctx,
		containers.ContainerCliName,
		containers.ContainerImageInspectArgs(redisImage)...,
	)
	if output, err := imageInspectCmd.CombinedOutput(); err != nil {
		t.Fatalf("inspect positional Redis image: %v\n%s", err, output)
	}

	// Start only the Redis service.
	if err := containers.UpContainers(ctx, options); err != nil {
		t.Fatalf("up positional Redis service: %v", err)
	}

	containerInspectCmd := exec.CommandContext(
		ctx,
		containers.ContainerCliName,
		containers.ContainerInspectArgs(redisContainer)...,
	)
	output, err := containerInspectCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("inspect positional Redis container: %v\n%s", err, output)
	}
	if !bytes.Contains(output, []byte(`"state" : "running"`)) {
		t.Fatalf("positional Redis container is not running:\n%s", output)
	}

	// Execute redis-cli inside the selected Redis service.
	execOptions := options
	execOptions.Args = []string{"redis", "redis-cli", "ping"}
	execOptions.Replica = 1
	if err := containers.ExecContainers(ctx, execOptions); err != nil {
		t.Fatalf("exec positional Redis service: %v", err)
	}

	// Restart only the Redis service.
	if err := containers.RestartContainers(ctx, options); err != nil {
		t.Fatalf("restart positional Redis service: %v", err)
	}

	containerInspectCmd = exec.CommandContext(
		ctx,
		containers.ContainerCliName,
		containers.ContainerInspectArgs(redisContainer)...,
	)
	output, err = containerInspectCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("inspect restarted positional Redis container: %v\n%s", err, output)
	}
	if !bytes.Contains(output, []byte(`"state" : "running"`)) {
		t.Fatalf("restarted positional Redis container is not running:\n%s", output)
	}

	// Inspect the selected Redis service.
	if err := containers.InspectContainers(ctx, options); err != nil {
		t.Fatalf("inspect positional Redis service: %v", err)
	}

	// List only the selected project service.
	if err := containers.PsContainers(ctx, options); err != nil {
		t.Fatalf("ps positional Redis service: %v", err)
	}

	// Get logs only from the Redis service.
	logOptions := options
	logOptions.Number = 5
	if err := containers.LogsContainers(ctx, logOptions); err != nil {
		t.Fatalf("logs positional Redis service: %v", err)
	}

	// Stop only the Redis service.
	if err := containers.DownContainers(ctx, options); err != nil {
		t.Fatalf("down positional Redis service: %v", err)
	}

	containerInspectCmd = exec.CommandContext(
		ctx,
		containers.ContainerCliName,
		containers.ContainerInspectArgs(redisContainer)...,
	)
	output, err = containerInspectCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("inspect stopped positional Redis container: %v\n%s", err, output)
	}
	if !bytes.Contains(output, []byte(`"state" : "stopped"`)) {
		t.Fatalf("positional Redis container is not stopped:\n%s", output)
	}

	// List every runtime container, including stopped containers.
	allOptions := options
	allOptions.All = true
	if err := containers.PsContainers(ctx, allOptions); err != nil {
		t.Fatalf("ps all containers: %v", err)
	}
}
