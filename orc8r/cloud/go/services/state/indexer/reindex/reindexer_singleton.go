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

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/metrics"
	state_types "magma/orc8r/cloud/go/services/state/types"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/util"

	"github.com/golang/glog"
)

// This Reindexer runs as though it is a singleton
type reindexerSingleton struct {
	versioner Versioner
	store     Store
}

func NewReindexerSingleton(store Store, versioner Versioner) Reindexer {
	return &reindexerSingleton{store: store, versioner: versioner}
}

func (r *reindexerSingleton) Run(ctx context.Context) {
	glog.V(2).Infof("called once")
	batches := r.getReindexBatches(ctx)
	for {
		if isCanceled(ctx) {
			glog.Warning("State reindexing async job canceled")
			return
		}
		err := r.reindexJobs(ctx, batches)

		clock.Sleep(failedJobSleep)

		if err == nil {
			continue
		}

		if err == merrors.ErrNotFound {
			glog.V(2).Infof("Failed to claim state reindex job from queue: %s", err)
		} else {
			glog.Errorf("Failed to get or complete state reindex job from queue: %s", err)
		}
		glog.V(2).Infof("sleep a minute")

		clock.Sleep(failedJobSleep)

	}
	r.RunUnsafe(ctx, "", nil)
}

// If no job available, returns ErrNotFound from magma/orc8r/lib/go/errors.
func (r *reindexerSingleton) reindexJobs(ctx context.Context, batches []reindexBatch) error {
	indexerID := ""
	glog.Infof("called once")


	jobs, err := r.getJobs(indexerID)
	glog.Infof("Reindex for indexer '%s' with reindex jobs: %+v, err: %s", indexerID, jobs, err)
	if err != nil || len(jobs) == 0 {
		return err
	}
	if err != nil {
		return wrap(err, ErrDefault, "")
	}
	if jobs == nil {
		return merrors.ErrNotFound
	}


	for _, j := range jobs {
		err = r.reindexJob(j, indexerID, ctx, batches, nil)
		if err != nil {
			return wrap(err, ErrDefault, indexerID)
		}
	}
	return nil
}

func (r * reindexerSingleton) reindexJob(job *Job, indexerID string, ctx context.Context, batches []reindexBatch, sendUpdate func(string)) error {
	defer TestHookReindexDone()
	glog.Infof("running job: %s", job)

	start := clock.Now()

	// glog.Infof("Reindex for indexer '%s', execute job %+v", indexerID, job)
	jobErr := executeJob(ctx, job, batches)

	// TODO add indexer fail count increase
	err := r.versioner.SetIndexerActualVersion(job.Idx.GetID(), job.To)
	if err != nil {
		return fmt.Errorf("error completing state reindex job %+v with job err <%s>: %s", job, jobErr, err)
	}
	glog.V(2).Infof("Completed state reindex job %+v with job err %+v", job, jobErr)
	duration := clock.Since(start).Seconds()
	metrics.ReindexDuration.WithLabelValues(job.Idx.GetID()).Set(duration)
	glog.Infof("Attempt at state reindex job %+v took %f seconds", job, duration)

	TestHookReindexSuccess()
	if sendUpdate != nil {
		sendUpdate(fmt.Sprintf("indexer %s successfully reindexed from version %d to version %d", job.Idx.GetID(), job.From, job.To))
	}
	return nil
}

// TODO cleanup
func (r *reindexerSingleton) RunUnsafe(ctx context.Context, indexerID string, sendUpdate func(string)) error {
	batches := r.getReindexBatches(ctx)
	glog.Infof("Reindex for indexer '%s' with state batches: %+v", indexerID, batches)
	jobs, err := r.getJobs(indexerID)
	if err != nil || len(jobs) == 0 {
		return err
	}
	glog.Infof("Reindex for indexer '%s' with reindex jobs: %+v", indexerID, jobs)

	for _, j := range jobs {
		err = r.reindexJob(j, indexerID, ctx, batches, sendUpdate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *reindexerSingleton) GetIndexerVersions() ([]*indexer.Versions, error) {
	return r.versioner.GetIndexerVersions()
}

// getJobs gets all required reindex jobs.
// If indexer ID is non-empty, only gets job for that indexer.
func (r *reindexerSingleton) getJobs(indexerID string) ([]*Job, error) {
	idxs, err := getIndexers(indexerID)

	if err != nil {
		return nil, err
	}

	var ret []*Job
	for _, x := range idxs {
		v, err := GetIndexerVersion(r.versioner, x.GetID())
		glog.Infof("indexerVersion %s, %s", v, err)

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

// Get network-segregated reindex batches with capped number of state IDs per batch.
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
