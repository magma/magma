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

package gx

import (
	"time"

	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/lte/cloud/go/protos"
)

type MonitoringLevel uint8

const (
	SessionLevel MonitoringLevel = 0x0
	RuleLevel    MonitoringLevel = 0x1
)

type MonitoringSupport uint8

const UsageMonitoringDisabled MonitoringSupport = 0x0

type MonitoringReport uint8

const UsageMonitoringReport MonitoringReport = 0x0

type EventTrigger uint32

const (
	RevalidationTimeout      EventTrigger = 17
	UsageReportTrigger       EventTrigger = 33
	PCRF91UsageReportTrigger EventTrigger = 26
)

// CreditControlRequest represents a call over gx
type CreditControlRequest struct {
	SessionID               string
	Type                    credit_control.CreditRequestType
	IMSI                    string
	RequestNumber           uint32
	IPAddr                  string
	IPv6Addr                string
	SpgwIPV4                string
	Apn                     string
	Msisdn                  []byte
	Imei                    string
	PlmnID                  string
	UserLocation            []byte
	GcID                    string
	Qos                     *QosRequestInfo
	UsageReports            []*UsageReport
	HardwareAddr            []byte
	IPCANType               credit_control.IPCANType
	RATType                 credit_control.RATType
	TgppCtx                 *protos.TgppContext
	EventTrigger            EventTrigger
	AccessTimezone          *protos.Timezone
	ChargingCharacteristics string
}

type QosRequestInfo struct {
	ApnAggMaxBitRateUL         uint32
	ApnAggMaxBitRateDL         uint32
	ApnExtendedAggMaxBitRateUL uint32
	ApnExtendedAggMaxBitRateDL uint32
	QosClassIdentifier         uint32
	PriLevel                   uint32
	PreCapability              uint32
	PreVulnerability           uint32
}

// CreditControlAnswer represents the gx CCA message we're expecting
type CreditControlAnswer struct {
	ResultCode             uint32
	ExperimentalResultCode uint32
	SessionID              string
	OriginHost             string
	RequestNumber          uint32
	RuleInstallAVP         []*RuleInstallAVP
	RuleRemoveAVP          []*RuleRemoveAVP
	UsageMonitors          []*UsageMonitoringInfo
	EventTriggers          []EventTrigger
	RevalidationTime       *time.Time
	Qos                    *QosInformation
}

type UsageReport struct {
	MonitoringKey []byte
	Level         MonitoringLevel
	InputOctets   uint64
	OutputOctets  uint64
	TotalOctets   uint64
}

// RedirectInformation represents Information needed for redirection setup
type RedirectInformation struct {
	RedirectSupport       uint32 `avp:"Redirect-Support"`
	RedirectAddressType   uint32 `avp:"Redirect-Address-Type"`
	RedirectServerAddress string `avp:"Redirect-Server-Address"`
}

// Flow-Information ::= < AVP Header: 1058 >
//  [ Flow-Description ]
//  [ Packet-Filter-Identifier ]
//  [ Packet-Filter-Usage ]
//  [ ToS-Traffic-Class ]
//  [ Security-Parameter-Index ]
//  [ Flow-Label ]
//  [ Flow-Direction ]
// Only Flow-Description is supported right now
type FlowInformation struct {
	FlowDescription string `avp:"Flow-Description"`
}

type RuleDefinition struct {
	RuleName            string               `avp:"Charging-Rule-Name"`
	RatingGroup         *uint32              `avp:"Rating-Group"`
	Precedence          uint32               `avp:"Precedence"`
	MonitoringKey       []byte               `avp:"Monitoring-Key"`
	FlowDescriptions    []string             `avp:"Flow-Description"`
	FlowInformations    []*FlowInformation   `avp:"Flow-Information"`
	RedirectInformation *RedirectInformation `avp:"Redirect-Information"`
	Qos                 *QosInformation      `avp:"QoS-Information"`
	ServiceIdentifier   *uint32              `avp:"Service-Identifier"`
}

