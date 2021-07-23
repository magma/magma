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

package orc8r

import "hash/fnv"

// JitterUint32 returns a jitter of the given uint32 value that is deterministic
// based on the given key.
func JitterUint32(keyn uint32, key string, maxJitterMultiplier float32) uint32 {
	// FNV-1 is a non-cryptographic hash function that is fast and very simple to implement
	h := fnv.New32a()
	_, err := h.Write([]byte(key))
	if err != nil {
		return 0
	}
	multiplier := float32(h.Sum32()%100) / 100.0
	maxJitter := float32(keyn) * maxJitterMultiplier
	return uint32(multiplier * maxJitter)
}
