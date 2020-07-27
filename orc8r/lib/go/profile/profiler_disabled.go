// +build !with_profiler

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
