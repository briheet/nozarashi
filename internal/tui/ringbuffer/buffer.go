package ringbuffer

import "sync"

// Ring Buffer helps in coordinating between the
// poller service and the tui consumption.
type RingBuffer[T any] struct {
	// Mutex for safe accesses
	mu sync.RWMutex

	// Datatype is usually defined in specs as Containers data
	data []T

	// Capacity is the fixed capacity of the buffer
	capacity int
}

// NewRingBuffer allocates memory for a fixed-size ring buffer
func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	return &RingBuffer[T]{
		data:     make([]T, capacity),
		capacity: capacity,
	}
}

// Push adds an element to the tail of the buffer
// If the buffer is full, it overwrites the oldest element
func (r *RingBuffer[T]) Push(item T) {}

// Pop removes and returns the oldest element from the head of the buffer
func (r *RingBuffer[T]) Pop(item T) {}
