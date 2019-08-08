/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package generate

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// MagmaSwaggerConfig is the Go struct version of our custom Swagger spec file.
// The only difference is the magma-gen-meta field, which specifies
// dependencies (other swagger specs that the spec has refs to), a desired
// filename that this file should be ref'd with from dependent files, and
// a list of Go struct types and filenames that this file produces when
// models are generated.
type MagmaSwaggerConfig struct {
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
	Tags        []TagDefinition
	Paths       map[string]interface{}
	Responses   map[string]interface{}
	Parameters  map[string]interface{}
	Definitions map[string]interface{}
}

func (msc MagmaSwaggerConfig) ToSwaggerConfig() SwaggerConfig {
	ret := SwaggerConfig{}
	ret.Swagger = msc.Swagger
	ret.Info = msc.Info
	ret.BasePath = msc.BasePath
	ret.Consumes = msc.Consumes
	ret.Produces = msc.Produces
	ret.Schemes = msc.Schemes
	ret.Tags = msc.Tags
	ret.Paths = msc.Paths
	ret.Responses = msc.Responses
	ret.Parameters = msc.Parameters
	ret.Definitions = msc.Definitions
	return ret
}

// SwaggerConfig is the Go struct version of a OAI/Swagger 2.0 YAML spec file.
type SwaggerConfig struct {
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

const tmpGenDir = "tmpgen"

// GenerateModels parses the magma-gen-meta key of the given swagger YML file,
// copies the files that the target file depends on into the current working
// directory, shells out to `swagger generate models`, then cleans up the
// dependency files.
func GenerateModels(targetFilepath string, templateFilepath string, rootDir string) error {
	absTargetFilepath, err := filepath.Abs(targetFilepath)
	if err != nil {
		return errors.Wrapf(err, "target filepath %s is invalid", targetFilepath)
	}

	allConfigs, err := ParseSwaggerDependencyTree(absTargetFilepath, rootDir)
	if err != nil {
		return errors.WithStack(err)
	}

	defer func() {
		// we always want to do the cleanup step
		r := recover()
		_ = os.RemoveAll("tmpgen")
		if r != nil {
			// repanic after cleaning up
			panic(r)
		}
	}()

	// For each dependency, strip the magma-gen-meta and write the result to
	// the filename specified by `dependent-filename`
	err = StripAndWriteSwaggerConfigs(allConfigs)
	if err != nil {
		return errors.WithStack(err)
	}

	// Shell out to go-swagger
	targetConfig := allConfigs[absTargetFilepath]
	outputDir := filepath.Join(os.Getenv("MAGMA_ROOT"), targetConfig.MagmaGenMeta.OutputDir)
	cmd := exec.Command(
		"swagger", "generate", "model",
		"-f", filepath.Join(tmpGenDir, targetConfig.MagmaGenMeta.TempGenFilename),
		"-t", outputDir,
		"-C", templateFilepath,
	)
	stdOutBuffer := &strings.Builder{}
	stdErrBuffer := &strings.Builder{}
	cmd.Stdout = stdOutBuffer
	cmd.Stderr = stdErrBuffer

	err = cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "failed to generate models; stdout:\n%s\nstderr:\n%s", stdOutBuffer.String(), stdErrBuffer.String())
	}

	return nil
}

// ParseSwaggerDependencyTree parses the entire dependency tree of a magma
// swagger spec file specified by the rootFilepath parameter.
// The returned value maps between the absolute specified dependency filepath
// and the parsed struct for the dependency file.
func ParseSwaggerDependencyTree(rootFilepath string, rootDir string) (map[string]MagmaSwaggerConfig, error) {
	absRootFilepath, err := filepath.Abs(rootFilepath)
	if err != nil {
		return nil, errors.Wrapf(err, "root filepath %s is invalid", rootFilepath)
	}

	targetConfig, err := readSwaggerSpec(absRootFilepath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	type mscAndPath struct {
		MagmaSwaggerConfig
		path string
	}

	// Do a BFS to parse the entire dependency tree of swagger config files
	// into structs
	openedFiles := map[string]bool{absRootFilepath: true}
	allConfigs := map[string]MagmaSwaggerConfig{}
	configsToVisit := []mscAndPath{{MagmaSwaggerConfig: targetConfig, path: absRootFilepath}}
	for len(configsToVisit) > 0 {
		nextConfig := configsToVisit[0]
		configsToVisit = configsToVisit[1:]
		allConfigs[nextConfig.path] = nextConfig.MagmaSwaggerConfig

		for _, dependencyPath := range nextConfig.MagmaGenMeta.Dependencies {
			absDependencyPath, err := filepath.Abs(filepath.Join(rootDir, dependencyPath))
			if err != nil {
				return nil, errors.Wrapf(err, "dependency filepath %s is invalid", dependencyPath)
			}
			if _, alreadyOpened := openedFiles[absDependencyPath]; alreadyOpened {
				continue
			}
			openedFiles[absDependencyPath] = true

			dependencyConfig, err := readSwaggerSpec(absDependencyPath)
			if err != nil {
				return nil, errors.Wrap(err, "failed to read dependency tree of swagger configs")
			}
			configsToVisit = append(
				configsToVisit,
				mscAndPath{
					MagmaSwaggerConfig: dependencyConfig,
					path:               absDependencyPath,
				},
			)
		}
	}
	return allConfigs, nil
}

func StripAndWriteSwaggerConfigs(allConfigs map[string]MagmaSwaggerConfig) error {
	err := os.Mkdir(tmpGenDir, 0777)
	if err != nil {
		return errors.Wrap(err, "could not create temporary gen directory")
	}

	for path, msc := range allConfigs {
		sanitized := msc.ToSwaggerConfig()
		marshaledSanitized, err := yaml.Marshal(sanitized)
		if err != nil {
			return errors.Wrapf(err, "could not re-marshal swagger config %s", path)
		}

		err = ioutil.WriteFile(filepath.Join(tmpGenDir, msc.MagmaGenMeta.TempGenFilename), marshaledSanitized, 0666) // \m/
		if err != nil {
			return errors.Wrapf(err, "could not write dependency swagger config %s", msc.MagmaGenMeta.TempGenFilename)
		}
	}
	return nil
}

func readSwaggerSpec(filepath string) (MagmaSwaggerConfig, error) {
	fileContents, err := ioutil.ReadFile(filepath)
	if err != nil {
		return MagmaSwaggerConfig{}, errors.Wrapf(err, "could not open target file %s", filepath)
	}

	config := MagmaSwaggerConfig{}
	err = yaml.Unmarshal(fileContents, &config)
	if err != nil {
		return MagmaSwaggerConfig{}, errors.Wrapf(err, "could not parse target file %s as yml", filepath)
	}
	return config, nil
}
