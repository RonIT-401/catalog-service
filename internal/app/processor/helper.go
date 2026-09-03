package processor

import (
	"context"
	"sync"
)

func Wrap(ctx context.Context, wg *sync.WaitGroup, cb func(context.Context)) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		select {
		case <-ctx.Done():
			return
		default:
			cb(ctx)
		}
	}()
}
