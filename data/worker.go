package data

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type worker struct {
	DataFetcher
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewWorker(fetcher DataFetcher) JobWorker {
	ctx, cancel := context.WithCancel(context.Background())
	w := &worker{
		DataFetcher: fetcher,
		ctx:         ctx,
		cancel:      cancel,
	}
	w.wg.Add(1)
	return w
}

func (w *worker) Run() {
	defer w.wg.Done()
	for {
		select {
		case <-w.ctx.Done():
			fmt.Println("Worker stopped")
			return
		default:
			fmt.Println("Worker running...")
			time.Sleep(1 * time.Second) // Simulate work
		}
	}
}

// Stop cancels the Goroutine
func (w *worker) Stop() {
	w.cancel()  // Signal the Goroutine to stop
	w.wg.Wait() // Wait for cleanup
	fmt.Println("Worker completely stopped")
}
