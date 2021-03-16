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
	"database/sql"
	"fmt"
	"sort"
	"time"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/sqorc"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/Masterminds/squirrel"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

const (
	// queueTableName is the name of the SQL table acting as the reindex job queue.
	queueTableName = "reindex_job_queue"

	// versionTableName is the name of the SQL table acting as the source of truth for indexer versions.
	versionTableName = "indexer_versions"

	// Job queue columns
	idCol         = "indexer_id"
	fromCol       = "from_version"
	toCol         = "to_version"
	statusCol     = "status"
	attemptsCol   = "attempts"
	errorCol      = "error"
	lastChangeCol = "last_status_change"

	// Version tracker columns
	idColVersions      = "indexer_id"
	actualColVersions  = "version_actual"
	desiredColVersions = "version_desired"

	// defaultJobTimeout after which reindex jobs are considered failed.
	defaultJobTimeout = 5 * time.Minute
)

var (
	// TestHookGet is an empty hook function for test coordination.
	// This should only be set by test code.
	TestHookGet = func() {}
)

// sqlJobQueue wraps a Postgres table to provide a job queue for state reindex jobs.
//
// sqlJobQueue stores the "actual" versions of state indexers, compares to the desired versions
// denoted in each registered indexer, then creates reindex jobs based on those discrepancies.
//
// Note that the version is only "actual" with quotes since indexers may only be in-transition to
// their denoted version, but, if so, a relevant reindex job will be present in the reindex job queue.
//
// Job queue columns:
//	- indexer_id 			-- ID of indexer needing a reindex
//	- status 				-- available, in_progress, complete, or error
//	- last_status_change	-- Unix time since last update to status
//	- attempts 				-- number of attempts at completing the reindex
//	- error 				-- non-empty string if the reindex job was completed with err
// Indexer versions columns:
//	- indexer_id			-- ID of an indexer
//	- version_desired		-- most-recently-updated desired indexer version
//	- version_actual		-- actual version of the indexer
//
// Notes:
//	- Reindex jobs are assumed to take less than 5 minutes (defaultJobTimeout) to complete. For jobs that
//	  happen to take longer, multiple controller instances may try to complete the job concurrently,
//	  under the assumption that previous jobs failed. Individual indexers should handle this gracefully.
//	  Last writer wins for storing error strings.
//	- As with other SQL usages in magma, multiple concurrent calls to Initialize can cause a race condition in Postgres's
//	  DDL table creation, which will return an error.
//	- Indexer versions (uint32) are stored in Postgres default integer types (int32). While this isn't expected to
//	  be an issue, future updates to this type should consider the possibility of a sufficiently-large version being misinterpreted
//	  by a SQL WHERE clause.
type sqlJobQueue struct {
	maxAttempts uint
	db          *sql.DB
	builder     sqorc.StatementBuilder
}

// NewSQLJobQueue returns a new SQL-backed implementation of an unordered job queue.
// The job queue is safe for use across goroutines and processes.
//
// maxAttempts is the max number of times to attempt reindexing the indexer.
//
// Populating the job queue is an exactly-once operation. We handle this in two parts
//	- Populate <= 1 time
//		- The job queue jobs are written as part of a tx that checks the "stored"
//		  indexer versions, and these stored versions are updated the the "desired" versions during the same tx,
//		  ensuring no more than one controller instance will write to the job queue per code push.
//		- There is a small race condition where multiple callers may both log that they successfully updated the job queue,
//		  but this is inconsequential since the condition (a) requires near-simultaneous calls and (b) actually results in the
//		  exact same jobs being written.
//	- Populate >= 1 time
//		- This work is best suited for a future where we have a message broker in the orc8r,
//		  so for now each controller warning-logs either success or failure to write to the job queue, and manual
//		  inspection of the logs would be required (thankfully, we also have tests to ensure this doesn't happen in the expected case).
//
// Only provides Postgres support due to use of the non-standard "FOR UPDATE SKIP LOCKED" clause.
func NewSQLJobQueue(maxAttempts uint, db *sql.DB, builder sqorc.StatementBuilder) JobQueue {
	return &sqlJobQueue{maxAttempts: maxAttempts, db: db, builder: builder}
}

func (s *sqlJobQueue) Initialize() error {
	err := s.initVersionTable()
	if err != nil {
		return err
	}
	return s.initQueueTable()
}

