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

//+build !windows

package server

import (
	"errors"
	"net"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/magma/magma/src/go/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func removeNoError(string) error {
	return nil
}

func osFileNoError(string) (os.FileInfo, error) {
	return nil, nil
}

func netDialNoError(network, address string, timeout time.Duration) (net.Conn, error) {
	return nil, nil
}

func netDialError(network, address string, timeout time.Duration) (net.Conn, error) {
	return nil, errors.New("any error will do")
}

// TestCleanupExistingUnusedSock ensures that should an un-used existing socket exist, we clean it up
// and log a warning.
func TestCleanupExistingUnusedSock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, logBuffer := testutil.NewTestLogger()

	err := cleanupUnixSocket(logger, osFileNoError, removeNoError, netDialError, "arbitrary_file_path.sock")

	assert.Nil(t, err)
	assert.Equal(
		t,
		"WARN\tRemoving existing socket file; previous unclean shutdown?\n",
		logBuffer.String())
}

// TestCleanupExistingConnectedSock validates that should an existing server be hosting on a unix
// domain socket, we do not clean it up but instead log an error.
func TestCleanupExistingConnectedSock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, _ := testutil.NewTestLogger()

	err := cleanupUnixSocket(logger, osFileNoError, removeNoError, netDialNoError, "arbitrary_file_path.sock")

	assert.NotNil(t, err)
	wantErrMsg := "existing listener on socket file"
	assert.Containsf(t,
		err.Error(), wantErrMsg, "expected error containing %q, got %s", wantErrMsg, err)
}
