// +build with_profiler

/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package profile provides CPU & memory profiling helper functions
// profiling is enabled by with_profiler build tag
package profile

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

const (
	memoryProfFileFmt = "%s/memory_%s.pprof"
	cpuProfFileFmt    = "%s/cpu_%s.pprof"
	timeFmt           = "0102_15.04.05"
)

var profileDir = flag.String("profiles_dir", os.TempDir(), "Destination directory for profiler files")

// MemWrite creates memory profile file, invokes GC and collects & writes memory profile
func MemWrite() error {
	os.MkdirAll(*profileDir, os.ModeDir)
	fname := fmt.Sprintf(memoryProfFileFmt, *profileDir, time.Now().Format(timeFmt))
	f, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("could not create memory profile in '%s': %v", fname, err)
	}
	defer f.Close()
	runtime.GC() // get up-to-date statistics
	if err := pprof.Lookup("heap").WriteTo(f, 1); err != nil {
		return fmt.Errorf("could not write memory profile in '%s': %v", fname, err)
	}
	return nil
}

// CpuStart creates a new CPU profile file and starts CPU profiling
// CpuStart will only return a non-nil file pointer on success, so
// it should be safe to do:
//     f, _ := profile.CpuStart()
//     ... // run profiled logic
//     profile.CpuStop(f)
func CpuStart() (*os.File, error) {
	os.MkdirAll(*profileDir, os.ModeDir)
	fname := fmt.Sprintf(cpuProfFileFmt, *profileDir, time.Now().Format(timeFmt))
	f, err := os.Create(fname)
	if err != nil {
		return nil, fmt.Errorf("could not create CPU profile in '%s': %v", fname, err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		f.Close()
		return nil, fmt.Errorf("could not start CPU profile: ", err)
	}
	return f, nil
}

// CpuStop collects and writes started CPU profile and closes the profile file
func CpuStop(f *os.File) error {
	pprof.StopCPUProfile()
	if f != nil {
		return f.Close()
	}
	return nil
}
