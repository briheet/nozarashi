package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/briheet/nozarashi/internal/specs"
)

func TestRenderContainerDetails(t *testing.T) {
	containers := &specs.Containers{
		Items: []specs.Container{
			{
				ID: "example-api-1",
				Configuration: specs.ContainerConfiguration{
					CreationDate: "2026-07-25T00:00:00Z",
					Image:        specs.ContainerImage{Reference: "example-api"},
					Platform: specs.ContainerPlatform{
						OperatingSystem: "linux",
						Architecture:    "arm64",
					},
					Resources: specs.ContainerResources{
						CPUs:          4,
						MemoryInBytes: 1024 * 1024 * 1024,
					},
					InitProcess: specs.ContainerInitProcess{
						Executable:       "./api",
						Arguments:        []string{"serve"},
						Environment:      []string{"PORT=8080"},
						WorkingDirectory: "/app",
					},
					Networks: []specs.ContainerNetwork{
						{Network: "example-default"},
					},
					Mounts: []specs.ContainerMount{
						{
							Destination: "/data",
							Type: specs.ContainerMountType{
								Volume: specs.ContainerVolumeMount{Name: "example-data"},
							},
						},
					},
					PublishedPorts: []specs.ContainerPublishedPort{
						{
							HostAddress:   "0.0.0.0",
							HostPort:      8080,
							ContainerPort: 8080,
							Protocol:      "tcp",
						},
					},
				},
				Status: specs.ContainerStatus{
					State:       "running",
					StartedDate: "2026-07-25T00:00:01Z",
					Networks: []specs.ContainerNetworkStatus{
						{Network: "example-default", IPv4Address: "192.168.64.2/24"},
					},
				},
			},
		},
	}

	var output bytes.Buffer
	if err := renderContainerDetails(t.Context(), &output, containers); err != nil {
		t.Fatalf("render container details: %v", err)
	}

	for _, expected := range []string{
		"FIELD",
		"example-api-1",
		"./api serve",
		"example-default (192.168.64.2/24)",
		"example-data -> /data",
		"0.0.0.0:8080 -> 8080/tcp",
		"PORT=8080",
	} {
		if !strings.Contains(output.String(), expected) {
			t.Fatalf("rendered container details missing %q:\n%s", expected, output.String())
		}
	}
}