// PopulateJobs tries to add necessary reindex jobs to the job queue.
//
// The population is performed atomically, so max 1 controller instance will be successful per push with indexer version updates.
// However, verifying that at least 1 controller instance was successful is left to manual inspection of the logs.
//
// The tx spans two tables--reindex_job_queue and indexer_versions. If anything fails during the tx, log the error and assume
// it was due to serializing the update with other controller instances.
func (s *sqlJobQueue) PopulateJobs() (bool, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		jobs, err := s.getNewJobs(tx)
		if err != nil {
			return false, err
		}
		if len(jobs) == 0 {
			glog.Info("All desired and actual indexer versions equal, not populating job queue")
			return false, nil
		}

		err = s.addJobs(tx, jobs)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	ret, err := sqorc.ExecInTx(s.db, &sql.TxOptions{Isolation: sql.LevelSerializable}, nil, txFn)
	if err != nil {
		glog.Warningf("Failed to populate reindex job queue; ignore if another controller instance succeeded: %s", err)
		return false, nil
	}
	updated := ret.(bool)

	if updated {
		// Info-level so can compare that at least one controller writes to the job queue
		glog.Info("Succeeded in updating reindex job queue and overwriting new indexer versions")
	}

	return updated, nil
}

func (s *sqlJobQueue) ClaimAvailableJob() (*Job, error) {
	job, err := s.claimAvailableJob()
	if err == merrors.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	idx, err := indexer.GetIndexer(job.id)
	if err != nil {
		failed := &Job{Idx: idx, From: job.from, To: job.to}
		completeErr := s.CompleteJob(failed, errors.Wrap(err, "error claiming available job"))
		if completeErr != nil {
			glog.Errorf("Error completing job after failing to claim it: %+v", completeErr)
		}
		return nil, err
	}

	claimedJob := &Job{Idx: idx, From: job.from, To: job.to}
	return claimedJob, nil
}

func (s *sqlJobQueue) CompleteJob(job *Job, withErr error) error {
	if job == nil {
		return errors.New("job cannot be nil")
	}

	txFn := func(tx *sql.Tx) (interface{}, error) {
		var errVal string
		var statusVal Status

		switch withErr {
		case nil:
			errVal = ""
			statusVal = StatusComplete
		default:
			errVal = withErr.Error()
			statusVal = StatusAvailable
		}

		// Only update indexer actual versions on successful job completion
		if withErr == nil {
			err := s.setIndexerActualVersionImpl(tx, job.Idx.GetID(), job.To)
			if err != nil {
				return nil, err
			}
		}

		_, err := s.builder.Update(queueTableName).
			Set(statusCol, statusVal).
			Set(errorCol, errVal).
			Set(lastChangeCol, clock.Now().Unix()).
			Where(squirrel.Eq{idCol: job.Idx.GetID()}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrapf(err, "update reindex job status to complete %+v", job)
		}

		return nil, nil
	}

	_, err := sqorc.ExecInTx(s.db, &sql.TxOptions{Isolation: sql.LevelSerializable}, nil, txFn)
	return err
}

func (s *sqlJobQueue) GetJobInfos() (map[string]JobInfo, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		rows, err := s.selectAll().From(queueTableName).RunWith(tx).Query()
		if err != nil {
			return nil, errors.Wrap(err, "select all reindex job infos")
		}
		defer sqorc.CloseRowsLogOnError(rows, "GetJobInfos")

		jobs, err := scanJobs(rows)
		if err != nil {
			return nil, err
		}

		return jobs, nil
	}

	txRet, err := sqorc.ExecInTx(s.db, nil, nil, txFn)
	if err != nil {
		return nil, err
	}
	jobs := txRet.(map[string]*reindexJob)

	infos := map[string]JobInfo{}
	for id, job := range jobs {
		infos[id] = JobInfo{IndexerID: job.id, Status: job.status, Error: job.getError(s.maxAttempts), Attempts: job.attempts}
	}

	return infos, nil
}

func (s *sqlJobQueue) GetIndexerVersions() ([]*indexer.Versions, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		return s.getIndexerVersionsImpl(tx)
	}
	txRet, err := sqorc.ExecInTx(s.db, &sql.TxOptions{Isolation: sql.LevelSerializable}, nil, txFn)
	if err != nil {
		return nil, err
	}
	ret := txRet.([]*indexer.Versions)
	return ret, nil
}

func (s *sqlJobQueue) SetIndexerActualVersion(indexerID string, version indexer.Version) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		return nil, s.setIndexerActualVersionImpl(tx, indexerID, version)
	}
	_, err := sqorc.ExecInTx(s.db, &sql.TxOptions{Isolation: sql.LevelSerializable}, nil, txFn)
	return err
}

