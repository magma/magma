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

func NewFileLocker(fs DirectoryClient) (*FileLocker, error) {
	fileLocks := make(map[string]*sync.RWMutex)

	_, err := fs.Stat()
	if err != nil {
		_ = fs.Mkdir(0766)
	}

	files, err := fs.ReadDir()
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		fullPath := fs.Dir() + "/" + f.Name()
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
		defer f.selfMutex.Unlock()
		mtx, ok := f.fileLocks[filename]
		if !ok {
			f.fileLocks[filename] = &sync.RWMutex{}
			f.fileLocks[filename].Lock()
			return
		}
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
		defer f.selfMutex.Unlock()
		mtx, ok := f.fileLocks[filename]
		if !ok {
			f.fileLocks[filename] = &sync.RWMutex{}
			f.fileLocks[filename].RLock()
			return
		}
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

// DirectoryClient provides the necessary functions to read and modify a single
// directory for the FileLocker to operate
type DirectoryClient interface {
	Stat() (os.FileInfo, error)
	Mkdir(perm os.FileMode) error
	ReadDir() ([]os.FileInfo, error)
	Dir() string
}

type dirClient struct {
	rulesDir string
}

func (f *dirClient) Stat() (os.FileInfo, error) {
	return os.Stat(f.rulesDir)
}

func (f *dirClient) Mkdir(perm os.FileMode) error {
	return os.Mkdir(f.rulesDir, perm)
}

func (f *dirClient) ReadDir() ([]os.FileInfo, error) {
	return ioutil.ReadDir(f.rulesDir)
}

func (f *dirClient) Dir() string {
	return f.rulesDir
}

func NewDirectoryClient(rulesDir string) DirectoryClient {
	return &dirClient{
		rulesDir: rulesDir,
	}
}
