package containers

import (
	"testing"

	"github.com/briheet/nozarashi/internal/specs"
)

func TestSelectServiceNodes(t *testing.T) {
	graph := &specs.Graph{
		Project: &specs.Specs{
			Project: specs.ProjectSpecs{Name: "example"},
			Services: map[string]specs.ServiceSpecs{
				"api":    {},
				"worker": {},
			},
		},
		Nodes: []*specs.ServiceNode{
			{ServiceName: "example-api"},
			{ServiceName: "example-worker"},
		},
	}

	if err := selectServiceNodes(graph, []string{"worker"}); err != nil {
		t.Fatalf("select worker service: %v", err)
	}

	if len(graph.Nodes) != 1 || graph.Nodes[0].ServiceName != "example-worker" {
		t.Fatalf("unexpected selected services: %#v", graph.Nodes)
	}
}

func TestSelectServiceNodesKeepsAllWithoutArguments(t *testing.T) {
	graph := &specs.Graph{
		Nodes: []*specs.ServiceNode{
			{ServiceName: "example-api"},
			{ServiceName: "example-worker"},
		},
	}

	if err := selectServiceNodes(graph, nil); err != nil {
		t.Fatalf("select all services: %v", err)
	}

	if len(graph.Nodes) != 2 {
		t.Fatalf("expected every service to remain selected: %#v", graph.Nodes)
	}
}

func TestSelectServiceNodesRejectsUnknownService(t *testing.T) {
	graph := &specs.Graph{
		Project: &specs.Specs{
			Project:  specs.ProjectSpecs{Name: "example"},
			Services: map[string]specs.ServiceSpecs{"api": {}},
		},
	}

	if err := selectServiceNodes(graph, []string{"worker"}); err == nil {
		t.Fatal("expected unknown service error")
	}
}
