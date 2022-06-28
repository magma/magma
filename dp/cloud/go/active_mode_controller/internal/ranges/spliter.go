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

func Split(ranges []Range, points []int) ([]Range, []Point) {
	var newRanges []Range
	var newPoints []Point
	i := 0
	for _, r := range ranges {
		end := r.End
		for ; i < len(points) && points[i] <= end; i++ {
			if points[i] < r.Begin {
				continue
			}
			newPoints = append(newPoints, Point{
				Pos:   points[i],
				Value: r.Value,
			})
			r.End = points[i] - 1
			newRanges = addRangeIfValid(newRanges, r)
			r.Begin = points[i] + 1
		}
		r.End = end
		newRanges = addRangeIfValid(newRanges, r)
	}
	return newRanges, newPoints
}

func addRangeIfValid(ranges []Range, r Range) []Range {
	if r.Begin > r.End {
		return ranges
	}
	return append(ranges, r)
}
