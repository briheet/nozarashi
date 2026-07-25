package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/briheet/nozarashi/internal/specs"
)

func TestRenderLogs(t *testing.T) {
	logs := &specs.Logs{
		Containers: []specs.ContainerLogs{
			{
				Name:  "example-api-1",
				Lines: []string{"server started", "request completed"},
			},
			{
				Name:  "example-worker-1",
				Lines: []string{"job completed"},
			},
		},
	}

	var output bytes.Buffer
	if err := renderLogs(t.Context(), &output, logs); err != nil {
		t.Fatalf("render logs: %v", err)
	}

	for _, expected := range []string{
		"CONTAINER",
		"LOG",
		"example-api-1",
		"server started",
		"request completed",
		"example-worker-1",
		"job completed",
	} {
		if !strings.Contains(output.String(), expected) {
			t.Fatalf("rendered logs missing %q:\n%s", expected, output.String())
		}
	}
}
