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

// package sim implements EAP-SIM provider
package sim

import (
	"fmt"
	"time"

	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
)

const (
	TYPE           = uint8(protos.EapType_SIM)
	MIN_PACKET_LEN = eap.EapSubtype

	EapSimServiceName = "eap_sim"
)

const (
	// SIM Attributes
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
	AT_CLIENT_ERROR_CODE // 22
)

const (
	// SIM Notification Codes
	NOTIFICATION_FAILURE uint16 = 16384
)

const (
	// IMSI Consts
	MinImsiLen = 6
	MaxImsiLen = 16
)

type Subtype uint8

const (
	// SIM Subtypes
	SubtypeStart            Subtype = 10
	SubtypeChallenge        Subtype = 11
	SubtypeNotification     Subtype = 12
	SubtypeReauthentication Subtype = 13
	SubtypeClientError      Subtype = 14
)

type SimState int16

const (
	// Processing/handling States
	StateNone          SimState = iota
	StateCreated                // newly created
	StateIdentity               // Valid permanent identity received
	StateChallenge              // Auth Challenge was returned to UE
	StateAuthenticated          // UE is successfully authenticated
	StateRedirected             // UE is redirected to another Auth method, cache this state to prevent redirection loop
)

const (
	ATT_HDR_LEN = 4
	AUTN_LEN    = 16
	RAND_LEN    = 16
	RandAutnLen = RAND_LEN + AUTN_LEN
	MAC_LEN     = 16

	DefaultChallengeTimeout            = time.Second * 20
	DefaultErrorNotificationTimeout    = time.Second * 10
	DefaultSessionTimeout              = time.Hour * 12
	DefaultSessionAuthenticatedTimeout = time.Second * 5
	GsmTripletsNumber                  = 3

	Version byte = 1 // SIM's Supported Version
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
	if l == MaxImsiLen && i[0] != '1' {
		return fmt.Errorf("Invalid IMSI %s", i)
	}
	for idx, c := range i {
		if c < '0' || c > '9' {
			return fmt.Errorf("Unexpected IMSI byte 0x%X (%c) at index %d", c, c, idx)
		}
	}
	return nil
}
