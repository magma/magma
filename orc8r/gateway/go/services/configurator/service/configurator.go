/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package service implements the core of configurator
package service

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"sync"
	"time"

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
		log.Printf("gateway mconfig streaming error: %v", e)
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
	u := updates[0]
	// Validate the received mconfig payload
	cfg := &protos.GatewayConfigs{}
	err := protos.UnmarshalMconfig(u.GetValue(), cfg)
	if err != nil {
		log.Printf("error unmarshaling mconfig update for GW %s: %v", u.GetKey(), err)
		return false // re-establish stream on error
	}
	mdCfg, ok := cfg.GetConfigsByKey()["magmad"]
	if !ok {
		log.Printf("invalid mconfig update for GW %s - missing magmad configuration", u.GetKey())
		return false
	}
	if err = ptypes.UnmarshalAny(mdCfg, new(mconfig.MagmaD)); err != nil {
		log.Printf("invalid magmad mconfig GW %s: %v", u.GetKey(), err)
		return false
	}
	c.Lock()
	updateChan := c.updateChan
	oldCfgJson, err := SaveConfigs(u.GetValue(), updateChan != nil)
	if err != err {
		log.Printf("error saving new gateway mconfig: %v", err)
		c.Unlock()
		return false
	}
	// everything succeeded, update digest for the next config stream
	c.latestConfigDigest = nil
	digest := cfg.GetMetadata().GetDigest().GetMd5HexDigest()
	if len(digest) > 0 {
		c.latestConfigDigest = &protos.GatewayConfigsDigest{Md5HexDigest: digest}
	} else {
		if digest, err = cfg.GetMconfigDigest(); err == nil {
			c.latestConfigDigest = &protos.GatewayConfigsDigest{Md5HexDigest: digest}
		} else {
			log.Printf("error encoding mconfig digest: %v", err)
		}
	}
	c.Unlock()

	// Prepare a list of service names with changed mconfigs and
	// notify with the list of updated service configs if requested (c.updateChan != nil)
	if updateChan != nil {
		updatedServices := UpdateCompletion{}
		newCfg := &rawMconfigMsg{ConfigsByKey: map[string]json.RawMessage{}}
		oldCfg := &rawMconfigMsg{ConfigsByKey: map[string]json.RawMessage{}}
		json.Unmarshal(u.GetValue(), newCfg)
		json.Unmarshal(oldCfgJson, oldCfg)
		newMap, oldMap := newCfg.ConfigsByKey, oldCfg.ConfigsByKey
		for sn, val := range newMap {
			if len(sn) > 0 {
				oldVal, found := oldMap[sn]
				if !found { // old config didn't have it, service needs to be updated
					updatedServices = append(updatedServices, sn)
				} else if !bytes.Equal(val, oldVal) { // non-matching values
					updatedServices = append(updatedServices, sn)
				}
			}
		}
		// find all removed configs
		for sn, _ := range oldMap {
			if _, found := newMap[sn]; !found {
				updatedServices = append(updatedServices, sn)
			}
		}
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
		log.Printf("error marshaling latest mconfig digest: %v", err)
	}
	return nil
}

type rawMconfigMsg struct {
	ConfigsByKey map[string]json.RawMessage
}
