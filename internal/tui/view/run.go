package view

import (
	"context"
	"time"

	"github.com/briheet/nozarashi/internal/containers"
	"github.com/briheet/nozarashi/internal/specs"
	"github.com/briheet/nozarashi/internal/tui/model"
	"github.com/briheet/nozarashi/internal/tui/poller"
	"github.com/briheet/nozarashi/internal/tui/ringbuffer"
	tea "github.com/charmbracelet/bubbletea"
)

func Run(ctx context.Context, opts containers.ContainerOptions) error {
	// First check this containers system is running
	if err := containers.StatusSystemContainers(ctx); err != nil {
		return err
	}

	// One time alloc memory for ring buffer
	// Make sure to not allow any allocs after this
	memBuf := ringbuffer.NewRingBuffer[specs.Containers](1000)

	// Get the initial model for tui
	// Pass in shared buffer and polling time
	model := model.InitialModel(memBuf, time.Millisecond*250)

	// Poller for polling container specs and pushing in ring buffer
	// Pass in shared buffer and polling time
	poller := poller.NewPoller(memBuf, time.Millisecond*200)

	// Create a new program
	program := tea.NewProgram(teaModel{m: model})

	if _, err := program.Run(); err != nil {
		return err
	}

	return nil
}
