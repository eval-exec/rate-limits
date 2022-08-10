package leaky_bucket_test

import (
	"context"
	"github.com/eval-exec/rate-limits/leaky_bucket"
	"sync/atomic"
	"testing"
	"time"
)

func TestLeakyBucket_Request(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	qpsLimit := uint64(10000)
	leakyBucket := leaky_bucket.New(ctx, cancelFunc, qpsLimit)
	defer leakyBucket.Close()

	{
		now := time.Now()
		leakyBucket.Request()
		t.Log(time.Since(now))
	}

	allReqCount, reqOkCount := uint64(0), uint64(0)

	now := time.Now()
	for time.Since(now) <= time.Second {
		allReqCount++
		if leakyBucket.Request() {
			atomic.AddUint64(&reqOkCount, 1)
		}
	}

	if reqOkCount > qpsLimit {
		t.Errorf("request success count is %d/%d, but qps limit is %d", reqOkCount, allReqCount, qpsLimit)
	}
	t.Logf("request ok count: %d/%d, qps limit: %d", reqOkCount, allReqCount, qpsLimit)

}
