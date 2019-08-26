/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package ocstats

import (
	"errors"
	"io"
	"testing"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewHandlerCreation(t *testing.T) {
	core, o := observer.New(zap.WarnLevel)
	logger := zap.New(core)

	var opts prometheus.Options
	handler, closer, err := NewHandler(
		WithNamespace("test"),
		WithLogger(logger),
		WithProcessCollector(),
		WithGoCollector(),
		func(opt *prometheus.Options) error {
			opts = *opt
			return nil
		},
	)
	require.NotNil(t, handler)
	assert.NotNil(t, closer)
	assert.NoError(t, err)

	assert.Equal(t, "test", opts.Namespace)
	opts.OnError(io.ErrUnexpectedEOF)
	assert.Equal(t, 1, o.Len())
}

func TestNewHandlerBadOption(t *testing.T) {
	_, _, err := NewHandler(func(*prometheus.Options) error {
		return errors.New("bad option")
	})
	assert.EqualError(t, err, "applying option: bad option")
}
