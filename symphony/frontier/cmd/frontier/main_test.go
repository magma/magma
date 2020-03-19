// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigParsing(t *testing.T) {
	tests := []struct {
		args   []string
		env    map[string]string
		expect func(*testing.T, *cliFlags, error)
	}{
		{
			args: []string{
				"--proxy-target", "http://proxy.me",
				"--inventory-target", "http://inventory.me",
			},
			env: map[string]string{
				"KEY_PAIRS": "lRKUN5SKyFOqvn81vIrvX7ppRsaeC36F,eyTzd21GdaVzKKznwYHSeOYX3DnKXuzI",
			},
			expect: func(t *testing.T, cfg *cliFlags, err error) {
				require.NoError(t, err)
				assert.Equal(t, "http://proxy.me", cfg.ProxyTarget.String())
				assert.Equal(t, "http://inventory.me", cfg.InventoryTarget.String())
				require.Len(t, cfg.KeyPairs, 2)
				assert.EqualValues(t, "lRKUN5SKyFOqvn81vIrvX7ppRsaeC36F", cfg.KeyPairs[0])
				assert.EqualValues(t, "eyTzd21GdaVzKKznwYHSeOYX3DnKXuzI", cfg.KeyPairs[1])
			},
		},
		{
			args: []string{
				"--proxy-target", "http://proxy.me",
				"--inventory-target", "http://inventory.me",
			},
			env: map[string]string{
				"KEY_PAIRS": "2sqGIo70vqONlkW58lq3nScxsDlGZTvR",
			},
			expect: func(t *testing.T, cfg *cliFlags, err error) {
				require.NoError(t, err)
				require.Len(t, cfg.KeyPairs, 1)
				assert.EqualValues(t, "2sqGIo70vqONlkW58lq3nScxsDlGZTvR", cfg.KeyPairs[0])
			},
		},
		{
			args: []string{
				"--proxy-target", "http://proxy.me",
				"--inventory-target", "http://inventory.me",
			},
			expect: func(t *testing.T, _ *cliFlags, err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tc := range tests {
		for k, v := range tc.env {
			err := os.Setenv(k, v)
			require.NoError(t, err)
		}
		var cfg cliFlags
		_, err := flags.NewParser(&cfg, flags.HelpFlag|flags.PassDoubleDash).
			ParseArgs(tc.args)
		tc.expect(t, &cfg, err)
		for k := range tc.env {
			err := os.Unsetenv(k)
			require.NoError(t, err)
		}
	}
}

func TestURLUnmarshal(t *testing.T) {
	tests := []struct {
		name   string
		target string
		expect func(*testing.T, target, error)
	}{
		{
			name:   "OK",
			target: "http://example.com",
			expect: func(t *testing.T, target target, err error) {
				require.NoError(t, err)
				assert.Equal(t, "http", target.Scheme)
				assert.Equal(t, "example.com", target.Host)
			},
		},
		{
			name:   "NoScheme",
			target: "example.com",
			expect: func(t *testing.T, _ target, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:   "BadUrl",
			target: string([]byte{0x7f}),
			expect: func(t *testing.T, _ target, err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var cfg struct {
				Target target `long:"target"`
			}
			_, err := flags.ParseArgs(&cfg,
				[]string{"--target", tc.target},
			)
			tc.expect(t, cfg.Target, err)
		})
	}
}
