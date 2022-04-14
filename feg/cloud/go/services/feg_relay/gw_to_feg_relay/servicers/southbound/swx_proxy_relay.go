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

// SwxProxyServer implementation
//
// Authenticate retrieves authentication vectors from the HSS using MAR/MAA
func (s *RelayRouter) Authenticate(
	c context.Context, r *protos.AuthenticationRequest) (*protos.AuthenticationAnswer, error) {

	client, ctx, cancel, err := s.getSwxClient(c, r.GetUserName())
	if err != nil {
		return nil, err
	}
	defer cancel()
	return client.Authenticate(ctx, r)
}

// Register the AAA server serving a user to the HSS using SAR/SAA
func (s *RelayRouter) Register(
	ctx context.Context, r *protos.RegistrationRequest) (*protos.RegistrationAnswer, error) {

	client, ctx, cancel, err := s.getSwxClient(ctx, r.GetUserName())
	if err != nil {
		return nil, err
	}
	defer cancel()
	return client.Register(ctx, r)
}

// Deregister the AAA server serving a user to the HSS using SAR/SAA
func (s *RelayRouter) Deregister(
	ctx context.Context, r *protos.RegistrationRequest) (*protos.RegistrationAnswer, error) {

	client, ctx, cancel, err := s.getSwxClient(ctx, r.GetUserName())
	if err != nil {
		return nil, err
	}
	defer cancel()
	return client.Deregister(ctx, r)
}

func (s *RelayRouter) getSwxClient(
	c context.Context, imsi string) (protos.SwxProxyClient, context.Context, context.CancelFunc, error) {

	conn, ctx, cancel, err := s.GetFegServiceConnection(c, imsi, FegSwxProxy)
	if err != nil {
		return nil, nil, nil, err
	}
	return protos.NewSwxProxyClient(conn), ctx, cancel, nil
}
