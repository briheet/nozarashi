package model

import (
	"os"

	"github.com/briheet/nozarashi/internal/specs"
	"github.com/briheet/nozarashi/internal/tui/ringbuffer"
)

// Base model
type Model struct {
	// Data specific
	RingBuffer *ringbuffer.RingBuffer[Snapshot]
	Updates    <-chan error

	Cursor         int
	ResourceCursor int
	ActivePanel    int
	LogOffset      int

	// Window specifics
	Width  int
	Height int

	Hostname string
	Err      error
}

// Snapshot type buffer
type Snapshot struct {
	Containers specs.Containers
	Images     []specs.Image
	Volumes    []specs.Volume
	Networks   []specs.Network
	Logs       map[string][]string
	Stats      map[string]specs.ContainerStats
}

// Model init
func InitialModel(ringBuffer *ringbuffer.RingBuffer[Snapshot], updates <-chan error) *Model {
	hostname, _ := os.Hostname()
	return &Model{
		RingBuffer:  ringBuffer,
		Updates:     updates,
		ActivePanel: 1,
		Hostname:    hostname,
	}
}
