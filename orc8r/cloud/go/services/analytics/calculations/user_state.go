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

package calculations

import (
	"magma/orc8r/cloud/go/services/analytics/query_api"
	"magma/orc8r/lib/go/metrics"

	"github.com/golang/glog"
)

// UserStateManager contains the user state in the deployment to enforce minUserThreshold constraint
type UserStateManager interface {
	Update()
	GetTotalUsers() int
	GetTotalUsersInNetwork(networkID string) int
	GetTotalUsersInGateway(networkID string, gatewayID string) int
}

// userNetworkState contains the user state per network
type userNetworkState struct {
	totalUsersPerNetwork int
	usersGatewayTable    map[string]int
}

// UserStateManager contains the user state in the deployment to enforce minUserThreshold constraint
type userStateManagerImpl struct {
	promAPIClient query_api.PrometheusAPI
	expr          string

	// total number of  users across the deployment
	totalUsers        int
	usersNetworkTable map[string]*userNetworkState
}

// NewUserStateManager construct a new user state manager
func NewUserStateManager(promAPIClient query_api.PrometheusAPI, expr string) UserStateManager {
	return &userStateManagerImpl{
		promAPIClient: promAPIClient,
		expr:          expr,
	}
}

// GetTotalUsers get the total number of users in a deployment
func (u *userStateManagerImpl) GetTotalUsers() int {
	return u.totalUsers
}

// GetTotalUsersInNetwork get the total number of users in the network
func (u *userStateManagerImpl) GetTotalUsersInNetwork(networkID string) int {
	if networkState, ok := u.usersNetworkTable[networkID]; ok {
		return networkState.totalUsersPerNetwork
	}
	return 0
}

// GetTotalUsersInGateway get the total number of users in the network
func (u *userStateManagerImpl) GetTotalUsersInGateway(networkID string, gatewayID string) int {
	networkState, ok := u.usersNetworkTable[networkID]
	if !ok {
		return 0
	}

	if gatewayState, ok := networkState.usersGatewayTable[gatewayID]; ok {
		return gatewayState
	}
	return 0
}

// Update the existing user state
func (u *userStateManagerImpl) Update() {
	if u.expr == "" {
		// Nothing to do
		return
	}
	vec, err := query_api.QueryPrometheusVector(u.promAPIClient, u.expr)
	if err != nil {
		glog.Errorf("failed querying user state metric %v", err)
		return
	}
	u.totalUsers = 0
	u.usersNetworkTable = make(map[string]*userNetworkState)
	for _, v := range vec {
		// Get labels from query result
		var networkID string
		var gatewayID string
		for label, value := range v.Metric {
			if string(label) == metrics.NetworkLabelName {
				networkID = string(value)
			} else if string(label) == metrics.GatewayLabelName {
				gatewayID = string(value)
			}
		}
		userNetworkTable, ok := u.usersNetworkTable[networkID]
		if !ok {
			userNetworkTable = &userNetworkState{
				totalUsersPerNetwork: 0,
				usersGatewayTable:    make(map[string]int),
			}
			u.usersNetworkTable[networkID] = userNetworkTable
		}
		userCount := int(v.Value)
		userNetworkTable.totalUsersPerNetwork += userCount
		userNetworkTable.usersGatewayTable[gatewayID] = userCount
		u.totalUsers += userCount
	}
}
