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
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/pkg/errors"
)

func (m *CwfNetwork) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m *NetworkCarrierWifiConfigs) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m *GatewayCwfConfigs) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	set := make(map[string][]uint32)
	for _, peer := range m.AllowedGrePeers {
		for _, key := range set[string(peer.IP)] {
			if swag.Uint32Value(peer.Key) == key {
				return errors.New(fmt.Sprintf("Found duplicate peer %s with key %d", string(peer.IP), key))
			}
		}
		set[string(peer.IP)] = append(set[string(peer.IP)], swag.Uint32Value(peer.Key))
	}
	return nil
}

func (m *CwfSubscriberDirectoryRecord) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m *CarrierWifiNetworkClusterStatus) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m *CarrierWifiGatewayHealthStatus) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}
