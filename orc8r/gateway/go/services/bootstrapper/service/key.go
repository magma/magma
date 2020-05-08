/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package service implements the core of bootstrapper

package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"magma/gateway/config"
	"magma/orc8r/lib/go/security/key"
)

// GetChallengeKey reads and returns the bootstrapper challenge key if present
// or creates and returns the key if the key file doesn't exist
func GetChallengeKey() (privKey interface{}, err error) {
	challengeKeyFile := config.GetMagmadConfigs().BootstrapConfig.ChallengeKey
	privKey, err = key.ReadKey(challengeKeyFile)
	if err == nil {
		return // all good, return the key
	}
	log.Printf("Bootstrapper ReadKey(%s) error: %v", challengeKeyFile, err)

	if os.IsNotExist(err) { // file doesn't exist, try to create it
		privKey, err = key.GenerateKey(PrivateKeyType, 0)
		if err != nil {
			err = fmt.Errorf("Bootstrapper Generate Key error: %v", err)
			return
		}
		dir := filepath.Dir(challengeKeyFile)
		if len(dir) > 3 {
			os.MkdirAll(dir, os.ModePerm)
		}
		err = key.WriteKey(challengeKeyFile, privKey)
		if err != nil {
			err = fmt.Errorf("Bootstrapper Write Key (%s) error: %v", challengeKeyFile, err)
			return
		}
		privKey, err = key.ReadKey(challengeKeyFile)
		if err != nil {
			err = fmt.Errorf(
				"Bootstrapper Failed to read recently created key from (%s) error: %v", challengeKeyFile, err)
			return
		}
	}
	return
}
