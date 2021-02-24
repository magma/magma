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

// Package wrappers provides semantic wrappers around the state service's
// client API.
package wrappers

import (
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/lib/go/errors"
)

// GetGatewayStatus returns the status for an indicated gateway.
func GetGatewayStatus(networkID string, deviceID string) (*models.GatewayStatus, error) {
	st, err := state.GetState(networkID, orc8r.GatewayStateType, deviceID, serdes.State)
	if err != nil {
		return nil, err
	}
	if st.ReportedState == nil {
		return nil, errors.ErrNotFound
	}
	return fillInGatewayStatusState(st), nil
}

// GetGatewayStatuses returns the status for indicated gateways, keyed by
// device ID.
func GetGatewayStatuses(networkID string, deviceIDs []string) (map[string]*models.GatewayStatus, error) {
	stateIDs := types.MakeIDs(orc8r.GatewayStateType, deviceIDs...)
	res, err := state.GetStates(networkID, stateIDs, serdes.State)
	if err != nil {
		return map[string]*models.GatewayStatus{}, err
	}

	ret := make(map[string]*models.GatewayStatus, len(res))
	for stateID, st := range res {
		ret[stateID.DeviceID] = fillInGatewayStatusState(st)
	}
	return ret, nil
}

func fillInGatewayStatusState(st types.State) *models.GatewayStatus {
	if st.ReportedState == nil {
		return nil
	}
	gwStatus := st.ReportedState.(*models.GatewayStatus)
	gwStatus.CheckinTime = st.TimeMs
	gwStatus.CertExpirationTime = st.CertExpirationTime
	gwStatus.HardwareID = st.ReporterID
	return gwStatus
}