func (s *sqlJobQueue) initQueueTable() error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := s.builder.CreateTable(queueTableName).
			IfNotExists().
			Column(idCol).Type(sqorc.ColumnTypeText).NotNull().PrimaryKey().EndColumn().
			Column(fromCol).Type(sqorc.ColumnTypeInt).NotNull().EndColumn().
			Column(toCol).Type(sqorc.ColumnTypeInt).NotNull().EndColumn().
			Column(statusCol).Type(sqorc.ColumnTypeText).Default(fmt.Sprintf("'%s'", StatusAvailable)).NotNull().EndColumn().
			Column(lastChangeCol).Type(sqorc.ColumnTypeInt).NotNull().EndColumn().
			Column(attemptsCol).Type(sqorc.ColumnTypeInt).Default(0).NotNull().EndColumn().
			Column(errorCol).Type(sqorc.ColumnTypeText).Default("''").NotNull().EndColumn().
			RunWith(tx).
			Exec()
		return nil, errors.Wrap(err, "initialize reindex job queue table")
	}
	_, err := sqorc.ExecInTx(s.db, &sql.TxOptions{Isolation: sql.LevelRepeatableRead}, nil, txFn)
	return err
}

func (s *sqlJobQueue) initVersionTable() error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := s.builder.CreateTable(versionTableName).
			IfNotExists().
			Column(idColVersions).Type(sqorc.ColumnTypeText).NotNull().PrimaryKey().EndColumn().
			Column(actualColVersions).Type(sqorc.ColumnTypeInt).Default(0).NotNull().EndColumn().
			Column(desiredColVersions).Type(sqorc.ColumnTypeInt).NotNull().EndColumn().
			RunWith(tx).
			Exec()
		return nil, errors.Wrap(err, "initialize indexer versions table")
	}
	_, err := sqorc.ExecInTx(s.db, &sql.TxOptions{Isolation: sql.LevelRepeatableRead}, nil, txFn)
	return err
}

// addJobs adds reindex jobs to the table.
func (s *sqlJobQueue) addJobs(tx *sql.Tx, newJobs []*reindexJob) error {
	jobsToInsert, err := s.getComposedJobs(tx, newJobs)
	if err != nil {
		return err
	}

	_, err = s.builder.Delete(queueTableName).RunWith(tx).Exec()
	if err != nil {
		return errors.Wrap(err, "add reindex jobs, delete existing table contents")
	}

	builder := s.builder.Insert(queueTableName).Columns(idCol, fromCol, toCol, lastChangeCol)
	for _, job := range jobsToInsert {
		builder = builder.Values(job.id, job.from, job.to, clock.Now().Unix())
	}

	_, err = builder.RunWith(tx).Exec()
	if err != nil {
		return errors.Wrapf(err, "add reindex jobs, insert new jobs %+v", jobsToInsert)
	}

	return nil
}

// getNewJobs returns slice of reindex jobs to run.
// Includes jobs where desired != actual.
func (s *sqlJobQueue) getNewJobs(tx *sql.Tx) ([]*reindexJob, error) {
	versions, err := s.getIndexerVersionsImpl(tx)
	if err != nil {
		return nil, err
	}
	var jobs []*reindexJob
	for _, v := range versions {
		if v.Desired == v.Actual {
			continue
		}
		jobs = append(jobs, &reindexJob{id: v.IndexerID, from: v.Actual, to: v.Desired})
	}
	return jobs, nil
}

// Venn diagram of indexer IDs in old and new jobs
//	- old_only:	indexer ID only present in old jobs -- existing job incomplete, but no new job needed
//	- new_only:	indexer ID only present in new jobs -- existing job not found, and new job needed
//	- both:		indexer ID present in both old and new jobs -- existing job incomplete, and new job also needed
// {    old_only    [    both    }    new_only    ]
func (s *sqlJobQueue) getComposedJobs(tx *sql.Tx, newJobs []*reindexJob) ([]*reindexJob, error) {
	oldJobs, err := s.getExistingIncompleteJobs(tx)
	if err != nil {
		return nil, err
	}

	insertJobs := map[string]*reindexJob{}

	// Include all new jobs -- on conflict with an old job, log and replace with new job
	for _, job := range newJobs {
		if prev, exists := oldJobs[job.id]; exists && !job.isSameVersions(prev) {
			glog.Warningf("Replacing existing reindex job %+v with %+v", prev, job)
		}
		insertJobs[job.id] = job
	}

	// Include remaining old jobs -- add back all incomplete old jobs that haven't been superseded by a new job
	for _, job := range oldJobs {
		if _, hasNewerVersion := insertJobs[job.id]; !hasNewerVersion {
			insertJobs[job.id] = job
		}
	}

	var ret []*reindexJob
	for _, j := range insertJobs {
		ret = append(ret, j)
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i].id < ret[j].id }) // for deterministic testing

	return ret, nil
}

