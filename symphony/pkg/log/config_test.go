// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"bytes"
	"log"
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestConfigParse(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		var config Config
		_, err := flags.ParseArgs(&config, nil)
		assert.NoError(t, err)
		assert.Equal(t, "info", config.Level.String())
		assert.Equal(t, "console", config.Format.String())
	})
	t.Run("OK", func(t *testing.T) {
		var config Config
		_, err := flags.ParseArgs(&config, []string{
			"--level", "error", "--format", "json",
		})
		assert.NoError(t, err)
		assert.Equal(t, "error", config.Level.String())
		assert.Equal(t, "json", config.Format.String())
	})
	t.Run("BadLevel", func(t *testing.T) {
		var cfg Config
		_, err := flags.ParseArgs(&cfg, []string{
			"--level", "foo",
		})
		assert.Error(t, err)
	})
	t.Run("BadFormat", func(t *testing.T) {
		var cfg Config
		_, err := flags.ParseArgs(&cfg, []string{
			"--format", "bar",
		})
		assert.Error(t, err)
	})
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "Production",
			config: Config{
				Level:  AllowedLevel(zap.InfoLevel),
				Format: "json",
			},
		},
		{
			name: "Development",
			config: Config{
				Level:  AllowedLevel(zap.DebugLevel),
				Format: "console",
			},
		},
		{
			name:   "Nop",
			config: Config{},
		},
		{
			name: "BadFormat",
			config: Config{
				Format: "fmt",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			logger, err := New(tc.config)
			if !tc.wantErr {
				assert.NotNil(t, logger)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestMustNew(t *testing.T) {
	var config Config
	assert.NotPanics(t, func() { _ = MustNew(config) })
	config.Format = "baz"
	assert.Panics(t, func() { _ = MustNew(config) })
}

func TestProvider(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	logger, restorer, err := Provider(Config{})
	require.NoError(t, err)
	defer restorer()
	assert.Equal(t, logger.Background(), zap.L())
	log.Println("suppressed message")
	assert.Zero(t, buf.Len())
}
