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

package generate

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"magma/orc8r/cloud/go/swagger"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// MagmaSwaggerSpec is the Go struct version of our custom Swagger spec file.
// The only difference is the magma-gen-meta field, which specifies
// dependencies (other swagger specs that the spec has refs to), a desired
// filename that this file should be ref'd with from dependent files, and
// a list of Go struct types and filenames that this file produces when
// models are generated.
type MagmaSwaggerSpec struct {
	Swagger      string
	MagmaGenMeta MagmaGenMeta `yaml:"magma-gen-meta"`
	Info         struct {
		Title       string
		Description string
		Version     string
	}
	BasePath    string `yaml:"basePath"`
	Consumes    []string
	Produces    []string
	Schemes     []string
	Tags        []swagger.TagDefinition
	Paths       map[string]interface{}
	Responses   map[string]interface{}
	Parameters  map[string]interface{}
	Definitions map[string]interface{}
}

func (m MagmaSwaggerSpec) ToSwaggerSpec() swagger.Spec {
	s := swagger.Spec{
		Swagger:     m.Swagger,
		Info:        m.Info,
		BasePath:    m.BasePath,
		Consumes:    m.Consumes,
		Produces:    m.Produces,
		Schemes:     m.Schemes,
		Tags:        m.Tags,
		Paths:       m.Paths,
		Responses:   m.Responses,
		Parameters:  m.Parameters,
		Definitions: m.Definitions,
	}
	return s
}

type MagmaGenMeta struct {
	// GoPackage is the target Go package that the models defined in this spec
	// will be generated into (including the trailing `models`)
	// This will be referenced by the tool when adding imports to generated
	// files during the modification step of generation.
	GoPackage string `yaml:"go-package"`

	// OutputDir specifies the desired output directory relative to MAGMA_ROOT
	// that the models in this package should be generated into.
	OutputDir string `yaml:"output-dir"`

	// Dependencies is a list of swagger spec files that this file references.
	// This should be a list of filepaths relative to MAGMA_ROOT.
	Dependencies []string

	// TempGenFilename will be the filename that the swaggergen tool renames
	// this spec file to when copying it into directories of dependent
	// specs during generation. Dependent specs should reference definitions
	// in this spec by referencing from this filename.
	TempGenFilename string `yaml:"temp-gen-filename"`

	// Types is a list of Go struct names and generated filenames that this
	// spec file produces. This will be referenced by the tool when replacing
	// references.
	Types []MagmaGenType
}

type MagmaGenType struct {
	GoStructName string `yaml:"go-struct-name"`
	Filename     string
}

// GenerateModels parses the magma-gen-meta key of the given swagger YAML file,
// copies the files that the target file depends on into the current working
// directory, shells out to `swagger generate models`, then cleans up the
// dependency files.
func GenerateModels(targetFilepath string, configFilepath string, rootDir string, specs map[string]MagmaSwaggerSpec) error {
	absTargetFilepath, err := filepath.Abs(targetFilepath)
	if err != nil {
		return errors.Wrapf(err, "target filepath %s is invalid", targetFilepath)
	}

	tmpGenDir, err := ioutil.TempDir(".", "tmpgen")
	if err != nil {
		return errors.Wrap(err, "could not create temporary gen directory")
	}
	defer os.RemoveAll(tmpGenDir)

	// For each dependency, strip the magma-gen-meta and write the result to
	// the filename specified by `dependent-filename`
	err = StripAndWriteSwaggerSpecs(specs, tmpGenDir)
	if err != nil {
		return errors.WithStack(err)
	}

	// Shell out to go-swagger
	targetSpec := specs[absTargetFilepath]
	outputDir := filepath.Join(rootDir, targetSpec.MagmaGenMeta.OutputDir)
	absConfigFilepath, err := filepath.Abs(configFilepath)
	if err != nil {
		return err
	}

	cmd := exec.Command(
		"swagger", "generate", "model",
		"--spec", filepath.Join(tmpGenDir, targetSpec.MagmaGenMeta.TempGenFilename),
		"--target", outputDir,
		"--config-file", absConfigFilepath,
	)
	stdoutBuf := &strings.Builder{}
	stderrBuf := &strings.Builder{}
	cmd.Stdout = stdoutBuf
	cmd.Stderr = stderrBuf

	err = cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "failed to generate models; stdout:\n%s\nstderr:\n%s", stdoutBuf.String(), stderrBuf.String())
	}

	return nil
}

