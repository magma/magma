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
)

func (m *SymphonyNetwork) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *SymphonyAgent) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *MutableSymphonyAgent) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *ManagedDevices) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *SymphonyDevice) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *MutableSymphonyDevice) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *SymphonyDeviceName) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *SymphonyDeviceConfig) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *SymphonyDeviceAgent) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *SymphonyDeviceState) ValidateModel() error {
	return m.Validate(strfmt.Default)
}
