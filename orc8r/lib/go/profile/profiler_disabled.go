//
// Copyright (c) Facebook, Inc. and its affiliates.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

// +build !with_profiler

// Package profile provides CPU & memory profiling helper functions
// profiling is enabled by with_profiler build tag
package profile

import "os"

// empty stubs for disabled profiler builds

// MemWrite stub
func MemWrite() error {
	return nil
}

// CpuStart stub
func CpuStart() (*os.File, error) {
	return nil, nil
}

// CpuStop stub
func CpuStop(f *os.File) error {
	return nil
}
