package render

import (
	"context"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/briheet/nozarashi/internal/specs"
)

// This function wraps over containers logs and uses text/tabwriter to pretty print
func RenderLogs(ctx context.Context, logs *specs.Logs) error {
	// os.Stdout now, can switch to something else in future
	return renderLogs(ctx, os.Stdout, logs)
}

// Renders container logs to the provided output.
func renderLogs(ctx context.Context, output io.Writer, logs *specs.Logs) error {
	writer := tabwriter.NewWriter(output, 0, 4, 2, ' ', 0)

	if _, err := fmt.Fprintln(writer, "CONTAINER\tLOG"); err != nil {
		return fmt.Errorf("write logs header: %w", err)
	}

	for _, containerLogs := range logs.Containers {
		if len(containerLogs.Lines) == 0 {
			if _, err := fmt.Fprintf(writer, "%s\t\n", containerLogs.Name); err != nil {
				return fmt.Errorf("write logs for container %q: %w", containerLogs.Name, err)
			}
			continue
		}

		for _, line := range containerLogs.Lines {
			if err := ctx.Err(); err != nil {
				return fmt.Errorf("render logs: %w", err)
			}

			if _, err := fmt.Fprintf(writer, "%s\t%s\n", containerLogs.Name, line); err != nil {
				return fmt.Errorf("write logs for container %q: %w", containerLogs.Name, err)
			}
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush logs: %w", err)
	}

	return nil
}
