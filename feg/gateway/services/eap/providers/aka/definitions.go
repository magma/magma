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

// package aka implements EAP-AKA provider
package aka

import (
	"fmt"
	"time"

	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
)

const (
	TYPE           = uint8(protos.EapType_AKA)
	MIN_PACKET_LEN = eap.EapSubtype

	EapAkaServiceName = "eap_aka"
)

const (
	// AKA Attributes
	AT_RAND eap.AttrType = iota + 1
	AT_AUTN
	AT_RES
	AT_AUTS
	_
	AT_PADDING
	AT_NONCE_MT
	_
	_
	AT_PERMANENT_ID_REQ
	AT_MAC
	AT_NOTIFICATION
	AT_ANY_ID_REQ
	AT_IDENTITY
	AT_VERSION_LIST
	AT_SELECTED_VERSION
	AT_FULLAUTH_ID_REQ
	_
	AT_COUNTER
	AT_COUNTER_TOO_SMALL
	AT_NONCE_S
	AT_CLIENT_ERROR_CODE              // 22
	AT_IV                eap.AttrType = 129
	AT_ENCR_DATA         eap.AttrType = 130
	AT_NEXT_PSEUDONYM    eap.AttrType = 132
	AT_NEXT_REAUTH_ID    eap.AttrType = 133
	AT_CHECKCODE         eap.AttrType = 134
	AT_RESULT_IND        eap.AttrType = 135
)

const (
	// AKA Notification Codes
	NOTIFICATION_FAILURE_AUTH   uint16 = 0
	NOTIFICATION_FAILURE        uint16 = 16384
	NOTIFICATION_SUCCESS        uint16 = 32768
	NOTIFICATION_ACCESS_DENIED  uint16 = 1026
	NOTIFICATION_NOT_SUBSCRIBED uint16 = 1031
)

const (
	// IMSI Consts
	MinImsiLen = 6
	MaxImsiLen = 16
)

type Subtype uint8

const (
	// AKA Subtypes
	_ Subtype = iota
	SubtypeChallenge
	SubtypeAuthenticationReject
	_
	SubtypeSynchronizationFailure
	SubtypeIdentity
	SubtypeNotification     Subtype = 12
	SubtypeReauthentication Subtype = 13
	SubtypeClientError      Subtype = 14
)

type AkaState int16

const (
	// Processing/handling States
	StateNone          AkaState = iota
	StateCreated                // newly created
	StateIdentity               // Valid permanent identity received
	StateChallenge              // Auth Challenge was returned to UE
	StateAuthenticated          // UE is successfully authenticated
	StateRedirected             // UE is redirected to another Auth method (SIM, AKA', etc.)
)

const (
	ATT_HDR_LEN = 4
	AUTN_LEN    = 16
	RAND_LEN    = 16
	RandAutnLen = RAND_LEN + AUTN_LEN
	MAC_LEN     = 16

	AT_RAND_ATTR_LEN = AUTN_LEN + ATT_HDR_LEN
	AT_AUTN_ATTR_LEN = RAND_LEN + ATT_HDR_LEN
	AT_MAC_ATTR_LEN  = MAC_LEN + ATT_HDR_LEN

	DefaultChallengeTimeout            = time.Second * 20
	DefaultErrorNotificationTimeout    = time.Second * 10
	DefaultSessionTimeout              = time.Hour * 12
	DefaultSessionAuthenticatedTimeout = time.Second * 5
)

type IMSI string

func (i IMSI) Validate() error {
	l := len(i)
	if l > MaxImsiLen {
		return fmt.Errorf("IMSI %s is too long: %d", i, l)
	}
	if l < MinImsiLen {
		return fmt.Errorf("IMSI %s is too short: %d", i, l)
	}
	if l == MaxImsiLen && i[0] != '0' {
		return fmt.Errorf("Invalid IMSI %s", i)
	}
	for idx, c := range i {
		if c < '0' || c > '9' {
			return fmt.Errorf("Unexpected IMSI byte 0x%X (%c) at index %d", c, c, idx)
		}
	}
	return nil
}
