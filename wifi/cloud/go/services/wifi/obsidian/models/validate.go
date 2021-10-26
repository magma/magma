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
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
)

func (m *WifiNetwork) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *NetworkWifiConfigs) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *WifiGateway) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *MutableWifiGateway) ValidateModel(context.Context) error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}

	// Custom validation only for wifi and device
	var res []error
	if err := m.Wifi.ValidateModel(context.Background()); err != nil {
		res = append(res, err)
	}
	if err := m.Device.ValidateModel(context.Background()); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *GatewayWifiConfigs) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *WifiMesh) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *MeshName) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *MeshWifiConfigs) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}
