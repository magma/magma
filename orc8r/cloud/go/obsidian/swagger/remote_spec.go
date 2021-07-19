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

package swagger

import (
	"context"
	"strings"
	"time"

	"magma/orc8r/cloud/go/obsidian/swagger/protos"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
)

// RemoteSpec identifies a remote spec
type RemoteSpec struct {
	// service name of the RemoteSpec
	// should always be lowercase to match service registry convention
	service string
}

// NewRemoteSpec constructs a endpoint to communicate with the spec servicer.
func NewRemoteSpec(serviceName string) RemoteSpec {
	return RemoteSpec{service: strings.ToLower(serviceName)}
}

// GetPartialSpec returns the partial spec associated to the service as a
// YAML string.
func (s *RemoteSpec) GetPartialSpec() (string, error) {
	c, err := s.getClient()
	if err != nil {
		return "", err
	}

	res, err := c.GetPartialSpec(context.Background(), &protos.PartialSpecRequest{})
	if err != nil {
		return "", err
	}

	return res.SwaggerSpec, nil
}

// GetStandaloneSpec returns the standalone spec associated to the service as
// a YAML string.
func (s *RemoteSpec) GetStandaloneSpec() (string, error) {
	c, err := s.getClient()
	if err != nil {
		return "", err
	}

	res, err := c.GetStandaloneSpec(context.Background(), &protos.StandaloneSpecRequest{})
	if err != nil {
		return "", err
	}

	return res.SwaggerSpec, nil
}

// GetService returns the service name.
func (s *RemoteSpec) GetService() string {
	return s.service
}

// Shorten timeout so that the caller doesn't time out (60s) if individual
// services time out.
const specClientDefaultTimeout = 5 * time.Second

func (s *RemoteSpec) getClient() (protos.SwaggerSpecClient, error) {
	conn, err := registry.GetConnectionWithTimeout(s.service, specClientDefaultTimeout)
	if err != nil {
		initErr := merrors.NewInitError(err, s.service)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewSwaggerSpecClient(conn), nil
}
