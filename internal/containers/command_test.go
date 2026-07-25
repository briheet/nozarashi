package containers

import (
	"slices"
	"testing"

	"github.com/briheet/nozarashi/internal/specs"
)

func TestContainerRunArgs(t *testing.T) {
	service := &specs.ServiceSpecs{
		EntryPoint: []string{"/usr/bin/env", "sh"},
		Command:    []string{"server", "--listen"},
		Environment: map[string]string{
			"PORT": "8080",
			"MODE": "test",
		},
		Volumes: []specs.VolumeMount{
			{Name: "data", Path: "/data"},
		},
		Networks: []string{"default"},
		Ports: []specs.PortSpec{
			{Host: 8080, Container: 80},
		},
	}

	got := ContainerRunArgs(
		"example",
		"example-api-1",
		"example-api",
		"linux/arm64",
		service,
	)
	want := []string{
		"run",
		"--detach",
		"--name",
		"example-api-1",
		"--platform",
		"linux/arm64",
		"--env",
		"MODE=test",
		"--env",
		"PORT=8080",
		"--volume",
		"example-data:/data",
		"--network",
		"example-default",
		"--publish",
		"8080:80",
		"--entrypoint",
		"/usr/bin/env",
		"example-api",
		"sh",
		"server",
		"--listen",
	}

	if !slices.Equal(got, want) {
		t.Fatalf("unexpected container run arguments:\ngot:  %q\nwant: %q", got, want)
	}
}

func TestContainerStopArgs(t *testing.T) {
	got := ContainerStopArgs("example-api-1")
	want := []string{"stop", "example-api-1"}

	if !slices.Equal(got, want) {
		t.Fatalf("unexpected container stop arguments:\ngot:  %q\nwant: %q", got, want)
	}
}

func TestContainerDeleteArgs(t *testing.T) {
	got := ContainerDeleteArgs("example-api-1")
	want := []string{"delete", "example-api-1"}

	if !slices.Equal(got, want) {
		t.Fatalf("unexpected container delete arguments:\ngot:  %q\nwant: %q", got, want)
	}
}

func TestContainerImageDeleteArgs(t *testing.T) {
	got := ContainerImageDeleteArgs("example-api")
	want := []string{"image", "delete", "example-api"}

	if !slices.Equal(got, want) {
		t.Fatalf("unexpected image delete arguments:\ngot:  %q\nwant: %q", got, want)
	}
}

func TestContainerLogsArgs(t *testing.T) {
	got := ContainerLogsArgs("example-api-1", 25)
	want := []string{"logs", "-n", "25", "example-api-1"}

	if !slices.Equal(got, want) {
		t.Fatalf("unexpected container logs arguments:\ngot:  %q\nwant: %q", got, want)
	}
}

func TestContainerLogsArgsWithoutNumber(t *testing.T) {
	got := ContainerLogsArgs("example-api-1", 0)
	want := []string{"logs", "example-api-1"}

	if !slices.Equal(got, want) {
		t.Fatalf("unexpected container logs arguments:\ngot:  %q\nwant: %q", got, want)
	}
}

func TestContainerListArgs(t *testing.T) {
	got := ContainerListArgs(false)
	want := []string{"list", "--format", "json"}

	if !slices.Equal(got, want) {
		t.Fatalf("unexpected container list arguments:\ngot:  %q\nwant: %q", got, want)
	}
}

func TestContainerListAllArgs(t *testing.T) {
	got := ContainerListArgs(true)
	want := []string{"list", "--all", "--format", "json"}

	if !slices.Equal(got, want) {
		t.Fatalf("unexpected container list all arguments:\ngot:  %q\nwant: %q", got, want)
	}
}

func TestContainerExecArgs(t *testing.T) {
	got := ContainerExecArgs(
		"example-api-1",
		ContainerOptions{
			Args:        []string{"api", "sh", "-c", "echo ready"},
			Interactive: true,
			TTY:         true,
		},
	)
	want := []string{
		"exec",
		"--interactive",
		"--tty",
		"example-api-1",
		"sh",
		"-c",
		"echo ready",
	}

	if !slices.Equal(got, want) {
		t.Fatalf("unexpected container exec arguments:\ngot:  %q\nwant: %q", got, want)
	}
}
