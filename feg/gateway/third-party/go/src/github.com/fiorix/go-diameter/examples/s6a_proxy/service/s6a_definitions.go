// package servce implements S6a GRPC proxy service which sends AIR, ULR messages over diameter connection,
// waits (blocks) for diameter's AIAs, ULAs & returns their RPC representation
package service

import (
	"github.com/fiorix/go-diameter/v4/diam/datatype"
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

const VENDOR_3GPP = uint32(10415)

type EUtranVector struct {
	RAND  datatype.OctetString `avp:"RAND"`
	XRES  datatype.OctetString `avp:"XRES"`
	AUTN  datatype.OctetString `avp:"AUTN"`
	KASME datatype.OctetString `avp:"KASME"`
}

type ExperimentalResult struct {
	VendorId               uint32 `avp:"Vendor-Id"`
	ExperimentalResultCode uint32 `avp:"Experimental-Result-Code"`
}

type AuthenticationInfo struct {
	EUtranVector EUtranVector `avp:"E-UTRAN-Vector"`
}

type AIA struct {
	SessionID          string                    `avp:"Session-Id"`
	ResultCode         uint32                    `avp:"Result-Code"`
	OriginHost         datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm        datatype.DiameterIdentity `avp:"Origin-Realm"`
	AuthSessionState   int32                     `avp:"Auth-Session-State"`
	ExperimentalResult ExperimentalResult        `avp:"Experimental-Result"`
	AIs                []AuthenticationInfo      `avp:"Authentication-Info"`
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
//				}}
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
	ContextIdentifier       uint32                  `avp:"Context-Identifier"`
	PDNType                 int32                   `avp:"PDN-Type"`
	ServiceSelection        string                  `avp:"Service-Selection"`
	EPSSubscribedQoSProfile EPSSubscribedQoSProfile `avp:"EPS-Subscribed-QoS-Profile"`
	AMBR                    AMBR                    `avp:"AMBR"`
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
}
