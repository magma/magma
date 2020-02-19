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

	"magma/orc8r/lib/go/security/key"
)

func (b *Bootstrapper) getChallengeKey() (privKey interface{}, err error) {
	privKey, err = key.ReadKey(b.ChallengeKeyFile)
	if err == nil {
		return // all good, return the key
	}
	log.Printf("Bootstrapper ReadKey(%s) error: %v", b.ChallengeKeyFile, err)
	privKey, err = key.GenerateKey(PrivateKeyType, 0)
	if err != nil {
		err = fmt.Errorf("Bootstrapper Generate Key error: %v", err)
		return
	}
	dir := filepath.Dir(b.ChallengeKeyFile)
	if len(dir) > 3 {
		os.MkdirAll(dir, os.ModePerm)
	}
	err = key.WriteKey(b.ChallengeKeyFile, privKey)
	if err != nil {
		err = fmt.Errorf("Bootstrapper Write Key (%s) error: %v", b.ChallengeKeyFile, err)
		return
	}
	privKey, err = key.ReadKey(b.ChallengeKeyFile)
	if err != nil {
		err = fmt.Errorf(
			"Bootstrapper Failed to read recently created key from (%s) error: %v", b.ChallengeKeyFile, err)
		return
	}
	return
}
