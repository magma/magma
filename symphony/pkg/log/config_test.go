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

func TestConfigParsing(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		var cfg Config
		_, err := flags.ParseArgs(&cfg, nil)
		assert.NoError(t, err)
		assert.EqualValues(t, zap.InfoLevel, cfg.Level)
		assert.Equal(t, "console", cfg.Format)
	})
	t.Run("ok", func(t *testing.T) {
		var cfg Config
		_, err := flags.ParseArgs(&cfg, []string{
			"--level", "error", "--format", "json",
		})
		assert.NoError(t, err)
		assert.EqualValues(t, zap.ErrorLevel, cfg.Level)
		assert.Equal(t, "json", cfg.Format)
	})
	t.Run("bad level", func(t *testing.T) {
		var cfg Config
		_, err := flags.ParseArgs(&cfg, []string{
			"--level", "garbage",
		})
		assert.Error(t, err)
	})
}

func TestLoggerBuild(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "Production",
			cfg:  Config{Level(zap.InfoLevel), "json"},
		},
		{
			name: "Development",
			cfg:  Config{Level(zap.DebugLevel), "console"},
		},
		{
			name: "Discard",
			cfg:  Config{},
		},
		{
			name:    "NoFormat",
			cfg:     Config{Level(zap.WarnLevel), ""},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			logger, err := tc.cfg.Build()
			if !tc.wantErr {
				assert.NotNil(t, logger)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestNewLogger(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	logger, restorer, err := New(Config{})
	require.NoError(t, err)
	defer restorer()
	assert.Equal(t, logger.Background(), zap.L())
	log.Println("suppressed message")
	assert.Zero(t, buf.Len())
}
