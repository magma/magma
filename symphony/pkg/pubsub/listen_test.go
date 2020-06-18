// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pubsub_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/facebookincubator/symphony/pkg/pubsub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	gcpubsub "gocloud.dev/pubsub"
	"gocloud.dev/pubsub/mempubsub"
)

func TestNewListener(t *testing.T) {
	ctx := context.Background()
	cfg := pubsub.ListenerConfig{
		Subscriber: pubsub.SubscriberFunc(func(context.Context) (*gcpubsub.Subscription, error) {
			return nil, nil
		}),
		Logger: zaptest.NewLogger(t),
		Tenant: pointer.ToString("test"),
		Events: []string{t.Name()},
		Handler: pubsub.HandlerFunc(func(context.Context, string, string, []byte) error {
			return nil
		}),
	}
	listener, err := pubsub.NewListener(ctx, cfg)
	assert.NoError(t, err)
	assert.NotNil(t, listener)

	t.Run("NoEvents", func(t *testing.T) {
		config := cfg
		config.Events = nil
		_, err := pubsub.NewListener(ctx, config)
		assert.Error(t, err)
	})
	t.Run("NoSubscription", func(t *testing.T) {
		config := cfg
		config.Subscriber = pubsub.SubscriberFunc(func(context.Context) (*gcpubsub.Subscription, error) {
			return nil, errors.New("no subscription")
		})
		_, err := pubsub.NewListener(ctx, config)
		assert.Error(t, err)
	})
	t.Run("NoLogger", func(t *testing.T) {
		config := cfg
		config.Logger = nil
		listener, err := pubsub.NewListener(ctx, config)
		assert.NoError(t, err)
		assert.NotNil(t, listener)
	})
}

type testHandler struct {
	mock.Mock
}

func (t *testHandler) Handle(ctx context.Context, tenant string, name string, data []byte) error {
	return t.Called(ctx, tenant, name, data).Error(0)
}

func TestSubscribeAndListen(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	topic := mempubsub.NewTopic()
	defer topic.Shutdown(context.Background())

	subscription := mempubsub.NewSubscription(topic, time.Second)
	subscriber := pubsub.SubscriberFunc(func(context.Context) (*gcpubsub.Subscription, error) {
		return subscription, nil
	})

	const tenant = "test-tenant"
	testBody := []byte("test-body")
	var h testHandler
	h.On("Handle", ctx, tenant, t.Name(), mock.AnythingOfType("[]uint8")).
		Run(func(args mock.Arguments) {
			body := args.Get(3).([]byte)
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
		err := pubsub.SubscribeAndListen(ctx,
			pubsub.ListenerConfig{
				Subscriber: subscriber,
				Logger:     zaptest.NewLogger(t),
				Tenant:     pointer.ToString(tenant),
				Events:     []string{t.Name()},
				Handler:    &h,
			})
		require.True(t, errors.Is(err, context.Canceled))
	}()
	defer wg.Wait()

	msgs := []gcpubsub.Message{
		{
			Metadata: map[string]string{
				pubsub.NameHeader: t.Name(),
			},
		},
		{
			Metadata: map[string]string{
				pubsub.TenantHeader: tenant,
			},
		},
		{
			Metadata: map[string]string{
				pubsub.TenantHeader: tenant,
				pubsub.NameHeader:   t.Name(),
			},
			Body: testBody,
		},
	}
	for i := range msgs {
		err := topic.Send(ctx, &msgs[i])
		require.NoError(t, err)
	}
}
