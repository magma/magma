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
	"io/ioutil"

	swagger_lib "magma/orc8r/cloud/go/swagger"

	"github.com/golang/glog"
)

var (
	commonSpecPath = "/etc/magma/configs/orc8r/swagger_specs/common/swagger-common.yml"
)

// GetCombinedSwaggerSpecs polls every servicer registered with
// a Swagger spec and merges them together to return a combined spec.
func GetCombinedSwaggerSpecs() (string, error) {
	servicers, err := GetSpecServicers()
	if err != nil {
		return "", err
	}

	// Retrieve common spec
	data, err := ioutil.ReadFile(commonSpecPath)
	if err != nil {
		glog.Fatalf("Error retrieving common Swagger spec %+v", err)
	}
	yamlCommon := string(data)

	// Retrieve specs
	var yamlSpecs []string
	for _, s := range servicers {
		yamlSpec, err := s.GetSpec()
		if err != nil {
			glog.Errorf("Invalid response from spec servicer \n %+v", err)
		} else {
			yamlSpecs = append(yamlSpecs, yamlSpec)
		}
	}

	// Combine specs.
	combined, warnings, err := swagger_lib.Combine(yamlCommon, yamlSpecs)
	if err != nil {
		glog.Fatal(err)
	}
	if warnings != nil {
		glog.Infof("Warnings: %+v \n", warnings)
	}

	// Merge specs
	return combined, nil
}
