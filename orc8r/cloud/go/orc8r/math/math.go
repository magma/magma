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

package math

import "hash/fnv"

// JitterUint32 jitters the given uint32 value deterministically, based on the
// given key.
func JitterUint32(n uint32, key string, maxMultiplier float32) uint32 {
	hash := getFNV1Hash(key)
	multiplier := maxMultiplier * float32(hash%100) / 100.0
	jittered := float32(n) * (1 + multiplier)
	// Check for integer overflow
	if uint32(jittered) < n {
		return n
	}

	return uint32(jittered)
}

// JitterInt64 jitters the given int64 value deterministically, based on the
// given key.
func JitterInt64(n int64, key string, maxMultiplier float64) int64 {
	hash := getFNV1Hash(key)
	multiplier := maxMultiplier * float64(hash%100) / 100.0
	jittered := float64(n) * (1 + multiplier)
	// Check for integer overflow
	if int64(jittered) < n {
		return n
	}

	return int64(jittered)
}

func getFNV1Hash(key string) uint32 {
	// FNV-1 is a non-cryptographic hash function that is fast and very simple to implement
	h := fnv.New32a()
	_, err := h.Write([]byte(key))
	if err != nil {
		return 0
	}
	return h.Sum32()
}
