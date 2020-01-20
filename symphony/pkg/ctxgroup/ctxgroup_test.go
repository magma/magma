// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ctxgroup

import (
	"context"
	"errors"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
)

func TestParallel(t *testing.T) {
	var (
		v = make([]int, 128)
		g = WithContext(context.Background())
	)
	for i := 1; i <= len(v); i++ {
		i := i
		g.Go(func(context.Context) error {
			v[i-1] = i
			return nil
		})
	}
	assert.NoError(t, g.Wait())
	var sum int
	for i := range v {
		sum += v[i]
	}
	assert.Equal(t, (1+len(v))*len(v)/2, sum)
}

func TestLimited(t *testing.T) {
	ctx := context.Background()
	t.Run("Execution", func(t *testing.T) {
		limit := int64(runtime.NumCPU())
		g := WithContext(ctx, MaxConcurrency(limit))
		sem := semaphore.NewWeighted(limit)
		for i := int64(0); i < limit*4; i++ {
			g.Go(func(context.Context) error {
				if ok := sem.TryAcquire(1); !ok {
					return errors.New("acquiring semaphore")
				}
				sem.Release(1)
				return nil
			})
		}
		assert.NoError(t, g.Wait())
	})
	t.Run("BadLimit", func(t *testing.T) {
		assert.Panics(t, func() {
			WithContext(ctx, MaxConcurrency(0))
		})
	})
}

func TestZero(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")

	tests := []struct {
		errs []error
	}{
		{errs: []error{}},
		{errs: []error{nil}},
		{errs: []error{err1}},
		{errs: []error{err1, nil}},
		{errs: []error{err1, nil, err2}},
	}
	for _, tt := range tests {
		var (
			g = WithContext(context.Background())
			e error
		)
		for _, err := range tt.errs {
			err := err
			g.Go(func(ctx context.Context) error { return err })
			if e == nil && err != nil {
				e = err
			}
			assert.Equal(t, e, g.Wait())
		}
	}
}

func TestWithContext(t *testing.T) {
	g := WithContext(context.Background())
	g.Go(func(ctx context.Context) error {
		return errors.New("execution failure")
	})
	var err error
	g.Go(func(ctx context.Context) error {
		<-ctx.Done()
		err = ctx.Err()
		return err
	})
	_ = g.Wait()
	assert.EqualError(t, err, context.Canceled.Error())
}
