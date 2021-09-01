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

// package servce implements S6a GRPC proxy service which sends AIR, ULR messages over diameter connection,
// waits (blocks) for diameter's AIAs, ULAs & returns their RPC representation
package servicers

import (
	"github.com/fiorix/go-diameter/v4/diam/datatype"
)

const (
	// 3GPP 29.273 5.2.3.6
	RadioAccessTechnologyType_EUTRAN = 1004
)

// Definitions for AIA, see sample below:
//
//Authentication-Information-Answer (AIA)
//{Code:318,Flags:0x40,Version:0x1,Length:556,ApplicationId:16777251,HopByHopId:0x9105bf89,EndToEndId:0x16c85bed}
//	Session-Id {Code:263,Flags:0x40,Length:28,VendorId:0,Value:UTF8String{session;3420619691},Padding:2}
//	Authentication-Info {Code:1413,Flags:0xc0,Length:144,VendorId:10415,Value:Grouped{
//		E-UTRAN-Vector {Code:1414,Flags:0xc0,Length:132,VendorId:10415,Value:Grouped{
//			RAND {Code:1447,Flags:0xc0,Length:28,VendorId:10415,Value:OctetString{0xf122047125e8372054d2b31643878866},Padding:0},
//			XRES {Code:1448,Flags:0xc0,Length:20,VendorId:10415,Value:OctetString{0x2c6e30243a103f0e},Padding:0},
//			AUTN {Code:1449,Flags:0xc0,Length:28,VendorId:10415,Value:OctetString{0xe73b5091a37080007d1726fd84830ecc},Padding:0},
//			KASME {Code:1450,Flags:0xc0,Length:44,VendorId:10415,Value:OctetString{0x08083ce5b62fdbe542ba0a19c415411cfaf1db35b8832b1f8a9c7cb525824c21},Padding:0},
//		}}
//	}}
//	Authentication-Info {Code:1413,Flags:0xc0,Length:144,VendorId:10415,Value:Grouped{
//		E-UTRAN-Vector {Code:1414,Flags:0xc0,Length:132,VendorId:10415,Value:Grouped{
//			RAND {Code:1447,Flags:0xc0,Length:28,VendorId:10415,Value:OctetString{0x12c7eb54f10c4007f65e14315545ed25},Padding:0},
//			XRES {Code:1448,Flags:0xc0,Length:20,VendorId:10415,Value:OctetString{0x22aeae2a4713ee62},Padding:0},
//			AUTN {Code:1449,Flags:0xc0,Length:28,VendorId:10415,Value:OctetString{0xfb97e19addee80002be44eee2df02059},Padding:0},
//			KASME {Code:1450,Flags:0xc0,Length:44,VendorId:10415,Value:OctetString{0x342a6173dda12c7902d2048d70fd83806a5e66b6fced874ccddfa106c9d4e03f},Padding:0},
//		}}
//	}}
//	Authentication-Info {Code:1413,Flags:0xc0,Length:144,VendorId:10415,Value:Grouped{
//		E-UTRAN-Vector {Code:1414,Flags:0xc0,Length:132,VendorId:10415,Value:Grouped{
//			RAND {Code:1447,Flags:0xc0,Length:28,VendorId:10415,Value:OctetString{0x23ea3e0ebd90b06b87e07554ac65d85d},Padding:0},
//			XRES {Code:1448,Flags:0xc0,Length:20,VendorId:10415,Value:OctetString{0x4c4f47cf85b84db9},Padding:0},
//			AUTN {Code:1449,Flags:0xc0,Length:28,VendorId:10415,Value:OctetString{0xcc49b7b25775800011079582097b2e48},Padding:0},
//			KASME {Code:1450,Flags:0xc0,Length:44,VendorId:10415,Value:OctetString{0x671c99d4aeca35a90c4bb26028df37a5151322c837c86189635da3ac24979d43},Padding:0},
//		}}
//	}}
//	Auth-Session-State {Code:277,Flags:0x40,Length:12,VendorId:0,Value:Enumerated{0}}
//	Origin-Host {Code:264,Flags:0x40,Length:28,VendorId:0,Value:DiameterIdentity{hss.openair4G.eur},Padding:3}
//	Origin-Realm {Code:296,Flags:0x40,Length:24,VendorId:0,Value:DiameterIdentity{openair4G.eur},Padding:3}
//	Result-Code {Code:268,Flags:0x40,Length:12,VendorId:0,Value:Unsigned32{2001}}
//

