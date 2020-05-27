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
	"unsafe"

	"magma/gateway/config"
	"magma/orc8r/lib/go/protos"
)

var (
	localConfig   unsafe.Pointer
	cfgMu         sync.RWMutex
	lastFileInfo  os.FileInfo
	lastFilePath  string
	refreshTicker *time.Ticker
)

func init() {
	cfgMu.Lock()
	refreshTicker = time.NewTicker(MconfigRefreshInterval)
	cfgMu.Unlock()
	go func() {
		for {
			<-refreshTicker.C
			cfgPath, err := RefreshConfigs()
			if err == nil {
				log.Print("Mconfig refresh succeeded from: ", cfgPath)
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
func RefreshConfigs() (string, error) {
	// get dynamic config path
	configPath := ConfigFilePath()
	err := RefreshConfigsFrom(configPath)
	if err != nil {
		log.Printf("Cannot load configs from %s: %v", configPath, err)
		configPath = DefaultConfigFilePath()
		err = RefreshConfigsFrom(configPath)
	}
	return configPath, err
}

// ConfigFilePath returns current GW mconfig file path
func ConfigFilePath() string {
	return filepath.Join(configFileDir(), MconfigFileName)
}

// DefaultConfigFilePath returns default GW mconfig file path
func DefaultConfigFilePath() string {
	return filepath.Join(staticConfigFileDir(), MconfigFileName)
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
		mcdir = config.GetMagmadConfigs().DynamicMconfigDir
		if len(mcdir) == 0 {
			mcdir = DefaultDynamicConfigFileDir
		}
	}
	return mcdir
}

func staticConfigFileDir() string {
	mcdir := config.GetMagmadConfigs().StaticMconfigDir
	if len(mcdir) == 0 {
		mcdir = DefaultConfigFileDir
	}
	return mcdir
}

func loadFromFile(path string) error {
	cont, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	mc := new(protos.GatewayConfigs)
	err = protos.UnmarshalMconfig(cont, mc)
	if err != nil {
		return err
	}
	if len(mc.GetConfigsByKey()) == 0 {
		return fmt.Errorf("Empty Managed Gateway Configs")
	}
	atomic.StorePointer(&localConfig, (unsafe.Pointer)(mc))
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

// Info returns last used mconfig file information
func Info() (fullPath string, fileInfo os.FileInfo) {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return lastFilePath, lastFileInfo
}
