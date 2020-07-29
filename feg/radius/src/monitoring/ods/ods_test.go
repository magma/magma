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

package ods

import (
	"context"
	"fbc/cwf/radius/config"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.uber.org/zap"
)

// TestAnalyticsModulesAuthenticate tests the Analytics module handling of the Authenticate RADIUS packet
func TestSendOdsCounters(t *testing.T) {
	// Arrange
	logger, _ := zap.NewDevelopment()
	Init(&config.Ods{
		ReportingPeriod: config.Duration{time.Second},
		GraphURL:        "http://127.0.0.1:1234/ods",
		Entity:          "entity",
		Category:        "123",
		Prefix:          "lalala",
	}, logger)

	tg, _ := tag.NewKey("test")
	ctr := stats.Int64("ctr", "Counter", stats.UnitDimensionless)
	view.Register(
		&view.View{
			Name:        "ctr_view",
			Measure:     ctr,
			Description: "Counter View",
			Aggregation: view.Count(),
		},
	)

	var gotMetrics = make(chan bool, 1)

	http.HandleFunc("/ods", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			gotMetrics <- false
			return
		}
		t.Logf("got event: %s", body)
		gotMetrics <- true
	})

	go func() {
		if err := http.ListenAndServe(":1234", nil); err != nil {
			panic(err)
		}
	}()

	// Act
	stats.RecordWithTags(
		context.Background(),
		[]tag.Mutator{tag.Upsert(tg, "test")},
		ctr.M(1),
	)
	time.Sleep(time.Duration(time.Second))

	// Assert
	timeout := time.NewTimer(5 * time.Second)
	select {
	case success := <-gotMetrics:
		require.Equal(t, true, success)
	case <-timeout.C:
		require.Fail(t, "timed out waiting for metrics to propagate")
	}
}
