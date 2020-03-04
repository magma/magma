/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config

import (
	"os"
	"sync/atomic"
	"time"
)

const MinFreshnessCheckInterval = time.Minute

type StructuredConfig interface {
	UpdateFromYml() (cfgPtr StructuredConfig, file, owFile string)
	FreshnessCheckInterval() time.Duration
}

type AtomicStore atomic.Value

// GetCurrentCfg returns current system/process wide configuration
// the returned pointer is guaranteed to not to be <nil> and MUST be used only in READ ONLY fashion.
// An attempt to modify the returned struct fields may lead to unpredictable results
// If a caller needs mutable copy of configuration, it must copy the returned object
func (val *AtomicStore) GetCurrent(defaultCfgFactory func() StructuredConfig) StructuredConfig {
	now := time.Now()
	cn, ok := (*atomic.Value)(val).Load().(*cfgNode)
	if (!ok) || cn == nil || cn.cfgPtr == nil {
		cfg, fp, owfp := defaultCfgFactory().UpdateFromYml()
		fcInterval := cfg.FreshnessCheckInterval()
		if fcInterval < MinFreshnessCheckInterval {
			fcInterval = MinFreshnessCheckInterval
		}
		cn = &cfgNode{
			cfgPtr:                 cfg,
			filePath:               fp,
			owFilePath:             owfp,
			freshnessCheckInterval: fcInterval,
			notAfterTime:           now.Add(fcInterval),
		}
		(*atomic.Value)(val).Store(cn)
		return cfg
	}
	if cn.notAfterTime.After(now) { // refresh, if needed
		update := false
		if len(cn.filePath) > 0 {
			fi, serr := os.Stat(cn.filePath)
			update = serr == nil && now.Sub(fi.ModTime()) < cn.freshnessCheckInterval
		}
		if (!update) && len(cn.owFilePath) > 0 {
			fi, serr := os.Stat(cn.owFilePath)
			update = serr == nil && now.Sub(fi.ModTime()) < cn.freshnessCheckInterval
		}
		if update {
			cfg := cn.cfgPtr
			cfg, fp, owfp := cfg.UpdateFromYml()
			cn = &cfgNode{
				cfgPtr: cfg, filePath: fp,
				owFilePath:             owfp,
				freshnessCheckInterval: cn.freshnessCheckInterval,
				notAfterTime:           now.Add(cn.freshnessCheckInterval),
			}
			(*atomic.Value)(val).Store(cn)
			return cfg
		}
	}
	return cn.cfgPtr
}

// OverwriteFor unconditionally swaps configs refered by val with the given config
func (val *AtomicStore) Overwrite(cfg StructuredConfig) {
	now := time.Now()
	cn := &cfgNode{
		cfgPtr:                 cfg,
		freshnessCheckInterval: cfg.FreshnessCheckInterval(),
		notAfterTime:           now.Add(cfg.FreshnessCheckInterval()),
	}
	(*atomic.Value)(val).Store(cn)
}

type cfgNode struct {
	cfgPtr StructuredConfig
	filePath,
	owFilePath string
	freshnessCheckInterval time.Duration
	notAfterTime           time.Time // optimization to avoid calculation on every check
}
