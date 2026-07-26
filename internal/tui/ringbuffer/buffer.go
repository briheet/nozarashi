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

	// Head points to the oldest element in the buffer
	head int

	// Size is the number of elements currently stored
	size int
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
func (r *RingBuffer[T]) Push(item T) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.capacity == 0 {
		return
	}

	tail := (r.head + r.size) % r.capacity
	r.data[tail] = item

	if r.size < r.capacity {
		r.size++
		return
	}

	r.head = (r.head + 1) % r.capacity
}

// Latest returns the newest element without removing it from the buffer
func (r *RingBuffer[T]) Latest() (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.size == 0 {
		var zero T
		return zero, false
	}

	tail := (r.head + r.size - 1) % r.capacity
	return r.data[tail], true
}
