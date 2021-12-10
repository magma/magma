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
	"time"

	"magma/orc8r/cloud/go/clock"
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
)


// Job required to carry out a reindex job.
type Job struct {
	Idx  indexer.Indexer
	From indexer.Version
	To   indexer.Version
}

func (j *Job) String() string {
	return fmt.Sprintf("{id: %s, from: %d, to: %d}", j.Idx.GetID(), j.From, j.To)
}

// JobInfo provides information about a job's progress.
type JobInfo struct {
	IndexerID string
	Status    Status
	Error     string
	Attempts  uint
}

// reindexJob is the internal representation of a reindex job.
type reindexJob struct {
	// Indexer-relevant
	id   string
	from indexer.Version
	to   indexer.Version
	// Job-relevant
	status     Status
	attempts   uint
	error      string
	lastChange time.Time
}

func (j *reindexJob) String() string {
	return fmt.Sprintf("{id: %s, from: %d, to: %d}", j.id, j.from, j.to)
}

func (j *reindexJob) isSameVersions(job *reindexJob) bool {
	return j.from == job.from && j.to == job.to
}

// getError for the reindex job.
// Only returns err if the reindex job has unsuccessfully passed the passed max
// number of reindex attempts.
func (j *reindexJob) getError(maxAttempts uint) string {
	now := clock.Now()
	timeoutThreshold := now.Add(-defaultJobTimeout)

	tooManyAttempts := j.attempts >= maxAttempts
	stalled := j.status == StatusAvailable ||
		(j.status == StatusInProgress && j.lastChange.Before(timeoutThreshold))

	if tooManyAttempts && stalled {
		return j.error
	}
	return ""
}
