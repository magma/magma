/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"context"
	"fmt"
	"strconv"

	"magma/feg/cloud/go/services/feg_relay/gw_to_feg_relay"
	"magma/orc8r/cloud/go/services/directoryd"
)

type teidType string

const (
	// types of TEIDs
	ControlPlaneTeid = "ControlPlaneTeid"
	UserPlaneTeid    = "UserPlaneTeid"
)

// getUniqueSgwTeid gets a unique Sgw C TEID making sure it is not being used checking on directoryd indexes
func GetUniqueSgwTeid(ctx context.Context, tType teidType) (uint32, error) {
	gw, err := gw_to_feg_relay.RetrieveGatewayIdentity(ctx)
	if err != nil {
		return 0, fmt.Errorf("Couldnt retrieve GatewayIdentity: %s", err)
	}
	var teid string
	switch tType {
	case ControlPlaneTeid:
		teid, err = directoryd.GetNewSgwCTeid(ctx, gw.NetworkId)
	case UserPlaneTeid:
		teid, err = directoryd.GetNewSgwUTeid(ctx, gw.NetworkId)
	default:
		err = fmt.Errorf("getUniqueSgwTeid: TeidType not found: %s", teid)
	}
	if err != nil {
		return 0, err
	}
	teidUint64, err := strconv.ParseUint(teid, 10, 32)
	return uint32(teidUint64), err
}
