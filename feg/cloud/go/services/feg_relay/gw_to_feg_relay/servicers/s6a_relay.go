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

package servicers

import (
	"context"

	"magma/feg/cloud/go/protos"
)

// S6AProxyServer implementation
//
// AuthenticationInformation sends AIR over diameter connection,
// waits (blocks) for AIA & returns its RPC representation
func (s *RelayRouter) AuthenticationInformation(
	c context.Context, r *protos.AuthenticationInformationRequest) (*protos.AuthenticationInformationAnswer, error) {

	client, ctx, cancel, err := s.getS6aClient(c, r.GetUserName())
	if err != nil {
		return nil, err
	}
	defer cancel()
	ret, err := client.AuthenticationInformation(ctx, r)
	return ret, err
}

// UpdateLocation sends ULR (Code 316) over diameter connection,
// waits (blocks) for ULA & returns its RPC representation
func (s *RelayRouter) UpdateLocation(
	ctx context.Context, r *protos.UpdateLocationRequest) (*protos.UpdateLocationAnswer, error) {

	client, ctx, cancel, err := s.getS6aClient(ctx, r.GetUserName())
	if err != nil {
		return nil, err
	}
	defer cancel()
	return client.UpdateLocation(ctx, r)
}

// PurgeUE sends PUR (Code 321) over diameter connection,
// waits (blocks) for PUA & returns its RPC representation
func (s *RelayRouter) PurgeUE(ctx context.Context, r *protos.PurgeUERequest) (*protos.PurgeUEAnswer, error) {

	client, ctx, cancel, err := s.getS6aClient(ctx, r.GetUserName())
	if err != nil {
		return nil, err
	}
	defer cancel()
	return client.PurgeUE(ctx, r)
}

func (s *RelayRouter) getS6aClient(
	c context.Context, imsi string) (protos.S6AProxyClient, context.Context, context.CancelFunc, error) {

	conn, ctx, cancel, err := s.GetFegServiceConnection(c, imsi, FegS6aProxy)
	if err != nil {
		return nil, nil, nil, err
	}
	return protos.NewS6AProxyClient(conn), ctx, cancel, nil
}
