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

const redisImage = "docker.io/library/redis:7-alpine"

const (
	redisVolume    = "nozarashi-integration-redis-data"
	redisNetwork   = "nozarashi-integration-default"
	redisContainer = "nozarashi-integration-redis-1.nozarashi-integration"
)

// TestRedisServiceLifecycle checks the complete up and down flow against Apple Containers.
func TestRedisServiceLifecycle(t *testing.T) {
	// Skip when the Apple Container CLI is not available.
	if _, err := exec.LookPath(containers.ContainerCliName); err != nil {
		t.Skipf("%s is not available: %v", containers.ContainerCliName, err)
	}

	// Resolve the config fixture independently from the test working directory.
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve integration test path")
	}
	configPath := filepath.Join(filepath.Dir(testFile), "..", "data", "config.toml")

	// Run the lifecycle flow with a bounded integration timeout.
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Minute)
	defer cancel()

	options := containers.ContainerOptions{FilePath: configPath}

	// Pull or build every configured service image.
	if err := containers.BuildContainers(ctx, options); err != nil {
		t.Fatalf("run Redis build flow: %v", err)
	}

	// Create every configured project volume and network.
	if err := containers.CreateResources(ctx, options); err != nil {
		t.Fatalf("run Redis create flow: %v", err)
	}

	if err := containers.UpContainers(ctx, options); err != nil {
		t.Fatalf("run Redis up flow: %v", err)
	}

	// Confirm the Redis image is available in the local image store.
	imageInspectCmd := exec.CommandContext(
		ctx,
		containers.ContainerCliName,
		containers.ContainerImageInspectArgs(redisImage)...,
	)
	if output, err := imageInspectCmd.CombinedOutput(); err != nil {
		t.Fatalf("inspect Redis image: %v\n%s", err, output)
	}

	// Confirm the project volume is available.
	volumeInspectCmd := exec.CommandContext(
		ctx,
		containers.ContainerCliName,
		containers.ContainerVolumeInspectArgs(redisVolume)...,
	)
	if output, err := volumeInspectCmd.CombinedOutput(); err != nil {
		t.Fatalf("inspect Redis volume: %v\n%s", err, output)
	}

	// Confirm the project network is available.
	networkInspectCmd := exec.CommandContext(
		ctx,
		containers.ContainerCliName,
		containers.ContainerNetworkInspectArgs(redisNetwork)...,
	)
	if output, err := networkInspectCmd.CombinedOutput(); err != nil {
		t.Fatalf("inspect Redis network: %v\n%s", err, output)
	}

	// Confirm the Redis service container is running.
	containerInspectCmd := exec.CommandContext(
		ctx,
		containers.ContainerCliName,
		containers.ContainerInspectArgs(redisContainer)...,
	)
	output, err := containerInspectCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("inspect Redis container: %v\n%s", err, output)
	}
	if !bytes.Contains(output, []byte(`"state" : "running"`)) {
		t.Fatalf("Redis container is not running:\n%s", output)
	}

	// Confirm the service configuration reached the runtime.
	for _, expected := range [][]byte{
		[]byte(`"name" : "nozarashi-integration-redis-data"`),
		[]byte(`"network" : "nozarashi-integration-default"`),
		[]byte(`"containerPort" : 6379`),
		[]byte(`"hostPort" : 6379`),
	} {
		if !bytes.Contains(output, expected) {
			t.Fatalf("Redis container is missing %q:\n%s", expected, output)
		}
	}

	// Stop the Redis service.
	if err := containers.DownContainers(ctx, options); err != nil {
		t.Fatalf("run Redis down flow: %v", err)
	}

	// Confirm the Redis service container is stopped.
	containerInspectCmd = exec.CommandContext(
		ctx,
		containers.ContainerCliName,
		containers.ContainerInspectArgs(redisContainer)...,
	)
	output, err = containerInspectCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("inspect stopped Redis container: %v\n%s", err, output)
	}
	if !bytes.Contains(output, []byte(`"state" : "stopped"`)) {
		t.Fatalf("Redis container is not stopped:\n%s", output)
	}

	// Destroy every runtime object declared by the project.
	if err := containers.DestroyContainers(ctx, options); err != nil {
		t.Fatalf("run Redis destroy flow: %v", err)
	}

	// Confirm the container, network, volume and image were deleted.
	for name, args := range map[string][]string{
		"container": containers.ContainerInspectArgs(redisContainer),
		"network":   containers.ContainerNetworkInspectArgs(redisNetwork),
		"volume":    containers.ContainerVolumeInspectArgs(redisVolume),
		"image":     containers.ContainerImageInspectArgs(redisImage),
	} {
		inspectCmd := exec.CommandContext(ctx, containers.ContainerCliName, args...)
		if output, err := inspectCmd.CombinedOutput(); err == nil {
			t.Fatalf("%s still exists after destroy:\n%s", name, output)
		}
	}
}
