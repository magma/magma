/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ranges

func Select(ranges []Range, n int, mod int) (Point, bool) {
	sum := 0
	for i := range ranges {
		sum += countDivisible(ranges[i].Begin, ranges[i].End, mod)
	}
	if sum == 0 {
		return Point{}, false
	}
	i := 0
	for n %= sum; n >= 0; i++ {
		n -= countDivisible(ranges[i].Begin, ranges[i].End, mod)
	}
	i--
	return Point{
		Value: ranges[i].Value,
		Pos:   (ranges[i].End/mod + n + 1) * mod,
	}, true
}

func countDivisible(a int, b int, mod int) int {
	x := b/mod - a/mod
	if a%mod == 0 {
		x++
	}
	return x
}
