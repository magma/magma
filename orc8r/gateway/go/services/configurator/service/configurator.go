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

// package service implements the core of configurator
package service

import (
	"bytes"
	"encoding/json"
	"io"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"

	"magma/gateway/config"
	"magma/gateway/streamer"
	"magma/orc8r/lib/go/definitions"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/protos/mconfig"
)

// Configurator - magma configurator implementation
type Configurator struct {
	sync.RWMutex
	streamerClient     streamer.Client
	latestConfigDigest *protos.GatewayConfigsDigest
	updateChan         chan interface{}
}

// UpdateCompletion is a type sent to update channel with every successful mconfig update
// it includes a list of service names of services with changed configs
type UpdateCompletion []string

// NewConfigurator returns a new Configuration & attempts to feel in its configs from magmad.yml
func NewConfigurator(updateCompletionChan chan interface{}) *Configurator {
	return &Configurator{
		streamerClient: streamer.NewStreamerClient(nil),
		updateChan:     updateCompletionChan}
}

// Start registers mconfig listener and starts streaming
// It'll only return on error & will block in the streaming loop otherwise
func (c *Configurator) Start() error {
	c.RLock()
	cl := c.streamerClient
	c.RUnlock()
	return cl.Stream(c)
}

// Gateway Streamer Listener Interface Implementation
func (c *Configurator) GetName() string {
	return definitions.MconfigStreamName
}

func (c *Configurator) ReportError(e error) error {
	if e != io.EOF {
		glog.Errorf("gateway mconfig streaming error: %v", e)
		sleepDuration := time.Second * time.Duration(config.GetMagmadConfigs().ConfigStreamErrorRetryInterval)
		time.Sleep(sleepDuration)
	}
	return nil // continue streaming anyway
}

type anyResolver struct{}

func (c *Configurator) Update(ub *protos.DataUpdateBatch) bool {
	updates := ub.GetUpdates()
	if len(updates) == 0 {
		return true // keep waiting for new configs
	}
	// There should be only one ...
	update := updates[0]
	// Validate the received mconfig payload
	cfg := &protos.GatewayConfigs{}
	err := protos.UnmarshalMconfig(update.GetValue(), cfg)
	if err != nil {
		glog.Errorf("error unmarshaling mconfig update for GW %s: %v", update.GetKey(), err)
		return false // re-establish stream on error
	}
	mdCfg, ok := cfg.GetConfigsByKey()["magmad"]
	if !ok {
		glog.Errorf("invalid mconfig update for GW %s - missing magmad configuration", update.GetKey())
		return false
	}
	if err = ptypes.UnmarshalAny(mdCfg, new(mconfig.MagmaD)); err != nil {
		glog.Errorf("invalid magmad mconfig GW %s: %v", update.GetKey(), err)
		return false
	}
	// find out if any of the service configs changed
	updatedServices := UpdateCompletion{}
	newCfg := &rawMconfigMsg{ConfigsByKey: map[string]json.RawMessage{}}
	oldCfg := &rawMconfigMsg{ConfigsByKey: map[string]json.RawMessage{}}
	json.Unmarshal(update.GetValue(), newCfg)

	c.Lock() // lock on all file operations

	if oldCfgJson, err := readCfg(); err == nil {
		json.Unmarshal(oldCfgJson, oldCfg)
	}
	newMap, oldMap := newCfg.ConfigsByKey, oldCfg.ConfigsByKey
	for sn, val := range newMap {
		if len(sn) > 0 {
			oldVal, found := oldMap[sn]
			if found && bytes.Equal(val, oldVal) {
				continue
			}
			// old config didn't have it or had non matching values
			updatedServices = append(updatedServices, sn)
		}
	}
	// find all removed configs
	for sn, _ := range oldMap {
		if _, found := newMap[sn]; !found {
			updatedServices = append(updatedServices, sn)
		}
	}
	if len(updatedServices) > 0 {
		glog.V(1).Infof("changes detected in configs for services: %v", updatedServices)
		err = SaveConfigs(update.GetValue())
		if err != nil {
			glog.Errorf("error saving new gateway mconfig: %v", err)
			c.Unlock()
			return false
		}
		// check if we need to update static copy of configs & update them
		updateStaticConfigs(update.GetValue())
	} else {
		glog.V(1).Info("no changes in cloud provided configs")
	}
	// everything succeeded, update digest for the next config stream
	c.latestConfigDigest = nil
	digest := cfg.GetMetadata().GetDigest().GetMd5HexDigest()
	if len(digest) > 0 {
		c.latestConfigDigest = &protos.GatewayConfigsDigest{Md5HexDigest: digest}
	} else {
		// TODO(hcgatewood): GetMconfigDigest isn't supposed to be used in the gateway, move this fn to cloud
		if digest, err = cfg.GetMconfigDigest(); err == nil {
			c.latestConfigDigest = &protos.GatewayConfigsDigest{Md5HexDigest: digest}
		} else {
			glog.Errorf("error encoding mconfig digest: %v", err)
		}
	}
	updateChan := c.updateChan

	c.Unlock() // unlock before possibly blocking on updateChan

	// Prepare a list of service names with changed mconfigs and
	// notify with the list of updated service configs if requested (c.updateChan != nil)
	if updateChan != nil {
		updateChan <- updatedServices
	}
	return false
}

func (c *Configurator) GetExtraArgs() *any.Any {
	c.RLock()
	defer c.RUnlock()
	if c.latestConfigDigest != nil {
		extra, err := ptypes.MarshalAny(c.latestConfigDigest)
		if err == nil {
			return extra
		}
		glog.Errorf("error marshaling latest mconfig digest: %v", err)
	}
	return nil
}

type rawMconfigMsg struct {
	ConfigsByKey map[string]json.RawMessage
}
