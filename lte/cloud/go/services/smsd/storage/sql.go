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
	"context"
	"database/sql"
	"fmt"
	"sort"
	"time"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

const (
	smsTable = "smsd_messages"
	imsiIdx  = "smsd_sms_imsi_idx"

	pkCol        = "pk"
	nidCol       = "network_id"
	deliveredCol = "is_delivered"
	imsiCol      = "imsi"
	sourceCol    = "src_msisdn"
	messageCol   = "message"
	createdCol   = "time_created_sec"
	errorCol     = "error_message"
	attemptsCol  = "num_attempts"
	// TODO: save time last sent (delivery response received from AGW)?

	refsTable     = "smsd_refs"
	refSmsCol     = "sms_id"
	refCol        = "ref_num"
	refCreatedCol = "ref_created_sec"
)

const (
	// How many times we'll try to send the same SMS before marking it as failed
	maxRetries = 3

	defaultTimeoutThreshold = 6 * time.Minute
)

var allCols = []string{pkCol, deliveredCol, imsiCol, sourceCol, messageCol, createdCol, errorCol, attemptsCol, refCol, refCreatedCol}

func NewSQLSMSStorage(db *sql.DB, sqlBuilder sqorc.StatementBuilder, counter SMSReferenceCounter, idGenerator storage.IDGenerator) SMSStorage {
	return &sqlSMSStorage{
		db:          db,
		builder:     sqlBuilder,
		counter:     counter,
		idGenerator: idGenerator,
	}
}

type sqlSMSStorage struct {
	db          *sql.DB
	builder     sqorc.StatementBuilder
	counter     SMSReferenceCounter
	idGenerator storage.IDGenerator
}