type EUtranVector struct {
	ItemNumber datatype.Unsigned32  `avp:"Item-Number"`
	RAND       datatype.OctetString `avp:"RAND"`
	XRES       datatype.OctetString `avp:"XRES"`
	AUTN       datatype.OctetString `avp:"AUTN"`
	KASME      datatype.OctetString `avp:"KASME"`
}

type UtranVector struct {
	ItemNumber datatype.Unsigned32  `avp:"Item-Number"`
	RAND       datatype.OctetString `avp:"RAND"`
	XRES       datatype.OctetString `avp:"XRES"`
	AUTN       datatype.OctetString `avp:"AUTN"`
	CK         datatype.OctetString `avp:"Confidentiality-Key"`
	IK         datatype.OctetString `avp:"Integrity-Key"`
}

type GeranVector struct {
	ItemNumber datatype.Unsigned32  `avp:"Item-Number"`
	RAND       datatype.OctetString `avp:"RAND"`
	SRES       datatype.OctetString `avp:"SRES"`
	Kc         datatype.OctetString `avp:"Kc"`
}

type ExperimentalResult struct {
	VendorId               uint32 `avp:"Vendor-Id"`
	ExperimentalResultCode uint32 `avp:"Experimental-Result-Code"`
}

type AuthenticationInfo struct {
	EUtranVectors []EUtranVector `avp:"E-UTRAN-Vector"`
	UtranVectors  []UtranVector  `avp:"UTRAN-Vector"`
	GeranVectors  []GeranVector  `avp:"GERAN-Vector"`
}

type AIA struct {
	SessionID          string                    `avp:"Session-Id"`
	ResultCode         uint32                    `avp:"Result-Code"`
	OriginHost         datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm        datatype.DiameterIdentity `avp:"Origin-Realm"`
	AuthSessionState   int32                     `avp:"Auth-Session-State"`
	ExperimentalResult ExperimentalResult        `avp:"Experimental-Result"`
	AI                 AuthenticationInfo        `avp:"Authentication-Info"`
}

