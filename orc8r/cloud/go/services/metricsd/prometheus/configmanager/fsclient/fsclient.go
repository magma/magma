/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package fsclient

import (
	"io/ioutil"
	"os"
)

type FSClient interface {
	WriteFile(filename string, data []byte, perm os.FileMode) error
	ReadFile(filename string) ([]byte, error)
	Stat(filename string) (os.FileInfo, error)
}

type fsclient struct{}

func NewFSClient() FSClient {
	return &fsclient{}
}

func (f *fsclient) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}

func (f *fsclient) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (f *fsclient) Stat(filename string) (os.FileInfo, error) {
	return os.Stat(filename)
}