func (s *sqlSMSStorage) Init() (err error) {
	tx, err := s.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		err = errors.Wrap(err, "table initialization failed")
	}

	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = fmt.Errorf("%s; rollback error: %s", err, rollbackErr)
			}
		}
	}()

	_, err = s.builder.CreateTable(smsTable).
		IfNotExists().
		Column(nidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(pkCol).Type(sqorc.ColumnTypeText).PrimaryKey().EndColumn().
		Column(deliveredCol).Type(sqorc.ColumnTypeBool).NotNull().Default(false).EndColumn().
		Column(imsiCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(sourceCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(messageCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(createdCol).Type(sqorc.ColumnTypeInt).NotNull().EndColumn().
		Column(errorCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(attemptsCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		RunWith(tx).
		Exec()
	if err != nil {
		err = errors.Wrap(err, "table initialization failed")
		return
	}

	// index on (nid, imsi)
	_, err = s.builder.CreateIndex(imsiIdx).
		IfNotExists().
		On(smsTable).
		Columns(nidCol, imsiCol).
		RunWith(tx).
		Exec()
	if err != nil {
		err = errors.Wrap(err, "failed to create sms imsi index")
		return
	}

	_, err = s.builder.CreateTable(refsTable).
		IfNotExists().
		Column(refSmsCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(refCol).Type(sqorc.ColumnTypeInt).NotNull().EndColumn().
		Column(refCreatedCol).Type(sqorc.ColumnTypeInt).NotNull().EndColumn().
		PrimaryKey(refSmsCol, refCol).
		ForeignKey(smsTable, map[string]string{refSmsCol: pkCol}, sqorc.ColumnOnDeleteCascade).
		RunWith(tx).
		Exec()
	if err != nil {
		err = errors.Wrap(err, "failed to create sms ref number table")
		return
	}

	return
}

func (s *sqlSMSStorage) GetSMSs(networkID string, pks []string, imsis []string, onlyWaiting bool, startTime, endTime *time.Time) ([]*SMS, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		/*
			SELECT * FROM smsd_messages
			LEFT JOIN smsd_refs ON smsd_refs.sms_id = smsd_messages.pk
			[[ WHERE (smsd_messages.network_id = {networkID} AND smsd_messages.imsi IN {imsis} AND NOT smsd_messages.is_delivered AND smsd_messages.time_created_sec > ... AND smsd_messages.time_created_sec < ... ]]
		*/
		builder := s.builder.Select(allCols...).
			From(smsTable).
			LeftJoin(fmt.Sprintf("%s ON %s=%s", refsTable, getFQColName(refsTable, refSmsCol), getFQColName(smsTable, pkCol))).
			Where(sq.Eq{getFQColName(smsTable, nidCol): networkID}).
			RunWith(tx)
		if !funk.IsEmpty(pks) {
			builder = builder.Where(sq.Eq{getFQColName(smsTable, pkCol): pks})
		}
		if !funk.IsEmpty(imsis) {
			builder = builder.Where(sq.Eq{getFQColName(smsTable, imsiCol): imsis})
		}
		if onlyWaiting {
			builder = builder.Where(sq.Eq{getFQColName(smsTable, deliveredCol): false})
		}
		if startTime != nil {
			builder = builder.Where(sq.Gt{getFQColName(smsTable, createdCol): startTime.Unix()})
		}
		if endTime != nil {
			builder = builder.Where(sq.Lt{getFQColName(smsTable, createdCol): endTime.Unix()})
		}

		rows, err := builder.Query()
		if err != nil {
			return nil, errors.Wrap(err, "failed to load messages")
		}
		defer sqorc.CloseRowsLogOnError(rows, "GetSMSs")

		return scanMessages(rows)
	}

	retMap, err := sqorc.ExecInTx(s.db, nil, nil, txFn)
	if err != nil {
		return []*SMS{}, err
	}

	retCasted := retMap.(map[string]tSmsByPk)
	var ret []*SMS
	for _, smsByPk := range retCasted {
		for _, msg := range smsByPk {
			ret = append(ret, msg)
		}
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i].Pk < ret[j].Pk })
	return ret, nil
}

func (s *sqlSMSStorage) GetSMSsToDeliver(networkID string, imsis []string, timeoutThreshold time.Duration) ([]*SMS, error) {
	if funk.IsEmpty(imsis) {
		return []*SMS{}, nil
	}
	if timeoutThreshold <= 0 {
		timeoutThreshold = defaultTimeoutThreshold
	}

	// Find all messages for these imsis which don't have refs or have expired
	// refs. In other words, load messages which haven't yet been sent or are
	// in-flight and have expired.
	//
	// Then for each message, allocate enough ref numbers to cover all the SMSs
	// that the message will encode into. Ref nums can only range 0-255 and
	// have to be unique for each IMSI, so we'll have to do another SELECT in
	// order to figure out what ref nums have already been allocated for this
	// set of IMSIs.
	//
	// We will also clear out ref numbers for messages that have exceeded the
	// retry limit but are still considered in-flight, because if we don't ever
	// receive a delivery digest for those messages, their corresponding refs
	// will never be garbage collected.
	txFn := func(tx *sql.Tx) (interface{}, error) {
		timeCreated := clock.Now().Unix()
		timeoutSecs := clock.Now().Add(-timeoutThreshold).Unix()
		updatedTimeCreatedTS, err := ptypes.TimestampProto(time.Unix(timeCreated, 0))
		if err != nil {
			return nil, errors.Wrap(err, "failed to create timestamp")
		}

		err = garbageCollectExpiredRefs(tx, s.builder, networkID, imsis, timeoutSecs)
		if err != nil {
			return nil, err
		}

		smsByImsi, err := loadMessagesToSend(tx, s.builder, networkID, imsis, timeoutSecs)
		if err != nil {
			return nil, err
		}

		refsMasksByImsi, err := loadRefsMasks(tx, s.builder, networkID, imsis)
		if err != nil {
			return nil, err
		}

		// Gather all refs that we need to create and fill messages with the
		// ref nums that will be allocated to them.
		var outputMessages []*SMS
		refsToCreateByPk := map[string][]byte{}
		for _, imsi := range imsis {
			newRefs := getRefsAndMessagesToEncode(smsByImsi[imsi], refsMasksByImsi[imsi], s.counter)

			for pk, refs := range newRefs {
				msg := smsByImsi[imsi][pk]
				// Update these fields in-place because the DB update happens
				// after the read that populates these structs
				msg.RefNums = refs
				msg.LastDeliveryAttemptTime = updatedTimeCreatedTS
				msg.AttemptCount += 1

				refsToCreateByPk[pk] = refs
				outputMessages = append(outputMessages, msg)
			}
		}

		err = persistNewRefNums(tx, s.builder, refsToCreateByPk, timeCreated)
		if err != nil {
			return nil, err
		}

		return outputMessages, nil
	}

	ret, err := sqorc.ExecInTx(s.db, nil, nil, txFn)
	if err != nil {
		return nil, err
	}

	retCasted := ret.([]*SMS)
	sort.Slice(retCasted, func(i, j int) bool { return retCasted[i].Pk < retCasted[j].Pk })
	return retCasted, nil
}

func (s *sqlSMSStorage) CreateSMS(networkID string, sms MutableSMS) (string, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		pk := s.idGenerator.New()
		timeCreated := clock.Now().Unix()

		_, err := s.builder.Insert(smsTable).
			Columns(pkCol, nidCol, imsiCol, sourceCol, messageCol, createdCol).
			Values(pk, networkID, sms.Imsi, sms.SourceMsisdn, sms.Message, timeCreated).
			RunWith(tx).
			Exec()
		if err != nil {
			return "", errors.Wrap(err, "failed to create SMS")
		}
		return pk, nil
	}

	iPK, err := sqorc.ExecInTx(s.db, nil, nil, txFn)
	if err != nil {
		return "", err
	}
	return iPK.(string), nil
}

func (s *sqlSMSStorage) DeleteSMSs(networkID string, pks []string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := s.builder.Delete(smsTable).
			Where(sq.Eq{nidCol: networkID, pkCol: pks}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrap(err, "failed to delete SMSs")
		}
		return nil, nil
	}

	_, err := sqorc.ExecInTx(s.db, nil, nil, txFn)
	return err
}

func (s *sqlSMSStorage) ReportDelivery(networkID string, deliveredMessages map[string][]SMSRef, failedMessages map[string][]SMSFailureReport) error {
	// TODO: what should we do in this case when we get a delivery digest for
	//  a message we don't know about?
	//  We probably don't want to error out the whole call
	txFn := func(tx *sql.Tx) (interface{}, error) {
		var allImsis []string
		for imsi := range deliveredMessages {
			allImsis = append(allImsis, imsi)
		}
		for imsi := range failedMessages {
			allImsis = append(allImsis, imsi)
		}
		allImsis = funk.UniqString(allImsis)

		// We need to map IMSI-ref_num back to message pk first
		pksByRef, err := loadPksByRefs(tx, s.builder, networkID, allImsis)
		if err != nil {
			return nil, err
		}

		err = markMessagesAsDelivered(tx, s.builder, networkID, deliveredMessages, pksByRef)
		if err != nil {
			return nil, err
		}

		err = processFailedMessages(tx, s.builder, networkID, failedMessages, pksByRef)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	_, err := sqorc.ExecInTx(s.db, nil, nil, txFn)
	return err
}
