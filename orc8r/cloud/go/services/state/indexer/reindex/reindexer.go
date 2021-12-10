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
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/metrics"
	state_types "magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/lib/go/util"
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
			err = errors.Wrap(err, "get states")
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

func wrap(err error, sentinel Error, indexerID string) error {
	var wrap string
	switch indexerID {
	case "":
		wrap = string(sentinel)
	default:
		wrap = fmt.Sprintf("%s for idx %s", sentinel, indexerID)
	}
	return errors.Wrap(err, wrap)
}

// reindexer specific implementation

// TODO(reginawang3495): Remove "Desired" from Indexer.Versions once reindexer_queue and queue are removed

type reindexer struct {
	Versioner
	store Store
}

const reindexLoopInterval = time.Minute

func NewReindexer(store Store, versioner Versioner) Reindexer {
	return &reindexer{store: store, Versioner: versioner}
}

func (r *reindexer) Run(ctx context.Context) {
	// indexerID being "" means that we are reindexing all indexers
	const indexerID = ""
	// TODO(reginawang3495) will support non-nil sendUpdate from Run after removing queue impl
	var sendUpdate func(string) = nil

	for {
		if isCanceled(ctx) {
			glog.Warning("State reindexing async job canceled")
			return
		}

		err := r.findAndReindexJobs(ctx, indexerID, sendUpdate)
		if err != nil {
			glog.Errorf("Failed to getJobs for indexerID %s: %v", indexerID, err)
		}

		clock.Sleep(reindexLoopInterval)
		glog.Infof("Sleeping for %v minute(s) before looking for new jobs to reindex", reindexLoopInterval.Minutes())
	}
}

// TODO(reginawang3495): remove RunUnsafe when jobQueue reindexer is removed
func (r *reindexer) RunUnsafe(ctx context.Context, indexerID string, sendUpdate func(string)) error {
	jobs, err := r.getJobs(indexerID)
	if err != nil {
		return wrap(err, ErrDefault, indexerID)
	}
	if jobs == nil {
		return nil
	}

	batches := r.getReindexBatches(ctx)
	return r.reindexJobs(ctx, jobs, batches, sendUpdate)
}

func (r *reindexer) findAndReindexJobs(ctx context.Context, indexerID string, sendUpdate func(string)) error {
	jobs, err := r.getJobs(indexerID)
	if err != nil {
		return err
	}
	if jobs != nil {
		batches := r.getReindexBatches(ctx)
		r.reindexJobs(ctx, jobs, batches, sendUpdate)
	}
	return nil
}

// getReindexBatches gets network-segregated reindex batches with capped number of state IDs per batch.
func (r *reindexer) getReindexBatches(ctx context.Context) []reindexBatch {
	var idsByNetwork state_types.IDsByNetwork
	for {
		if isCanceled(ctx) {
			return nil
		}
		ids, err := r.store.GetAllIDs()
		if err != nil {
			err = wrap(err, ErrDefault, "")
			glog.Errorf("Failed to get all state IDs for state indexer reindexing, will retry: %v", err)
			clock.Sleep(failedGetBatchesSleep)
			continue
		}
		idsByNetwork = ids
		break
	}

	var current, rest state_types.IDs
	var batches []reindexBatch
	for networkID, ids := range idsByNetwork {
		rest = ids
		for len(rest) > 0 {
			count := util.MinInt(numStatesToReindexPerCall, len(rest))
			current, rest = rest[:count], rest[count:]
			batches = append(batches, reindexBatch{networkID: networkID, stateIDs: current})
		}
	}
	return batches
}

func (r *reindexer) reindexJobs(ctx context.Context, jobs []*Job, batches []reindexBatch, sendUpdate func(string)) error {
	errs := &multierror.Error{}
	for _, j := range jobs {
		err := r.reindexJob(j, ctx, batches, sendUpdate)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs.ErrorOrNil()
}

func (r *reindexer) reindexJob(job *Job, ctx context.Context, batches []reindexBatch, sendUpdate func(string)) error {
	defer TestHookReindexDone()

	start := clock.Now()

	jobErr := executeJob(ctx, job, batches)

	if jobErr == nil {
		err := r.SetIndexerActualVersion(job.Idx.GetID(), job.To)
		if err != nil {
			return fmt.Errorf("error completing state reindex job %+v: %s", job, err)
		}
	}

	glog.V(2).Infof("Completed state reindex job %+v with job err %+v", job, jobErr)

	duration := clock.Since(start).Seconds()
	metrics.ReindexDuration.WithLabelValues(job.Idx.GetID()).Set(duration)
	glog.Infof("Attempt at state reindex job %+v took %f seconds", job, duration)

	if jobErr == nil {
		// TODO(reginawang3495): refactor out at the end of milestone by using sendUpdate instead of this test hook
		TestHookReindexSuccess()
	}
	if sendUpdate != nil {
		sendUpdate(fmt.Sprintf("indexer %s successfully reindexed from version %d to version %d", job.Idx.GetID(), job.From, job.To))
	}
	return nil
}

// getJobs gets all required reindex jobs.
// If indexer ID is non-empty, only gets job for that indexer.
func (r *reindexer) getJobs(indexerID string) ([]*Job, error) {
	idxs, err := getIndexersFromRegistry(indexerID)

	if err != nil {
		return nil, err
	}

	var ret []*Job
	for _, x := range idxs {
		// Get Indexer Version saved in db
		v, err := GetIndexerVersion(r.Versioner, x.GetID())

		if err != nil {
			return nil, err
		}
		if v == nil {
			return nil, fmt.Errorf("indexer %s version not tracked", x.GetID())
		}
		if x.GetVersion() != v.Actual {
			ret = append(ret, &Job{Idx: x, From: v.Actual, To: x.GetVersion()})
		}
	}

	return ret, nil
}