// Definitions for ULA, see sample below:
//
//Update-Location-Answer (ULA)
//{Code:316,Flags:0x40,Version:0x1,Length:516,ApplicationId:16777251,HopByHopId:0x22910d0a,EndToEndId:0x8d330652}
//	Session-Id {Code:263,Flags:0x40,Length:24,VendorId:0,Value:UTF8String{session;89988919},Padding:0}
//	ULA-Flags {Code:1406,Flags:0xc0,Length:16,VendorId:10415,Value:Unsigned32{1}}
//	Subscription-Data {Code:1400,Flags:0xc0,Length:380,VendorId:10415,Value:Grouped{
//		MSISDN {Code:701,Flags:0xc0,Length:20,VendorId:10415,Value:OctetString{0x33638060010f},Padding:2},
//		Access-Restriction-Data {Code:1426,Flags:0xc0,Length:16,VendorId:10415,Value:Unsigned32{47}},
//		Subscriber-Status {Code:1424,Flags:0xc0,Length:16,VendorId:10415,Value:Enumerated{0}},
//		Network-Access-Mode {Code:1417,Flags:0xc0,Length:16,VendorId:10415,Value:Enumerated{2}},
//		AMBR {Code:1435,Flags:0xc0,Length:44,VendorId:10415,Value:Grouped{
//			Max-Requested-Bandwidth-UL {Code:516,Flags:0xc0,Length:16,VendorId:10415,Value:Unsigned32{50000000}},
//			Max-Requested-Bandwidth-DL {Code:515,Flags:0xc0,Length:16,VendorId:10415,Value:Unsigned32{100000000}},
//		}}
//		APN-Configuration-Profile {Code:1429,Flags:0xc0,Length:240,VendorId:10415,Value:Grouped{
//			Context-Identifier {Code:1423,Flags:0xc0,Length:16,VendorId:10415,Value:Unsigned32{0}},
//			All-APN-Configurations-Included-Indicator {Code:1428,Flags:0xc0,Length:16,VendorId:10415,Value:Enumerated{0}},
//			APN-Configuration {Code:1430,Flags:0xc0,Length:196,VendorId:10415,Value:Grouped{
//				Context-Identifier {Code:1423,Flags:0xc0,Length:16,VendorId:10415,Value:Unsigned32{0}},
//				PDN-Type {Code:1456,Flags:0xc0,Length:16,VendorId:10415,Value:Enumerated{0}},
//				Service-Selection {Code:493,Flags:0xc0,Length:20,VendorId:10415,Value:UTF8String{oai.ipv4},Padding:0},
//				EPS-Subscribed-QoS-Profile {Code:1431,Flags:0xc0,Length:88,VendorId:10415,Value:Grouped{
//					QoS-Class-Identifier {Code:1028,Flags:0xc0,Length:16,VendorId:10415,Value:Enumerated{9}},
//					Allocation-Retention-Priority {Code:1034,Flags:0x80,Length:60,VendorId:10415,Value:Grouped{
//						Priority-Level {Code:1046,Flags:0x80,Length:16,VendorId:10415,Value:Unsigned32{15}},
//						Pre-emption-Capability {Code:1047,Flags:0x80,Length:16,VendorId:10415,Value:Enumerated{1}},
//						Pre-emption-Vulnerability {Code:1048,Flags:0x80,Length:16,VendorId:10415,Value:Enumerated{0}},
//					}}
//				}},
//              TGPP-Charging-Characteristics {Code: 13,Flags:0x80,VendorId:10415,Value:UTF8String{12}},
//				AMBR {Code:1435,Flags:0xc0,Length:44,VendorId:10415,Value:Grouped{
//					Max-Requested-Bandwidth-UL {Code:516,Flags:0xc0,Length:16,VendorId:10415,Value:Unsigned32{50000000}},
//					Max-Requested-Bandwidth-DL {Code:515,Flags:0xc0,Length:16,VendorId:10415,Value:Unsigned32{100000000}},
//				}}
//			}}
//		}}
//		Subscribed-Periodic-RAU-TAU-Timer {Code:1619,Flags:0x80,Length:16,VendorId:10415,Value:Unsigned32{120}},
//	}}
//	Auth-Session-State {Code:277,Flags:0x40,Length:12,VendorId:0,Value:Enumerated{0}}
//	Origin-Host {Code:264,Flags:0x40,Length:28,VendorId:0,Value:DiameterIdentity{hss.openair4G.eur},Padding:3}
//	Origin-Realm {Code:296,Flags:0x40,Length:24,VendorId:0,Value:DiameterIdentity{openair4G.eur},Padding:3}
//	Result-Code {Code:268,Flags:0x40,Length:12,VendorId:0,Value:Unsigned32{2001}}
//
type AMBR struct {
	MaxRequestedBandwidthUL uint32 `avp:"Max-Requested-Bandwidth-UL"`
	MaxRequestedBandwidthDL uint32 `avp:"Max-Requested-Bandwidth-DL"`
	ExtendMaxRequestedBwUL  uint32 `avp:"Extended-Max-Requested-BW-UL"`
	ExtendMaxRequestedBwDL  uint32 `avp:"Extended-Max-Requested-BW-DL"`
}

type AllocationRetentionPriority struct {
	PriorityLevel           uint32 `avp:"Priority-Level"`
	PreemptionCapability    int32  `avp:"Pre-emption-Capability"`
	PreemptionVulnerability int32  `avp:"Pre-emption-Vulnerability"`
}

