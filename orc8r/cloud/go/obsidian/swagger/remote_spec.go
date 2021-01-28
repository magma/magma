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

	"magma/orc8r/cloud/go/obsidian/swagger/protos"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
)

type remoteSpec struct {
	// service name of the remoteSpec
	// should always be lowercase to match service registry convention
	service string
}

func NewRemoteSpec(serviceName string) remoteSpec {
	return remoteSpec{service: strings.ToLower(serviceName)}
}

func (s *remoteSpec) GetSpec() (string, error) {
	c, err := s.getClient()
	if err != nil {
		return "", err
	}

	res, err := c.GetSpec(context.Background(), &protos.GetSpecRequest{})
	if err != nil {
		return "", err
	}

	return res.SwaggerSpec, nil
}

func (s *remoteSpec) getClient() (protos.SwaggerSpecClient, error) {
	conn, err := registry.GetConnection(s.service)
	if err != nil {
		initErr := merrors.NewInitError(err, s.service)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewSwaggerSpecClient(conn), nil
}
