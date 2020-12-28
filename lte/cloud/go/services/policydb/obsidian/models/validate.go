/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package models

import (
	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
)

func (m BaseNames) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m RuleNames) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *PolicyRule) ValidateModel() error {
	for _, flow := range m.FlowList {
		if flow.Match != nil {
			errMatch := flow.Match.ValidateModel()
			if errMatch != nil {
				return errMatch
			}
		}
	}
	return m.Validate(strfmt.Default)
}

func (m *FlowMatch) ValidateModel() error {
	if (m.IPV4Dst != "" || m.IPV4Src != "") && (m.IPSrc != nil || m.IPDst != nil) {
		return errors.New("Invalid Argument: Can't mix old ipv4_src/ipv4_dst type with the new ip_src/ip_dst")
	}
	return m.Validate(strfmt.Default)
}

func (m *RatingGroup) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *MutableRatingGroup) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *NetworkSubscriberConfig) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *PolicyQosProfile) ValidateModel() error {
	return m.Validate(strfmt.Default)
}
