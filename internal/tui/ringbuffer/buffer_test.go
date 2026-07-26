package ringbuffer

import (
	"sync"
	"testing"
)

func TestRingBuffer(t *testing.T) {
	buffer := NewRingBuffer[int](2)
	buffer.Push(1)
	buffer.Push(2)
	buffer.Push(3)

	item, ok := buffer.Latest()
	if !ok || item != 3 {
		t.Fatalf("unexpected latest item: %d, %t", item, ok)
	}

	item, ok = buffer.Latest()
	if !ok || item != 3 {
		t.Fatalf("unexpected preserved item: %d, %t", item, ok)
	}

	emptyBuffer := NewRingBuffer[int](2)
	if _, ok := emptyBuffer.Latest(); ok {
		t.Fatal("expected empty buffer")
	}
}

func TestRingBufferConcurrentAccess(t *testing.T) {
	buffer := NewRingBuffer[int](100)
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	go func() {
		defer waitGroup.Done()
		for value := range 1000 {
			buffer.Push(value)
		}
	}()

	go func() {
		defer waitGroup.Done()
		for range 1000 {
			buffer.Latest()
		}
	}()

	waitGroup.Wait()
}
