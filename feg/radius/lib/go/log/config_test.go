/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package log

import (
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
)

func TestConfigParsing(t *testing.T) {
	var cfg Config
	_, err := flags.ParseArgs(&cfg, []string{
		"--level", "error", "--format", "console",
	})
	assert.NoError(t, err)
	assert.Equal(t, "error", cfg.Level)
	assert.Equal(t, "console", cfg.Format)
}

func TestLoggerBuild(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "Production",
			cfg:  Config{"info", "json"},
		},
		{
			name: "Development",
			cfg:  Config{"debug", "console"},
		},
		{
			name: "Discard",
			cfg:  Config{},
		},
		{
			name:    "BadLevel",
			cfg:     Config{"bad", "console"},
			wantErr: true,
		},
		{
			name:    "NoFormat",
			cfg:     Config{"warn", ""},
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
