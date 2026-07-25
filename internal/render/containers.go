package render

import (
	"context"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/briheet/nozarashi/internal/specs"
)

// RenderContainers uses text/tabwriter to print containers.
func RenderContainers(ctx context.Context, containers *specs.Containers) error {
	// os.Stdout now, can switch to something else in future.
	return renderContainers(ctx, os.Stdout, containers)
}

// Renders containers to the provided output.
func renderContainers(ctx context.Context, output io.Writer, containers *specs.Containers) error {
	writer := tabwriter.NewWriter(output, 0, 4, 2, ' ', 0)

	if _, err := fmt.Fprintln(writer, "ID\tIMAGE\tOS\tARCH\tSTATE\tIP\tCPUS\tMEMORY\tSTARTED"); err != nil {
		return fmt.Errorf("write containers header: %w", err)
	}

	for _, container := range containers.Items {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("render containers: %w", err)
		}

		ipAddress := ""
		if len(container.Status.Networks) > 0 {
			ipAddress = container.Status.Networks[0].IPv4Address
		}

		memory := fmt.Sprintf(
			"%d MB",
			container.Configuration.Resources.MemoryInBytes/(1024*1024),
		)

		if _, err := fmt.Fprintf(
			writer,
			"%s\t%s\t%s\t%s\t%s\t%s\t%d\t%s\t%s\n",
			container.ID,
			container.Configuration.Image.Reference,
			container.Configuration.Platform.OperatingSystem,
			container.Configuration.Platform.Architecture,
			container.Status.State,
			ipAddress,
			container.Configuration.Resources.CPUs,
			memory,
			container.Status.StartedDate,
		); err != nil {
			return fmt.Errorf("write container %q: %w", container.ID, err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush containers: %w", err)
	}

	return nil
}
