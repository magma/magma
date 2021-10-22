package reindex

import (
	"context"
	"fmt"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/services/state/indexer"
	state_types "magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/lib/go/util"

	"github.com/golang/glog"
)

// This Reindexer runs as though it is a singleton
type reindexerSingleton struct {
	versioner Versioner
	store Store
}

func NewReindexerSingleton(store Store, versioner Versioner) Reindexer {
	return &reindexerSingleton{store: store, versioner: versioner}
}

func (r *reindexerSingleton) Run(ctx context.Context) {
	r.RunUnsafe(ctx, "", nil)
}

func (r *reindexerSingleton) RunUnsafe(ctx context.Context, indexerID string, sendUpdate func(string)) error {
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
		err = r.versioner.SetIndexerActualVersion(j.Idx.GetID(), j.To)
		if err != nil {
			return err
		}
		if sendUpdate != nil {
			sendUpdate(fmt.Sprintf("indexer %s successfully reindexed from version %d to version %d", j.Idx.GetID(), j.From, j.To))
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


