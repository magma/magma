/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

// SwaggerConfig is used only for marshalling / unmarshalling from yml
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

// TagDefinition is used only for marshalling / unmarshalling from yml
type TagDefinition struct {
	Description string
	Name        string
}

func main() {
	// Parse options
	inpPtr := flag.String("inp", "", "Input folder")
	comPtr := flag.String("common", "", "Common Definitions path")
	outPtr := flag.String("out", "", "Output path")
	flag.Parse()

	// Now do the rest of the stuff
	fmt.Printf("Reading swagger configs from folder:\n%s\n\n", *inpPtr)
	inpConfigArr := readFromInpFolder(*inpPtr)

	fmt.Printf("Reading common config definitions from file:\n%s\n\n", *comPtr)
	commonConfig := readCommonConfig(*comPtr)

	fmt.Printf("Attempting to combine configs together..\n\n")
	outConfig := combineConfigs(commonConfig, inpConfigArr)

	fmt.Printf("Writing combined swagger config to file:\n%s\n\n", *outPtr)
	writeOutConfig(outConfig, *outPtr)
}

// For a list of input file paths, unmarshal the file contents
func readFromInpFolder(inpPath string) []SwaggerConfig {
	filePathArr := getFilePaths(inpPath)
	fileContentArr := getContentsOfFiles(filePathArr)
	editedContentArr := makeAllYmlReferencesLocal(fileContentArr)
	configArr := unmarshalArrToSwagger(editedContentArr)
	return configArr
}

// Unmarshal the common swagger config from its file path
func readCommonConfig(inpPath string) SwaggerConfig {
	contents := getFileContents(inpPath)
	return unmarshalToSwagger(contents)
}

func combineConfigs(common SwaggerConfig, inpArr []SwaggerConfig) SwaggerConfig {
	out := SwaggerConfig{}
	out.Swagger = common.Swagger
	out.Info = common.Info
	out.BasePath = common.BasePath
	out.Consumes = common.Consumes
	out.Produces = common.Produces
	out.Schemes = common.Schemes

	var allTags [][]TagDefinition
	var allInpPaths []map[string]interface{}
	var allInpResponses []map[string]interface{}
	var allInpParameters []map[string]interface{}
	var allInpDefinitions []map[string]interface{}

	for _, inpConfig := range inpArr {
		allTags = append(allTags, inpConfig.Tags)
		allInpPaths = append(allInpPaths, inpConfig.Paths)
		allInpResponses = append(allInpResponses, inpConfig.Responses)
		allInpParameters = append(allInpParameters, inpConfig.Parameters)
		allInpDefinitions = append(allInpDefinitions, inpConfig.Definitions)
	}

	out.Tags = combineTags(common.Tags, allTags)
	out.Paths = combineSubConfig(common.Paths, allInpPaths)
	out.Responses = combineSubConfig(common.Responses, allInpResponses)
	out.Parameters = combineSubConfig(common.Parameters, allInpParameters)
	out.Definitions = combineSubConfig(common.Definitions, allInpDefinitions)

	return out
}

func combineSubConfig(common map[string]interface{}, inpArr []map[string]interface{}) map[string]interface{} {
	outSubConfig := make(map[string]interface{})
	for _, subConfig := range inpArr {
		keyArr := keysOfMap(subConfig)
		for _, key := range keyArr {
			outSubConfig[key] = subConfig[key]
		}
	}
	// Use the common ones
	keyArr := keysOfMap(common)
	for _, key := range keyArr {
		outSubConfig[key] = common[key]
	}
	return outSubConfig
}

func combineTags(common []TagDefinition, inpArr [][]TagDefinition) []TagDefinition {
	outTags := make(map[string]string)
	for _, tagArr := range inpArr {
		for _, tag := range tagArr {
			outTags[tag.Name] = tag.Description
		}
	}
	for _, tag := range common {
		outTags[tag.Name] = tag.Description
	}
	keys := keysOfStrMap(outTags)
	uniqueOutTags := make([]TagDefinition, len(keys))
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		uniqueOutTags[i] = TagDefinition{Name: key, Description: outTags[key]}
	}
	return uniqueOutTags
}

func writeOutConfig(outConfig SwaggerConfig, outPath string) {
	strConfig := marshalFromSwagger(outConfig)
	f, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(strConfig)
	f.Sync()
}

func unmarshalArrToSwagger(fileContentArr []string) []SwaggerConfig {
	var configArr []SwaggerConfig
	for _, fileContent := range fileContentArr {
		config := unmarshalToSwagger(fileContent)
		configArr = append(configArr, config)
	}
	return configArr
}

func unmarshalToSwagger(fileContents string) SwaggerConfig {
	config := SwaggerConfig{}
	err := yaml.Unmarshal([]byte(fileContents), &config)
	if err != nil {
		panic(err)
	}
	return config
}

func marshalFromSwagger(config SwaggerConfig) string {
	d, err := yaml.Marshal(&config)
	if err != nil {
		panic(err)
	}
	return string(d)
}

// Change all cross-file references to local intra-file references
func makeAllYmlReferencesLocal(fileContentArr []string) []string {
	outArr := make([]string, 0, len(fileContentArr))

	// match on any yml reference to file_name_here.foo.bar#/baz
	// and change those references to #/baz (strip the prefix)
	// e.g. $ref: 'foo_bar_baz.blah#/asdf' -> $ref: '#/asdf'
	ymlRefRe := regexp.MustCompile(`(\$ref:\s*)['"].+(#/.+)['"]`)
	for _, fileContent := range fileContentArr {
		outArr = append(outArr, ymlRefRe.ReplaceAllString(fileContent, "$1'$2'"))
	}
	return outArr
}

// Get the text contents from an array of file paths
func getContentsOfFiles(filePathArr []string) []string {
	var fileContentArr []string
	for _, filePath := range filePathArr {
		contents := getFileContents(filePath)
		fileContentArr = append(fileContentArr, contents)
	}
	return fileContentArr
}

// Get the text contents from a file path
func getFileContents(filePath string) string {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return string(data)
}

// Return swagger config file paths from a folder path
func getFilePaths(inpPath string) []string {
	var filePathArr []string
	filepath.Walk(inpPath, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".yml") {
			filePathArr = append(filePathArr, path)
		}
		return nil
	})
	return filePathArr
}

func keysOfMap(inp map[string]interface{}) []string {
	keys := reflect.ValueOf(inp).MapKeys()
	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
	}
	// Sort before returning for more deterministic behavior
	sort.Strings(strkeys)
	return strkeys
}

func keysOfStrMap(inp map[string]string) []string {
	keys := reflect.ValueOf(inp).MapKeys()
	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
	}
	// Sort before returning for more deterministic behavior
	sort.Strings(strkeys)
	return strkeys
}
