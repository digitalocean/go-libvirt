package event

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testEvent struct{ id int32 }

func (e testEvent) GetCallbackID() int32 {
	return e.id
}

func testEvents(count int) []Event {
	ev := make([]Event, count)

	for i := 0; i < count; i++ {
		ev[i] = testEvent{int32(i)}
	}

	return ev
}

// TestStreamSequential tests the stream objects by first writing some entries
// the stream, then reading them back.
func TestStreamSequential(t *testing.T) {
	const evCount = 100000

	s := NewStream(1, 2)
	defer s.Shutdown()
	events := testEvents(evCount)

	// send a bunch of "events", and make sure we receive them all.
	fmt.Println("sending test events to queue")
	for i, e := range events {
		assert.Equal(t, int32(i), e.GetCallbackID())
		s.Push(e)
		assert.Equal(t, i+1, s.Len())
	}

	assert.Equal(t, evCount, s.Len())

	fmt.Println("reading events from queue")
	for _, e := range events {
		qe, ok := <-s.Recv()
		assert.True(t, ok)
		assert.Equal(t, e, qe)
	}

	assert.Zero(t, s.Len())
}

// TestStreamParallel tests the stream object with a pair of goroutines, one
// writing to the queue, the other reading from it.
func TestStreamParallel(t *testing.T) {
	const evCount = 10000

	s := NewStream(1, 2)
	defer s.Shutdown()
	events := testEvents(evCount)
	wg := sync.WaitGroup{}
	wg.Add(2)

	// Start 2 goroutines, one to send, the other to receive.
	go func() {
		defer wg.Done()
		for _, e := range events {
			s.Push(e)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < evCount; i++ {
			e := <-s.Recv()
			assert.Equal(t, int32(i), e.GetCallbackID())
		}
	}()

	wg.Wait()
	assert.Zero(t, s.Len())
}
