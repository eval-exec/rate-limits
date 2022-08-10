package token_bucket_test

import (
	"context"
	"github.com/eval-exec/rate-limits/token_bucket"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTokenBucket(t *testing.T) {
	re := require.New(t)
	ctx, cancelFunc := context.WithCancel(context.Background())

	tokenBucket := token_bucket.New(ctx, cancelFunc, 10000)
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
