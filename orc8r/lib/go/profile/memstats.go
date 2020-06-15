//
// Copyright (c) Facebook, Inc. and its affiliates.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

// Package profile provides CPU & memory profiling helper functions
package profile

import (
	"bytes"
	"fmt"
	"runtime"

	"github.com/golang/glog"
)

// LogMemStats collects the process memory stats and logs them out @ INFO level
func LogMemStats() {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	glog.Info(MemStatsToString(memStats))
}

// MemStatsToString returns a string with formatted runtime.MemStats
func MemStatsToString(s *runtime.MemStats) string {
	if s == nil {
		return ""
	}
	var b = new(bytes.Buffer)
	fmt.Fprintf(b, "Allocated:%9d; Objects#:%6d; Stack:%7d; VM:%9d\nBySize:",
		s.Alloc, s.HeapObjects, s.StackSys, s.Sys)
	for i := len(s.BySize) - 1; i > 0; i-- {
		fmt.Fprintf(b, " %d*%d;", s.BySize[i].Size, s.BySize[i].Mallocs-s.BySize[i].Frees)
	}
	return b.String()
}
