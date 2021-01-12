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
	"time"

	"magma/lte/cloud/go/sms_ll"
)

type SMSRef = byte

type SMSFailureReport struct {
	// TODO: this should really be []SMSRef but I'm not 100% sure if this
	//  struct will remain a valid map key
	Ref          SMSRef
	ErrorMessage string
}

// SMSStorage is the storage interface for managing SMS messages and their
// delivery.
// SMS's are intended to be immutable upon creation (delete or read only).
type SMSStorage interface {
	// Init performs on-start initialization work such as table creation.
	Init() error

	// GetSMSs returns all SMS messages in a time window matching the
	// provided IMSI and status filters.
	//
	// If pks is non-empty, this will fetch only the specified messages, as
	// long as they are in the specified network.
	// If imsis is non-empty, this will query for all SMS messages tracked in
	// the system, otherwise this method will additionally filter on
	// destination IMSI.
	// If onlyWaiting is true, this will only query for messages that need to
	// be delivered.
	// startTime defaults to epoch if nil
	// endTime defaults to current time if nil
	GetSMSs(networkID string, pks []string, imsis []string, onlyWaiting bool, startTime, endTime *time.Time) ([]*SMS, error)

	// GetSMSsToDeliver will return a collection of messages that need to be
	// delivered. For each IMSI requested, this will return a maximum of 256
	// pending messages.
	//
	// Only messages which are not in-flight will be returned by this call.
	// This will generate reference numbers for each message to be sent.
	// An in-flight message will be returned again if its reference number
	// has been generated longer than timeoutThreshold ago (default 30m).
	//
	// WARNING: Concurrent calls to this method with intersecting IMSI sets
	// will result in undefined behavior.
	GetSMSsToDeliver(networkID string, imsis []string, timeoutThreshold time.Duration) ([]*SMS, error)

	// CreateSMS creates a new SMS message. The auto-generated pk for the
	// message is returned.
	CreateSMS(networkID string, sms MutableSMS) (string, error)

	// DeleteSMSs deletes messages by pk. Semantics are all or nothing.
	DeleteSMSs(networkID string, pks []string) error

	// ReportDelivery reports delivery status of a set of SMSs
	// Map keys for both arguments are IMSIs
	ReportDelivery(networkID string, deliveredMessages map[string][]SMSRef, failedMessages map[string][]SMSFailureReport) error
}

// SMSReferenceCounter is a functional interface that wraps the logic to
// determine how many SMS messages a message string will be encoded into.
type SMSReferenceCounter interface {
	GetReferenceNumberCount(message string) uint16
}

type DefaultSMSReferenceCounter struct{}

func (*DefaultSMSReferenceCounter) GetReferenceNumberCount(message string) uint16 {
	return uint16(sms_ll.GetMessageCount(message))
}
