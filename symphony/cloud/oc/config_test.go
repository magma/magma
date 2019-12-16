// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oc

import (
	"os"
	"testing"

	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		args   []string
		env    map[string]string
		expect func(*testing.T, *Options, error)
	}{
		{
			expect: func(t *testing.T, opts *Options, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 1.0, opts.SamplingProbability)
				assert.Nil(t, opts.Jaeger)
				assert.Equal(t, jaeger.Options{}, JaegerOptions(*opts))
			},
		},
		{
			args: []string{
				"--jaeger", `{"Process":{"ServiceName":"test"}}`,
			},
			env: map[string]string{
				"SAMPLING_PROBABILITY": "0.5",
			},
			expect: func(t *testing.T, opts *Options, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 0.5, opts.SamplingProbability)
				assert.Equal(t, "test", opts.Jaeger.Process.ServiceName)
			},
		},
		{
			args: []string{
				"--jaeger", `{"AgentEndpoint:"localhost:6831"}`,
			},
			expect: func(t *testing.T, _ *Options, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range tests {
		for k, v := range tc.env {
			err := os.Setenv(k, v)
			require.NoError(t, err)
		}
		var opts Options
		_, err := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash).ParseArgs(tc.args)
		tc.expect(t, &opts, err)
		for k := range tc.env {
			err := os.Unsetenv(k)
			require.NoError(t, err)
		}
	}
}
