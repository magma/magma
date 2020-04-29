/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package service implements the core of configurator
package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"magma/gateway/mconfig"
)

// SaveConfigFile saves new gateway.configs and returns old configuration if any
func SaveConfigs(cfgJson []byte, readOldCfg bool) (oldCfgJson []byte, err error) {
	if len(cfgJson) == 0 {
		return oldCfgJson, fmt.Errorf("empty gateway mconfigs")
	}
	mconfigPath := mconfig.ConfigFilePath()
	newMconfigPath := mconfigPath + ".new"
	oldMconfigPath := mconfigPath + ".old"
	err = ioutil.WriteFile(newMconfigPath, cfgJson, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to save mconfigs into %s: %v", newMconfigPath, err)
	}
	oerr := os.Rename(mconfigPath, oldMconfigPath)
	err = os.Rename(newMconfigPath, mconfigPath)
	if err != nil {
		err = fmt.Errorf("failed to move mconfigs from %s to %s: %v", newMconfigPath, mconfigPath, err)
		if oerr == nil { // roll back if previous rename succeeded
			os.Rename(oldMconfigPath, mconfigPath)
		}
	} else {
		log.Printf("successfully updated mconfig in %s", mconfigPath)
		if readOldCfg && oerr == nil {
			if oldCfgJson, oerr = ioutil.ReadFile(oldMconfigPath); oerr != nil {
				oldCfgJson = nil
			}
		}
	}
	return oldCfgJson, err
}