// QoS per service date flow message
type QosInformation struct {
	BearerIdentifier   string  `avp:"Bearer-Identifier"`
	MaxReqBwUL         *uint32 `avp:"Max-Requested-Bandwidth-UL"`
	MaxReqBwDL         *uint32 `avp:"Max-Requested-Bandwidth-DL"`
	ExtendedMaxReqBwUL *uint32 `avp:"Extended-Max-Requested-BW-UL"`
	ExtendedMaxReqBwDL *uint32 `avp:"Extended-Max-Requested-BW-DL"`
	GbrDL              *uint32 `avp:"Guaranteed-Bitrate-DL"`
	GbrUL              *uint32 `avp:"Guaranteed-Bitrate-UL"`
	Qci                *uint32 `avp:"QoS-Class-Identifier"`
}

// RuleInstallAVP represents a policy rule to install. It can hold one of
// rule name, base name or definition
type RuleInstallAVP struct {
	RuleNames            []string          `avp:"Charging-Rule-Name"`
	RuleBaseNames        []string          `avp:"Charging-Rule-Base-Name"`
	RuleDefinitions      []*RuleDefinition `avp:"Charging-Rule-Definition"`
	RuleActivationTime   *time.Time        `avp:"Rule-Activation-Time"`
	RuleDeactivationTime *time.Time        `avp:"Rule-Deactivation-Time"`
}

type RuleRemoveAVP struct {
	RuleNames     []string `avp:"Charging-Rule-Name"`
	RuleBaseNames []string `avp:"Charging-Rule-Base-Name"`
}

type UsageMonitoringInfo struct {
	MonitoringKey      []byte                             `avp:"Monitoring-Key"`
	GrantedServiceUnit *credit_control.GrantedServiceUnit `avp:"Granted-Service-Unit"`
	Level              MonitoringLevel                    `avp:"Usage-Monitoring-Level"`
	Support            *MonitoringSupport                 `avp:"Usage-Monitoring-Support"`
	Report             *MonitoringReport                  `avp:"Usage-Monitoring-Report"`
}

// CCADiameterMessage is a gx CCA message as defined in 3GPP 29.212
type CCADiameterMessage struct {
	SessionID          string `avp:"Session-Id"`
	RequestNumber      uint32 `avp:"CC-Request-Number"`
	ResultCode         uint32 `avp:"Result-Code"`
	OriginHost         string `avp:"Origin-Host"`
	ExperimentalResult struct {
		VendorId               uint32 `avp:"Vendor-Id"`
		ExperimentalResultCode uint32 `avp:"Experimental-Result-Code"`
	} `avp:"Experimental-Result"`
	RequestType      uint32                 `avp:"CC-Request-Type"`
	RuleInstalls     []*RuleInstallAVP      `avp:"Charging-Rule-Install"`
	RuleRemovals     []*RuleRemoveAVP       `avp:"Charging-Rule-Remove"`
	UsageMonitors    []*UsageMonitoringInfo `avp:"Usage-Monitoring-Information"`
	EventTriggers    []EventTrigger         `avp:"Event-Trigger"`
	RevalidationTime *time.Time             `avp:"Revalidation-Time"`
	Qos              *QosInformation        `avp:"QoS-Information"`
}

//<RA-Request> ::= 	< Diameter Header: 258, REQ, PXY >
//					< Session-Id >
//					[ DRMP ]
//					{ Auth-Application-Id }
//					{ Origin-Host }
//					{ Origin-Realm }
//					{ Destination-Realm }
//					{ Destination-Host }
//					{ Re-Auth-Request-Type }
//					[ Session-Release-Cause ]
//					[ Origin-State-Id ]
//					[ OC-Supported-Features ]
//					*[ Event-Trigger ]
//					[ Event-Report-Indication ]
//					*[ Charging-Rule-Remove ]
//					*[ Charging-Rule-Install ]
//					[ Default-EPS-Bearer-QoS ]
//					*[ QoS-Information ]
//					[ Default-QoS-Information ]
//					[ Revalidation-Time ]
//					*[ Usage-Monitoring-Information ]
//					[ PCSCF-Restoration-Indication ] 0*4[ Conditional-Policy-Information ]
//					[ Removal-Of-Access ]
//					[ IP-CAN-Type ]
//					[ PRA-Install ]
//					[ PRA-Remove ]
//					*[ Proxy-Info ]
//					*[ Route-Record ]
//					*[ AVP ]
type PolicyReAuthRequest struct {
	SessionID        string                 `avp:"Session-Id"`
	OriginHost       string                 `avp:"Origin-Host"`
	RulesToRemove    []*RuleRemoveAVP       `avp:"Charging-Rule-Remove"`
	RulesToInstall   []*RuleInstallAVP      `avp:"Charging-Rule-Install"`
	Qos              *QosInformation        `avp:"QoS-Information"`
	UsageMonitors    []*UsageMonitoringInfo `avp:"Usage-Monitoring-Information"`
	EventTriggers    []EventTrigger         `avp:"Event-Trigger"`
	RevalidationTime *time.Time             `avp:"Revalidation-Time"`
}

