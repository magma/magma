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

	"github.com/golang/glog"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/metrics"
	state_types "magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/lib/go/merrors"
	"magma/orc8r/lib/go/util"
)

type reindexerQueue struct {
	queue JobQueue
	store Store
}

func NewReindexerQueue(queue JobQueue, store Store) Reindexer {
	return &reindexerQueue{queue: queue, store: store}
}

func (r *reindexerQueue) Run(ctx context.Context) {
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

func (r *reindexerQueue) RunUnsafe(ctx context.Context, indexerID string, sendUpdate func(string)) error {
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

func (r *reindexerQueue) GetIndexerVersions() ([]*indexer.Versions, error) {
	return r.queue.GetIndexerVersions()
}

// If no job available, returns ErrNotFound from magma/orc8r/lib/go/merrors.
func (r *reindexerQueue) claimAndReindexOne(ctx context.Context, batches []reindexBatch) error {
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

	err = r.queue.CompleteJob(job, jobErr) // Marking attempt as done, not successful
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
func (r *reindexerQueue) getJobs(indexerID string) ([]*Job, error) {
	idxs, err := getIndexersFromRegistry(indexerID)
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

func (r *reindexerQueue) reportStatusMetrics() {
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

// Get network-segregated reindex batches with capped number of state IDs per batch.
func (r *reindexerQueue) getReindexBatches(ctx context.Context) []reindexBatch {
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
