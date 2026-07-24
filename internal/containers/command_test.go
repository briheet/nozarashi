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
