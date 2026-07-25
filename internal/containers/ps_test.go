package containers

import (
	"testing"

	"github.com/briheet/nozarashi/internal/specs"
)

func TestSelectProjectContainers(t *testing.T) {
	containers := &specs.Containers{
		Items: []specs.Container{
			{ID: "example-api-1"},
			{ID: "example-worker-1"},
			{ID: "another-project-api-1"},
		},
	}
	graph := &specs.Graph{
		Nodes: []*specs.ServiceNode{
			{
				ServiceName: "example-worker",
				Spec:        &specs.ServiceSpecs{Replicas: 1},
			},
		},
	}

	selected := selectProjectContainers(containers, graph, false)

	if len(selected.Items) != 1 || selected.Items[0].ID != "example-worker-1" {
		t.Fatalf("unexpected selected containers: %#v", selected.Items)
	}
}

func TestSelectProjectContainersKeepsAll(t *testing.T) {
	containers := &specs.Containers{
		Items: []specs.Container{
			{ID: "example-api-1"},
			{ID: "another-project-api-1"},
		},
	}

	selected := selectProjectContainers(containers, &specs.Graph{}, true)

	if len(selected.Items) != 2 {
		t.Fatalf("expected every runtime container: %#v", selected.Items)
	}
}
