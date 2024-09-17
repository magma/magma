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
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/go-openapi/strfmt"
)

const echoKeyType = "ECHO"
const ecdsaKeyType = "SOFTWARE_ECDSA_SHA256"

func (m *Network) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *NetworkDNSConfig) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *NetworkFeatures) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *NetworkSentryConfig) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *StateConfig) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m NetworkDNSRecords) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m DNSConfigRecord) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *MagmadGateway) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *GatewayDevice) ValidateModel(context.Context) error {
	if err := m.Key.ValidateModel(context.Background()); err != nil {
		return err
	}
	return m.Validate(strfmt.Default)
}

func (m *ChallengeKey) ValidateModel(context.Context) error {
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

func (m *MagmadGatewayConfigs) ValidateModel(context.Context) error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m TierID) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *ReleaseChannel) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *Tier) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *TierName) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *TierVersion) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *TierGateways) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *TierImages) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *TierImage) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *GatewayStatus) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}

func (m *GatewayVpnConfigs) ValidateModel(context.Context) error {
	return m.Validate(strfmt.Default)
}
