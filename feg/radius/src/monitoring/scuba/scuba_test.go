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

package scuba

import (
	"fbc/cwf/radius/config"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestAnalyticsModulesAuthenticate tests the Analytics module handling of the Authenticate RADIUS packet
func TestSendOdsCounters(t *testing.T) {
	// Arrange
	logger, _ := zap.NewDevelopment()
	Initialize(&config.Scuba{
		MessageQueueSize: 1,
		FlushIntervalSec: 1,
		BatchSize:        2,
		GraphURL:         "http://127.0.0.1:4321/scuba",
		AccessToken:      "at",
	}, logger)

	var gotLog = make(chan bool, 1)

	http.HandleFunc("/scuba", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			gotLog <- false
			return
		}
		body := string(bodyBytes)
		require.Contains(t, body, "this+is+a+log")
		require.Contains(t, body, "perfpipe_some_table")
		require.Contains(t, body, "xwf_json_to_any_scuba")
		gotLog <- true
	})

	go func() {
		if err := http.ListenAndServe(":4321", nil); err != nil {
			panic(err)
		}
	}()

	// Act
	logger, err := NewLogger("some_table")
	if err != nil {
		require.Fail(t, "failed to create logger")
	}
	logger.Info("this is a log")

	// Assert
	timeout := time.NewTimer(5 * time.Second)
	select {
	case success := <-gotLog:
		require.Equal(t, true, success)
	case <-timeout.C:
		require.Fail(t, "timed out waiting for metrics to propagate")
	}
}
