/*
 *  Copyright 2020 The Magma Authors.
 *
 *  This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package storage

import (
	"database/sql"
	"fmt"
	"sort"
	"time"

	"magma/orc8r/cloud/go/sqorc"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type tSmsByPk = map[string]*SMS

func garbageCollectExpiredRefs(tx *sql.Tx, builder sqorc.StatementBuilder, networkID string, imsis []string, timeoutSecs int64) error {
	/*
		DELETE FROM smsd_refs
		WHERE sms_id IN (
			SELECT sms_id FROM smsd_refs
			INNER JOIN smsd_messages on smsd_refs.sms_id = smsd_messages.pk
			WHERE network_id = {nid} AND imsi IN {imsis} AND ref_created_sec < {timeout} AND num_attempts >= {limit}
		)
	*/
	subSelect, selectArgs, _ := builder.Select(refSmsCol).
		From(refsTable).
		JoinClause(fmt.Sprintf("INNER JOIN %s ON %s=%s", smsTable, getFQColName(refsTable, refSmsCol), getFQColName(smsTable, pkCol))).
		Where(sq.And{
			sq.Eq{
				getFQColName(smsTable, nidCol):  networkID,
				getFQColName(smsTable, imsiCol): imsis,
			},
			sq.Lt{getFQColName(refsTable, refCreatedCol): timeoutSecs},
			sq.GtOrEq{getFQColName(smsTable, attemptsCol): maxRetries},
		}).
		ToSql()
	_, err := builder.Delete(refsTable).
		Where(fmt.Sprintf("%s IN (%s)", refSmsCol, subSelect), selectArgs...).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to garbage collect old SMS refs")
	}
	return nil
}

func loadMessagesToSend(tx *sql.Tx, builder sqorc.StatementBuilder, networkID string, imsis []string, timeoutSecs int64) (map[string]tSmsByPk, error) {
	/*
		SELECT * FROM smsd_messages
		LEFT OUTER JOIN smsd_refs ON smsd_messages.pk = smsd_refs.sms_id
		WHERE
			smsd_messages.network_id = {nid}
			AND
			smsd_messages.imsi IN {imsis}
			AND
			(smsd_refs.sms_id is NULL OR smsd_refs.ref_created_sec < {timeout})
			AND
			NOT smsd_messages.is_delivered
			AND
			smsd_messages.num_attempts < 3
	*/
	rows, err := builder.Select(allCols...).
		From(smsTable).
		JoinClause(fmt.Sprintf("LEFT OUTER JOIN %s ON %s=%s", refsTable, getFQColName(smsTable, pkCol), getFQColName(refsTable, refSmsCol))).
		Where(
			sq.And{
				sq.Eq{
					getFQColName(smsTable, nidCol):  networkID,
					getFQColName(smsTable, imsiCol): imsis,
				},
				sq.Or{
					sq.Eq{getFQColName(refsTable, refSmsCol): nil},
					sq.Lt{getFQColName(refsTable, refCreatedCol): timeoutSecs},
				},
				sq.Eq{getFQColName(smsTable, deliveredCol): false},
				sq.Lt{getFQColName(smsTable, attemptsCol): maxRetries},
			},
		).
		RunWith(tx).
		Query()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load SMSs to deliver")
	}
	defer sqorc.CloseRowsLogOnError(rows, "loadMessagesToSend")

	smsByImsi, err := scanMessages(rows)
	if err != nil {
		return nil, err
	}
	return smsByImsi, nil
}

func loadRefsMasks(tx *sql.Tx, builder sqorc.StatementBuilder, networkID string, imsis []string) (map[string]*[256]bool, error) {
	/*
		SELECT smsd_messages.imsi, smsd_refs.ref_num FROM smsd_refs
		INNER JOIN smsd_messages ON smsd_refs.sms_id = smsd_messages.pk
		WHERE smsd_messages.network_id = {nid} AND smsd_messages.imsi IN {imsis}
	*/
	rows, err := builder.Select(getFQColName(smsTable, imsiCol), getFQColName(refsTable, refCol)).
		From(refsTable).
		JoinClause(fmt.Sprintf("INNER JOIN %s ON %s=%s", smsTable, getFQColName(refsTable, refSmsCol), getFQColName(smsTable, pkCol))).
		Where(sq.Eq{
			getFQColName(smsTable, nidCol):  networkID,
			getFQColName(smsTable, imsiCol): imsis,
		}).
		RunWith(tx).
		Query()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load allocated ref nums")
	}
	defer sqorc.CloseRowsLogOnError(rows, "loadRefsMasks")

	ret, err := scanRefs(rows)
	if err != nil {
		return nil, err
	}
	// Populate the return value with an all-false mask for all IMSIs which are
	// not present as keys (i.e. there are no in-flight messages for this IMSI)
	for _, imsi := range imsis {
		if _, imsiFound := ret[imsi]; !imsiFound {
			ret[imsi] = &[256]bool{}
		}
	}
	return ret, nil
}

