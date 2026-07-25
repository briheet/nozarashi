package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/briheet/nozarashi/internal/specs"
)

func TestRenderContainers(t *testing.T) {
	containers := &specs.Containers{
		Items: []specs.Container{
			{
				ID: "example-api-1",
				Configuration: specs.ContainerConfiguration{
					Image:    specs.ContainerImage{Reference: "example-api"},
					Platform: specs.ContainerPlatform{OperatingSystem: "linux", Architecture: "arm64"},
					Resources: specs.ContainerResources{
						CPUs:          4,
						MemoryInBytes: 1024 * 1024 * 1024,
					},
				},
				Status: specs.ContainerStatus{
					State:       "running",
					StartedDate: "2026-07-25T00:00:00Z",
					Networks: []specs.ContainerNetworkStatus{
						{IPv4Address: "192.168.64.2/24"},
					},
				},
			},
		},
	}

	var output bytes.Buffer
	if err := renderContainers(t.Context(), &output, containers); err != nil {
		t.Fatalf("render containers: %v", err)
	}

	for _, expected := range []string{
		"ID",
		"IMAGE",
		"example-api-1",
		"example-api",
		"linux",
		"arm64",
		"running",
		"192.168.64.2/24",
		"1024 MB",
	} {
		if !strings.Contains(output.String(), expected) {
			t.Fatalf("rendered containers missing %q:\n%s", expected, output.String())
		}
	}
}
