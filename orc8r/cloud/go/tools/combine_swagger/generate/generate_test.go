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

package generate_test

import (
	"io/ioutil"
	"os"
	"testing"

	"magma/orc8r/cloud/go/tools/combine_swagger/generate"
	swaggergen "magma/orc8r/cloud/go/tools/swaggergen/generate"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateStandaloneSpec(t *testing.T) {
	goldenFilePath := "../testdata/standalone.yml.golden"
	targetFilePath := "../testdata/configs/importer2.yml"
	specTargetPath := "../testdata/test_result.yml"

	os.Remove(specTargetPath)
	defer os.Remove(specTargetPath)

	specs, err := swaggergen.ParseSwaggerDependencyTree(targetFilePath, "../testdata")
	assert.NoError(t, err)

	err = generate.GenerateSpec(targetFilePath, specs, specTargetPath)
	assert.NoError(t, err)

	actual, err := ioutil.ReadFile(specTargetPath)
	assert.NoError(t, err)

	expected, err := ioutil.ReadFile(goldenFilePath)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}
