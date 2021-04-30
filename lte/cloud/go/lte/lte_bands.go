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

package lte

import "fmt"

// LTEBand struct for converting EARFCN to Band
type LTEBand struct {
	ID            uint32
	Mode          DuplexMode
	CountEarfcn   uint32
	StartEarfcnDl uint32
	StartEarfcnUl uint32
}

// DuplexMode of LTE Band
type DuplexMode int32

const (
	// TDDMode
	TDDMode DuplexMode = iota
	// FDDMode
	FDDMode
)

var bands = [...]LTEBand{
	// FDDMode
	{ID: 1, Mode: FDDMode, StartEarfcnDl: 0, StartEarfcnUl: 18000, CountEarfcn: 600},
	{ID: 2, Mode: FDDMode, StartEarfcnDl: 600, StartEarfcnUl: 18600, CountEarfcn: 600},
	{ID: 3, Mode: FDDMode, StartEarfcnDl: 1200, StartEarfcnUl: 19200, CountEarfcn: 750},
	{ID: 4, Mode: FDDMode, StartEarfcnDl: 1950, StartEarfcnUl: 19950, CountEarfcn: 450},
	{ID: 5, Mode: FDDMode, StartEarfcnDl: 2400, StartEarfcnUl: 20400, CountEarfcn: 250},
	{ID: 6, Mode: FDDMode, StartEarfcnDl: 2650, StartEarfcnUl: 20560, CountEarfcn: 100},
        {ID: 7, Mode: FDDMode, StartEarfcnDl: 2750, StartEarfcnUl: 20750, CountEarfcn: 700},
	{ID: 8, Mode: FDDMode, StartEarfcnDl: 3450, StartEarfcnUl: 21450, CountEarfcn: 350},
	{ID: 9, Mode: FDDMode, StartEarfcnDl: 3800, StartEarfcnUl: 21800, CountEarfcn: 350},
	{ID: 10, Mode: FDDMode, StartEarfcnDl: 4150, StartEarfcnUl: 22150, CountEarfcn: 600},
	{ID: 11, Mode: FDDMode, StartEarfcnDl: 4750, StartEarfcnUl: 22750, CountEarfcn: 260},
        {ID: 12, Mode: FDDMode, StartEarfcnDl: 5010, StartEarfcnUl: 23010, CountEarfcn: 170},
        {ID: 13, Mode: FDDMode, StartEarfcnDl: 5180, StartEarfcnUl: 23180, CountEarfcn: 100},
        {ID: 14, Mode: FDDMode, StartEarfcnDl: 5280, StartEarfcnUl: 23280, CountEarfcn: 450},
        {ID: 17, Mode: FDDMode, StartEarfcnDl: 5730, StartEarfcnUl: 23730, CountEarfcn: 120},
        {ID: 18, Mode: FDDMode, StartEarfcnDl: 5850, StartEarfcnUl: 23850, CountEarfcn: 150},
        {ID: 19, Mode: FDDMode, StartEarfcnDl: 6000, StartEarfcnUl: 24000, CountEarfcn: 150},
        {ID: 20, Mode: FDDMode, StartEarfcnDl: 6150, StartEarfcnUl: 24150, CountEarfcn: 300},
        {ID: 21, Mode: FDDMode, StartEarfcnDl: 6450, StartEarfcnUl: 24450, CountEarfcn: 150},
        {ID: 22, Mode: FDDMode, StartEarfcnDl: 6600, StartEarfcnUl: 24600, CountEarfcn: 900},
        {ID: 23, Mode: FDDMode, StartEarfcnDl: 7500, StartEarfcnUl: 25500, CountEarfcn: 200},
	{ID: 24, Mode: FDDMode, StartEarfcnDl: 7700, StartEarfcnUl: 25700, CountEarfcn: 340},
        {ID: 25, Mode: FDDMode, StartEarfcnDl: 8040, StartEarfcnUl: 26040, CountEarfcn: 650},
        {ID: 26, Mode: FDDMode, StartEarfcnDl: 8690, StartEarfcnUl: 26690, CountEarfcn: 350},
        {ID: 27, Mode: FDDMode, StartEarfcnDl: 9040, StartEarfcnUl: 27040, CountEarfcn: 170},
	{ID: 28, Mode: FDDMode, StartEarfcnDl: 9210, StartEarfcnUl: 27210, CountEarfcn: 59376},
	{ID: 71, Mode: FDDMode, StartEarfcnDl: 68586, StartEarfcnUl: 133122, CountEarfcn: 59426},
	// TDDMode
	{ID: 33, Mode: TDDMode, StartEarfcnDl: 36000, CountEarfcn: 200},
	{ID: 34, Mode: TDDMode, StartEarfcnDl: 36200, CountEarfcn: 150},
	{ID: 35, Mode: TDDMode, StartEarfcnDl: 36350, CountEarfcn: 600},
	{ID: 36, Mode: TDDMode, StartEarfcnDl: 36950, CountEarfcn: 600},
	{ID: 37, Mode: TDDMode, StartEarfcnDl: 37550, CountEarfcn: 200},
	{ID: 38, Mode: TDDMode, StartEarfcnDl: 37750, CountEarfcn: 500},
	{ID: 39, Mode: TDDMode, StartEarfcnDl: 38250, CountEarfcn: 400},
	{ID: 40, Mode: TDDMode, StartEarfcnDl: 38650, CountEarfcn: 1000},
	{ID: 41, Mode: TDDMode, StartEarfcnDl: 39650, CountEarfcn: 1940},
	{ID: 42, Mode: TDDMode, StartEarfcnDl: 41590, CountEarfcn: 2000},
	{ID: 43, Mode: TDDMode, StartEarfcnDl: 43590, CountEarfcn: 2000},
	{ID: 44, Mode: TDDMode, StartEarfcnDl: 45590, CountEarfcn: 1000},
	{ID: 45, Mode: TDDMode, StartEarfcnDl: 46590, CountEarfcn: 200},
	{ID: 46, Mode: TDDMode, StartEarfcnDl: 46790, CountEarfcn: 7750},
	{ID: 47, Mode: TDDMode, StartEarfcnDl: 54540, CountEarfcn: 700},
	{ID: 48, Mode: TDDMode, StartEarfcnDl: 55240, CountEarfcn: 1500},
	{ID: 49, Mode: TDDMode, StartEarfcnDl: 56740, CountEarfcn: 1500},
	{ID: 50, Mode: TDDMode, StartEarfcnDl: 58240, CountEarfcn: 850},
	{ID: 51, Mode: TDDMode, StartEarfcnDl: 59090, CountEarfcn: 50},
	{ID: 52, Mode: TDDMode, StartEarfcnDl: 59140, CountEarfcn: 1000},
	// Adding Band #53 require changes in the python code cause it's
	// start_freq_dl is float value.
	//{ID: 53, Mode: TDDMode, StartEarfcnDl: 60140, CountEarfcn: 115},
}

// EarfcnDLInRange checks that an EARFCN-DL belongs to a band
func (band LTEBand) EarfcnDLInRange(earfcndl uint32) bool {
	return band.StartEarfcnDl <= earfcndl && earfcndl < band.StartEarfcnDl+band.CountEarfcn
}

// EarfcnULInRange checks that an EARFCN-UL belongs to a band
func (band LTEBand) EarfcnULInRange(earfcnul uint32) bool {
	if band.Mode == FDDMode {
		return band.StartEarfcnUl <= earfcnul && earfcnul < band.StartEarfcnUl+band.CountEarfcn
	}
	return band.EarfcnDLInRange(earfcnul)
}

// GetBand for a EARFCN-UL
func GetBand(earfcndl uint32) (*LTEBand, error) {
	for _, band := range bands {
		if band.EarfcnDLInRange(earfcndl) {
			return &band, nil
		}
	}
	return nil, fmt.Errorf("Invalid EARFCNDL: no matching band")
}
