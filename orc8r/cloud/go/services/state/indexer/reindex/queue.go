/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package reindex

import (
	"fmt"

	"magma/orc8r/cloud/go/services/state/indexer"
)

// Status of a reindex job.
type Status string

const (
	// StatusAvailable indicates the job is not in progress and has not been successfully completed.
	// Available jobs will be considered "errored" when attempts >= max attempts.
	StatusAvailable Status = "available"
	// StatusInProgress indicates a job is being processed.
	// In-progress jobs can be claimed by a new caller if last_status_change is more than defaultTimeout ago.
	StatusInProgress Status = "in_progress"
	// StatusComplete indicates a job has been completed successfully.
	StatusComplete Status = "complete"

	// DefaultMaxAttempts is the default max number of attempts at a reindex job before it's considered failed.
	DefaultMaxAttempts = 3
)

// TODO(hcgatewood): remove this later in the diff stack
type Version struct {
	IndexerID string
	Actual    indexer.Version
	Desired   indexer.Version
}

// Job required to carry out a reindex job.
type Job struct {
	Idx  indexer.Indexer
	From indexer.Version
	To   indexer.Version
}

// JobInfo provides information about a job's progress.
type JobInfo struct {
	Status   Status
	Attempts uint
}

// JobQueue is a static, unordered job queue containing state indexers.
// State indexers are added to the queue for reindexing at initialization, and removed from the queue after the reindex is complete.
// ClaimAvailableJob should be polled periodically, as jobs may become available at any time.
type JobQueue interface {
	// Initialize the queue. Call before other methods.
	Initialize() error

	// PopulateJobs populates the queue with the necessary jobs read from the indexer registry.
	// Returns true if job queue was updated with new jobs.
	PopulateJobs() (bool, error)

	// ClaimAvailableJob claims the returned indexer, to safely perform reindexing operations.
	// Returns nil if no job available.
	ClaimAvailableJob() (*Job, error)

	// CompleteJob indicates completion of the reindexing operation and returns ownership of the job to the queue.
	CompleteJob(job *Job, withErr error) error

	// GetAllErrors returns all job errors, keyed by indexer ID.
	// A job only returns an error when its been attempted at least the max number of attempts.
	GetAllErrors() (map[string]string, error)

	// GetAllJobInfo provides information about job progress, keyed by indexer ID.
	GetAllJobInfo() (map[string]JobInfo, error)
}

func GetError(queue JobQueue, indexerID string) (string, error) {
	errs, err := queue.GetAllErrors()
	if err != nil {
		return "", err
	}
	errVal, ok := errs[indexerID]
	if !ok {
		return "", nil
	}
	return errVal, nil
}

// GetStatus returns the job status of the job for a particular reindex job.
func GetStatus(queue JobQueue, indexerID string) (Status, error) {
	statuses, err := GetAllStatuses(queue)
	if err != nil {
		return "", err
	}
	status, ok := statuses[indexerID]
	if !ok {
		return "", nil
	}
	return status, nil
}

// GetAllStatuses returns all job statuses, keyed by indexer ID.
func GetAllStatuses(queue JobQueue) (map[string]Status, error) {
	infos, err := queue.GetAllJobInfo()
	if err != nil {
		return nil, err
	}
	statuses := map[string]Status{}
	for id, info := range infos {
		statuses[id] = info.Status
	}
	return statuses, nil
}

func (j *Job) String() string {
	return fmt.Sprintf("{id: %s, from: %d, to: %d}", j.Idx.GetID(), j.From, j.To)
}
