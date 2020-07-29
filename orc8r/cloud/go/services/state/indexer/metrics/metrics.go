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

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	// IndexerNameLabel values contain the name of the relevant indexer.
	IndexerNameLabel = "indexerID"

	// IndexerVersionLabel values contain the version of the relevant indexer.
	// Values should derive from indexer.Version.
	IndexerVersionLabel = "indexerVersion"

	// SourceLabel values indicates whether the metric was produced during
	// indexing or reindexing operations.
	// Values should derive from MetricSource.
	SourceLabel = "metricSource"
	// SourceValueIndex indicates the metric originated during normal indexing operations.
	SourceValueIndex = "index"
	// SourceValueReindex indicates the metric originated during a reindex job.
	SourceValueReindex = "reindex"

	// ReindexStatusSuccess indicates the reported job as a whole has completed successfully.
	ReindexStatusSuccess = 1
	// ReindexStatusInProcess indicates the reported job as a whole is currently active.
	ReindexStatusInProcess = 0
	// ReindexStatusIncomplete indicates the reported job as a whole is incomplete.
	// This is the default status.
	ReindexStatusIncomplete = -1
)

var (
	IndexErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "stateindexer_indexerrors_count",
			Help: "Number of per-state errors reported by indexers",
		},
		[]string{IndexerNameLabel, IndexerVersionLabel, SourceLabel},
	)
	IndexerVersion = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "stateindexer_indexerversion",
			Help: "Indexer's actual version (may be stale)",
		},
		[]string{IndexerNameLabel},
	)
	ReindexStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "stateindexer_reindexjob_status",
			Help: "Status of indexer's most-recent reindex job: 1 succeeded, 0 in-progress, -1 considered failed",
		},
		[]string{IndexerNameLabel},
	)
	ReindexAttempts = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "stateindexer_reindexjob_attempts",
			Help: "Number of times indexer's most-recent reindex job has been attempted",
		},
		[]string{IndexerNameLabel},
	)
	ReindexDuration = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "stateindexer_reindexjob_seconds",
			Help: "Duration of indexer's most-recent reindex job attempt",
		},
		[]string{IndexerNameLabel},
	)
)
