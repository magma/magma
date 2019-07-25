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
`
	TestOverrideYML = `---
# TEST YML

foo: 1234

la: override
`
	TestConfigDir         = "/tmp/"
	TestConfigOverrideDir = "/tmp/overrides"
)

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

	os.Remove(filepath.Join(TestConfigDir, serviceFileName))
	os.Remove(filepath.Join(TestConfigOverrideDir, serviceFileName))
	os.Remove(TestConfigOverrideDir)

}
