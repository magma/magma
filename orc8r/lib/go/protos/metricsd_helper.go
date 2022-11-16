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

// metricsd_helper.go adds some useful conversions between metric/label enum
// values and names

package protos

import (
	dto "github.com/prometheus/client_model/go"
)

// GetDecodedLabel converts a Metric to a list of LabelPair
func GetDecodedLabel(m *dto.Metric) []*dto.LabelPair {
	var newLabels []*dto.LabelPair
	for _, labelPair := range m.GetLabel() {
		labelName := labelPair.GetName()
		labelValue := labelPair.GetValue()
		newLabels = append(
			newLabels,
			&dto.LabelPair{Name: &labelName, Value: &labelValue})
	}
	return newLabels
}
