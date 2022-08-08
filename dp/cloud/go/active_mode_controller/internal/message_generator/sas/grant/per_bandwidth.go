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

package grant

import "math/bits"

func SelectGrantsWithRedundancy(available uint32, grants uint32, pref []uint32, minWidth int, maxWidth int, index int) uint32 {
	newGrants := selectTwoGrants(grants, pref, minWidth, maxWidth, index)
	if newGrants != 0 {
		return newGrants
	}
	available &^= toCloseMask(grants, minWidth)
	available |= grants
	newGrants = selectTwoGrants(available, pref, minWidth, maxWidth, index)
	if newGrants != 0 {
		return grants | newGrants
	}
	return selectOneGrant(grants, pref, index)
}

func SelectGrantsWithoutRedundancy(available uint32, grants uint32, pref []uint32, index int) uint32 {
	if grants != 0 {
		return selectOneGrant(grants, pref, index)
	}
	return selectOneGrant(available, pref, index)
}

func selectTwoGrants(available uint32, pref []uint32, minWidth int, maxWidth int, index int) uint32 {
	for _, x := range pref {
		y := available & okMask(x, minWidth, maxWidth)
		if available&x == x && y != 0 {
			return x | nearest(x, y)
		}
	}
	a, m := available, uint32(0)
	for a != 0 {
		x := a & -a
		y := available & okMask(x, minWidth, maxWidth)
		if y != 0 {
			m |= x
		}
		a -= x
	}
	if m != 0 {
		x := nth(m, index)
		y := available & okMask(x, minWidth, maxWidth)
		return x | nearest(x, y)
	}
	return 0
}

func selectOneGrant(available uint32, pref []uint32, index int) uint32 {
	for _, x := range pref {
		if available&x == x {
			return x
		}
	}
	if available != 0 {
		return nth(available, index)
	}
	return 0
}

func okMask(x uint32, i int, j int) uint32 {
	return x<<(j+1) - x<<(i+1) + (x-1)>>i - (x-1)>>j
}

func nearest(x uint32, y uint32) uint32 {
	left, right := y&^(x-1), y&(x-1)
	l, m, r := bits.TrailingZeros32(left), bits.TrailingZeros32(x), bits.Len32(right)-1
	if right == 0 || (left != 0 && l-m < m-r) {
		return 1 << l
	}
	return 1 << r
}

func toCloseMask(x uint32, i int) uint32 {
	y := uint32(0)
	for ; i > 0; i-- {
		y |= x<<i | x>>i
	}
	return y
}

func nth(x uint32, i int) uint32 {
	i %= bits.OnesCount32(x)
	y := uint32(0)
	for ; i >= 0; i-- {
		y = x & -x
		x -= y
	}
	return y
}
