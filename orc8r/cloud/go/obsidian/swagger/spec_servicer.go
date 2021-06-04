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

	"magma/orc8r/cloud/go/obsidian/swagger/protos"
	"magma/orc8r/cloud/go/obsidian/swagger/spec"

	"github.com/golang/glog"
)

type specServicer struct {
	partialSpec    string
	standaloneSpec string
}

// NewSpecServicer constructs a spec servicer.
func NewSpecServicer(partialSpec string, standaloneSpec string) protos.SwaggerSpecServer {
	return &specServicer{partialSpec: partialSpec, standaloneSpec: standaloneSpec}
}

func NewSpecServicerWithLoader(specs spec.Loader, service string) protos.SwaggerSpecServer {
	// Swallow errors because the service should continue to run even if it
	// can't find its Swagger spec Loader.
	partial, err := specs.GetPartialSpec(service)
	if err != nil {
		glog.Errorf("Error retrieving Swagger partial spec of service %s: %+v", service, err)
	}
	standalone, err := specs.GetStandaloneSpec(service)
	if err != nil {
		glog.Errorf("Error retrieving Swagger standalone spec of service %s: %+v", service, err)
	}
	return NewSpecServicer(partial, standalone)
}

// NewSpecServicerFromFile initializes a specServicer given a service name.
func NewSpecServicerFromFile(service string) protos.SwaggerSpecServer {
	return NewSpecServicerWithLoader(spec.GetDefaultLoader(), service)
}

func (s *specServicer) GetPartialSpec(ctx context.Context, request *protos.PartialSpecRequest) (*protos.PartialSpecResponse, error) {
	return &protos.PartialSpecResponse{SwaggerSpec: s.partialSpec}, nil
}

func (s *specServicer) GetStandaloneSpec(ctx context.Context, request *protos.StandaloneSpecRequest) (*protos.StandaloneSpecResponse, error) {
	return &protos.StandaloneSpecResponse{SwaggerSpec: s.standaloneSpec}, nil
}
