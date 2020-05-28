// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"
)

func TestConfigParse(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		app := kingpin.New(t.Name(), "")
		config := AddFlags(app)
		_, err := app.Parse(nil)
		assert.NoError(t, err)
		assert.Equal(t, "info", config.Level.String())
		assert.Equal(t, "console", config.Format.String())
	})
	t.Run("OK", func(t *testing.T) {
		app := kingpin.New(t.Name(), "")
		config := AddFlags(app)
		_, err := app.Parse([]string{
			"--" + LevelFlagName, "error",
			"--" + FormatFlagName, "json",
		})
		assert.NoError(t, err)
		assert.Equal(t, "error", config.Level.String())
		assert.Equal(t, "json", config.Format.String())
	})
	t.Run("BadLevel", func(t *testing.T) {
		app := kingpin.New(t.Name(), "")
		_ = AddFlags(app)
		_, err := app.Parse([]string{
			"--" + LevelFlagName, "foo",
		})
		assert.Error(t, err)
	})
	t.Run("BadFormat", func(t *testing.T) {
		app := kingpin.New(t.Name(), "")
		_ = AddFlags(app)
		_, err := app.Parse([]string{
			"--" + FormatFlagName, "bar",
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
