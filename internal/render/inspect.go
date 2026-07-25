package render

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/briheet/nozarashi/internal/specs"
)

// RenderContainerDetails uses text/tabwriter to print container details.
func RenderContainerDetails(ctx context.Context, containers *specs.Containers) error {
	// os.Stdout now, can switch to something else in future.
	return renderContainerDetails(ctx, os.Stdout, containers)
}

// Renders detailed container fields to the provided output.
func renderContainerDetails(
	ctx context.Context,
	output io.Writer,
	containers *specs.Containers,
) error {
	writer := tabwriter.NewWriter(output, 0, 4, 2, ' ', 0)

	for containerIndex, container := range containers.Items {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("render container details: %w", err)
		}

		if _, err := fmt.Fprintln(writer, "FIELD\tVALUE"); err != nil {
			return fmt.Errorf("write container details header: %w", err)
		}

		command := strings.Join(
			append(
				[]string{container.Configuration.InitProcess.Executable},
				container.Configuration.InitProcess.Arguments...,
			),
			" ",
		)
		memory := fmt.Sprintf(
			"%d MB",
			container.Configuration.Resources.MemoryInBytes/(1024*1024),
		)

		fields := [][2]string{
			{"ID", container.ID},
			{"IMAGE", container.Configuration.Image.Reference},
			{"STATE", container.Status.State},
			{"CREATED", container.Configuration.CreationDate},
			{"STARTED", container.Status.StartedDate},
			{"OS", container.Configuration.Platform.OperatingSystem},
			{"ARCH", container.Configuration.Platform.Architecture},
			{"CPUS", fmt.Sprintf("%d", container.Configuration.Resources.CPUs)},
			{"MEMORY", memory},
			{"COMMAND", command},
			{"WORKDIR", container.Configuration.InitProcess.WorkingDirectory},
		}

		for _, field := range fields {
			if _, err := fmt.Fprintf(writer, "%s\t%s\n", field[0], field[1]); err != nil {
				return fmt.Errorf("write container %q details: %w", container.ID, err)
			}
		}

		for _, network := range container.Configuration.Networks {
			value := network.Network
			for _, status := range container.Status.Networks {
				if status.Network == network.Network && status.IPv4Address != "" {
					value = fmt.Sprintf("%s (%s)", value, status.IPv4Address)
					break
				}
			}

			if _, err := fmt.Fprintf(writer, "NETWORK\t%s\n", value); err != nil {
				return fmt.Errorf("write container %q network: %w", container.ID, err)
			}
		}

		for _, mount := range container.Configuration.Mounts {
			source := mount.Type.Volume.Name
			if source == "" {
				source = mount.Source
			}

			if _, err := fmt.Fprintf(
				writer,
				"MOUNT\t%s -> %s\n",
				source,
				mount.Destination,
			); err != nil {
				return fmt.Errorf("write container %q mount: %w", container.ID, err)
			}
		}

		for _, port := range container.Configuration.PublishedPorts {
			if _, err := fmt.Fprintf(
				writer,
				"PORT\t%s:%d -> %d/%s\n",
				port.HostAddress,
				port.HostPort,
				port.ContainerPort,
				port.Protocol,
			); err != nil {
				return fmt.Errorf("write container %q port: %w", container.ID, err)
			}
		}

		for _, environment := range container.Configuration.InitProcess.Environment {
			if _, err := fmt.Fprintf(writer, "ENV\t%s\n", environment); err != nil {
				return fmt.Errorf("write container %q environment: %w", container.ID, err)
			}
		}

		if containerIndex < len(containers.Items)-1 {
			if _, err := fmt.Fprintln(writer); err != nil {
				return fmt.Errorf("separate container details: %w", err)
			}
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush container details: %w", err)
	}

	return nil
}
