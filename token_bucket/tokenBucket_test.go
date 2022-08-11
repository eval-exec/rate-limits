package token_bucket_test

import (
	"context"
	"github.com/eval-exec/rate-limits/token_bucket"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestLongPeriodRequest(t *testing.T) {

	re := require.New(t)
	ctx, cancelFunc := context.WithCancel(context.Background())

	qpsLimit := uint64(100)
	tokenBucket := token_bucket.New(ctx, cancelFunc, qpsLimit)
	defer tokenBucket.Close()

	timer := time.NewTimer(5 * time.Second)

	memo := make([]time.Time, 1000000)
	idx := 0
loop:
	for {
		select {
		case <-timer.C:
			break loop
		default:
			if tokenBucket.GetToken() {
				memo[idx] = time.Now()
				idx++
			}
		}
		time.Sleep(time.Duration(rand.Uint64()%10+1) * time.Nanosecond)
	}
	memo = memo[:idx]

	end := 0
	for i := 0; i < len(memo); i++ {
		for end < len(memo) && memo[end].Sub(memo[i]) < time.Second {
			end++
		}
		tQps := uint64(end - i)
		re.LessOrEqualf(tQps, qpsLimit+200, "got qps:%d, but qpsLimit is %d in (%s, %s)", tQps, qpsLimit, memo[i], memo[end-1])
	}

}

func TestTokenBucket(t *testing.T) {
	re := require.New(t)
	ctx, cancelFunc := context.WithCancel(context.Background())

	tokenBucket := token_bucket.New(ctx, cancelFunc, 10000)
	defer tokenBucket.Close()

	timer := time.NewTimer(time.Second)
	for i := 0; i < 10000; i++ {
		select {
		case <-timer.C:
			return
		default:
			re.True(tokenBucket.GetToken())
		}
	}

}

func TestTokenBucketFail(t *testing.T) {
	re := require.New(t)

	ctx, cancelFunc := context.WithCancel(context.Background())

	tokenBucket := token_bucket.New(ctx, cancelFunc, 10000)
	defer tokenBucket.Close()

	timer := time.NewTimer(time.Second)

	for i := 0; i < 10001; i++ {
		select {
		case <-timer.C:
			return
		default:
			if i < 10000 {
				re.True(tokenBucket.GetToken())
			} else {
				re.False(tokenBucket.GetToken())
			}
		}
	}

}
