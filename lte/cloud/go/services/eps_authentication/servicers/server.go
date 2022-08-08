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
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fegprotos "magma/feg/cloud/go/protos"
	lteprotos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/eps_authentication/storage"
	"magma/orc8r/lib/go/protos"
)

type EPSAuthServer struct {
	store storage.SubscriberDBStorage
}

// NewEPSAuthServer returns a Server with the provided store.
func NewEPSAuthServer(store storage.SubscriberDBStorage) (*EPSAuthServer, error) {
	if store == nil {
		return nil, fmt.Errorf("Cannot initialize eps authentication server with nil store")
	}
	return &EPSAuthServer{store: store}, nil
}

// lookupSubscriber returns a subscriber's data or an error.
func (srv *EPSAuthServer) lookupSubscriber(user, nid string) (*lteprotos.SubscriberData, fegprotos.ErrorCode, error) {
	subscriber, err := srv.store.GetSubscriberData(&lteprotos.SubscriberID{Id: user}, &protos.NetworkID{Id: nid})
	if err != nil {
		if status.Convert(err).Code() == codes.NotFound {
			return nil, fegprotos.ErrorCode_USER_UNKNOWN, err
		}
		return nil, fegprotos.ErrorCode_AUTHENTICATION_DATA_UNAVAILABLE, err
	}
	return subscriber, fegprotos.ErrorCode_SUCCESS, nil
}

// lookupSubscriberProfile returns a subscriber's data & profile or an error.
func (srv *EPSAuthServer) lookupSubscriberProfile(
	userName, networkID string) (*lteprotos.SubscriberData, map[string]string, []string, fegprotos.ErrorCode, error) {

	subscriber, staticIps, subApns, err := srv.store.GetSubscriberDataProfile(
		&lteprotos.SubscriberID{Id: userName}, &protos.NetworkID{Id: networkID})
	if err != nil {
		if status.Convert(err).Code() == codes.NotFound {
			return nil, staticIps, subApns, fegprotos.ErrorCode_USER_UNKNOWN, err
		}
		return nil, staticIps, subApns, fegprotos.ErrorCode_AUTHENTICATION_DATA_UNAVAILABLE, err
	}
	return subscriber, staticIps, subApns, fegprotos.ErrorCode_SUCCESS, nil
}
