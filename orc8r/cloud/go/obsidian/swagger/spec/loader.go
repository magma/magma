/*
 Copyright 2021 The Magma Authors.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package spec

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Loader provide access to swagger spec YAML files.
type Loader interface {
	GetCommonSpec() (string, error)
	GetPartialSpec(service string) (string, error)
	GetStandaloneSpec(service string) (string, error)
}

// fsLoader implements Loader via the filesystem at a specific base dir.
type fsLoader struct {
	baseDir string
}

const commonSpecRelPath = "common/swagger-common.yml"

// GetCommonSpec returns ./common/swagger-common.yml
func (f *fsLoader) GetCommonSpec() (string, error) {
	return readEntireFileString(filepath.Join(f.baseDir, commonSpecRelPath))
}

// GetPartialSpec returns ./partial/$service.swagger.v1.yml
func (f *fsLoader) GetPartialSpec(service string) (string, error) {
	return readEntireFileString(partial.path(f.baseDir, service))
}

// GetStandaloneSpec returns ./standalone/$service.swagger.v1.yml
func (f *fsLoader) GetStandaloneSpec(service string) (string, error) {
	return readEntireFileString(standalone.path(f.baseDir, service))
}

// NewFSLoader returns Loader backed by fsLoader using the specified base path.
// See also GetDefaultLoader which embeds the default base path and is preferred
// when appropriate.
func NewFSLoader(baseDir string) Loader {
	return &fsLoader{baseDir: baseDir}
}

const defaultSpecBasePath = "/etc/magma/swagger/specs"

// GetDefaultLoader returns Loader with the default production base path.
func GetDefaultLoader() Loader {
	return NewFSLoader(defaultSpecBasePath)
}

type specType string

const (
	partial    specType = "partial"
	standalone specType = "standalone"
)

func (t specType) path(base string, service string) string {
	return filepath.Join(base, string(t), fmt.Sprintf("%s.swagger.v1.yml", strings.ToLower(service)))
}

// readEntireFileString wraps ioutil.ReadFile and returns a string iff the
// entire file can be read without error.
func readEntireFileString(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