// Given a set of messages and a reference number mask, returns any new ref
// nums that need to be allocated (keyed by message pk).
// If a message already as assigned ref nums, we will still return the existing
// refs because the creation time on those will need to be advanced.
//
// No guarantee that all messages will be assigned refs - if there are > 256
// messages that need to be sent or are in flight, we will truncate some of the
// input.
//
// No matter what, all messages with ref nums already allocated will be
// returned as part of the output.
func getRefsAndMessagesToEncode(messages tSmsByPk, refsMask *[256]bool, refCounter SMSReferenceCounter) map[string][]byte {
	if funk.IsEmpty(messages) || refsMask == nil {
		return nil
	}

	// Sort map keys for deterministic testing
	sortedPks := funk.Keys(messages).([]string)
	sort.Strings(sortedPks)

	refsToAllocate := map[string][]byte{}
	for _, pk := range sortedPks {
		msg := messages[pk]
		numFreeRefs := getNumFreeRefs(refsMask)

		// Easy case - refs already allocated, this message has timed out.
		// Re-use the existing refs.
		// We still return the refs as something to allocate because we need
		// to advance the persisted creation time for the ref num
		if !funk.IsEmpty(msg.RefNums) {
			refsToAllocate[pk] = msg.RefNums
			continue
		}

		numRefsNeeded := refCounter.GetReferenceNumberCount(msg.Message)
		if numRefsNeeded > numFreeRefs {
			// Not enough available refs to encode this message right now
			continue
		}
		refsToAllocate[pk] = allocateRefs(refsMask, numRefsNeeded)
	}
	return refsToAllocate
}

// Write new refs and increment attempt count for all messages
func persistNewRefNums(tx *sql.Tx, builder sqorc.StatementBuilder, refsByPk map[string][]byte, timeCreated int64) error {
	// INSERT INTO smsd_refs (sms_id, ref_num, ref_created_sec) VALUES ($1, $2, $3)
	// ON CONFLICT (sms_id, ref_num) DO UPDATE SET ref_created_sec = $4
	sc := sq.NewStmtCache(tx)
	defer sqorc.ClearStatementCacheLogOnError(sc, "persistNewRefNums")
	for pk, refs := range refsByPk {
		for _, refNum := range refs {
			_, err := builder.Insert(refsTable).
				Columns(refSmsCol, refCol, refCreatedCol).
				Values(pk, refNum, timeCreated).
				OnConflict(
					[]sqorc.UpsertValue{{Column: refCreatedCol, Value: timeCreated}},
					refSmsCol, refCol,
				).
				RunWith(sc).
				Exec()
			if err != nil {
				return errors.Wrap(err, "failed to create new SMS ref numbers")
			}
		}
	}

	// UPDATE smsd_messages SET num_attempts = num_attempts + 1 WHERE pk IN {pks}
	allPks := funk.Keys(refsByPk).([]string)
	_, err := builder.Update(smsTable).
		Set(attemptsCol, sq.Expr(fmt.Sprintf("%s+1", attemptsCol))).
		Where(sq.Eq{pkCol: allPks}).
		RunWith(sc).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to increment sms attempt counts")
	}

	return nil
}

type imsiAndRef struct {
	imsi string
	ref  byte
}

