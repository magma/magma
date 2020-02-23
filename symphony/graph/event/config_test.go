// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gocloud.dev/pubsub/mempubsub"
)

func TestConfigFlag(t *testing.T) {
	const goodURL = "file://test"
	var cfg Config
	err := cfg.UnmarshalFlag(goodURL)
	assert.NoError(t, err)
	assert.Equal(t, goodURL, cfg.String())

	var badURL = string([]byte{0x7f})
	err = cfg.UnmarshalFlag(badURL)
	assert.Error(t, err)
}

func TestProvide(t *testing.T) {
	cfg := Config{url: mempubsub.Scheme + "://" + uuid.New().String()}
	emitter, subscriber, err := Provide(context.Background(), cfg)
	assert.NoError(t, err)
	assert.NotNil(t, emitter)
	assert.NotNil(t, subscriber)

	cfg.url = string([]byte{0x7f})
	_, _, err = Provide(context.Background(), cfg)
	assert.Error(t, err)
}
