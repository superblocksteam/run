package waitgroup

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWaitAlive(t *testing.T) {
	testCases := []struct {
		name         string
		blockOnClose bool
	}{
		{
			name:         "no block on close",
			blockOnClose: false,
		},
		{
			name:         "block on close",
			blockOnClose: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.True(t, New(&sync.WaitGroup{}, tc.blockOnClose).Alive())
		})
	}
}

func TestWaitNoBlockOnClose(t *testing.T) {
	wg := sync.WaitGroup{}
	done := make(chan struct{})
	runnable := New(&wg, false)

	go func() {
		runnable.Run(context.Background())
		close(done)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
	}()

	time.Sleep(50 * time.Millisecond) // Ensure goroutine has time to execute
	runnable.Close(context.Background())
	<-done
}

func TestWaitBlockOnClose(t *testing.T) {
	var done bool
	wg := sync.WaitGroup{}
	runnable := New(&wg, true)

	go func() {
		runnable.Run(context.Background())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		done = true
	}()

	time.Sleep(50 * time.Millisecond) // Ensure goroutine has time to execute
	runnable.Close(context.Background())
	assert.True(t, done)
}
