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

// Package eap (EAP Authenticator) provides interface to supported & registered EAP Authenticator Providers
package eap

import "sync/atomic"

const BitsPerUint = 32

// Duplicate detector provides atomic "locking" to support EAP Identifier based duplicate prevention
type duplicateDetector struct {
	mask  []uint32
	bitLn uint
}

// NewDuplicateDetector returns a new detector supporting bitLen possible combinations
func NewDuplicateDetector(bitLen uint) *duplicateDetector {
	ln := bitLen / BitsPerUint
	if bitLen%BitsPerUint != 0 {
		ln += 1
	}
	return &duplicateDetector{make([]uint32, ln), bitLen}
}

// Add "locks" identifier/type pair while in flight & prevents duplicates to be accepted while the previous message
// is in flight. Returns true if the addition is successful, false if the same identifier (bit) is already in flight
// All "out of bound" entities are considered to be "locked" and the API will act according to this convention
func (dd *duplicateDetector) Add(bit uint) bool {
	if bit >= dd.bitLn {
		return false // out of bounds - Add will always fail for out of bounds request
	}
	bitIdx := bit % BitsPerUint
	ptr := &dd.mask[bit/BitsPerUint]
	bitMask := uint32(1) << bitIdx
	for {
		mask := atomic.LoadUint32(ptr)
		if (mask & bitMask) != 0 {
			return false
		}
		if atomic.CompareAndSwapUint32(ptr, mask, mask|bitMask) {
			return true
		}
	}
}

// Remove "unlocks" identifier/type pair while in flight & prevents duplicates to be accepted while the previous message
// is in flight. Returns true if the bit was locked & swapped
func (dd *duplicateDetector) Remove(bit uint) bool {
	if bit >= dd.bitLn {
		return false // out of bounds - Remove will always fail for out of bounds request
	}
	bitIdx := bit % BitsPerUint
	ptr := &dd.mask[bit/BitsPerUint]
	bitMask := uint32(1) << bitIdx
	for {
		mask := atomic.LoadUint32(ptr)
		if (mask & bitMask) == 0 {
			return false
		}
		if atomic.CompareAndSwapUint32(ptr, mask, mask&(^bitMask)) {
			return true
		}
	}
}

// Check verifies if a request with given identity (bit) is in flight and returns true if it is
func (dd *duplicateDetector) Check(bit uint) bool {
	if bit >= dd.bitLn {
		return true // out of bounds - act as if it's "locked"
	}
	return atomic.LoadUint32(&dd.mask[bit/BitsPerUint])&(uint32(1)<<bit%BitsPerUint) != 0
}
