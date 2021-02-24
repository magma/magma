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
	"net"

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
		for _, key := range set[peer.IP] {
			if swag.Uint32Value(peer.Key) == key {
				return errors.New(fmt.Sprintf("Found duplicate peer %s with key %d", peer.IP, key))
			}
		}
		set[peer.IP] = append(set[peer.IP], swag.Uint32Value(peer.Key))
	}
	return nil
}

func (m *CwfSubscriberDirectoryRecord) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m *CarrierWifiHaPairStatus) ValidateModel() error {
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

func (m *CwfHaPair) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	if m.GatewayID1 == m.GatewayID2 {
		return fmt.Errorf("GatewayID1 and GatewayID2 cannot be the same")
	}
	return m.Config.ValidateModel()
}

func (m *MutableCwfHaPair) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	if m.GatewayID1 == m.GatewayID2 {
		return fmt.Errorf("GatewayID1 and GatewayID2 cannot be the same")
	}
	return m.Config.ValidateModel()
}

func (m *CwfHaPairConfigs) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	_, _, err := net.ParseCIDR(m.TransportVirtualIP)
	if err != nil {
		return fmt.Errorf("Transport virtual IP must be specified in CIDR format (e.g. '10.10.10.11/24')")
	}
	return nil
}
