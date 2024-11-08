package waitgroup

import (
	"context"
	"sync"

	"github.com/superblocksteam/run"
)

type wait struct {
	blockOnClose bool
	done         chan struct{}
	wg           *sync.WaitGroup

	run.ForwardCompatibility
}

// NewWait returns a runnable that ensures that the wait group completes.
// This is useful when you want to wait for dynamically created tasks
// (i.e. async api executions) to complete before exiting.
func New(wg *sync.WaitGroup, blockOnClose bool) run.Runnable {
	return &wait{
		blockOnClose: blockOnClose,
		wg:           wg,
		done:         make(chan struct{}),
	}
}

func (w *wait) Run(context.Context) error {
	<-w.done
	w.wg.Wait()

	return nil
}

func (w *wait) Name() string { return "wait group reaper" }

func (w *wait) Close(context.Context) error {
	if w.blockOnClose {
		w.wg.Wait()
	}

	close(w.done)
	return nil
}
