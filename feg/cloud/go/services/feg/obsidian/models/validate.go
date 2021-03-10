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
	"regexp"

	"magma/orc8r/cloud/go/services/configurator"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/pkg/errors"
)

func (m FegNetworkID) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	if !swag.IsZero(m) {
		exists, err := configurator.DoesNetworkExist(string(m))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to search for network %s", string(m)))
		}
		if !exists {
			return errors.New(fmt.Sprintf("Network: %s does not exist", string(m)))
		}
	}
	return nil
}

func (m *FederatedNetworkConfigs) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	nid := *m.FegNetworkID
	if !swag.IsZero(nid) {
		exists, err := configurator.DoesNetworkExist(nid)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to search for network %s", nid))
		}
		if !exists {
			return fmt.Errorf("Network: %s does not exist", nid)
		}
	}
	return nil
}

func (m *DiameterClientConfigs) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m *DiameterServerConfigs) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m *EapAkaTimeouts) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m *GatewayFederationConfigs) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

var nhRouteRegex = regexp.MustCompile(`^(\d{5,6})$`)

func (m *NetworkFederationConfigs) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	if m != nil {
		for k, v := range m.NhRoutes {
			if !nhRouteRegex.Match([]byte(k)) {
				return fmt.Errorf("invalid NH route PLMNID: %s for serving FeG network: %s", k, v)
			}
		}
	}
	return nil
}

func (m *SubscriptionProfile) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}
