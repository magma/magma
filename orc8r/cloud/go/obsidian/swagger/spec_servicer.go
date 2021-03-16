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
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"magma/orc8r/cloud/go/obsidian/swagger/protos"

	"github.com/golang/glog"
)

type specServicer struct {
	partialSpec    string
	standaloneSpec string
}

// NewSpecServicerFromFile intializes a spec servicer
// given a service name.
func NewSpecServicerFromFile(service string) protos.SwaggerSpecServer {
	service = strings.ToLower(service)
	partialPath, standalonePath := getSpecPaths(service)
	partial, err := ioutil.ReadFile(partialPath)
	if err != nil {
		// Swallow error because the service should continue to
		// run even if it can't find its partial Swagger spec file.
		glog.Errorf("Error retrieving Swagger Spec of service %s: %+v", service, err)
		return NewSpecServicer("", "")
	}
	standalone, err := ioutil.ReadFile(standalonePath)
	if err != nil {
		// Swallowing ReadFile error because the service should continue to
		// run even if it can't find its standalone Swagger spec file.
		glog.Errorf("Error retrieving Swagger Spec of service %s: %+v", service, err)
		return NewSpecServicer("", "")
	}

	return NewSpecServicer(string(partial), string(standalone))
}

// NewSpecServicer constructs a spec servicer.
func NewSpecServicer(partialSpec string, standaloneSpec string) protos.SwaggerSpecServer {
	return &specServicer{partialSpec: partialSpec, standaloneSpec: standaloneSpec}
}

func (s *specServicer) GetPartialSpec(ctx context.Context, request *protos.PartialSpecRequest) (*protos.PartialSpecResponse, error) {
	return &protos.PartialSpecResponse{SwaggerSpec: s.partialSpec}, nil
}

func (s *specServicer) GetStandaloneSpec(ctx context.Context, request *protos.StandaloneSpecRequest) (*protos.StandaloneSpecResponse, error) {
	return &protos.StandaloneSpecResponse{SwaggerSpec: s.standaloneSpec}, nil
}

// getSpecPaths returns the filepath on the production image
// that contains the service's Swagger spec
func getSpecPaths(service string) (string, string) {
	specDir := "/etc/magma/swagger/specs"
	partialSpecPath := filepath.Join(specDir, "partial", fmt.Sprintf("%s.swagger.v1.yml", service))
	standaloneSpecPath := filepath.Join(specDir, "standalone", fmt.Sprintf("%s.swagger.v1.yml", service))
	return partialSpecPath, standaloneSpecPath
}
