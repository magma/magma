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

package census

import (
	"net/http"

	"go.uber.org/zap"
)

// Init Initialize views and Prometheus exporter
func Init(config Config, logger *zap.Logger) {
	// Create metrics server
	census, err := config.Build(logger)
	if err != nil {
		logger.Error("Failed building census", zap.Error(err))
		return
	}
	http.Handle("/metrics", census.StatsHandler)
	go func() {
		defer census.Close()
		http.ListenAndServe(":9100", nil)
	}()
}
