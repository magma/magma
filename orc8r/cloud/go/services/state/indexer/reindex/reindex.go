/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
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
	"magma/orc8r/cloud/go/services/state/servicers"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/util"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

type Error string

const (
	// ErrDefault is the default error reported for reindex failure.
	ErrDefault Error = "state reindex error"
	// ErrReindexPerState indicates a Reindex error occured for specific keys. Not included in any job errors.
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
	testHookReindexComplete = func() {}
)

type reindexBatch struct {
	networkID string
	stateIDs  []state.ID
}

// Run to progressively complete required reindex jobs.
// Periodically polls the reindex job queue for reindex jobs, attempts to
// complete the job, and writes back any encountered errors.
// Returns only upon context cancellation, which can optionally be nil.
func Run(ctx context.Context, queue JobQueue, store servicers.StateServiceInternal) {
	batches := getReindexBatches(store)
	for {
		if isCanceled(ctx) {
			glog.Warning("State reindexing async job canceled")
			return
		}

		reportStatusMetrics(queue)

		err := reindexOne(queue, batches)
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

// If no job available, returns ErrNotFound from magma/orc8r/lib/go/errors.
func reindexOne(queue JobQueue, batches []reindexBatch) (err error) {
	job, err := queue.ClaimAvailableJob()
	if err != nil {
		return wrap(err, ErrDefault, "")
	}
	if job == nil {
		return merrors.ErrNotFound
	}

	start := clock.Now()
	id := job.Idx.GetID()
	subs := job.Idx.GetSubscriptions()

	defer func() {
		glog.V(2).Infof("Attempting to complete state reindex job %+v with job err %+v", job, err)
		completeErr := queue.CompleteJob(job, err)
		if completeErr != nil {
			glog.Errorf("Failed to complete state reindex job %+v with job err <%s>: %s", job, err, completeErr)
		}

		duration := clock.Since(start).Seconds()
		metrics.ReindexDuration.WithLabelValues(id).Set(duration)
		glog.Infof("Attempt at state reindex job %+v took %f seconds", job, duration)

		testHookReindexComplete()
	}()

	isFirst := job.From == 0
	err = job.Idx.PrepareReindex(job.From, job.To, isFirst)
	if err != nil {
		return wrap(err, ErrPrepare, id)
	}

	for _, b := range batches {
		ids := filterIDs(subs, b.stateIDs)
		if len(ids) == 0 {
			continue
		}

		// Convert IDs to states -- silently ignore not-found (stale) state IDs
		statesByID, err := state.GetStates(b.networkID, ids)
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

func reportStatusMetrics(queue JobQueue) {
	infos, err := queue.GetAllJobInfo()
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
			glog.Errorf("Report reindex metrics failed to get indexer %s from registry: %v", id, err)
			continue
		}
		metrics.IndexerVersion.WithLabelValues(id).Set(float64(idx.GetVersion()))
	}
}

// Get network-segregated reindex batches with capped number of state IDs per batch.
func getReindexBatches(store servicers.StateServiceInternal) []reindexBatch {
	var idsByNetwork state.IDsByNetwork
	for {
		ids, err := store.GetAllIDs()
		if err == nil {
			idsByNetwork = ids
			break
		}
		err = wrap(err, ErrDefault, "")
		glog.Errorf("Failed to get all state IDs for state indexer reindexing, will retry: %v", err)
		clock.Sleep(failedGetBatchesSleep)
	}

	var current, rest []state.ID
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

func isCanceled(ctx context.Context) bool {
	if ctx == nil {
		return false
	}

	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// filterIDs to the subset that match at least one subscription.
func filterIDs(subs []indexer.Subscription, ids []state.ID) []state.ID {
	var ret []state.ID
	for _, id := range ids {
		for _, s := range subs {
			if s.Match(id) {
				ret = append(ret, id)
				break
			}
		}
	}
	return ret
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
