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
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/go-openapi/strfmt"
)

const echoKeyType = "ECHO"
const ecdsaKeyType = "SOFTWARE_ECDSA_SHA256"

func (m *Network) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *NetworkDNSConfig) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *NetworkFeatures) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m NetworkDNSRecords) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m DNSConfigRecord) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *MagmadGateway) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *GatewayDevice) ValidateModel() error {
	if err := m.Key.ValidateModel(); err != nil {
		return err
	}
	return m.Validate(strfmt.Default)
}

func (m *ChallengeKey) ValidateModel() error {
	switch m.KeyType {
	case echoKeyType:
		if m.Key != nil {
			return errors.New("ECHO mode should not have key value")
		}
		return nil
	case ecdsaKeyType:
		if m.Key == nil {
			return fmt.Errorf("No key supplied")
		}
		_, err := x509.ParsePKIXPublicKey(*m.Key)
		if err != nil {
			return fmt.Errorf("Failed to parse key: %s", err)
		}
		return nil
	default:
		return fmt.Errorf("Unknown key type %s", m.KeyType)
	}
}

func (m *MagmadGatewayConfigs) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m TierID) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *ReleaseChannel) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *Tier) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *TierName) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *TierVersion) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *TierGateways) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *TierImages) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *TierImage) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *GatewayStatus) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *GatewayVpnConfigs) ValidateModel() error {
	return m.Validate(strfmt.Default)
}
