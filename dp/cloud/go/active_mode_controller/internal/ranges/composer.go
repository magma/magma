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

import "magma/dp/cloud/go/active_mode_controller/internal/containers"

func ComposeForMidpoints(points []Point, length int, minValue int) []Range {
	if len(points) == 0 {
		return nil
	}
	var res []Range
	i, j := -1, 0
	last := points[len(points)-1].Pos
	pos := points[0].Pos - length
	mq := &containers.MinQueue{}
	mq.Push(points[0].Value)
	for {
		moveBegin := points[i+1].Pos <= points[j].Pos-length
		delta := 0
		if moveBegin {
			i++
			delta = points[i].Pos - pos
		} else {
			delta = points[j].Pos - pos - length
			j++
		}
		if v := mq.Top(); v >= minValue {
			res = addRange(res, Range{
				Begin: pos + length/2,
				End:   pos + delta + length/2,
				Value: v,
			})
		}
		pos += delta
		if !moveBegin && pos+length == last {
			break
		}
		if moveBegin {
			mq.Pop()
		} else {
			mq.Push(points[j].Value)
		}
	}
	return res
}

func addRange(ranges []Range, r Range) []Range {
	i := len(ranges) - 1
	if i >= 0 &&
		ranges[i].End == r.Begin &&
		r.Value == ranges[i].Value {
		ranges[i].End = r.End
		return ranges
	}
	return append(ranges, r)
}
