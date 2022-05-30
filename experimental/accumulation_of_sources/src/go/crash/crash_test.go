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

package crash_test

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/magma/magma/src/go/crash"
	mock_crash "github.com/magma/magma/src/go/crash/mock_crash"
	"github.com/stretchr/testify/assert"
)

//go:generate go run github.com/golang/mock/mockgen -package crash -destination mock_crash/mock_crash.go . Crash

func TestLevel_String(t *testing.T) {
	tests := []struct {
		level crash.Level
		want  string
	}{
		{
			level: crash.LevelDebug,
			want:  "debug",
		},
		{
			level: crash.LevelInfo,
			want:  "info",
		},
		{
			level: crash.LevelWarning,
			want:  "warning",
		},
		{
			level: crash.LevelError,
			want:  "error",
		},
		{
			level: crash.LevelFatal,
			want:  "fatal",
		},
	}

	for _, test := range tests {
		got := string(test.level)
		assert.Equal(t, test.want, got)
	}
}

// TestWrap_Panic confirms flush and recover are called properly when the function passed into wrap panics. Also that
// the panic is re-raised to end the program.
func TestWrap_Panic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	err := errors.New("panic")
	mockCrash := mock_crash.NewMockCrash(ctrl)
	mockCrash.EXPECT().Recover(err)
	mockCrash.EXPECT().Flush(time.Second * 5)

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("Wrap did not re-raise the panic")
			}
		}()
		crash.Wrap(mockCrash, func() {
			panic(err)
		})
	}()
}

// TestWrap_NoPanic confirms flush and recover are not called when the function doesn't panic.
func TestWrap_NoPanic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCrash := mock_crash.NewMockCrash(ctrl)

	crash.Wrap(mockCrash, func() {})
}
