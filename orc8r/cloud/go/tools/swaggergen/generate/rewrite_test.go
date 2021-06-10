/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package generate_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"magma/orc8r/cloud/go/tools/swaggergen/generate"

	"github.com/stretchr/testify/assert"
)

func TestRewriteGeneratedRefs(t *testing.T) {
	runRewriteTestCase(t, "../testdata/importer2.yml", "../testdata/importer2")
	runRewriteTestCase(t, "../testdata/importer.yml", "../testdata/importer")
	runRewriteTestCase(t, "../testdata/base.yml", "../testdata/base")
}

func runRewriteTestCase(t *testing.T, ymlFile string, outputDir string) {
	defer cleanupActualFiles(outputDir)

	specs, err := generate.ParseSwaggerDependencyTree(ymlFile, "../testdata")
	assert.NoError(t, err)

	err = generate.GenerateModels(ymlFile, "../testdata/config.yml", "../testdata", specs)
	assert.NoError(t, err)

	err = generate.RewriteGeneratedRefs(ymlFile, "../testdata", specs)
	assert.NoError(t, err)

	goldenFiles, actualFiles := []string{}, []string{}
	err = filepath.Walk(outputDir, func(path string, _ os.FileInfo, _ error) error {
		if strings.HasSuffix(path, "actual.golden") {
			goldenFiles = append(goldenFiles, strings.TrimSuffix(path, ".golden"))
		} else if strings.HasSuffix(path, ".actual") {
			actualFiles = append(actualFiles, path)
		}
		return nil
	})
	assert.NoError(t, err)
	sort.Strings(goldenFiles)
	sort.Strings(actualFiles)
	assert.Equal(t, goldenFiles, actualFiles)

	// Verify contents of actual vs golden
	for _, baseFilename := range goldenFiles {
		goldenFileContents, err := ioutil.ReadFile(baseFilename + ".golden")
		assert.NoError(t, err)
		actualFileContents, err := ioutil.ReadFile(baseFilename)
		assert.NoError(t, err)
		assert.Equal(t, goldenFileContents, actualFileContents)
	}
}
