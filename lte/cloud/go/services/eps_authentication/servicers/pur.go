/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package servicers

import (
	"errors"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/feg/cloud/go/protos"
	"magma/lte/cloud/go/services/eps_authentication/metrics"
	"magma/orc8r/cloud/go/identity"
)

func (srv *EPSAuthServer) PurgeUE(ctx context.Context, purge *protos.PurgeUERequest) (*protos.PurgeUEAnswer, error) {
	glog.V(2).Infof("received PUR from: %s", purge.GetUserName())
	metrics.PURequests.Inc()
	if err := validatePUR(purge); err != nil {
		glog.V(2).Infof("PUR is invalid: %v", err.Error())
		metrics.InvalidRequests.Inc()
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	networkID, err := identity.GetClientNetworkID(ctx)
	if err != nil {
		glog.V(2).Infof("could not lookup networkID: %v", err.Error())
		metrics.NetworkIDErrors.Inc()
		return nil, err
	}
	_, errorCode, err := srv.lookupSubscriber(purge.UserName, networkID)
	if err != nil {
		glog.V(2).Infof("failed to lookup subscriber '%s': %v", purge.UserName, err.Error())
		metrics.UnknownSubscribers.Inc()
		return &protos.PurgeUEAnswer{ErrorCode: errorCode}, err
	}
	return &protos.PurgeUEAnswer{ErrorCode: protos.ErrorCode_SUCCESS}, nil
}

// validatePUR returns an error iff the PUR is invalid.
func validatePUR(purge *protos.PurgeUERequest) error {
	if purge == nil {
		return errors.New("received a nil PurgeUERequest")
	}
	if len(purge.UserName) == 0 {
		return errors.New("user name was empty")
	}
	return nil
}
