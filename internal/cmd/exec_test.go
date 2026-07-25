package cmd

import (
	"slices"
	"testing"
)

func TestExecCmdPreservesCommandFlags(t *testing.T) {
	execCmd := ExecCmd()

	if err := execCmd.ParseFlags([]string{
		"--replica",
		"2",
		"redis",
		"sh",
		"-c",
		"echo ready",
	}); err != nil {
		t.Fatalf("parse exec flags: %v", err)
	}

	got := execCmd.Flags().Args()
	want := []string{"redis", "sh", "-c", "echo ready"}

	if !slices.Equal(got, want) {
		t.Fatalf("unexpected exec command arguments:\ngot:  %q\nwant: %q", got, want)
	}
}
