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

// CreditControl constants and structs to be used in sending/receiving messages
package gy

import (
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/lte/cloud/go/protos"
)

type FinalUnitAction uint8

const (
	Terminate      FinalUnitAction = 0x0
	Redirect       FinalUnitAction = 0x1
	RestrictAccess FinalUnitAction = 0x2
)

type UsedCreditsType int32

const (
	THRESHOLD UsedCreditsType = iota
	QHT
	FINAL                  // FINAL - UE disconnected, flow not in use
	QUOTA_EXHAUSTED        // UE hit credit limit
	VALIDITY_TIMER_EXPIRED // Credit expired
	OTHER_QUOTA_TYPE
	RATING_CONDITION_CHANGE
	FORCED_REAUTHORISATION
	POOL_EXHAUSTED
)

const (
	// 3GPP TS 29.274 RAT Types (for Gy)
	RAT_TYPE_WLAN   = "\x03"
	RAT_TYPE_EUTRAN = "\x06"
)

type UsedCredits struct {
	RatingGroup       uint32
	ServiceIdentifier *uint32
	InputOctets       uint64
	OutputOctets      uint64
	TotalOctets       uint64
	Type              UsedCreditsType
	RequestedUnits    *protos.RequestedUnits
}

type CreditControlRequest struct {
	SessionID               string
	Type                    credit_control.CreditRequestType
	IMSI                    string
	RequestNumber           uint32
	UeIPV4                  string
	SpgwIPV4                string
	Apn                     string
	Imei                    string
	PlmnID                  string
	GcID                    string
	UserLocation            []byte
	Msisdn                  []byte
	Qos                     *QosRequestInfo
	Credits                 []*UsedCredits
	RatType                 string
	TgppCtx                 *protos.TgppContext
	ChargingCharacteristics string
}

type QosRequestInfo struct {
	ApnAggMaxBitRateUL uint32
	ApnAggMaxBitRateDL uint32
}

type ReceivedCredits struct {
	ResultCode          uint32
	RatingGroup         uint32
	ServiceIdentifier   *uint32
	GrantedUnits        *credit_control.GrantedServiceUnit
	ValidityTime        uint32
	FinalUnitIndication *FinalUnitIndication
}

type CreditControlAnswer struct {
	ResultCode    uint32
	SessionID     string
	RequestNumber uint32
	OriginHost    string
	Credits       []*ReceivedCredits
}

type FinalUnitIndication struct {
	FinalAction    FinalUnitAction `avp:"Final-Unit-Action"`
	RedirectServer RedirectServer  `avp:"Redirect-Server"`
	RestrictRules  []string        `avp:"Filter-Id"`
}

type RedirectServer struct {
	RedirectAddressType   RedirectAddressType `avp:"Redirect-Address-Type"`
	RedirectServerAddress string              `avp:"Redirect-Server-Address"`
}

type RedirectAddressType uint8

const (
	IPV4Address RedirectAddressType = iota
	IPV6Address
	URL
	SIPURI
)

type MSCCDiameterMessage struct {
	ResultCode          uint32                            `avp:"Result-Code"`
	GrantedServiceUnit  credit_control.GrantedServiceUnit `avp:"Granted-Service-Unit"`
	ValidityTime        uint32                            `avp:"Validity-Time"`
	FinalUnitIndication *FinalUnitIndication              `avp:"Final-Unit-Indication"`
	RatingGroup         uint32                            `avp:"Rating-Group"`
	ServiceIdentifier   *uint32                           `avp:"Service-Identifier"`
}

type CCADiameterMessage struct {
	SessionID     string                 `avp:"Session-Id"`
	RequestNumber uint32                 `avp:"CC-Request-Number"`
	ResultCode    uint32                 `avp:"Result-Code"`
	OriginHost    string                 `avp:"Origin-Host"`
	RequestType   uint32                 `avp:"CC-Request-Type"`
	CreditControl []*MSCCDiameterMessage `avp:"Multiple-Services-Credit-Control"`
}

// ReAuthRequest is a diameter request received from the OCS to initiate a
// credit update
type ChargingReAuthRequest struct {
	SessionID         string  `avp:"Session-Id"`
	RatingGroup       *uint32 `avp:"Rating-Group"`
	ServiceIdentifier *uint32 `avp:"Service-Identifier"`
}

// ReAuthAnswer is a diameter answer sent back to the OCS after a credit update
// is initiated
type ChargingReAuthAnswer struct {
	SessionID  string `avp:"Session-Id"`
	ResultCode uint32 `avp:"Result-Code"`
}