// This will return a super-set of the refs that we actually care about but
// crafting it in another way is way too complicated.
func loadPksByRefs(tx *sql.Tx, builder sqorc.StatementBuilder, networkID string, imsis []string) (map[imsiAndRef]string, error) {
	/*
		SELECT smsd_messages.imsi, smsd_refs.sms_id, smsd_refs.ref_num FROM smsd_refs
		INNER JOIN smsd_messages ON smsd_messages.pk = smsd_refs.sms_id
		WHERE smsd_messages.network_id = {nid} AND smsd_messages.imsi IN {imsis}
	*/
	rows, err := builder.Select(getFQColName(smsTable, imsiCol), getFQColName(refsTable, refSmsCol), getFQColName(refsTable, refCol)).
		From(refsTable).
		JoinClause(fmt.Sprintf("INNER JOIN %s ON %s=%s", smsTable, getFQColName(smsTable, pkCol), getFQColName(refsTable, refSmsCol))).
		Where(sq.Eq{
			getFQColName(smsTable, nidCol):  networkID,
			getFQColName(smsTable, imsiCol): imsis,
		}).
		RunWith(tx).
		Query()
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for ref-pk mapping")
	}
	defer sqorc.CloseRowsLogOnError(rows, "loadPksByRefs")

	ret := map[imsiAndRef]string{}
	for rows.Next() {
		var imsi, pk string
		var ref int64

		err := rows.Scan(&imsi, &pk, &ref)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan ref-pk mapping")
		}
		ret[imsiAndRef{imsi: imsi, ref: byte(ref)}] = pk
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "sql rows err")
	}
	return ret, nil
}

func markMessagesAsDelivered(tx *sql.Tx, builder sqorc.StatementBuilder, networkID string, deliveredMessages map[string][]SMSRef, pksByRef map[imsiAndRef]string) error {
	// For delivered messages, mark them as such in the table and delete
	// all the refs that have been allocated for them.
	var deliveredPks []string
	for imsi, deliveredRefs := range deliveredMessages {
		for _, ref := range deliveredRefs {
			pk, found := pksByRef[imsiAndRef{imsi: imsi, ref: ref}]
			if found {
				deliveredPks = append(deliveredPks, pk)
				continue
			}
		}
	}
	deliveredPks = funk.UniqString(deliveredPks)

	_, err := builder.Update(smsTable).
		Set(deliveredCol, true).
		Set(errorCol, sql.NullString{Valid: false}).
		Where(sq.Eq{nidCol: networkID, pkCol: deliveredPks}).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to mark SMSs as delivered")
	}

	// Subquery to limit operation to this network
	/*
		DELETE FROM smsd_refs
		WHERE sms_id IN (
			SELECT pk FROM smsd_messages WHERE network_id = {nid} AND pk IN {pks}
		)
	*/
	subSelect, selectArgs, _ := builder.Select(pkCol).
		From(smsTable).
		Where(sq.Eq{nidCol: networkID, pkCol: deliveredPks}).
		ToSql()
	_, err = builder.Delete(refsTable).
		Where(sq.Expr(fmt.Sprintf("%s IN (%s)", refSmsCol, subSelect), selectArgs...)).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to clear refs for delivered messages")
	}

	return nil
}

func processFailedMessages(tx *sql.Tx, builder sqorc.StatementBuilder, networkID string, failedMessages map[string][]SMSFailureReport, pksByRef map[imsiAndRef]string) error {
	sc := sq.NewStmtCache(tx)
	defer sqorc.ClearStatementCacheLogOnError(sc, "processFailedMessages")

	// For failed messages, persist the error message and delete the refs if
	// the message's attempt count is over the retry limit
	var failedPks []string
	for imsi, failureReport := range failedMessages {
		for _, report := range failureReport {
			pk, found := pksByRef[imsiAndRef{imsi: imsi, ref: report.Ref}]
			if !found {
				continue
			}
			// track this to do the ref deletion later
			failedPks = append(failedPks, pk)

			// Set error message
			_, err := builder.Update(smsTable).
				Set(errorCol, sql.NullString{Valid: true, String: report.ErrorMessage}).
				Where(sq.Eq{nidCol: networkID, pkCol: pk}).
				RunWith(sc).
				Exec()
			if err != nil {
				return errors.Wrap(err, "failed to set error message on SMS")
			}
		}
	}
	failedPks = funk.UniqString(failedPks)

	/*
		DELETE FROM smsd_refs
		WHERE sms_id IN (
			SELECT pk FROM smsd_messages
			WHERE network_id = {nid} AND pk IN {pks} AND num_attempts >= 3
		)
	*/
	subSelect, selectArgs, _ := builder.Select(pkCol).
		From(smsTable).
		Where(
			sq.And{
				sq.Eq{nidCol: networkID, pkCol: failedPks},
				sq.GtOrEq{attemptsCol: maxRetries},
			},
		).
		ToSql()
	_, err := builder.Delete(refsTable).
		Where(sq.Expr(fmt.Sprintf("%s IN (%s)", refSmsCol, subSelect), selectArgs...)).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to delete refs for failed messages over retry threshold")
	}

	return nil
}

