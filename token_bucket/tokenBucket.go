package token_bucket

import (
	"context"
	"sync/atomic"
	"time"
)

func New(ctx context.Context, cancel context.CancelFunc, qps uint64) *TokenBucket {
	return &TokenBucket{
		ctx:    ctx,
		cancel: cancel,
		Qps:    qps,
		tokens: qps,
	}

}

type TokenBucket struct {
	ctx    context.Context
	cancel context.CancelFunc
	Qps    uint64
	tokens uint64
}

func (r *TokenBucket) addTokens() {

	tk := time.NewTicker(time.Second / time.Duration(r.Qps))
	defer tk.Stop()

	for {
		select {
		case <-tk.C:
			if r.tokens >= r.Qps {
				continue
			}
			atomic.AddUint64(&r.tokens, 1)
		case <-r.ctx.Done():
			return
		}
	}

}

func (r *TokenBucket) GetToken() bool {
retry:
	tokens := atomic.LoadUint64(&r.tokens)
	if tokens > 0 {
		if !atomic.CompareAndSwapUint64(&r.tokens, tokens, tokens-1) {
			goto retry
		}
		return true
	}
	return false
}

func (r *TokenBucket) Close() {
	r.cancel()
}