// ParseSwaggerDependencyTree parses the entire dependency tree of a magma
// swagger spec file specified by the rootFilepath parameter.
// The returned value maps between the absolute specified dependency filepath
// and the parsed struct for the dependency file.
func ParseSwaggerDependencyTree(rootFilepath string, rootDir string) (map[string]MagmaSwaggerSpec, error) {
	absRootFilepath, err := filepath.Abs(rootFilepath)
	if err != nil {
		return nil, errors.Wrapf(err, "root filepath %s is invalid", rootFilepath)
	}

	targetSpec, err := readSwaggerSpec(absRootFilepath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	type mscAndPath struct {
		MagmaSwaggerSpec
		path string
	}

	// Do a BFS to parse the entire dependency tree of swagger spec files
	// into structs
	openedFiles := map[string]bool{absRootFilepath: true}
	allSpecs := map[string]MagmaSwaggerSpec{}
	specsToVisit := []mscAndPath{{MagmaSwaggerSpec: targetSpec, path: absRootFilepath}}
	for len(specsToVisit) > 0 {
		nextSpec := specsToVisit[0]
		specsToVisit = specsToVisit[1:]
		allSpecs[nextSpec.path] = nextSpec.MagmaSwaggerSpec

		for _, dependencyPath := range nextSpec.MagmaGenMeta.Dependencies {
			absDependencyPath, err := filepath.Abs(filepath.Join(rootDir, dependencyPath))
			if err != nil {
				return nil, errors.Wrapf(err, "dependency filepath %s is invalid", dependencyPath)
			}
			if _, alreadyOpened := openedFiles[absDependencyPath]; alreadyOpened {
				continue
			}
			openedFiles[absDependencyPath] = true

			dependencySpec, err := readSwaggerSpec(absDependencyPath)
			if err != nil {
				return nil, errors.Wrap(err, "failed to read dependency tree of swagger specs")
			}
			specsToVisit = append(
				specsToVisit,
				mscAndPath{
					MagmaSwaggerSpec: dependencySpec,
					path:             absDependencyPath,
				},
			)
		}
	}
	return allSpecs, nil
}

func StripAndWriteSwaggerSpecs(specs map[string]MagmaSwaggerSpec, outDir string) error {
	for path, msc := range specs {
		sanitized := msc.ToSwaggerSpec()
		marshaledSanitized, err := yaml.Marshal(sanitized)
		if err != nil {
			return errors.Wrapf(err, "could not re-marshal swagger spec %s", path)
		}

		err = ioutil.WriteFile(filepath.Join(outDir, msc.MagmaGenMeta.TempGenFilename), marshaledSanitized, 0666)
		if err != nil {
			return errors.Wrapf(err, "could not write dependency swagger spec %s", msc.MagmaGenMeta.TempGenFilename)
		}
	}
	return nil
}

func readSwaggerSpec(filepath string) (MagmaSwaggerSpec, error) {
	fileContents, err := ioutil.ReadFile(filepath)
	if err != nil {
		return MagmaSwaggerSpec{}, errors.Wrapf(err, "could not open target file %s", filepath)
	}

	spec := MagmaSwaggerSpec{}
	err = yaml.Unmarshal(fileContents, &spec)
	if err != nil {
		return MagmaSwaggerSpec{}, errors.Wrapf(err, "could not parse target file %s as yml", filepath)
	}
	return spec, nil
}

// MarshalToYAML marshals a MagmaSwaggerSpec to a YAML-formatted string.
func (m *MagmaSwaggerSpec) MarshalToYAML() (string, error) {
	d, err := yaml.Marshal(&m)
	if err != nil {
		return "", err
	}
	return string(d), nil
}
