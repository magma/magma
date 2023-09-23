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

package reindex

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/glog"

	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/metrics"
	state_types "magma/orc8r/cloud/go/services/state/types"
)

type Error string

const (
	// ErrDefault is the default error reported for reindex failure.
	ErrDefault Error = "state reindex error"
	// ErrReindexPerState indicates a Reindex error occurred for specific keys. Not included in any job errors.
	ErrReindexPerState Error = "reindex error: per-state errors"

	// ErrPrepare is included in job error when error source is indexer PrepareReindex call.
	ErrPrepare Error = "state reindex error: error from PrepareReindex"
	// ErrReindex is included in job error when error source is indexer Reindex call.
	ErrReindex Error = "state reindex error: error from Reindex"
	// ErrComplete is included in job error when error source is indexer CompleteReindex call.
	ErrComplete Error = "state reindex error: error from CompleteReindex"

	failedJobSleep            = 1 * time.Minute
	failedGetBatchesSleep     = 5 * time.Second
	numStatesToReindexPerCall = 100
)

var (
	// TestHookReindexDone is a function called after each reindex in Run()
	// completes, regardless of success or failure.
	// This should only be set by test code.
	TestHookReindexDone = func() {}

	// TestHookReindexSuccess is a function called after each reindex in Run()
	// completes successfully.
	// This should only be set by test code.
	TestHookReindexSuccess = func() {}

	// TestHookReindexFailure is a function called whenever we fail
	// to establish a connection in the remote indexer.
	// This should only be set by test code.
	TestHookReindexFailure = func() {}
)

type Reindexer interface {
	// Run to progressively complete required reindex jobs.
	// Periodically polls the reindex job queue for reindex jobs, attempts to
	// complete the job, and writes back any encountered errors.
	// Returns only upon context cancellation, which can optionally be nil.
	Run(ctx context.Context)

	// RunUnsafe tries to complete all required reindex jobs.
	// If ID is non-empty, only tries to reindex specified indexer.
	// This function is intended for use only when automatic reindexing (via the
	// reindex queue) is disabled.
	// Arguments:
	//	- Loggable updates sent synchronously via sendUpdate
	// DO NOT use in parallel with Run().
	RunUnsafe(ctx context.Context, indexerID string, sendUpdate func(string)) error

	// GetIndexerVersions returns version info for all tracked indexers, keyed by indexer ID.
	GetIndexerVersions() ([]*indexer.Versions, error)
}

type reindexBatch struct {
	networkID string
	stateIDs  state_types.IDs
}

func executeJob(ctx context.Context, job *Job, batches []reindexBatch) error {
	id := job.Idx.GetID()
	stateTypes := job.Idx.GetTypes()

	isFirst := job.From == 0
	err := job.Idx.PrepareReindex(job.From, job.To, isFirst)
	if err != nil {
		return wrap(err, ErrPrepare, id)
	}

	for _, b := range batches {
		if isCanceled(ctx) {
			return wrap(err, ErrDefault, "context canceled")
		}
		ids := b.stateIDs.Filter(stateTypes...)
		if len(ids) == 0 {
			continue
		}

		// Convert IDs to states -- silently ignore not-found (stale) state IDs
		statesByID, err := state.GetSerializedStates(ctx, b.networkID, ids)
		if err != nil {
			err = fmt.Errorf("get states: %w", err)
			return wrap(err, ErrDefault, id)
		}

		errs, err := job.Idx.Index(b.networkID, statesByID)
		if err != nil {
			return wrap(err, ErrReindex, id)
		}
		if len(errs) == len(b.stateIDs) {
			err = errors.New("reindex call succeeded but all state IDs returned per-state reindex errors")
			return wrap(err, ErrReindex, id)
		} else if len(errs) != 0 {
			metrics.IndexErrors.WithLabelValues(id, getVersion(job), metrics.SourceValueReindex).Add(float64(len(errs)))
			glog.Warningf("%s: %s", ErrReindexPerState, errs)
		}
	}

	err = job.Idx.CompleteReindex(job.From, job.To)
	if err != nil {
		return wrap(err, ErrComplete, id)
	}

	return nil
}

func getIndexersFromRegistry(indexerID string) ([]indexer.Indexer, error) {
	idxs, err := indexer.GetIndexers()
	if err != nil {
		return nil, err
	}
	if indexerID == "" {
		return idxs, nil
	}

	for _, x := range idxs {
		if indexerID == x.GetID() {
			return []indexer.Indexer{x}, nil
		}
	}
	return nil, fmt.Errorf("indexer with ID %s not found in registry", indexerID)
}

func isCanceled(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	return ctx.Err() == context.Canceled
}

func getVersion(job *Job) string {
	return fmt.Sprint(job.Idx.GetVersion())
}

func getStatus(status Status) float64 {
	switch status {
	case StatusAvailable:
		return metrics.ReindexStatusIncomplete
	case StatusInProgress:
		return metrics.ReindexStatusInProcess
	case StatusComplete:
		return metrics.ReindexStatusSuccess
	}
	glog.Errorf("Unrecognized state reindexer job status: %s", status)
	return metrics.ReindexStatusIncomplete
}

func wrap(err error, sentinel Error, indexerID string) error {
	var wrap string
	switch indexerID {
	case "":
		wrap = string(sentinel)
	default:
		wrap = fmt.Sprintf("%s for idx %s", sentinel, indexerID)
	}
	return fmt.Errorf(wrap+": %w", err)
}
