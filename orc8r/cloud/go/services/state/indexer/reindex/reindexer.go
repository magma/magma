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

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/metrics"
	state_types "magma/orc8r/cloud/go/services/state/types"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/util"

	"github.com/golang/glog"
	"github.com/pkg/errors"
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

type reindexerImpl struct {
	queue JobQueue
	store Store
}

func NewReindexer(queue JobQueue, store Store) Reindexer {
	return &reindexerImpl{queue: queue, store: store}
}

func (r *reindexerImpl) Run(ctx context.Context) {
	batches := r.getReindexBatches(ctx)
	for {
		if isCanceled(ctx) {
			glog.Warning("State reindexing async job canceled")
			return
		}

		r.reportStatusMetrics()

		err := r.claimAndReindexOne(ctx, batches)
		if err == nil {
			continue
		}

		if err == merrors.ErrNotFound {
			glog.V(2).Infof("Failed to claim state reindex job from queue: %s", err)
		} else {
			glog.Errorf("Failed to get or complete state reindex job from queue: %s", err)
		}
		clock.Sleep(failedJobSleep)
	}
}

func (r *reindexerImpl) RunUnsafe(ctx context.Context, indexerID string, sendUpdate func(string)) error {
	batches := r.getReindexBatches(ctx)
	glog.Infof("Reindex for indexer '%s' with state batches: %+v", indexerID, batches)
	jobs, err := r.getJobs(indexerID)
	if err != nil || len(jobs) == 0 {
		return err
	}
	glog.Infof("Reindex for indexer '%s' with reindex jobs: %+v", indexerID, jobs)

	for _, j := range jobs {
		glog.Infof("Reindex for indexer '%s', execute job %+v", indexerID, j)
		err = executeJob(ctx, j, batches)
		if err != nil {
			return err
		}
		err = r.queue.SetIndexerActualVersion(j.Idx.GetID(), j.To)
		if err != nil {
			return err
		}
		if sendUpdate != nil {
			sendUpdate(fmt.Sprintf("indexer %s successfully reindexed from version %d to version %d", j.Idx.GetID(), j.From, j.To))
		}
	}
	return nil
}

func (r *reindexerImpl) GetIndexerVersions() ([]*indexer.Versions, error) {
	return r.queue.GetIndexerVersions()
}

// If no job available, returns ErrNotFound from magma/orc8r/lib/go/errors.
func (r *reindexerImpl) claimAndReindexOne(ctx context.Context, batches []reindexBatch) error {
	defer TestHookReindexDone()

	job, err := r.queue.ClaimAvailableJob()
	if err != nil {
		return wrap(err, ErrDefault, "")
	}
	if job == nil {
		return merrors.ErrNotFound
	}
	start := clock.Now()

	jobErr := executeJob(ctx, job, batches)

	err = r.queue.CompleteJob(job, jobErr)
	if err != nil {
		return fmt.Errorf("error completing state reindex job %+v with job err <%s>: %s", job, jobErr, err)
	}
	glog.V(2).Infof("Completed state reindex job %+v with job err %+v", job, jobErr)

	duration := clock.Since(start).Seconds()
	metrics.ReindexDuration.WithLabelValues(job.Idx.GetID()).Set(duration)
	glog.Infof("Attempt at state reindex job %+v took %f seconds", job, duration)

	TestHookReindexSuccess()
	return nil
}

// getJobs gets all required reindex jobs.
// If indexer ID is non-empty, only gets job for that indexer.
func (r *reindexerImpl) getJobs(indexerID string) ([]*Job, error) {
	idxs, err := getIndexers(indexerID)
	if err != nil {
		return nil, err
	}

	var ret []*Job
	for _, x := range idxs {
		v, err := GetIndexerVersion(r.queue, x.GetID())
		if err != nil {
			return nil, err
		}
		if v == nil {
			return nil, fmt.Errorf("indexer %s version not tracked", x.GetID())
		}
		if v.Actual != v.Desired {
			ret = append(ret, &Job{Idx: x, From: v.Actual, To: v.Desired})
		}
	}

	return ret, nil
}

func (r *reindexerImpl) reportStatusMetrics() {
	infos, err := r.queue.GetJobInfos()
	if err != nil {
		err = wrap(err, ErrDefault, "")
		glog.Errorf("Report reindex metrics failed to get all reindex job info: %v", err)
		return
	}

	for id, info := range infos {
		metrics.ReindexStatus.WithLabelValues(id).Set(getStatus(info.Status))
		metrics.ReindexAttempts.WithLabelValues(id).Set(float64(info.Attempts))
		if info.Status != StatusComplete {
			continue
		}
		idx, err := indexer.GetIndexer(id)
		if err != nil {
			glog.Errorf("Report reindex metrics failed to get indexer %s from registry with error: %v", id, err)
			continue
		}
		if idx == nil {
			glog.Errorf("Report reindex metrics failed to get indexer %s from registry", id)
			continue
		}
		metrics.IndexerVersion.WithLabelValues(id).Set(float64(idx.GetVersion()))
	}
}

type reindexBatch struct {
	networkID string
	stateIDs  state_types.IDs
}

// Get network-segregated reindex batches with capped number of state IDs per batch.
func (r *reindexerImpl) getReindexBatches(ctx context.Context) []reindexBatch {
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
		statesByID, err := state.GetSerializedStates(b.networkID, ids)
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

func getIndexers(indexerID string) ([]indexer.Indexer, error) {
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
	return errors.Wrap(err, wrap)
}