type EPSSubscribedQoSProfile struct {
	QoSClassIdentifier          int32                       `avp:"QoS-Class-Identifier"`
	AllocationRetentionPriority AllocationRetentionPriority `avp:"Allocation-Retention-Priority"`
}

type APNConfiguration struct {
	ContextIdentifier           uint32                  `avp:"Context-Identifier"`
	PDNType                     uint32                  `avp:"PDN-Type"`
	ServiceSelection            string                  `avp:"Service-Selection"`
	EPSSubscribedQoSProfile     EPSSubscribedQoSProfile `avp:"EPS-Subscribed-QoS-Profile"`
	AMBR                        AMBR                    `avp:"AMBR"`
	TgppChargingCharacteristics string                  `avp:"TGPP-Charging-Characteristics"`
}

type APNConfigurationProfile struct {
	ContextIdentifier                     uint32             `avp:"Context-Identifier"`
	AllAPNConfigurationsIncludedIndicator int32              `avp:"All-APN-Configurations-Included-Indicator"`
	APNConfigs                            []APNConfiguration `avp:"APN-Configuration"`
}

type SubscriptionData struct {
	MSISDN                        datatype.OctetString    `avp:"MSISDN"`
	AccessRestrictionData         uint32                  `avp:"Access-Restriction-Data"`
	SubscriberStatus              int32                   `avp:"Subscriber-Status"`
	NetworkAccessMode             int32                   `avp:"Network-Access-Mode"`
	AMBR                          AMBR                    `avp:"AMBR"`
	APNConfigurationProfile       APNConfigurationProfile `avp:"APN-Configuration-Profile"`
	SubscribedPeriodicRauTauTimer uint32                  `avp:"Subscribed-Periodic-RAU-TAU-Timer"`
	TgppChargingCharacteristics   string                  `avp:"TGPP-Charging-Characteristics"`
	RegionalSubscriptionZoneCode  []datatype.OctetString  `avp:"Regional-Subscription-Zone-Code"`
}

type ULA struct {
	SessionID          string                    `avp:"Session-Id"`
	ULAFlags           uint32                    `avp:"ULA-Flags"`
	SubscriptionData   SubscriptionData          `avp:"Subscription-Data"`
	AuthSessionState   int32                     `avp:"Auth-Session-State"`
	ResultCode         uint32                    `avp:"Result-Code"`
	OriginHost         datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm        datatype.DiameterIdentity `avp:"Origin-Realm"`
	ExperimentalResult ExperimentalResult        `avp:"Experimental-Result"`
	SupportedFeatures  []SupportedFeatures       `avp:"Supported-Features"`
}

type CLR struct {
	SessionID        string                    `avp:"Session-Id"`
	AuthSessionState int32                     `avp:"Auth-Session-State"`
	OriginHost       datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm      datatype.DiameterIdentity `avp:"Origin-Realm"`
	CancellationType int32                     `avp:"Cancellation-Type"`
	DestinationHost  datatype.DiameterIdentity `avp:"Destination-Host"`
	DestinationRealm datatype.DiameterIdentity `avp:"Destination-Realm"`
	UserName         string                    `avp:"User-Name"`
}

