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
	"io/ioutil"

	"magma/orc8r/cloud/go/obsidian/swagger/protos"

	"github.com/golang/glog"
)

type specServicer struct {
	spec string
}

func NewSpecServicerFromFile(path string, service string) protos.SwaggerSpecServer {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		// We are swallowing this error because a singular service spec failure
		// should not down the entire Swagger UI
		glog.Errorf("Error retrieving Swagger Spec of service %s: %+v", service, err)
		return NewSpecServicer("")
	}
	return NewSpecServicer(string(data))
}

func NewSpecServicer(spec string) protos.SwaggerSpecServer {
	return &specServicer{spec}
}

func (s *specServicer) GetSpec(ctx context.Context, request *protos.GetSpecRequest) (*protos.GetSpecResponse, error) {
	return &protos.GetSpecResponse{SwaggerSpec: s.spec}, nil
}
