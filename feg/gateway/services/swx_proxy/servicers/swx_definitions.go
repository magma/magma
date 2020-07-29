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

package servicers

import "github.com/fiorix/go-diameter/v4/diam/datatype"

const (
	// 3GPP 29.273 8.2.3.12:
	// register AAA server serving authenticated user to HSS
	ServerAssignmentType_REGISTRATION = 1
	// de-register AAA server serving authenticated user to HSS
	ServerAssignnmentType_USER_DEREGISTRATION = 5
	// request user profile from HSS
	ServerAssignmentType_AAA_USER_DATA_REQUEST = 12

	// 3GPP 29.273 5.2.3.6
	RadioAccessTechnologyType_WLAN = 0

	// 3GPP 29.273 8.1.2.1.1/2
	SipAuthScheme_EAP_AKA       = "EAP-AKA"
	SipAuthScheme_EAP_AKA_PRIME = "EAP-AKA'"

	// END_USER_E164 - Subscription-ID Type indicating that the identifier is
	// in international E.164 format (eg. MSISDN).
	// See IETF RFC 4006 section 8.47.
	END_USER_E164 = 0

	// 3GPP 29.273 8.2.3.4
	Non3GPPIPAccess_ENABLED = 0

	// Value of AVP auth-session-state indicating that no state is maintained
	// between calls.
	AuthSessionState_NO_STATE_MAINTAINED = 1
)

// 3GPP 29.273 8.2.2.1 - Multimedia Authentication Request
type MAR struct {
	SessionID           datatype.UTF8String         `avp:"Session-Id"`
	VendorSpecificAppId VendorSpecificApplicationId `avp:"Vendor-Specific-Application-Id"`
	OriginHost          datatype.DiameterIdentity   `avp:"Origin-Host"`
	OriginRealm         datatype.DiameterIdentity   `avp:"Origin-Realm"`
	AuthSessionState    datatype.UTF8String         `avp:"Auth-Session-State"`
	UserName            string                      `avp:"User-Name"`
	RATType             datatype.Enumerated         `avp:"RAT-Type"`
	AuthData            SIPAuthDataItem             `avp:"SIP-Auth-Data-Item"`
	NumberAuthItems     uint32                      `avp:"SIP-Number-Auth-Items"`
}

// 3GPP 29.273 8.2.2.1 - Multimedia Authentication Answer
type MAA struct {
	SessionID          string                    `avp:"Session-Id"`
	ResultCode         uint32                    `avp:"Result-Code"`
	OriginHost         datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm        datatype.DiameterIdentity `avp:"Origin-Realm"`
	AuthSessionState   int32                     `avp:"Auth-Session-State"`
	ExperimentalResult ExperimentalResult        `avp:"Experimental-Result"`
	SIPAuthDataItems   []SIPAuthDataItem         `avp:"SIP-Auth-Data-Item"`
	SIPNumberAuthItems uint32                    `avp:"SIP-Number-Auth-Items"`
	AAAServerName      datatype.DiameterIdentity `avp:"TGPP-AAA-Server-Name"`
}

// VendorSpecificApplicationId
type VendorSpecificApplicationId struct {
	VendorId          uint32 `avp:"Vendor-Id"`
	AuthApplicationId uint32 `avp:"Auth-Application-Id"`
	AcctApplicationId uint32 `avp:"Acct-Application-Id"`
}

type ExperimentalResult struct {
	VendorId               uint32 `avp:"Vendor-Id"`
	ExperimentalResultCode uint32 `avp:"Experimental-Result-Code"`
}

type SIPAuthDataItem struct {
	AuthScheme         string               `avp:"SIP-Authentication-Scheme"`
	Authenticate       datatype.OctetString `avp:"SIP-Authenticate"`
	Authorization      datatype.OctetString `avp:"SIP-Authorization"`
	ConfidentialityKey datatype.OctetString `avp:"Confidentiality-Key"`
	IntegrityKey       datatype.OctetString `avp:"Integrity-Key"`
}

// 3GPP 29.273 8.2.2.3 - Server Assignment Request
type SAR struct {
	SessionID            datatype.UTF8String         `avp:"Session-Id"`
	VendorSpecificAppId  VendorSpecificApplicationId `avp:"Vendor-Specific-Application-Id"`
	OriginHost           datatype.DiameterIdentity   `avp:"Origin-Host"`
	OriginRealm          datatype.DiameterIdentity   `avp:"Origin-Realm"`
	AuthSessionState     datatype.Unsigned32         `avp:"Auth-Session-State"`
	UserName             datatype.UTF8String         `avp:"User-Name"`
	ServerAssignmentType datatype.Enumerated         `avp:"Server-Assignment-Type"`
}

// 3GPP 29.273 8.2.2.3 - Server Assignment Answer
type SAA struct {
	SessionID          string                    `avp:"Session-Id"`
	AuthSessionState   int32                     `avp:"Auth-Session-State"`
	ResultCode         uint32                    `avp:"Result-Code"`
	OriginHost         datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm        datatype.DiameterIdentity `avp:"Origin-Realm"`
	ExperimentalResult ExperimentalResult        `avp:"Experimental-Result"`
	UserName           datatype.UTF8String       `avp:"User-Name"`
	UserData           Non3GPPUserData           `avp:"Non-3GPP-User-Data"`
	AAAServerName      datatype.DiameterIdentity `avp:"TGPP-AAA-Server-Name"`
}

type Non3GPPUserData struct {
	SubscriptionId  SubscriptionId      `avp:"Subscription-Id"`
	Non3GPPIPAccess datatype.Enumerated `avp:"Non-3GPP-IP-Access"`
}

type SubscriptionId struct {
	SubscriptionIdType datatype.Enumerated `avp:"Subscription-Id-Type"`
	SubscriptionIdData datatype.UTF8String `avp:"Subscription-Id-Data"`
}

// 3GPP 29.273 8.2.2.4 - Registration Termination Request
type RTR struct {
	SessionID            datatype.UTF8String         `avp:"Session-Id"`
	VendorSpecificAppId  VendorSpecificApplicationId `avp:"Vendor-Specific-Application-Id"`
	OriginHost           datatype.DiameterIdentity   `avp:"Origin-Host"`
	OriginRealm          datatype.DiameterIdentity   `avp:"Origin-Realm"`
	AuthSessionState     datatype.Unsigned32         `avp:"Auth-Session-State"`
	UserName             datatype.UTF8String         `avp:"User-Name"`
	DeregistrationReason DeregistrationReason        `avp:"Deregistration-Reason"`
}

// 3GPP 29.273 8.2.2.4 - Registration Termination Answer
type RTA struct {
	SessionID          string                    `avp:"Session-Id"`
	AuthSessionState   int32                     `avp:"Auth-Session-State"`
	ResultCode         uint32                    `avp:"Result-Code"`
	OriginHost         datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm        datatype.DiameterIdentity `avp:"Origin-Realm"`
	ExperimentalResult ExperimentalResult        `avp:"Experimental-Result"`
}

type DeregistrationReason struct {
	ReasonCode datatype.Enumerated `avp:"Reason-Code"`
	ReasonInfo datatype.UTF8String `avp:"Reason-Info"`
}
