// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pubsub_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/facebookincubator/symphony/pkg/pubsub"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gocloud.dev/pubsub/mempubsub"
)

func TestProvider(t *testing.T) {
	u := &url.URL{
		Scheme: mempubsub.Scheme,
		Host:   uuid.New().String(),
	}
	cfg := pubsub.Config{PubURL: u, SubURL: u}
	t.Run("Emitter", func(t *testing.T) {
		emitter, shutdown, err := pubsub.ProvideEmitter(context.Background(), cfg)
		assert.NoError(t, err)
		assert.NotNil(t, emitter)
		assert.NotNil(t, shutdown)
	})
	t.Run("Subscriber", func(t *testing.T) {
		subscriber := pubsub.ProvideSubscriber(cfg)
		assert.NotNil(t, subscriber)
	})
}