// Definitions for PU
//
// PUR is Go representation of Purge-UE-Request message
//
//	< Purge-UE-Request> ::=	< Diameter Header: 321, REQ, PXY, 16777251 >
//	< Session-Id >
//	[ DRMP ]
//	[ Vendor-Specific-Application-Id ]
//	{ Auth-Session-State }
//	{ Origin-Host }
//	{ Origin-Realm }
//	[ Destination-Host ]
//	{ Destination-Realm }
//	{ User-Name }
//	[ OC-Supported-Features ]
//	[ PUR-Flags ]
//	*[ Supported-Features ]
//	[ EPS-Location-Information ]
//	*[ AVP ]
//	*[ Proxy-Info ]
//	*[ Route-Record ]
type PUR struct {
	SessionID                   string                      `avp:"Session-Id"`
	DRMP                        uint32                      `avp:"DRMP"`
	VendorSpecificApplicationId VendorSpecificApplicationId `avp:"Vendor-Specific-Application-Id"`
	AuthSessionState            int32                       `avp:"Auth-Session-State"`
	OriginHost                  datatype.DiameterIdentity   `avp:"Origin-Host"`
	OriginRealm                 datatype.DiameterIdentity   `avp:"Origin-Realm"`
	DestinationHost             datatype.DiameterIdentity   `avp:"Destination-Host"`
	DestinationRealm            datatype.DiameterIdentity   `avp:"Destination-Realm"`
	UserName                    datatype.UTF8String         `avp:"User-Name"`
	OCSupportedFeatures         OCSupportedFeatures         `avp:"OC-Supported-Features"`
	PURFlags                    uint32                      `avp:"PUR-Flags"`
	SupportedFeatures           []SupportedFeatures         `avp:"Supported-Features"`
}

// PUA is Go representation of Purge-UE-Answer message
//
//	< Purge-UE-Answer> ::=	< Diameter Header: 321, PXY, 16777251 >
//	< Session-Id >
//	[ DRMP ]
//	[ Vendor-Specific-Application-Id ]
//	*[ Supported-Features ]
//	[ Result-Code ]
//	[ Experimental-Result ]
//	{ Auth-Session-State }
//	{ Origin-Host }
//	{ Origin-Realm }
//	[ OC-Supported-Features ]
//	[ OC-OLR ]
//	*[ Load ]
//	[ PUA-Flags ]
//	*[ AVP ]
//	[ Failed-AVP ]
//	*[ Proxy-Info ]
//	*[ Route-Record ]
//
type PUA struct {
	SessionID                   string                      `avp:"Session-Id"`
	DRMP                        uint32                      `avp:"DRMP"`
	VendorSpecificApplicationId VendorSpecificApplicationId `avp:"Vendor-Specific-Application-Id"`
	SupportedFeatures           []SupportedFeatures         `avp:"Supported-Features"`
	ResultCode                  uint32                      `avp:"Result-Code"`
	ExperimentalResult          ExperimentalResult          `avp:"Experimental-Result"`
	AuthSessionState            int32                       `avp:"Auth-Session-State"`
	OriginHost                  datatype.DiameterIdentity   `avp:"Origin-Host"`
	OriginRealm                 datatype.DiameterIdentity   `avp:"Origin-Realm"`
	OCSupportedFeatures         OCSupportedFeatures         `avp:"OC-Supported-Features"`
	OC_OLR                      OC_OLR                      `avp:"OC-OLR"`
	PUAFlags                    uint32                      `avp:"PUA-Flags"`
}

// VendorSpecificApplicationId -> Vendor-Specific-Application-Id AVP
type VendorSpecificApplicationId struct {
	VendorId          uint32 `avp:"Vendor-Id"`
	AuthApplicationId uint32 `avp:"Auth-Application-Id"`
	AcctApplicationId uint32 `avp:"Acct-Application-Id"`
}

// SupportedFeatures -> Supported-Features AVP
type SupportedFeatures struct {
	VendorId      uint32 `avp:"Vendor-Id"`
	FeatureListID uint32 `avp:"Feature-List-ID"`
	FeatureList   uint32 `avp:"Feature-List"`
}

// OCSupportedFeatures -> OC-Supported-Features AVP
type OCSupportedFeatures struct {
	OCFeatureVector uint64 `avp:"OC-Feature-Vector"`
}

// OC_OLR -> OC-OLR AVP
type OC_OLR struct {
	OCSequenceNumber      uint64 `avp:"OC-Sequence-Number"`
	OCReportType          uint32 `avp:"OC-Report-Type"`
	OCReductionPercentage uint32 `avp:"OC-Reduction-Percentage"`
	OCValidityDuration    uint32 `avp:"OC-Validity-Duration"`
}

