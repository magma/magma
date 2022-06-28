/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package signal_test

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"magma/dp/cloud/go/active_mode_controller/internal/signal"
)

const (
	timeout = time.Millisecond * 10
)

func TestRun(t *testing.T) {
	app, errs := whenAppWasStarted()
	thenAppWasStarted(t, app)

	whenSignalWasSent(t)
	thenAppWasShutdown(t, errs)
}

func whenAppWasStarted() (*stubApp, chan error) {
	errs := make(chan error)
	app := &stubApp{running: make(chan struct{})}
	go func() {
		errs <- signal.Run(context.Background(), app)
	}()
	return app, errs
}

func thenAppWasStarted(t *testing.T, s *stubApp) {
	select {
	case <-s.running:
		return
	case <-time.After(timeout):
		assert.Fail(t, "app failed to start")
	}
}

func whenSignalWasSent(t *testing.T) {
	proc, err := os.FindProcess(os.Getpid())
	assert.NoError(t, err)
	assert.NoError(t, proc.Signal(syscall.SIGTERM))
}

func thenAppWasShutdown(t *testing.T, errs chan error) {
	select {
	case err := <-errs:
		assert.NoError(t, err)
	case <-time.After(timeout):
		assert.Fail(t, "app failed to stop")
	}
}

type stubApp struct {
	running chan struct{}
}

func (s *stubApp) Run(ctx context.Context) error {
	close(s.running)
	<-ctx.Done()
	return nil
}
