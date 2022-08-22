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

package oc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opencensus.io/stats/view"
)

func TestViewerRegistration(t *testing.T) {
	want := Views{&view.View{Name: "test"}}
	err := RegisterViewer("test", want)
	assert.NoError(t, err)
	got := GetViewer("test")
	assert.Equal(t, want, got)
	assert.Panics(t, func() { MustRegisterViewer("test", nil) })
}
