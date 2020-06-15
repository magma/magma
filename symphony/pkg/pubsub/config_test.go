// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pubsub

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gocloud.dev/pubsub/mempubsub"
)

func TestConfigFlag(t *testing.T) {
	const goodURL = "file://test"
	var cfgURL URL
	err := cfgURL.Set(goodURL)
	assert.NoError(t, err)
	assert.Equal(t, goodURL, cfgURL.String())

	var badURL = string([]byte{0x7f})
	err = cfgURL.Set(badURL)
	assert.Error(t, err)
}

func TestProvider(t *testing.T) {
	cfg := newConfig(mempubsub.Scheme + "://" + uuid.New().String())
	t.Run("Emitter", func(t *testing.T) {
		emitter, shutdown, err := ProvideEmitter(context.Background(), cfg)
		assert.NoError(t, err)
		assert.NotNil(t, emitter)
		assert.NotNil(t, shutdown)
		_, _, err = ProvideEmitter(context.Background(), newConfig(string([]byte{0x7f})))
		assert.Error(t, err)
	})
	t.Run("Subscriber", func(t *testing.T) {
		subscriber := ProvideSubscriber(cfg)
		assert.NotNil(t, subscriber)
	})
}
