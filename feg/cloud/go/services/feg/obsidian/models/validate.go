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
	"fmt"
	"regexp"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"magma/orc8r/cloud/go/services/configurator"
)

func (m FegNetworkID) ValidateModel(ctx context.Context) error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	if !swag.IsZero(m) {
		exists, err := configurator.DoesNetworkExist(ctx, string(m))
		if err != nil {
			return fmt.Errorf("Failed to search for network %s: %w", string(m), err)
		}
		if !exists {
			return fmt.Errorf("Network: %s does not exist", string(m))
		}
	}
	return nil
}

func (m *FederatedNetworkConfigs) ValidateModel(ctx context.Context) error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	nid := *m.FegNetworkID
	if !swag.IsZero(nid) {
		exists, err := configurator.DoesNetworkExist(ctx, nid)
		if err != nil {
			return fmt.Errorf("Failed to search for network %s: %w", nid, err)
		}
		if !exists {
			return fmt.Errorf("Network: %s does not exist", nid)
		}
	}
	return nil
}

func (m *DiameterClientConfigs) ValidateModel(context.Context) error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m *DiameterServerConfigs) ValidateModel(context.Context) error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m *EapAkaTimeouts) ValidateModel(context.Context) error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

func (m *GatewayFederationConfigs) ValidateModel(context.Context) error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}

var nhRouteRegex = regexp.MustCompile(`^(\d{5,6})$`)

func (m *NetworkFederationConfigs) ValidateModel(context.Context) error {
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

func (m *SubscriptionProfile) ValidateModel(context.Context) error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}
	return nil
}
