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
	"github.com/pkg/errors"
)

var (
	commonSpecPath = "/etc/magma/configs/orc8r/swagger_specs/common/swagger-common.yml"
)

// GetCombinedSpec polls every servicer registered with
// a Swagger spec and merges them together to return a combined spec.
func GetCombinedSpec(yamlCommon string) (string, error) {
	servicers, err := GetSpecServicers()
	if err != nil {
		return "", err
	}

	var yamlSpecs []string
	for _, s := range servicers {
		yamlSpec, err := s.GetSpec()
		if err != nil {
			err = errors.Wrapf(err, "get Swagger spec from %s service", s.GetService())
			glog.Error(err)
		} else {
			yamlSpecs = append(yamlSpecs, yamlSpec)
		}
	}

	combined, warnings, err := swagger_lib.Combine(yamlCommon, yamlSpecs)
	if err != nil {
		return "", err
	}
	if warnings != nil {
		glog.Infof("Swagger spec traits were overwritten or unable to be read: %+v \n", warnings)
	}

	return combined, nil
}

// GetCommonSpec returns the YAML string of the Swagger Common Spec
func GetCommonSpec() (string, error) {
	data, err := ioutil.ReadFile(commonSpecPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
