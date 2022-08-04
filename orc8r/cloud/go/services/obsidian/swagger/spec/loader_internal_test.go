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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecType_Path(t *testing.T) {
	tests := map[string]struct {
		specType specType
		base     string
		service  string
		want     string
	}{
		"partial feg": {
			specType: partial,
			base:     "base/dir",
			service:  "feg",
			want:     "base/dir/partial/feg.swagger.v1.yml",
		},
		"standalone lte": {
			specType: standalone,
			base:     "another/base/dir",
			service:  "LTE",
			want:     "another/base/dir/standalone/lte.swagger.v1.yml",
		},
		"partial nprobe": {
			specType: partial,
			base:     "base/place",
			service:  "nprobe",
			want:     "base/place/partial/nprobe.swagger.v1.yml",
		},
		"standalone smsd": {
			specType: standalone,
			base:     "a/b/c",
			service:  "SMSD",
			want:     "a/b/c/standalone/smsd.swagger.v1.yml",
		},
	}

	for desc, test := range tests {
		t.Run(desc, func(t *testing.T) {
			assert.Equal(t, test.want, test.specType.path(test.base, test.service))
		})
	}
}

func TestNewSpecs(t *testing.T) {
	want := "a/b/c"
	s := NewFSLoader(want)
	f, ok := s.(*fsLoader)
	assert.True(t, ok)
	assert.Equal(t, want, f.baseDir)
}

func TestDefaultSpecs(t *testing.T) {
	s := GetDefaultLoader()
	f, ok := s.(*fsLoader)
	assert.True(t, ok)
	assert.Equal(t, "/etc/magma/swagger/specs", f.baseDir)
}
