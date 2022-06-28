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

package ranges_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/dp/cloud/go/active_mode_controller/internal/ranges"
)

func TestDecomposeOverlapping(t *testing.T) {
	const minValue = -200
	testData := []struct {
		name     string
		ranges   []ranges.Range
		expected []ranges.Point
	}{{
		name:     "Should convert empty list to empty list",
		ranges:   nil,
		expected: nil,
	}, {
		name: "Should convert single channel to points",
		ranges: []ranges.Range{
			{Begin: 35600, End: 35700, Value: 30},
		},
		expected: []ranges.Point{
			{Pos: 35600, Value: minValue},
			{Pos: 35700, Value: 30},
		},
	}, {
		name: "Should join connected ranges with same eirp",
		ranges: []ranges.Range{
			{Begin: 35600, End: 35700, Value: 30},
			{Begin: 35700, End: 35800, Value: 30},
		},
		expected: []ranges.Point{
			{Pos: 35600, Value: minValue},
			{Pos: 35800, Value: 30},
		},
	}, {
		name: "Should join connected ranges with different eirp",
		ranges: []ranges.Range{
			{Begin: 35600, End: 35700, Value: 30},
			{Begin: 35700, End: 35800, Value: 20},
		},
		expected: []ranges.Point{
			{Pos: 35600, Value: minValue},
			{Pos: 35700, Value: 30},
			{Pos: 35800, Value: 20},
		},
	}, {
		name: "Should handle nested ranges",
		ranges: []ranges.Range{
			{Begin: 35600, End: 35800, Value: 20},
			{Begin: 35650, End: 35750, Value: 30},
		},
		expected: []ranges.Point{
			{Pos: 35600, Value: minValue},
			{Pos: 35650, Value: 20},
			{Pos: 35750, Value: 30},
			{Pos: 35800, Value: 20},
		},
	}, {
		name: "Should handle overlapping ranges",
		ranges: []ranges.Range{
			{Begin: 35600, End: 35800, Value: 20},
			{Begin: 35700, End: 35900, Value: 30},
		},
		expected: []ranges.Point{
			{Pos: 35600, Value: minValue},
			{Pos: 35700, Value: 20},
			{Pos: 35900, Value: 30},
		},
	}, {
		name: "Should handle disjoint ranges",
		ranges: []ranges.Range{
			{Begin: 35600, End: 35700, Value: 30},
			{Begin: 35800, End: 35900, Value: 30},
		},
		expected: []ranges.Point{
			{Pos: 35600, Value: minValue},
			{Pos: 35700, Value: 30},
			{Pos: 35800, Value: minValue},
			{Pos: 35900, Value: 30},
		},
	}, {
		name: "Should handle complex case",
		ranges: []ranges.Range{
			{Begin: 35500, End: 36300, Value: 10},
			{Begin: 35600, End: 36000, Value: 20},
			{Begin: 35800, End: 36200, Value: 25},
			{Begin: 36000, End: 36100, Value: 30},
			{Begin: 36500, End: 36700, Value: 25},
			{Begin: 36600, End: 36800, Value: 20},
		},
		expected: []ranges.Point{
			{Pos: 35500, Value: minValue},
			{Pos: 35600, Value: 10},
			{Pos: 35800, Value: 20},
			{Pos: 36000, Value: 25},
			{Pos: 36100, Value: 30},
			{Pos: 36200, Value: 25},
			{Pos: 36300, Value: 10},
			{Pos: 36500, Value: minValue},
			{Pos: 36700, Value: 25},
			{Pos: 36800, Value: 20},
		},
	}}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			actual := ranges.DecomposeOverlapping(tt.ranges, minValue)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
