package poller

import (
	"time"

	"github.com/briheet/nozarashi/internal/specs"
	"github.com/briheet/nozarashi/internal/tui/ringbuffer"
)

type poller struct{}

func NewPoller(ringBuffer *ringbuffer.RingBuffer[specs.Containers], pollingTime time.Duration) *poller {
	return &poller{}
}
