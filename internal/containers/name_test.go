package containers

import "testing"

func TestServiceContainerName(t *testing.T) {
	got := serviceContainerName("example", "example-api", 1)
	want := "example-api-1.example"

	if got != want {
		t.Fatalf("unexpected service container name: got %q, want %q", got, want)
	}
}
