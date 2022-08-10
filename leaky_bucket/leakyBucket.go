package leaky_bucket

import (
	"context"
	"time"
)

func New(ctx context.Context, cancel context.CancelFunc, qps uint64) *LeakyBucket {
	l := &LeakyBucket{
		ctx:    ctx,
		cancel: cancel,
		Qps:    qps,
		bucket: make(chan struct{}, qps),
		hole:   make(chan struct{}, 1),
	}
	go func() {
		duration := time.Second / time.Duration(qps)

		lastDrop := time.Now()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if time.Since(lastDrop) < duration {
					continue
				}
				drip, ok := <-l.bucket
				if !ok {
					continue
				}
				select {
				case l.hole <- drip:
					lastDrop = time.Now()
				default:

				}
			}

		}

	}()
	return l
}

type LeakyBucket struct {
	ctx    context.Context
	cancel context.CancelFunc
	Qps    uint64
	bucket chan struct{}
	hole   chan struct{}
}

func (l *LeakyBucket) Request() bool {
	select {
	case l.bucket <- struct{}{}:
		<-l.hole
		return true
	default:
		return false
	}

}

func (l *LeakyBucket) Close() {
	l.cancel()
}
