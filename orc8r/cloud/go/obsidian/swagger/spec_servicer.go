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

	protos "magma/orc8r/cloud/go/obsidian/swagger/protos"
)

type specServicer struct {
	spec string
}

func NewSpecServicerWithPath(specPath string) (protos.SwaggerSpecServer, error) {
	data, err := ioutil.ReadFile(specPath)
	if err != nil {
		return nil, err
	}

	return NewSpecServicer(string(data)), nil
}

func NewSpecServicer(spec string) protos.SwaggerSpecServer {
	return &specServicer{spec}
}

func (s *specServicer) GetSpec(ctx context.Context, request *protos.GetSpecRequest) (*protos.GetSpecResponse, error) {
	ret := &protos.GetSpecResponse{SwaggerSpec: ""}
	ret.SwaggerSpec = s.spec

	return ret, nil
}
