/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package mconfig provides gateway Go support for cloud managed configuration (mconfig)
package mconfig

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"magma/orc8r/cloud/go/protos"
)

var (
	localConfig   atomic.Value // always *mconfig.GatewayConfigs, never should be nil
	cfgMu         sync.Mutex
	lastFileInfo  os.FileInfo
	lastFilePath  string
	refreshTicker *time.Ticker
)

func init() {
	localConfig.Store(new(protos.GatewayConfigs))
	RefreshConfigs()
	cfgMu.Lock()
	refreshTicker = time.NewTicker(MconfigRefreshInterval)
	cfgMu.Unlock()
	go func() {
		for {
			<-refreshTicker.C
			err := RefreshConfigs()
			if err == nil {
				log.Print("Mconfig refresh succeeded")
			} else {
				log.Printf("Mconfig refresh error: %v", err)
			}
		}
	}()
}

// RefreshConfigs checks if Managed Config File's path or content has changed
// and tries to reload mamaged configs from the file
// refreshConfigs is thread safe and can be safely called while current configs are in use by
// other threads/routines
func RefreshConfigs() error {
	dynamicConfigPath := configFilePath()
	err := RefreshConfigsFrom(dynamicConfigPath)
	if err != nil {
		log.Printf("Cannot load configs from %s: %v", dynamicConfigPath, err)
		err = RefreshConfigsFrom(defaultConfigFilePath())
	}
	return err
}

// RefreshConfigsFrom checks if Managed Config File mcpath has changed
// and tries to reload mamaged configs from the file
// RefreshConfigsFrom is thread safe and can be safely called while current configs are in use by
// other threads/routines
func RefreshConfigsFrom(mcpath string) error {
	cfgMu.Lock()
	defer cfgMu.Unlock()

	fi, err := os.Stat(mcpath)
	if err != nil {
		return fmt.Errorf("Managed Config File '%s' stat error: %v", mcpath, err)
	}
	if sameFile(lastFileInfo, fi) {
		return nil
	}
	err = loadFromFile(mcpath)
	if err != nil {
		return fmt.Errorf("Error Loading Managed Config File '%s': %v", mcpath, err)
	}
	lastFileInfo = fi
	lastFilePath = mcpath
	return nil
}

func sameFile(oldInfo, newInfo os.FileInfo) bool {
	return oldInfo != nil &&
		newInfo != nil &&
		os.SameFile(oldInfo, newInfo) &&
		oldInfo.ModTime() == newInfo.ModTime()
}

func configFileDir() string {
	mcdir := os.Getenv(ConfigFileDirEnv)
	if len(mcdir) == 0 {
		mcdir = DefaultDynamicConfigFileDir
	}
	return mcdir
}

func configFilePath() string {
	return filepath.Join(configFileDir(), MconfigFileName)
}

func defaultConfigFilePath() string {
	return filepath.Join(DefaultConfigFileDir, MconfigFileName)
}

func loadFromFile(path string) error {
	cont, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	mc := new(protos.GatewayConfigs)
	err = protos.Unmarshal(cont, mc)
	if err != nil {
		return err
	}
	if len(mc.GetConfigsByKey()) == 0 {
		return fmt.Errorf("Empty Managed Gateway Configs")
	}
	localConfig.Store(mc)
	return nil
}

// StopRefreshTicker stops refresh ticker
func StopRefreshTicker() {
	cfgMu.Lock()
	if refreshTicker != nil {
		refreshTicker.Stop()
	}
	cfgMu.Unlock()
}
