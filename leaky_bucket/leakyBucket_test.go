package leaky_bucket

import (
	"context"
	"reflect"
	"testing"
)

func TestLeakyBucket_Request(t *testing.T) {
	type fields struct {
		ctx    context.Context
		cancel context.CancelFunc
		Qps    uint64
		tokens uint64
		hole   chan struct {
		}
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LeakyBucket{
				ctx:    tt.fields.ctx,
				cancel: tt.fields.cancel,
				Qps:    tt.fields.Qps,
				tokens: tt.fields.tokens,
				hole:   tt.fields.hole,
			}
			if got := l.Request(); got != tt.want {
				t.Errorf("Request() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		ctx    context.Context
		cancel context.CancelFunc
		qps    uint64
	}
	tests := []struct {
		name string
		args args
		want *LeakyBucket
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.ctx, tt.args.cancel, tt.args.qps); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}