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

func getFNV1Hash(key string) uint32 {
	// FNV-1 is a non-cryptographic hash function that is fast and very simple to implement
	h := fnv.New32a()
	_, err := h.Write([]byte(key))
	if err != nil {
		return 0
	}
	return h.Sum32()
}

// JitterUint32 returns a jitter of the given uint32 value that is deterministic
// based on the given key.
func JitterUint32(n uint32, key string, maxMultiplier float32) uint32 {
	fnv1Hash := getFNV1Hash(key)
	multiplier := maxMultiplier * float32(fnv1Hash%100) / 100.0
	if multiplier >= 1 {
		return n
	}
	return uint32(multiplier * float32(n))
}

// JitterInt64 returns a jitter of the given int64 value that is deterministic
// based on the given key.
func JitterInt64(n int64, key string, maxMultiplier float64) int64 {
	fnv1Hash := getFNV1Hash(key)
	multiplier := maxMultiplier * float64(fnv1Hash%100) / 100.0
	if multiplier >= 1 {
		return n
	}
	return int64(multiplier * float64(n))
}
