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

import "gopkg.in/yaml.v2"

// Spec is the Go struct version of a OAI/Swagger 2.0 YAML spec file.
type Spec struct {
	Swagger string
	Info    struct {
		Title       string
		Description string
		Version     string
	}
	BasePath    string `yaml:"basePath"`
	Consumes    []string
	Produces    []string
	Schemes     []string
	Tags        []TagDefinition
	Paths       map[string]interface{}
	Responses   map[string]interface{}
	Parameters  map[string]interface{}
	Definitions map[string]interface{}
}

type TagDefinition struct {
	Description string
	Name        string
}

// MarshalBinary marshals the spec to bytes.
func (s Spec) MarshalBinary() ([]byte, error) {
	yamlSpec, err := marshalToYAML(s)
	if err != nil {
		return nil, nil
	}
	return []byte(yamlSpec), nil
}

// marshalToYAML marshals the passed Swagger spec to a YAML-formatted string.
func marshalToYAML(spec Spec) (string, error) {
	d, err := yaml.Marshal(&spec)
	if err != nil {
		return "", err
	}
	return string(d), nil
}
