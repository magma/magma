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

package combine

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Load a common swagger spec file and a directory of
// swagger specs to their YAML string format.
func Load(commonFilepath, inDir string) (string, []string, error) {
	yamlSpecs, err := loadSpecsFromInputDir(inDir)
	if err != nil {
		return "", nil, err
	}

	yamlCommon, err := readFile(commonFilepath)
	if err != nil {
		return "", nil, err
	}

	return yamlCommon, yamlSpecs, nil
}

// Write spec to filepath.
func Write(outSpec string, filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = f.WriteString(outSpec)
	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return err
	}

	return nil
}

// loadSpecsFromInputDir loads all input Swagger files' contents
// to string.
func loadSpecsFromInputDir(inDir string) ([]string, error) {
	filepaths := getFilepaths(inDir)
	contents, err := readFiles(filepaths)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

// readFiles maps the passed filepaths to their contents.
func readFiles(filepaths []string) ([]string, error) {
	var contents []string
	for _, path := range filepaths {
		s, err := readFile(path)
		if err != nil {
			return nil, err
		}
		contents = append(contents, s)
	}
	return contents, nil
}

// readFile returns the content of the passed filepath.
func readFile(filepath string) (string, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// getFilepaths returns the filepaths of each Swagger YAML file in or
// below inDir, in lexical order.
func getFilepaths(inDir string) []string {
	var filepaths []string
	filepath.Walk(inDir, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".yml") {
			filepaths = append(filepaths, path)
		}
		return nil
	})
	return filepaths
}
