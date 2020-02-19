/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TestYML = `---
# TEST YML

foo: 8443

bar: something

la: test

baz:
  - first
  - second
  - third

map1:
  key1: value_1
  key2: value_2

foo_int: 12345
`
	TestOverrideYML = `---
# TEST YML

foo: 1234

la: override

bla_bla_bla: bla_bla_bla value
`
	TestConfigDir         = "/tmp/"
	TestConfigOverrideDir = "/tmp/overrides"
)

type ConfigTestStruct struct {
	Foo    string `yaml:"foo"`
	Bar    string
	La     string
	Baz    []string `yaml:"baz"`
	BlaBla string   `yaml:"bla_bla_bla"`
	Fooint int      `yaml:"foo_int"`
}

func TestGetConfigWithoutOverride(t *testing.T) {
	serviceFileName := "test.yml"
	f, err := os.Create(filepath.Join(TestConfigDir, serviceFileName))
	assert.NoError(t, err)

	_, err = f.WriteString(TestYML)
	f.Close()
	err = os.Mkdir(TestConfigOverrideDir, 0700)
	assert.NoError(t, err)

	configMap, err := getServiceConfigImpl("", "test", TestConfigDir, "", TestConfigOverrideDir)
	assert.NoError(t, err)
	foo, err := configMap.GetIntParam("foo")
	assert.NoError(t, err)
	assert.Equal(t, 8443, foo)

	bar, err := configMap.GetStringParam("bar")
	assert.NoError(t, err)
	assert.Equal(t, "something", bar)

	baz, err := configMap.GetStringArrayParam("baz")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(baz))

	assert.Equal(t, "first", baz[0])
	assert.Equal(t, "second", baz[1])
	assert.Equal(t, "third", baz[2])

	mapParam, err := configMap.GetMapParam("map1")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(mapParam))
	assert.Equal(t, "value_1", mapParam["key1"])
	assert.Equal(t, "value_2", mapParam["key2"])

	os.Remove(filepath.Join(TestConfigDir, serviceFileName))
	os.Remove(TestConfigOverrideDir)
}

func TestGetConfigWithOverride(t *testing.T) {
	serviceFileName := "test_service.yml"
	f, err := os.Create(filepath.Join(TestConfigDir, serviceFileName))
	assert.NoError(t, err)

	_, err = f.WriteString(TestYML)
	f.Close()
	err = os.Mkdir(TestConfigOverrideDir, 0700)
	assert.NoError(t, err)
	f, err = os.Create(filepath.Join(TestConfigOverrideDir, serviceFileName))
	assert.NoError(t, err)

	_, err = f.WriteString(TestOverrideYML)
	f.Close()

	configMap, err := getServiceConfigImpl("", "test_service", TestConfigDir, "", TestConfigOverrideDir)
	assert.NoError(t, err)
	foo, err := configMap.GetIntParam("foo")
	assert.NoError(t, err)
	assert.Equal(t, 1234, foo)

	bar, err := configMap.GetStringParam("bar")
	assert.NoError(t, err)
	assert.Equal(t, "something", bar)

	la, err := configMap.GetStringParam("la")
	assert.NoError(t, err)
	assert.Equal(t, "override", la)

	var testStruct ConfigTestStruct
	err = GetStructuredServiceConfigExt(
		"", "test_service", TestConfigDir, "", TestConfigOverrideDir, &testStruct)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(testStruct.Baz))
	assert.Equal(t, "first", testStruct.Baz[0])
	assert.Equal(t, "second", testStruct.Baz[1])
	assert.Equal(t, "third", testStruct.Baz[2])
	assert.Equal(t, "1234", testStruct.Foo)
	assert.Equal(t, "something", testStruct.Bar)
	assert.Equal(t, "override", testStruct.La)
	assert.Equal(t, "bla_bla_bla value", testStruct.BlaBla)
	assert.Equal(t, 12345, testStruct.Fooint)

	var testIndirectStruct ConfigTestStruct
	testIndirectStruct.Fooint = 54321
	var testInterface interface{} = &testIndirectStruct
	err = GetStructuredServiceConfigExt(
		"", "test_service", TestConfigDir, "", TestConfigOverrideDir, testInterface)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(testIndirectStruct.Baz))
	assert.Equal(t, "first", testIndirectStruct.Baz[0])
	assert.Equal(t, "second", testIndirectStruct.Baz[1])
	assert.Equal(t, "third", testIndirectStruct.Baz[2])
	assert.Equal(t, "1234", testIndirectStruct.Foo)
	assert.Equal(t, "something", testIndirectStruct.Bar)
	assert.Equal(t, "override", testIndirectStruct.La)
	assert.Equal(t, "bla_bla_bla value", testIndirectStruct.BlaBla)
	assert.Equal(t, 12345, testStruct.Fooint)

	os.Remove(filepath.Join(TestConfigDir, serviceFileName))
	os.Remove(filepath.Join(TestConfigOverrideDir, serviceFileName))
	os.Remove(TestConfigOverrideDir)
}
