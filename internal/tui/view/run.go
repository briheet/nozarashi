package view

import (
	"context"
	"io"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/briheet/nozarashi/internal/containers"
	"github.com/briheet/nozarashi/internal/tui/model"
	"github.com/briheet/nozarashi/internal/tui/poller"
	"github.com/briheet/nozarashi/internal/tui/ringbuffer"
)

const pollInterval = time.Second

func Run(ctx context.Context, opts containers.ContainerOptions) error {
	// First check this containers system is running
	if err := containers.StatusSystemContainers(ctx, io.Discard); err != nil {
		return err
	}

	// One time alloc memory for ring buffer
	// Make sure to not allow any allocs after this
	memBuf := ringbuffer.NewRingBuffer[model.Snapshot](1000)

	// Notify the model after the poller writes data to the ring buffer
	updates := make(chan error, 1)

	// Get the initial model for tui
	// Pass in shared buffer and update channel
	baseModel := model.InitialModel(memBuf, updates)

	// Poller for polling container specs and pushing in ring buffer
	// Pass in shared buffer, update channel and polling time
	containerPoller := poller.NewPoller(memBuf, updates, pollInterval)

	ctx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	go func() {
		defer close(done)
		containerPoller.Poll(ctx)
	}()
	defer cleanup(cancel, done)

	// Create a new program
	program := tea.NewProgram(
		teaModel{m: baseModel},
		tea.WithContext(ctx),
	)

	if _, err := program.Run(); err != nil {
		return err
	}

	return nil
}

func cleanup(cancel context.CancelFunc, done <-chan struct{}) {
	cancel()
	<-done
}
