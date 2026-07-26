package containers

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/briheet/nozarashi/internal/specs"
)

// GetContainerStats gets one resource usage sample from Apple Container.
func GetContainerStats(ctx context.Context) ([]specs.ContainerStats, error) {
	statsCmd := exec.CommandContext(ctx, ContainerCliName, ContainerStatsArgs...)

	output, err := statsCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("get container stats: %s: %w", strings.TrimSpace(string(output)), err)
	}

	var stats []specs.ContainerStats
	if err := json.Unmarshal(output, &stats); err != nil {
		return nil, fmt.Errorf("decode container stats: %w", err)
	}

	return stats, nil
}