// RSR is Go representation of Reset-Request message
//
// < Reset-Request> ::= < Diameter Header: 322, REQ, PXY, 16777251 >
//
// < Session-Id >
// [ Vendor-Specific-Application-Id ]
// { Auth-Session-State }
// { Origin-Host }
// { Origin-Realm }
// { Destination-Host }
// { Destination-Realm }
// *[ Supported-Features ]
// *[ User-Id ]
// *[ AVP ]
// *[ Proxy-Info ]
// *[ Route-Record ]
type RSR struct {
	SessionID                   string                      `avp:"Session-Id"`
	VendorSpecificApplicationId VendorSpecificApplicationId `avp:"Vendor-Specific-Application-Id"`
	AuthSessionState            int32                       `avp:"Auth-Session-State"`
	OriginHost                  datatype.DiameterIdentity   `avp:"Origin-Host"`
	OriginRealm                 datatype.DiameterIdentity   `avp:"Origin-Realm"`
	DestinationHost             datatype.DiameterIdentity   `avp:"Destination-Host"`
	DestinationRealm            datatype.DiameterIdentity   `avp:"Destination-Realm"`
	SupportedFeatures           []SupportedFeatures         `avp:"Supported-Features"`
	UserId                      []datatype.UTF8String       `avp:"User-Id"`
}

// RequestedEUTRANAuthInfo contains the information needed for authentication requests
// for E-UTRAN.
type RequestedEUTRANAuthInfo struct {
	NumVectors        datatype.Unsigned32  `avp:"Number-Of-Requested-Vectors"`
	ImmediateResponse datatype.Unsigned32  `avp:"Immediate-Response-Preferred"`
	ResyncInfo        datatype.OctetString `avp:"Re-synchronization-Info"`
}

// RequestedUtranGeranAuthInfo contains the information needed for authentication requests
// for UTRAN/GERAN.
type RequestedUtranGeranAuthInfo struct {
	NumVectors        datatype.Unsigned32  `avp:"Number-Of-Requested-Vectors"`
	ImmediateResponse datatype.Unsigned32  `avp:"Immediate-Response-Preferred"`
	ResyncInfo        datatype.OctetString `avp:"Re-synchronization-Info"`
}

// AIR encapsulates all of the information contained in an authentication information request.
// This information is sent to fetch data in order to authenticate a subscriber.
type AIR struct {
	SessionID                   datatype.UTF8String       `avp:"Session-Id"`
	OriginHost                  datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm                 datatype.DiameterIdentity `avp:"Origin-Realm"`
	AuthSessionState            datatype.UTF8String       `avp:"Auth-Session-State"`
	UserName                    string                    `avp:"User-Name"`
	VisitedPLMNID               datatype.Unsigned32       `avp:"Visited-PLMN-Id"`
	RequestedEUTRANAuthInfo     RequestedEUTRANAuthInfo   `avp:"Requested-EUTRAN-Authentication-Info"`
	RequestedUtranGeranAuthInfo RequestedEUTRANAuthInfo   `avp:"Requested-UTRAN-GERAN-Authentication-Info"`
}

// ULR is an update location request. It is used to update location information in the HSS.
type ULR struct {
	SessionID         datatype.UTF8String       `avp:"Session-Id"`
	OriginHost        datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm       datatype.DiameterIdentity `avp:"Origin-Realm"`
	AuthSessionState  datatype.Unsigned32       `avp:"Auth-Session-State"`
	UserName          datatype.UTF8String       `avp:"User-Name"`
	VisitedPLMNID     datatype.Unsigned32       `avp:"Visited-PLMN-Id"`
	RATType           datatype.Unsigned32       `avp:"RAT-Type"`
	ULRFlags          datatype.Unsigned32       `avp:"ULR-Flags"`
	SupportedFeatures []SupportedFeatures       `avp:"Supported-Features"`
}