func scanMessages(rows *sql.Rows) (map[string]tSmsByPk, error) {
	smsByImsi := map[string]tSmsByPk{}
	for rows.Next() {
		var pk, imsi, srcMsisdn, message string
		var errorMessage sql.NullString
		var delivered bool
		var timeCreated, numAttempts int64
		var refNum, refCreated sql.NullInt64

		err := rows.Scan(&pk, &delivered, &imsi, &srcMsisdn, &message, &timeCreated, &errorMessage, &numAttempts, &refNum, &refCreated)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan sms row")
		}

		createdTs, err := ptypes.TimestampProto(time.Unix(timeCreated, 0))
		if err != nil {
			return nil, errors.Wrapf(err, "could not validate created time for sms %s", pk)
		}

		var attemptedTs *timestamp.Timestamp
		if refCreated.Valid {
			attemptedTs, err = ptypes.TimestampProto(time.Unix(refCreated.Int64, 0))
			if err != nil {
				return nil, errors.Wrapf(err, "could not validate attempted time for sms %s", pk)
			}
		}

		status := MessageStatus_WAITING
		switch {
		case delivered:
			status = MessageStatus_DELIVERED
		case numAttempts >= maxRetries:
			status = MessageStatus_FAILED
		}

		var refs []byte
		if refNum.Valid {
			// Technically, this could truncate the refNum value but we ensure
			// that this value will fit into 1 byte in application code
			refs = append(refs, byte(refNum.Int64))
		}

		// If we already scanned this message, we just need to append another
		// ref number for it (because left outer join).
		if _, subMapExists := smsByImsi[imsi]; !subMapExists {
			smsByImsi[imsi] = tSmsByPk{}
		}
		existingSMS, found := smsByImsi[imsi][pk]
		if found {
			existingSMS.RefNums = append(existingSMS.RefNums, refs...)
		} else {
			smsByImsi[imsi][pk] = &SMS{
				Pk:                      pk,
				Status:                  status,
				Imsi:                    imsi,
				SourceMsisdn:            srcMsisdn,
				Message:                 message,
				CreatedTime:             createdTs,
				LastDeliveryAttemptTime: attemptedTs,
				// Technically this could also truncate but then you've got
				// bigger problems
				AttemptCount:  uint32(numAttempts),
				DeliveryError: errorMessage.String,
				RefNums:       refs,
			}
		}
	}
	return smsByImsi, nil
}

// returns masks where true at index i means that ref#i has been assigned
func scanRefs(rows *sql.Rows) (map[string]*[256]bool, error) {
	ret := map[string]*[256]bool{}
	for rows.Next() {
		var imsi string
		var ref int64

		err := rows.Scan(&imsi, &ref)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan ref num")
		}

		if _, imsiArrExists := ret[imsi]; !imsiArrExists {
			ret[imsi] = &[256]bool{}
		}
		ret[imsi][ref] = true
	}
	return ret, nil
}

func getNumFreeRefs(refMask *[256]bool) uint16 {
	// 16 because 256 is a valid value for this
	numFreeRefs := uint16(0)
	for _, isAllocated := range refMask {
		if !isAllocated {
			numFreeRefs += 1
		}
	}
	return numFreeRefs
}

func allocateRefs(refMask *[256]bool, count uint16) []byte {
	ret := make([]byte, 0, count)
	for i, isAllocated := range refMask {
		if !isAllocated {
			ret = append(ret, byte(i))
			refMask[i] = true
		}

		if len(ret) == int(count) {
			break
		}
	}
	return ret
}

func getFQColName(table, col string) string {
	return fmt.Sprintf("%s.%s", table, col)
}