//<RA-Answer> ::= 	< Diameter Header: 258, PXY >
//					< Session-Id >
//					[ DRMP ]
//					{ Origin-Host }
//					{ Origin-Realm }
//					[ Result-Code ]
//					[ Experimental-Result ]
//					[ Origin-State-Id ]
//					[ OC-Supported-Features ]
//					[ OC-OLR ]
//					[ IP-CAN-Type ]
//					[ RAT-Type ]
//					[ AN-Trusted ]
//					0*2 [ AN-GW-Address ]
//					[ 3GPP-SGSN-MCC-MNC ]
//					[ 3GPP-SGSN-Address ]
//					[ 3GPP-SGSN-Ipv6-Address ]
//					[ RAI ]
//					[ 3GPP-User-Location-Info ]
//					[ User-Location-Info-Time ]
//					[ NetLoc-Access-Support ]
//					[ User-CSG-Information ]
//					[ 3GPP-MS-TimeZone ]
//					[ Default-QoS-Information ]
//					*[ Charging-Rule-Report]
//					[ Error-Message ]
//					[ Error-Reporting-Host ]
//					[ Failed-AVP ]
//					*[ Proxy-Info ]
//					*[ AVP ]
type PolicyReAuthAnswer struct {
	SessionID   string                `avp:"Session-Id"`
	ResultCode  uint32                `avp:"Result-Code"`
	RuleReports []*ChargingRuleReport `avp:"Charging-Rule-Report"`
}

//Charging-Rule-Report ::= < AVP Header: 1018 >
// 						  *[ Charging-Rule-Name ]
//                        *[ Charging-Rule-Base-Name ]
//                         [ Bearer-Identifier ]
//                         [ PCC-Rule-Status ]
//                         [ Rule-Failure-Code ]
//                         [ Final-Unit-Indication ]
//                        *[ RAN-NAS-Release-Cause ]
//                        *[ Content-Version ]
//                        *[ AVP ]
type ChargingRuleReport struct {
	RuleNames     []string        `avp:"Charging-Rule-Name"`
	RuleBaseNames []string        `avp:"Charging-Rule-Base-Name"`
	FailureCode   RuleFailureCode `avp:"Rule-Failure-Code"`
}

type RuleFailureCode uint32

const (
	UnknownRuleName                 RuleFailureCode = 1
	RatingGroupError                RuleFailureCode = 2
	ServiceIdentifierError          RuleFailureCode = 3
	GwPCEFMalfunction               RuleFailureCode = 4
	ResourcesLimitation             RuleFailureCode = 5
	MaxNrBearersReached             RuleFailureCode = 6
	UnknownBearerID                 RuleFailureCode = 7
	MissingBearerID                 RuleFailureCode = 8
	MissingFlowInformation          RuleFailureCode = 9
	ResourceAllocationFailure       RuleFailureCode = 10
	UnsuccessfulQoSValidation       RuleFailureCode = 11
	IncorrectFlowInformation        RuleFailureCode = 12
	PsToCsHandover                  RuleFailureCode = 13
	TDFApplicationIdentifierError   RuleFailureCode = 14
	NoBearerFound                   RuleFailureCode = 15
	FilterRestrictions              RuleFailureCode = 16
	ANGWFailed                      RuleFailureCode = 17
	MissingRedirectServerAddress    RuleFailureCode = 18
	CMEndUserServiceDenied          RuleFailureCode = 19
	CMCreditControlNotApplicable    RuleFailureCode = 20
	CMAuthorizationRejected         RuleFailureCode = 21
	CMUserUnknown                   RuleFailureCode = 22
	CMRatingFailed                  RuleFailureCode = 23
	RoutingRuleRejection            RuleFailureCode = 24
	UnknownRoutingAccessInformation RuleFailureCode = 25
	NoNBIFOMSupport                 RuleFailureCode = 26
)
