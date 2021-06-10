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

package generate

import (
	"os"
	"path/filepath"
	"strings"

	"magma/orc8r/cloud/go/obsidian/swagger/spec"
	"magma/orc8r/cloud/go/tools/swaggergen/generate"

	"github.com/pkg/errors"
)

// GenerateStandaloneSpecs generates standalone specs for all Swagger specs in
// a directory
func GenerateStandaloneSpecs(specDir string, rootDir string) error {
	outDir := filepath.Join(filepath.Dir(specDir), "standalone")
	specPaths := getFilepaths(specDir)
	for _, path := range specPaths {
		specs, err := generate.ParseSwaggerDependencyTree(path, rootDir)
		if err != nil {
			return err
		}

		outPath := filepath.Join(outDir, filepath.Base(path))
		if err != nil {
			return err
		}

		GenerateSpec(path, specs, outPath)
	}
	return nil
}

// GenerateSpec generates a standalone spec file for a particular
// Swagger spec.
func GenerateSpec(targetFilepath string, specs map[string]generate.MagmaSwaggerSpec, outPath string) error {
	absTargetFilepath, err := filepath.Abs(targetFilepath)
	if err != nil {
		return errors.Wrapf(err, "target filepath %s is invalid", targetFilepath)
	}

	var yamlSpecs []string
	var yamlCommon string
	for specFilepath, s := range specs {
		if specFilepath != absTargetFilepath {
			// Clearing fields so that standalone spec is not overpopulated
			// by unnecessary fields from its dependencies
			s.Paths = nil
			s.Tags = nil
		}

		yamlSpec, err := s.MarshalToYAML()
		if err != nil {
			return err
		}

		if s.MagmaGenMeta.TempGenFilename == "orc8r-swagger-common.yml" {
			yamlCommon = yamlSpec
		} else {
			yamlSpecs = append(yamlSpecs, yamlSpec)
		}
	}

	combined, warnings, err := spec.Combine(yamlCommon, yamlSpecs)
	if err != nil {
		return err
	}
	if warnings != nil {
		return warnings
	}

	err = write(combined, outPath)
	if err != nil {
		return err
	}

	return nil
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

// write spec to filepath.
func write(outSpec string, filepath string) error {
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
