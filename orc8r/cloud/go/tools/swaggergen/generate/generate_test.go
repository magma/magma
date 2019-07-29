/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
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
	"gopkg.in/yaml.v2"
)

func TestParseSwaggerDependencyTree(t *testing.T) {
	actual, err := generate.ParseSwaggerDependencyTree("../testdata/importer2.yml", os.Getenv("MAGMA_ROOT"))
	assert.NoError(t, err)

	expectedFiles := []string{"../testdata/base.yml", "../testdata/importer.yml", "../testdata/importer2.yml"}
	expected := parseExpectedFiles(t, expectedFiles)

	assert.Equal(t, expected, actual)
}

func TestParseSwaggerDependencyTree_Cycle(t *testing.T) {
	actual, err := generate.ParseSwaggerDependencyTree("../testdata/cycle1.yml", os.Getenv("MAGMA_ROOT"))
	assert.NoError(t, err)

	expectedFiles := []string{"../testdata/cycle1.yml", "../testdata/cycle2.yml"}
	expected := parseExpectedFiles(t, expectedFiles)

	assert.Equal(t, expected, actual)
}

func TestGenerateModels(t *testing.T) {
	runTestGenerateCase(t, "../testdata/importer2.yml", "../testdata/importer2")
	runTestGenerateCase(t, "../testdata/importer.yml", "../testdata/importer")
	runTestGenerateCase(t, "../testdata/base.yml", "../testdata/base")
}

// outputDir has to match what's specified in the yml
func runTestGenerateCase(t *testing.T, ymlFile string, outputDir string) {
	defer cleanupActualFiles(outputDir)

	err := generate.GenerateModels(ymlFile, "../testdata/template.yml", os.Getenv("MAGMA_ROOT"))
	assert.NoError(t, err)

	// Verify that generated files are the same as the expected golden files
	goldenFiles, actualFiles := []string{}, []string{}
	err = filepath.Walk(outputDir, func(path string, _ os.FileInfo, _ error) error {
		if strings.HasSuffix(path, "go.golden") {
			goldenFiles = append(goldenFiles, strings.TrimSuffix(path, ".golden"))
		} else if strings.HasSuffix(path, ".actual") {
			actualFiles = append(actualFiles, strings.TrimSuffix(path, ".actual")+".go")
		}
		return nil
	})
	assert.NoError(t, err)
	sort.Strings(goldenFiles)
	sort.Strings(actualFiles)
	assert.Equal(t, goldenFiles, actualFiles)

	// Verify contents of golden vs actual files
	for _, baseFilename := range goldenFiles {
		goldenFileContents, err := ioutil.ReadFile(baseFilename + ".golden")
		assert.NoError(t, err)
		actualFileContents, err := ioutil.ReadFile(strings.TrimSuffix(baseFilename, ".go") + ".actual")
		assert.NoError(t, err)
		assert.Equal(t, goldenFileContents, actualFileContents)
	}
}

func parseExpectedFiles(t *testing.T, files []string) map[string]generate.MagmaSwaggerConfig {
	expected := map[string]generate.MagmaSwaggerConfig{}
	for _, expectedPath := range files {
		expectedAbs, err := filepath.Abs(expectedPath)
		assert.NoError(t, err)
		contents, err := ioutil.ReadFile(expectedAbs)
		assert.NoError(t, err)

		expectedStruct := generate.MagmaSwaggerConfig{}
		err = yaml.Unmarshal(contents, &expectedStruct)
		assert.NoError(t, err)
		expected[expectedAbs] = expectedStruct
	}
	return expected
}

func cleanupActualFiles(outputDir string) {
	_ = recover()

	cleanupFiles := []string{}
	_ = filepath.Walk(outputDir, func(path string, _ os.FileInfo, _ error) error {
		if strings.HasSuffix(path, "actual") {
			cleanupFiles = append(cleanupFiles, path)
		}
		return nil
	})

	for _, cleanupFile := range cleanupFiles {
		_ = os.Remove(cleanupFile)
	}
}
