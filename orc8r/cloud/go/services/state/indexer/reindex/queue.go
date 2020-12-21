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
	DefaultMaxAttempts uint = 3
)

// Job required to carry out a reindex job.
type Job struct {
	Idx  indexer.Indexer
	From indexer.Version
	To   indexer.Version
}

// JobInfo provides information about a job's progress.
type JobInfo struct {
	IndexerID string
	Status    Status
	Error     string
	Attempts  uint
}

// JobQueue is a static, unordered job queue containing state indexers.
// State indexers are added to the queue for reindexing at initialization, and removed from the queue after the reindex is complete.
// ClaimAvailableJob should be polled periodically, as jobs may become available at any time.
type JobQueue interface {
	// Initialize the queue.
	// Call before other methods.
	Initialize() error

	// PopulateJobs populates the queue with the necessary jobs read from the indexer registry.
	// Returns true if job queue was updated with new jobs.
	PopulateJobs() (bool, error)

	// ClaimAvailableJob claims the returned indexer, to safely perform reindexing operations.
	// Returns nil if no job available.
	ClaimAvailableJob() (*Job, error)

	// CompleteJob indicates completion of the reindexing operation and returns ownership of the job to the queue.
	CompleteJob(job *Job, withErr error) error

	// GetJobInfos provides full information about job progress, keyed by indexer ID.
	// A job info only includes an error when its job has been attempted at least the max number of attempts.
	GetJobInfos() (map[string]JobInfo, error)

	// GetIndexerVersions returns version info for all tracked indexers, keyed by indexer ID.
	// Intended for use when automatic reindexing is disabled.
	GetIndexerVersions() ([]*indexer.Versions, error)

	// SetIndexerActualVersion sets the actual version of an indexer, post-reindex.
	// Intended for use when automatic reindexing is disabled.
	SetIndexerActualVersion(indexerID string, actual indexer.Version) error
}

// GetError returns the job error for a particular reindex job.
// A job only returns an error when its been attempted at least the max number of attempts.
func GetError(queue JobQueue, indexerID string) (string, error) {
	errs, err := GetErrors(queue)
	if err != nil {
		return "", err
	}
	errVal, ok := errs[indexerID]
	if !ok {
		return "", nil
	}
	return errVal, nil
}

// GetErrors returns all job errors, keyed by indexer ID.
// A job only returns an error when its been attempted at least the max number of attempts.
func GetErrors(queue JobQueue) (map[string]string, error) {
	infos, err := queue.GetJobInfos()
	if err != nil {
		return nil, err
	}
	jobErrsByID := map[string]string{}
	for id, info := range infos {
		if info.Error != "" {
			jobErrsByID[id] = info.Error
		}
	}
	return jobErrsByID, nil
}

// GetStatus returns the job status of the job for a particular reindex job.
func GetStatus(queue JobQueue, indexerID string) (Status, error) {
	statuses, err := GetStatuses(queue)
	if err != nil {
		return "", err
	}
	status, ok := statuses[indexerID]
	if !ok {
		return "", nil
	}
	return status, nil
}

// GetStatuses returns all job statuses, keyed by indexer ID.
func GetStatuses(queue JobQueue) (map[string]Status, error) {
	infos, err := queue.GetJobInfos()
	if err != nil {
		return nil, err
	}
	statuses := map[string]Status{}
	for id, info := range infos {
		statuses[id] = info.Status
	}
	return statuses, nil
}

// GetIndexerVersion gets the tracked indexer versions for an indexer ID.
// Returns nil if not found.
func GetIndexerVersion(queue JobQueue, indexerID string) (*indexer.Versions, error) {
	versions, err := queue.GetIndexerVersions()
	if err != nil {
		return nil, err
	}
	for _, v := range versions {
		if v.IndexerID == indexerID {
			return v, nil
		}
	}
	return nil, nil
}

func (j *Job) String() string {
	return fmt.Sprintf("{id: %s, from: %d, to: %d}", j.Idx.GetID(), j.From, j.To)
}

// EqualVersions returns true iff the slices are equal.
// Assumes the slices are already sorted. Any nil elements results in false.
func EqualVersions(a []*indexer.Versions, b []*indexer.Versions) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] == nil || b[i] == nil || *a[i] != *b[i] {
			return false
		}
	}
	return true
}
