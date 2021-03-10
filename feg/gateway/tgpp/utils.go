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

// Package tgpp implements 3GPP related utility function usable across 3GPP related protocols
package tgpp

import (
	"strconv"

	"github.com/golang/glog"
)

const (
	MinMncLen     = 2
	MaxMncLen     = 3
	DefaultMncLen = MinMncLen
)

// GetPlmnID returns PLMN ID part of given IMSI
func GetPlmnID(imsi string, mncLen int) []byte {
	if mncLen < MinMncLen {
		mncLen = MinMncLen
	} else if mncLen > MaxMncLen {
		mncLen = MaxMncLen
	}
	imsiBytes := [6]byte{}
	for i := 0; i < 6; i++ {
		v, err := strconv.Atoi(imsi[i : i+1])
		if err != nil {
			glog.Errorf("Invalid Digit '%s' in IMSI '%s': %v", imsi[i:i+1], imsi, err)
		}
		imsiBytes[i] = byte(v)
	}
	// see https://www.arib.or.jp/english/html/overview/doc/STD-T63v10_70/5_Appendix/Rel11/29/29272-bb0.pdf#page=73
	plmnId := [3]byte{
		imsiBytes[0] | (imsiBytes[1] << 4),
		imsiBytes[2] | (imsiBytes[5] << 4),
		imsiBytes[3] | (imsiBytes[4] << 4)}
	if mncLen < 3 {
		plmnId[1] |= 0xF0
	}
	return plmnId[:]
}

// DecodeMsisdn decdes TBCD encoded E.164 MSISDN into a string
func DecodeMsisdn(encoded []byte) string {
	res := make([]byte, 0, len(encoded)*2)
	for _, b := range encoded {
		n1, n2 := b&0xF, b>>4
		res = append(res, '0'+n1)
		if n2 != 0xF {
			res = append(res, '0'+n2)
		}
	}
	return string(res)
}
