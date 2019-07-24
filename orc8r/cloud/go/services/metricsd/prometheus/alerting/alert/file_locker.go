/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package alert

import (
	"io/ioutil"
	"os"
	"sync"
)

type FileLocker struct {
	fileLocks map[string]*sync.RWMutex
	selfMutex sync.Mutex
}

func NewFileLocker(rulesDir string) (*FileLocker, error) {
	fileLocks := make(map[string]*sync.RWMutex)

	_, err := os.Stat(rulesDir)
	if err != nil {
		_ = os.Mkdir(rulesDir, 0766)
	}

	files, err := ioutil.ReadDir(rulesDir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		fullPath := rulesDir + "/" + f.Name()
		fileLocks[fullPath] = &sync.RWMutex{}
	}
	return &FileLocker{
		fileLocks: fileLocks,
	}, nil
}

// Lock locks the mutex associated with the given filename for writing. If
// mutex does not exist in map yet, create one.
func (f *FileLocker) Lock(filename string) {
	mtx, ok := f.fileLocks[filename]
	if !ok {
		f.selfMutex.Lock()
		mtx, ok := f.fileLocks[filename]
		if !ok {
			f.fileLocks[filename] = &sync.RWMutex{}
			f.fileLocks[filename].Lock()
			return
		}
		f.selfMutex.Unlock()
		mtx.Lock()
		return
	}
	mtx.Lock()
}

// Unlock unlocks the mutex associated with the given filename for writing.
// No-op if mutex does not exist in map
func (f *FileLocker) Unlock(filename string) {
	if mutex, ok := f.fileLocks[filename]; ok {
		mutex.Unlock()
	}
}

// RLock locks the mutex associated with the given filename for reading.
// No-op if mutex does not exist in map
func (f *FileLocker) RLock(filename string) {
	mtx, ok := f.fileLocks[filename]
	if !ok {
		f.selfMutex.Lock()
		mtx, ok := f.fileLocks[filename]
		if !ok {
			f.fileLocks[filename] = &sync.RWMutex{}
			f.fileLocks[filename].RLock()
			return
		}
		f.selfMutex.Unlock()
		mtx.RLock()
		return
	}
	mtx.RLock()
}

// RUnlock unlocks the mutex associated with the given filename for reading.
// No-op if mutex does not exist in map
func (f *FileLocker) RUnlock(filename string) {
	if mutex, ok := f.fileLocks[filename]; ok {
		mutex.RUnlock()
	}
}
