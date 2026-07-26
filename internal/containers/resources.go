package containers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/briheet/nozarashi/internal/specs"
	"github.com/charmbracelet/x/ansi"
)

// GetImages lists images in the local Apple Container image store.
func GetImages(ctx context.Context) ([]specs.Image, error) {
	return getResourceList[specs.Image](ctx, "images", ContainerImageListArgs)
}

// GetVolumes lists volumes managed by Apple Container.
func GetVolumes(ctx context.Context) ([]specs.Volume, error) {
	return getResourceList[specs.Volume](ctx, "volumes", ContainerVolumeListArgs)
}

// GetNetworks lists networks managed by Apple Container.
func GetNetworks(ctx context.Context) ([]specs.Network, error) {
	return getResourceList[specs.Network](ctx, "networks", ContainerNetworkListArgs)
}

func getResourceList[T any](ctx context.Context, resource string, args []string) ([]T, error) {
	command := exec.CommandContext(ctx, ContainerCliName, args...)
	stderr := &bytes.Buffer{}
	command.Stderr = stderr

	output, err := command.Output()
	if err != nil {
		message := strings.TrimSpace(ansi.Strip(stderr.String()))
		if message != "" {
			return nil, fmt.Errorf("list %s: %s: %w", resource, message, err)
		}
		return nil, fmt.Errorf("list %s: %w", resource, err)
	}

	var items []T
	if err := json.Unmarshal(output, &items); err != nil {
		return nil, fmt.Errorf("decode %s list: %w", resource, err)
	}
	return items, nil
}
