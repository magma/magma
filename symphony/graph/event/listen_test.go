// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/mempubsub"
)

func TestNewListener(t *testing.T) {
	ctx := context.Background()
	cfg := ListenerConfig{
		Subscriber: SubscriberFunc(func(context.Context) (*pubsub.Subscription, error) {
			return nil, nil
		}),
		Logger: zaptest.NewLogger(t),
		Tenant: "test",
		Events: []string{t.Name()},
		Handler: HandlerFunc(func(context.Context, string, []byte) error {
			return nil
		}),
	}
	listener, err := NewListener(ctx, cfg)
	assert.NoError(t, err)
	assert.NotNil(t, listener)

	t.Run("NoTenant", func(t *testing.T) {
		config := cfg
		config.Tenant = ""
		_, err := NewListener(ctx, config)
		assert.Error(t, err)
	})
	t.Run("NoEvents", func(t *testing.T) {
		config := cfg
		config.Events = nil
		_, err := NewListener(ctx, config)
		assert.Error(t, err)
	})
	t.Run("NoSubscription", func(t *testing.T) {
		config := cfg
		config.Subscriber = SubscriberFunc(func(context.Context) (*pubsub.Subscription, error) {
			return nil, errors.New("no subscription")
		})
		_, err := NewListener(ctx, config)
		assert.Error(t, err)
	})
	t.Run("NoLogger", func(t *testing.T) {
		config := cfg
		config.Logger = nil
		listener, err := NewListener(ctx, config)
		assert.NoError(t, err)
		assert.NotNil(t, listener)
	})
}

type testHandler struct {
	mock.Mock
}

func (t *testHandler) Handle(ctx context.Context, name string, data []byte) error {
	return t.Called(ctx, name, data).Error(0)
}

func TestSubscribeAndListen(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	topic := mempubsub.NewTopic()
	defer topic.Shutdown(context.Background())

	subscription := mempubsub.NewSubscription(topic, time.Second)
	subscriber := SubscriberFunc(func(context.Context) (*pubsub.Subscription, error) {
		return subscription, nil
	})

	const tenant = "test-tenant"
	testBody := []byte("test-body")
	var h testHandler
	h.On("Handle", ctx, t.Name(), mock.AnythingOfType("[]uint8")).
		Run(func(args mock.Arguments) {
			body := args.Get(2).([]byte)
			assert.Equal(t, testBody, body)
			cancel()
		}).
		Return(nil).
		Once()
	defer h.AssertExpectations(t)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := SubscribeAndListen(ctx,
			ListenerConfig{
				Subscriber: subscriber,
				Logger:     zaptest.NewLogger(t),
				Tenant:     tenant,
				Events:     []string{t.Name()},
				Handler:    &h,
			})
		require.True(t, errors.Is(err, context.Canceled))
	}()
	defer wg.Wait()

	msgs := []pubsub.Message{
		{
			Metadata: map[string]string{
				NameHeader: t.Name(),
			},
		},
		{
			Metadata: map[string]string{
				TenantHeader: tenant,
			},
		},
		{
			Metadata: map[string]string{
				TenantHeader: tenant,
				NameHeader:   t.Name(),
			},
			Body: testBody,
		},
	}
	for i := range msgs {
		err := topic.Send(ctx, &msgs[i])
		require.NoError(t, err)
	}
}
