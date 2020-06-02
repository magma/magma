// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pubsub

import (
	"context"
	"flag"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/mempubsub"
)

func TestURLSubscriber(t *testing.T) {
	t.Run("Open", func(t *testing.T) {
		assert.Implements(t, (*Subscriber)(nil), new(URLSubscriber))

		var (
			ctx  = context.Background()
			name = mempubsub.Scheme + "://" + uuid.New().String()
		)
		topic, err := pubsub.OpenTopic(ctx, name)
		require.NoError(t, err)
		defer topic.Shutdown(ctx)

		subscription, err := NewURLSubscriber(name).Subscribe(ctx)
		require.NoError(t, err)
		subscription.Shutdown(ctx)
	})

	t.Run("Flag", func(t *testing.T) {
		var subscriber URLSubscriber
		assert.Implements(t, (*flag.Value)(nil), &subscriber)

		const goodURL = "file://test"
		err := subscriber.Set(goodURL)
		assert.NoError(t, err)
		assert.Equal(t, goodURL, subscriber.String())

		var badURL = string([]byte{0x7f})
		err = subscriber.Set(badURL)
		assert.Error(t, err)
	})
}

func TestNopSubscriber(t *testing.T) {
	subscriber := NewNopSubscriber()
	_, err := subscriber.Subscribe(context.Background())
	assert.EqualError(t, err, "nop subscriber")
}
