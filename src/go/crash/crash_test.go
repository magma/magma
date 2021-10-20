// Copyright 2021 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:generate go run github.com/golang/mock/mockgen -package crash -destination mock_crash/mock_crash.go . Crash

func TestLevel_String(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{
			level: LevelDebug,
			want:  "debug",
		},
		{
			level: LevelInfo,
			want:  "info",
		},
		{
			level: LevelWarning,
			want:  "warning",
		},
		{
			level: LevelError,
			want:  "error",
		},
		{
			level: LevelFatal,
			want:  "fatal",
		},
	}

	for _, test := range tests {
		got := string(test.level)
		assert.Equal(t, test.want, got)
	}
}
