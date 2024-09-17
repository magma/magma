/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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
	cn, ok := (*atomic.Value)(val).Load().(*cfgNode)
	if (!ok) || cn == nil || cn.cfgPtr == nil {
		cfg, fp, owfp := defaultCfgFactory().UpdateFromYml()
		fcInterval := cfg.FreshnessCheckInterval()
		if fcInterval < MinFreshnessCheckInterval {
			fcInterval = MinFreshnessCheckInterval
		}
		val.updateStore(cfg, fp, owfp, fcInterval)
		return cfg
	}
	if cn.notAfterTime.After(time.Now()) { // refresh, if needed
		update := false
		if len(cn.filePath) > 0 {
			fi, serr := os.Stat(cn.filePath)
			update = (serr == nil) && fi.ModTime().After(cn.updatedTime)
		}
		if (!update) && len(cn.owFilePath) > 0 {
			fi, serr := os.Stat(cn.owFilePath)
			update = (serr == nil) && fi.ModTime().After(cn.updatedTime)
		}
		if update {
			cfg := cn.cfgPtr
			cfg, fp, owfp := cfg.UpdateFromYml()
			val.updateStore(cfg, fp, owfp, cn.freshnessCheckInterval)
			return cfg
		}
	}
	return cn.cfgPtr
}

// Overwrite unconditionally swaps configs referred by val with the given config
func (val *AtomicStore) Overwrite(cfg StructuredConfig) {
	now := time.Now()
	cn := &cfgNode{
		cfgPtr:                 cfg,
		freshnessCheckInterval: cfg.FreshnessCheckInterval(),
		notAfterTime:           now.Add(cfg.FreshnessCheckInterval()),
	}
	(*atomic.Value)(val).Store(cn)
}

// updateStore stores new structured configs into atomic store cache
func (val *AtomicStore) updateStore(cfg StructuredConfig, cfgPath, cfgOverwritePath string, chkInterval time.Duration) {
	now := time.Now()
	cn := &cfgNode{
		cfgPtr:                 cfg,
		filePath:               cfgPath,
		owFilePath:             cfgOverwritePath,
		freshnessCheckInterval: chkInterval,
		updatedTime:            now,
		notAfterTime:           now.Add(chkInterval),
	}
	(*atomic.Value)(val).Store(cn)
}

type cfgNode struct {
	cfgPtr StructuredConfig
	filePath,
	owFilePath string
	freshnessCheckInterval time.Duration
	updatedTime            time.Time
	notAfterTime           time.Time // optimization to avoid calculation on every check
}
