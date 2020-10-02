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

package storage_test

import (
	"fmt"
	"testing"
	"time"

	"magma/lte/cloud/go/services/smsd/storage"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
)

func TestSQLSMSStorage_Integration(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:?_foreign.keys=1")
	if err != nil {
		t.Fatalf("Could not initialize sqlite DB: %s", err)
	}
	refCounter := &mockRefCounter{numRefs: 1}
	store := storage.NewSQLSMSStorage(db, sqorc.GetSqlBuilder(), refCounter, &mockIDGenerator{})

	err = store.Init()
	if err != nil {
		t.Fatalf("Could not initialize smsd tables: %s", err)
	}

	var frozenClock int64 = 1000
	clock.SetAndFreezeClock(t, time.Unix(frozenClock, 0))
	defer clock.UnfreezeClock(t)

	// Empty-case tests:
	// No SMSs when you list them
	// No SMSs to deliver
	actualMessages, err := store.GetSMSs("n1", nil, nil, false, nil, nil)
	assert.NoError(t, err)
	assert.Empty(t, actualMessages)

	actualMessages, err = store.GetSMSsToDeliver("n1", []string{"IMSI1", "IMSI2"}, 0)
	assert.NoError(t, err)
	assert.Empty(t, actualMessages)

	// Create 1 SMS each for 2 subs
	_, err = store.CreateSMS(
		"n1",
		storage.MutableSMS{
			Imsi:         "IMSI1",
			SourceMsisdn: "123",
			Message:      "hello world",
		},
	)
	assert.NoError(t, err)
	_, err = store.CreateSMS(
		"n1",
		storage.MutableSMS{
			Imsi:         "IMSI2",
			SourceMsisdn: "456",
			Message:      "goodbye world",
		},
	)
	assert.NoError(t, err)

	actualMessages, err = store.GetSMSs("n1", nil, nil, false, nil, nil)
	assert.NoError(t, err)
	expectedMessages := []*storage.SMS{
		{
			Pk:                      "1",
			Status:                  storage.MessageStatus_WAITING,
			Imsi:                    "IMSI1",
			SourceMsisdn:            "123",
			Message:                 "hello world",
			CreatedTime:             timestampProto(t, 1000),
			LastDeliveryAttemptTime: nil,
			AttemptCount:            0,
			DeliveryError:           "",
			RefNums:                 nil,
		},
		{
			Pk:                      "2",
			Status:                  storage.MessageStatus_WAITING,
			Imsi:                    "IMSI2",
			SourceMsisdn:            "456",
			Message:                 "goodbye world",
			CreatedTime:             timestampProto(t, 1000),
			LastDeliveryAttemptTime: nil,
			AttemptCount:            0,
			DeliveryError:           "",
			RefNums:                 nil,
		},
	}
	assert.Equal(t, expectedMessages, actualMessages)

	frozenClock += 100
	clock.SetAndFreezeClock(t, time.Unix(frozenClock, 0))

	// Ask for SMSs to deliver for IMSI1
	actualMessages, err = store.GetSMSsToDeliver("n1", []string{"IMSI1"}, 0)
	assert.NoError(t, err)
	expectedMessages = []*storage.SMS{
		{
			Pk:           "1",
			Status:       storage.MessageStatus_WAITING,
			Imsi:         "IMSI1",
			SourceMsisdn: "123",
			Message:      "hello world",
			CreatedTime:  timestampProto(t, 1000),
			// "Now"
			LastDeliveryAttemptTime: timestampProto(t, frozenClock),
			AttemptCount:            1,
			DeliveryError:           "",
			RefNums:                 []byte{0},
		},
	}
	assert.Equal(t, expectedMessages, actualMessages)

	// Should get nothing back from another call
	actualMessages, err = store.GetSMSsToDeliver("n1", []string{"IMSI1"}, 0)
	assert.NoError(t, err)
	assert.Empty(t, actualMessages)

	// Timeout this first message, we should get it back again with the same
	// ref but advanced delivery attempt time and attempt count
	frozenClock += 10000
	clock.SetAndFreezeClock(t, time.Unix(frozenClock, 0))

	actualMessages, err = store.GetSMSsToDeliver("n1", []string{"IMSI1"}, 0)
	assert.NoError(t, err)
	expectedMessages[0].LastDeliveryAttemptTime = timestampProto(t, frozenClock)
	expectedMessages[0].AttemptCount = 2
	assert.Equal(t, expectedMessages, actualMessages)

	// Double-check the all message list output
	actualMessages, err = store.GetSMSs("n1", nil, nil, false, nil, nil)
	assert.NoError(t, err)
	expectedAllMessages := []*storage.SMS{
		expectedMessages[0],
		{
			Pk:                      "2",
			Status:                  storage.MessageStatus_WAITING,
			Imsi:                    "IMSI2",
			SourceMsisdn:            "456",
			Message:                 "goodbye world",
			CreatedTime:             timestampProto(t, 1000),
			LastDeliveryAttemptTime: nil,
			AttemptCount:            0,
			DeliveryError:           "",
			RefNums:                 nil,
		},
	}
	assert.Equal(t, actualMessages, expectedAllMessages)

	// Report that delivery failed, should get it back again with the same ref
	// number again but attempt time and count advanced
	// Also include an unknown message
	err = store.ReportDelivery(
		"n1",
		map[string][]storage.SMSRef{"IMSI3": {0x1}},
		map[string][]storage.SMSFailureReport{
			"IMSI1": {{Ref: 0, ErrorMessage: "foobar"}},
		},
	)
	assert.NoError(t, err)
	actualMessages, err = store.GetSMSs("n1", nil, nil, false, nil, nil)
	assert.NoError(t, err)
	expectedAllMessages[0].DeliveryError = "foobar"
	assert.Equal(t, expectedAllMessages, actualMessages)

	// We need to timeout the message to try and re-send it, because we don't
	// delete the ref for a failed message that is not at the retry limit
	// This is on purpose because otherwise the select query for this would
	// have a crazy WHERE clause.
	frozenClock += 1000
	clock.SetAndFreezeClock(t, time.Unix(frozenClock, 0))

	actualMessages, err = store.GetSMSsToDeliver("n1", []string{"IMSI1"}, 0)
	assert.NoError(t, err)
	expectedMessages[0].LastDeliveryAttemptTime = timestampProto(t, frozenClock)
	expectedMessages[0].AttemptCount = 3
	assert.Equal(t, actualMessages, expectedMessages)

	// Mark this message as failed delivery again, time it out, and we should
	// no longer see it as a message that needs to be sent (exceeded retry)
	err = store.ReportDelivery(
		"n1",
		nil,
		map[string][]storage.SMSFailureReport{
			"IMSI1": {{Ref: 0, ErrorMessage: "barbaz"}},
		},
	)
	assert.NoError(t, err)
	frozenClock += 1000
	clock.SetAndFreezeClock(t, time.Unix(frozenClock, 0))
	actualMessages, err = store.GetSMSsToDeliver("n1", []string{"IMSI1"}, 0)
	assert.NoError(t, err)
	assert.Empty(t, actualMessages)
	actualMessages, err = store.GetSMSs("n1", nil, nil, false, nil, nil)
	assert.NoError(t, err)
	expectedAllMessages[0] = &storage.SMS{
		Pk:            "1",
		Status:        storage.MessageStatus_FAILED,
		Imsi:          "IMSI1",
		SourceMsisdn:  "123",
		Message:       "hello world",
		CreatedTime:   timestampProto(t, 1000),
		AttemptCount:  3,
		DeliveryError: "barbaz",
		// Note that LastDeliveryAttemptTime is unfilled because the ref is
		// gone and the ref row is where we save that timestamp.
	}
	assert.Equal(t, expectedAllMessages, actualMessages)

	// Create new SMSs: 1 for imsi2, 1 for imsi3
	_, err = store.CreateSMS(
		"n1",
		storage.MutableSMS{
			Imsi:         "IMSI2",
			SourceMsisdn: "789",
			Message:      "message 3",
		},
	)
	assert.NoError(t, err)
	_, err = store.CreateSMS(
		"n1",
		storage.MutableSMS{
			Imsi:         "IMSI3",
			SourceMsisdn: "123",
			Message:      "message 4",
		},
	)
	assert.NoError(t, err)

	// Allocate 2 ref nums per message this time
	refCounter.numRefs = 2
	actualMessages, err = store.GetSMSsToDeliver("n1", []string{"IMSI1", "IMSI2", "IMSI3"}, 0)
	assert.NoError(t, err)
	expectedMessages = []*storage.SMS{
		{
			Pk:                      "2",
			Status:                  storage.MessageStatus_WAITING,
			Imsi:                    "IMSI2",
			SourceMsisdn:            "456",
			Message:                 "goodbye world",
			CreatedTime:             timestampProto(t, 1000),
			LastDeliveryAttemptTime: timestampProto(t, frozenClock),
			AttemptCount:            1,
			DeliveryError:           "",
			RefNums:                 []byte{0x0, 0x1},
		},
		{
			Pk:                      "3",
			Status:                  storage.MessageStatus_WAITING,
			Imsi:                    "IMSI2",
			SourceMsisdn:            "789",
			Message:                 "message 3",
			CreatedTime:             timestampProto(t, frozenClock),
			LastDeliveryAttemptTime: timestampProto(t, frozenClock),
			AttemptCount:            1,
			DeliveryError:           "",
			RefNums:                 []byte{0x2, 0x3},
		},
		{
			Pk:                      "4",
			Status:                  storage.MessageStatus_WAITING,
			Imsi:                    "IMSI3",
			SourceMsisdn:            "123",
			Message:                 "message 4",
			CreatedTime:             timestampProto(t, frozenClock),
			LastDeliveryAttemptTime: timestampProto(t, frozenClock),
			AttemptCount:            1,
			DeliveryError:           "",
			RefNums:                 []byte{0x0, 0x1},
		},
	}
	assert.Equal(t, expectedMessages, actualMessages)

	// Mark both imsi2 as delivered, imsi3 as failed
	err = store.ReportDelivery(
		"n1",
		map[string][]storage.SMSRef{
			"IMSI2": {0x0, 0x1, 0x2, 0x3},
		},
		map[string][]storage.SMSFailureReport{
			"IMSI3": {
				{Ref: 0, ErrorMessage: "foo bar"},
				{Ref: 1, ErrorMessage: "foo bar"},
			},
		},
	)
	assert.NoError(t, err)

	// Another query should return only the message for imsi3
	frozenClock += 1000
	clock.SetAndFreezeClock(t, time.Unix(frozenClock, 0))
	actualMessages, err = store.GetSMSsToDeliver("n1", []string{"IMSI1", "IMSI2", "IMSI3"}, 0)
	assert.NoError(t, err)
	expectedMessages = []*storage.SMS{
		{
			Pk:                      "4",
			Status:                  storage.MessageStatus_WAITING,
			Imsi:                    "IMSI3",
			SourceMsisdn:            "123",
			Message:                 "message 4",
			CreatedTime:             timestampProto(t, frozenClock-1000),
			LastDeliveryAttemptTime: timestampProto(t, frozenClock),
			AttemptCount:            2,
			DeliveryError:           "foo bar",
			RefNums:                 []byte{0x0, 0x1},
		},
	}
	assert.Equal(t, expectedMessages, actualMessages)

	// Mark it as delivered
	err = store.ReportDelivery(
		"n1",
		map[string][]storage.SMSRef{
			"IMSI3": {0x0, 0x1},
		},
		nil,
	)
	assert.NoError(t, err)

	// No more messages to send
	frozenClock += 1000
	clock.SetAndFreezeClock(t, time.Unix(frozenClock, 0))
	actualMessages, err = store.GetSMSsToDeliver("n1", []string{"IMSI1", "IMSI2", "IMSI3"}, 0)
	assert.NoError(t, err)
	assert.Empty(t, actualMessages)

	// Create a new message for imsi1 and verify refs get re-used
	_, err = store.CreateSMS(
		"n1",
		storage.MutableSMS{
			Imsi:         "IMSI1",
			SourceMsisdn: "123",
			Message:      "message 5",
		},
	)
	assert.NoError(t, err)
	actualMessages, err = store.GetSMSsToDeliver("n1", []string{"IMSI1", "IMSI2", "IMSI3"}, 0)
	assert.NoError(t, err)
	expectedMessages = []*storage.SMS{
		{
			Pk:                      "5",
			Status:                  storage.MessageStatus_WAITING,
			Imsi:                    "IMSI1",
			SourceMsisdn:            "123",
			Message:                 "message 5",
			CreatedTime:             timestampProto(t, frozenClock),
			LastDeliveryAttemptTime: timestampProto(t, frozenClock),
			AttemptCount:            1,
			DeliveryError:           "",
			RefNums:                 []byte{0x0, 0x1},
		},
	}
	assert.Equal(t, expectedMessages, actualMessages)

	actualMessages, err = store.GetSMSs("n1", nil, nil, false, nil, nil)
	assert.NoError(t, err)
	expectedAllMessages = []*storage.SMS{
		{
			Pk:            "1",
			Status:        storage.MessageStatus_FAILED,
			Imsi:          "IMSI1",
			SourceMsisdn:  "123",
			Message:       "hello world",
			CreatedTime:   timestampProto(t, 1000),
			AttemptCount:  3,
			DeliveryError: "barbaz",
		},
		{
			Pk:           "2",
			Status:       storage.MessageStatus_DELIVERED,
			Imsi:         "IMSI2",
			SourceMsisdn: "456",
			Message:      "goodbye world",
			CreatedTime:  timestampProto(t, 1000),
			AttemptCount: 1,
		},
		{
			Pk:           "3",
			Status:       storage.MessageStatus_DELIVERED,
			Imsi:         "IMSI2",
			SourceMsisdn: "789",
			Message:      "message 3",
			CreatedTime:  timestampProto(t, 13100),
			AttemptCount: 1,
		},
		{
			Pk:           "4",
			Status:       storage.MessageStatus_DELIVERED,
			Imsi:         "IMSI3",
			SourceMsisdn: "123",
			Message:      "message 4",
			CreatedTime:  timestampProto(t, 13100),
			AttemptCount: 2,
		},
		{
			Pk:                      "5",
			Status:                  storage.MessageStatus_WAITING,
			Imsi:                    "IMSI1",
			SourceMsisdn:            "123",
			Message:                 "message 5",
			CreatedTime:             timestampProto(t, frozenClock),
			LastDeliveryAttemptTime: timestampProto(t, frozenClock),
			AttemptCount:            1,
			DeliveryError:           "",
			RefNums:                 []byte{0x0, 0x1},
		},
	}
	assert.Equal(t, expectedAllMessages, actualMessages)

	// Check PK filter
	actualMessages, err = store.GetSMSs("n1", []string{"1", "2", "3", "4"}, nil, false, nil, nil)
	assert.NoError(t, err)
	expectedAllMessages = expectedAllMessages[0:4]
	assert.Equal(t, expectedAllMessages, actualMessages)

	// Make sure things are gated across networks
	actualMessages, err = store.GetSMSs("n2", nil, nil, false, nil, nil)
	assert.NoError(t, err)
	assert.Empty(t, actualMessages)

	actualMessages, err = store.GetSMSsToDeliver("n2", []string{"IMSI1", "IMSI2", "IMSI3"}, 0)
	assert.NoError(t, err)
	assert.Empty(t, actualMessages)

	// Delete all messages
	err = store.DeleteSMSs("n1", []string{"1", "2", "3", "4", "5"})
	assert.NoError(t, err)
	actualMessages, err = store.GetSMSs("n1", nil, nil, false, nil, nil)
	assert.NoError(t, err)
	assert.Empty(t, actualMessages)
}

type mockRefCounter struct {
	numRefs uint16
}

func (m *mockRefCounter) GetReferenceNumberCount(message string) uint16 {
	return m.numRefs
}

type mockIDGenerator struct {
	curID int
}

func (m *mockIDGenerator) New() string {
	m.curID++
	return fmt.Sprintf("%d", m.curID)
}

func timestampProto(t *testing.T, unix int64) *timestamp.Timestamp {
	ret, err := ptypes.TimestampProto(time.Unix(unix, 0))
	assert.NoError(t, err)
	return ret
}