// getExistingIncompleteJobs returns a map {indexer ID -> reindex job} for incomplete jobs currently stored in the table.
func (s *sqlJobQueue) getExistingIncompleteJobs(tx *sql.Tx) (map[string]*reindexJob, error) {
	rows, err := s.selectAll().
		From(queueTableName).
		Where(squirrel.NotEq{statusCol: StatusComplete}).
		RunWith(tx).
		Query()
	if err != nil {
		return nil, errors.Wrap(err, "select existing incomplete reindex jobs")
	}
	defer sqorc.CloseRowsLogOnError(rows, "getExistingIncompleteJobs")

	jobs, err := scanJobs(rows)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

// If no job available, returns ErrNotFound from magma/orc8r/lib/go/errors.
func (s *sqlJobQueue) claimAvailableJob() (*reindexJob, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		now := clock.Now()
		timeoutThreshold := now.Add(-defaultJobTimeout)

		rows, err := s.selectAll().
			From(queueTableName).
			Where(
				squirrel.And{
					// Hasn't been claimed/attempted too many times
					squirrel.Lt{attemptsCol: s.maxAttempts},
					// Is available
					squirrel.Or{
						// Normal case: job is available
						squirrel.Eq{statusCol: StatusAvailable},
						// Timeout case: claim job that has been "executing" for too long
						squirrel.And{
							squirrel.Eq{statusCol: StatusInProgress},
							squirrel.Lt{lastChangeCol: timeoutThreshold.Unix()},
						},
					},
				},
			).
			Limit(1).
			Suffix("FOR UPDATE SKIP LOCKED").
			RunWith(tx).
			Query()

		if err != nil {
			return nil, errors.Wrap(err, "claim available reindex job, select available job")
		}
		defer sqorc.CloseRowsLogOnError(rows, "ClaimAvailableJob")

		job, err := scanJob(rows)
		if err != nil {
			return nil, err
		}

		// Set job status to in_progress
		_, err = s.builder.Update(queueTableName).
			Set(statusCol, StatusInProgress).
			Set(attemptsCol, job.attempts+1).
			Set(lastChangeCol, now.Unix()).
			Where(squirrel.Eq{idCol: job.id}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrapf(err, "claim available reindex job, update job status for %+v", job)
		}

		return job, nil
	}

	ret, err := sqorc.ExecInTx(s.db, nil, nil, txFn)
	if err != nil {
		return nil, err
	}
	job := ret.(*reindexJob)

	return job, nil
}

func (s *sqlJobQueue) selectAll() squirrel.SelectBuilder {
	return s.builder.Select(idCol, fromCol, toCol, statusCol, attemptsCol, errorCol, lastChangeCol)
}

func (s *sqlJobQueue) getIndexerVersionsImpl(tx *sql.Tx) ([]*indexer.Versions, error) {
	old, err := s.getTrackedVersions(tx)
	if err != nil {
		return nil, err
	}

	composed, err := getComposedVersions(old)
	if err != nil {
		return nil, err
	}
	if EqualVersions(composed, old) {
		return composed, nil
	}

	// Test hook after first db call so the tx has "officially" started by acquiring some locks
	TestHookGet()

	err = s.overwriteAllVersions(tx, composed)
	if err != nil {
		return nil, err
	}

	return composed, nil
}

func (s *sqlJobQueue) setIndexerActualVersionImpl(tx *sql.Tx, indexerID string, version indexer.Version) error {
	_, err := s.builder.Update(versionTableName).
		Set(actualColVersions, version).
		Where(squirrel.Eq{idColVersions: indexerID}).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrapf(err, "update indexer actual version for %s to %d", indexerID, version)
	}
	return nil
}

func (s *sqlJobQueue) getTrackedVersions(tx *sql.Tx) ([]*indexer.Versions, error) {
	var ret []*indexer.Versions

	rows, err := s.builder.Select(idColVersions, actualColVersions, desiredColVersions).
		From(versionTableName).
		RunWith(tx).
		Query()
	if err != nil {
		return nil, errors.Wrap(err, "get all indexer versions, select existing versions")
	}

	defer sqorc.CloseRowsLogOnError(rows, "GetAllIndexerVersions")

	var idVal string
	var actualVal, desiredVal int64 // int64 is driver's default int type, though these cols are actually int32 storing a uint32
	for rows.Next() {
		err = rows.Scan(&idVal, &actualVal, &desiredVal)
		if err != nil {
			return ret, errors.Wrap(err, "get all indexer versions, SQL row scan error")
		}
		v, err := newVersions(idVal, actualVal, desiredVal)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v)
	}

	err = rows.Err()
	if err != nil {
		return ret, errors.Wrap(err, "get all indexer versions, SQL rows error")
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i].IndexerID < ret[j].IndexerID }) // make deterministic
	return ret, nil
}

