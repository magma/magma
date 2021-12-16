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

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/services/state/indexer/metrics"
	state_types "magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/lib/go/util"
)

// TODO(#10295): Remove "Desired" from Indexer.Versions once reindexer_queue and queue are removed

// reindexerSingleton runs as though it is a singleton
type reindexerSingleton struct {
	Versioner
	store Store
}

const reindexLoopInterval = time.Minute

func NewReindexerSingleton(store Store, versioner Versioner) Reindexer {
	return &reindexerSingleton{store: store, Versioner: versioner}
}

func (r *reindexerSingleton) Run(ctx context.Context) {
	// indexerID being "" means that we are reindexing all indexers
	const indexerID = ""
	// TODO(#10293) will support non-nil sendUpdate from Run after removing queue impl
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

// TODO(#10293): remove RunUnsafe when jobQueue reindexer is removed
func (r *reindexerSingleton) RunUnsafe(ctx context.Context, indexerID string, sendUpdate func(string)) error {
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

func (r *reindexerSingleton) findAndReindexJobs(ctx context.Context, indexerID string, sendUpdate func(string)) error {
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
func (r *reindexerSingleton) getReindexBatches(ctx context.Context) []reindexBatch {
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

func (r *reindexerSingleton) reindexJobs(ctx context.Context, jobs []*Job, batches []reindexBatch, sendUpdate func(string)) error {
	errs := &multierror.Error{}
	for _, j := range jobs {
		err := r.reindexJob(j, ctx, batches, sendUpdate)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs.ErrorOrNil()
}

func (r *reindexerSingleton) reindexJob(job *Job, ctx context.Context, batches []reindexBatch, sendUpdate func(string)) error {
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
		// TODO(#10293): consider refactoring out by using sendUpdate instead of this test hook
		TestHookReindexSuccess()
	}
	if sendUpdate != nil {
		sendUpdate(fmt.Sprintf("indexer %s successfully reindexed from version %d to version %d", job.Idx.GetID(), job.From, job.To))
	}
	return nil
}

// getJobs gets all required reindex jobs.
// If indexer ID is non-empty, only gets job for that indexer.
func (r *reindexerSingleton) getJobs(indexerID string) ([]*Job, error) {
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
