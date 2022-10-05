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
	assert.NoError(t, err)
	f.Close()
	err = os.Mkdir(TestConfigOverrideDir, 0700)
	assert.NoError(t, err)

	configMap, err := getServiceConfigImpl("", "test", TestConfigDir, "", TestConfigOverrideDir)
	assert.NoError(t, err)
	foo, err := configMap.GetInt("foo")
	assert.NoError(t, err)
	assert.Equal(t, 8443, foo)

	bar, err := configMap.GetString("bar")
	assert.NoError(t, err)
	assert.Equal(t, "something", bar)

	baz, err := configMap.GetStrings("baz")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(baz))

	assert.Equal(t, "first", baz[0])
	assert.Equal(t, "second", baz[1])
	assert.Equal(t, "third", baz[2])

	mapParam, err := configMap.GetMap("map1")
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
	assert.NoError(t, err)
	f.Close()
	err = os.Mkdir(TestConfigOverrideDir, 0700)
	assert.NoError(t, err)
	f, err = os.Create(filepath.Join(TestConfigOverrideDir, serviceFileName))
	assert.NoError(t, err)

	_, err = f.WriteString(TestOverrideYML)
	assert.NoError(t, err)
	f.Close()

	// ensure case insensitivity
	for _, serviceName := range []string{"test_service", "TEST_SERVICE", "TesT_sERviCE"} {
		configMap, err := getServiceConfigImpl("", serviceName, TestConfigDir, "", TestConfigOverrideDir)
		assert.NoError(t, err)
		foo, err := configMap.GetInt("foo")
		assert.NoError(t, err)
		assert.Equal(t, 1234, foo)

		bar, err := configMap.GetString("bar")
		assert.NoError(t, err)
		assert.Equal(t, "something", bar)

		la, err := configMap.GetString("la")
		assert.NoError(t, err)
		assert.Equal(t, "override", la)
	}

	var testStruct ConfigTestStruct
	_, _, err = GetStructuredServiceConfigExt(
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
	_, _, err = GetStructuredServiceConfigExt(
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