func (s *sqlJobQueue) overwriteAllVersions(tx *sql.Tx, versions []*indexer.Versions) error {
	_, err := s.builder.Delete(versionTableName).RunWith(tx).Exec()
	if err != nil {
		return errors.Wrap(err, "overwrite all indexer versions, delete existing versions")
	}

	if len(versions) == 0 {
		return nil
	}

	builder := s.builder.Insert(versionTableName).Columns(idColVersions, actualColVersions, desiredColVersions)
	for _, v := range versions {
		builder = builder.Values(v.IndexerID, v.Actual, v.Desired)
	}
	_, err = builder.RunWith(tx).Exec()
	if err != nil {
		return errors.Wrapf(err, "overwrite all indexer desired versions, insert new versions %+v", versions)
	}

	return nil
}

// If no job available, returns ErrNotFound from magma/orc8r/lib/go/errors.
func scanJob(rows *sql.Rows) (*reindexJob, error) {
	jobs, err := scanJobs(rows)
	if err != nil {
		return nil, err
	}
	if len(jobs) == 0 {
		return nil, merrors.ErrNotFound
	}

	// Return one job
	for _, job := range jobs {
		return job, nil
	}
	return nil, errors.New("err returning one job") // to appease compiler
}

// Returns map of indexer ID to reindex job.
func scanJobs(rows *sql.Rows) (map[string]*reindexJob, error) {
	var err error
	jobs := map[string]*reindexJob{}

	for rows.Next() {
		job := &reindexJob{}
		var lastChangeVal int64
		err = rows.Scan(&job.id, &job.from, &job.to, &job.status, &job.attempts, &job.error, &lastChangeVal)
		if err != nil {
			return nil, errors.Wrap(err, "scan reindex job, SQL row scan error")
		}
		job.lastChange = time.Unix(lastChangeVal, 0)

		jobs[job.id] = job
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "scan reindex job, SQL rows error")
	}

	return jobs, nil
}

// getComposedVersions writes the composition of tracked (old) and local (new) indexers to store.
// Determining whether an indexer needs to be reindexed depends on three recorded version infos per indexer:
//	- new_desired	-- desired version from indexer registry
//	- old_desired	-- desired version from existing reindex jobs
//	- actual		-- actual version updated upon successful reindex job completion
func getComposedVersions(old []*indexer.Versions) ([]*indexer.Versions, error) {
	newv, err := getIndexerVersionsByID()
	if err != nil {
		return nil, err
	}
	composed := map[string]*indexer.Versions{}

	// Insert all old versions -- old_desired and actual values
	for _, v := range old {
		composed[v.IndexerID] = v
	}

	// Insert all new versions -- new_desired overwrite any existing old_desired
	for id, newDesired := range newv {
		if _, present := composed[id]; present {
			composed[id].Desired = newDesired
		} else {
			composed[id] = &indexer.Versions{IndexerID: id, Actual: 0, Desired: newDesired}
		}
	}

	ret := funk.Map(composed, func(k string, v *indexer.Versions) *indexer.Versions { return v }).([]*indexer.Versions)
	sort.Slice(ret, func(i, j int) bool { return ret[i].IndexerID < ret[j].IndexerID }) // make deterministic
	return ret, nil
}

// getIndexerVersionsByID returns a map of registered indexer IDs to their registered ("desired") versions.
func getIndexerVersionsByID() (map[string]indexer.Version, error) {
	indexers, err := indexer.GetIndexers()
	if err != nil {
		return nil, err
	}
	ret := map[string]indexer.Version{}
	for _, x := range indexers {
		ret[x.GetID()] = x.GetVersion()
	}
	return ret, nil
}

// newVersions returns a new indexer versions view.
// First checks the indexer versions fit in an indexer.Version.
func newVersions(indexerID string, actualVersion, desiredVersion int64) (*indexer.Versions, error) {
	actual, err := indexer.NewIndexerVersion(actualVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "new actual version for indexer %s", indexerID)
	}
	desired, err := indexer.NewIndexerVersion(desiredVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "new desired version for indexer %s", indexerID)
	}
	v := &indexer.Versions{
		IndexerID: indexerID,
		Actual:    actual,
		Desired:   desired,
	}
	return v, nil
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
