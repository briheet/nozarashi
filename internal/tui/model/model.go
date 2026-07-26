package model

import (
	"time"

	"github.com/briheet/nozarashi/internal/specs"
	"github.com/briheet/nozarashi/internal/tui/ringbuffer"
)

// Base model
type Model struct {
}

func InitialModel(membuf *ringbuffer.RingBuffer[specs.Containers], pollingTime time.Duration) *Model {
	return &Model{}
}
