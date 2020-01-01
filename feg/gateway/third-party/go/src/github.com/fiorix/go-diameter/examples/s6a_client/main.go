// Copyright 2013-2018 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter S6A client example.
package main

/* NOTE: If you are using OAI HSS for testing - Update oai_db database:
*
* -- Update Tables oai_db.mmeidentity
* mysql -u <user> -p
# 	enter your password when prompted
* mysql > use oai_db
* mysql > show tables;  # show all tables
* mysql > select * from mmeidentity; # show all entries in mmeidentity
*  -- Insert your mme identity if not present
* mysql > INSERT INTO mmeidentity (`idmmeidentity`,`mmehost`,`mmerealm`,`UE-reachability`) VALUES ('7','magma.openair4G.eur','openair4G.eur','0');
*  -- Note- Here “7 “is assumed to be one number which should not present be  in table .
*
*  -- Insert user in users table
*
* mysql > INSERT INTO users (`imsi`, `msisdn`, `imei`, `imei_sv`, `ms_ps_status`, `rau_tau_timer`, `ue_ambr_ul`,
* 	`ue_ambr_dl`, `access_restriction`, `mme_cap`, `mmeidentity_idmmeidentity`, `key`, `RFSP-Index`, `urrp_mme`,
* 	`sqn`, `rand`, `OPc`) VALUES ('001010000000001', '33638060010', NULL, NULL, 'PURGED', '120', '50000000',
* 	'100000000', '47', '0000000000', '3', 0x8BAF473F2F8FD09487CCCBD7097C6862, '1', '0', 0,
* 	0x00000000000000000000000000000000, '');
*
* mysql> INSERT INTO pdn (`id`, `apn`, `pdn_type`, `pdn_ipv4`, `pdn_ipv6`, `aggregate_ambr_ul`, `aggregate_ambr_dl`,
* 	`pgw_id`, `users_imsi`, `qci`, `priority_level`,`pre_emp_cap`,`pre_emp_vul`, `LIPA-Permissions`)
* 	VALUES ('60', 'oai.ipv4','IPV4', '0.0.0.0', '0:0:0:0:0:0:0:0', '50000000', '100000000', '3', '001010000000001',
* 	'9', '15', 'DISABLED', 'ENABLED', 'LIPA-ONLY');
*/

import (
	"errors"
	"flag"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/fiorix/go-diameter/v4/diam/sm"
	"github.com/fiorix/go-diameter/v4/diam/sm/smpeer"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	addr            = flag.String("addr", "192.168.60.145:3868", "address in form of ip:port to connect to")
	host            = flag.String("diam_host", "magma-oai.openair4G.eur", "diameter identity host")
	realm           = flag.String("diam_realm", "openair4G.eur", "diameter identity realm")
	networkType     = flag.String("network_type", "sctp", "protocol type tcp/sctp/tcp4/tcp6/sctp4/sctp6")
	retries         = flag.Uint("retries", 3, "Maximum number of retransmits")
	watchdog        = flag.Uint("watchdog", 5, "Diameter watchdog interval in seconds. 0 to disable watchdog.")
	vendorID        = flag.Uint("vendor", 10415, "Vendor ID")
	appID           = flag.Uint("app", 16777251, "AuthApplicationID")
	ueIMSI          = flag.String("imsi", "001010000000001", "Client (UE) IMSI")
	plmnID          = flag.String("plmnid", "\x00\xF1\x10", "Client (UE) PLMN ID")
	vectors         = flag.Uint("vectors", 3, "Number Of Requested Auth Vectors")
	completionSleep = flag.Uint("sleep", 10, "After Completion Sleep Time (seconds)")
)

func main() {

	flag.Parse()
	if len(*addr) == 0 {
		flag.Usage()
	}

	cfg := &sm.Settings{
		OriginHost:       datatype.DiameterIdentity(*host),
		OriginRealm:      datatype.DiameterIdentity(*realm),
		VendorID:         datatype.Unsigned32(*vendorID),
		ProductName:      "go-diameter-s6a",
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
		HostIPAddresses: []datatype.Address{
			datatype.Address(net.ParseIP("127.0.0.1")),
		},
	}

	// Create the state machine (it's a diam.ServeMux) and client.
	mux := sm.New(cfg)

	cli := &sm.Client{
		Dict:               dict.Default,
		Handler:            mux,
		MaxRetransmits:     *retries,
		RetransmitInterval: time.Second,
		EnableWatchdog:     *watchdog != 0,
		WatchdogInterval:   time.Duration(*watchdog) * time.Second,
		SupportedVendorID: []*diam.AVP{
			diam.NewAVP(avp.SupportedVendorID, avp.Mbit, 0, datatype.Unsigned32(*vendorID)),
		},
		VendorSpecificApplicationID: []*diam.AVP{
			diam.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(*appID)),
					diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(*vendorID)),
				},
			}),
		},
	}

	// Set message handlers.
	done := make(chan struct{}, 1000)
	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.AuthenticationInformation, Request: false},
		handleAuthenticationInformationAnswer(done))

	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.UpdateLocation, Request: false},
		handleUpdateLocationAnswer(done))

	// Catch All
	mux.HandleIdx(diam.ALL_CMD_INDEX, handleAll())

	// Print error reports.
	go printErrors(mux.ErrorReports())

	conn, err := cli.DialNetwork(*networkType, *addr)
	if err != nil {
		log.Fatal(err)
	}
	err = sendAIR(conn, cfg)
	if err != nil {
		log.Fatal(err)
	}
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		log.Fatal("Authentication Information timeout")
	}
	err = sendULR(conn, cfg)
	if err != nil {
		log.Fatal(err)
	}
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		log.Fatal("Update Location timeout")
	}

	// Sleep after completion to observe DWR/As going in the background
	time.Sleep(time.Duration(*completionSleep) * time.Second)
}

func printErrors(ec <-chan *diam.ErrorReport) {
	for err := range ec {
		log.Println(err)
	}
}

// Create & send Authentication-Information Request
func sendAIR(c diam.Conn, cfg *sm.Settings) error {
	meta, ok := smpeer.FromContext(c.Context())
	if !ok {
		return errors.New("peer metadata unavailable")
	}
	sid := "session;" + strconv.Itoa(int(rand.Uint32()))
	m := diam.NewRequest(diam.AuthenticationInformation, diam.TGPP_S6A_APP_ID, dict.Default)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, cfg.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, cfg.OriginRealm)
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, meta.OriginRealm)
	m.NewAVP(avp.DestinationHost, avp.Mbit, 0, meta.OriginHost)
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(*ueIMSI))
	m.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(0))
	m.NewAVP(avp.VisitedPLMNID, avp.Vbit|avp.Mbit, uint32(*vendorID), datatype.OctetString(*plmnID))
	m.NewAVP(avp.RequestedEUTRANAuthenticationInfo, avp.Vbit|avp.Mbit, uint32(*vendorID), &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(
				avp.NumberOfRequestedVectors, avp.Vbit|avp.Mbit, uint32(*vendorID), datatype.Unsigned32(*vectors)),
			diam.NewAVP(
				avp.ImmediateResponsePreferred, avp.Vbit|avp.Mbit, uint32(*vendorID), datatype.Unsigned32(0)),
		},
	})
	log.Printf("\nSending AIR to %s\n%s\n", c.RemoteAddr(), m)
	_, err := m.WriteTo(c)
	return err
}

const ULR_FLAGS = 1<<1 | 1<<5

type EUtranVector struct {
	RAND  datatype.OctetString `avp:"RAND"`
	XRES  datatype.OctetString `avp:"XRES"`
	AUTN  datatype.OctetString `avp:"AUTN"`
	KASME datatype.OctetString `avp:"KASME"`
}

type ExperimentalResult struct {
	ExperimentalResultCode datatype.Unsigned32 `avp:"Experimental-Result-Code"`
}

type AuthenticationInfo struct {
	EUtranVector EUtranVector `avp:"E-UTRAN-Vector"`
}

type AIA struct {
	SessionID          datatype.UTF8String       `avp:"Session-Id"`
	ResultCode         datatype.Unsigned32       `avp:"Result-Code"`
	OriginHost         datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm        datatype.DiameterIdentity `avp:"Origin-Realm"`
	AuthSessionState   datatype.UTF8String       `avp:"Auth-Session-State"`
	ExperimentalResult ExperimentalResult        `avp:"Experimental-Result"`
	AIs                []AuthenticationInfo      `avp:"Authentication-Info"`
}

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
	ContextIdentifier                     uint32           `avp:"Context-Identifier"`
	AllAPNConfigurationsIncludedIndicator int32            `avp:"All-APN-Configurations-Included-Indicator"`
	APNConfiguration                      APNConfiguration `avp:"APN-Configuration"`
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

// Create & send Update-Location Request
func sendULR(c diam.Conn, cfg *sm.Settings) error {
	meta, ok := smpeer.FromContext(c.Context())
	if !ok {
		return errors.New("peer metadata unavailable")
	}
	sid := "session;" + strconv.Itoa(int(rand.Uint32()))
	m := diam.NewRequest(diam.UpdateLocation, diam.TGPP_S6A_APP_ID, dict.Default)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, cfg.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, cfg.OriginRealm)
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, meta.OriginRealm)
	m.NewAVP(avp.DestinationHost, avp.Mbit, 0, meta.OriginHost)
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(*ueIMSI))
	m.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(0))
	m.NewAVP(avp.RATType, avp.Mbit, uint32(*vendorID), datatype.Enumerated(1004))
	m.NewAVP(avp.ULRFlags, avp.Vbit|avp.Mbit, uint32(*vendorID), datatype.Unsigned32(ULR_FLAGS))
	m.NewAVP(avp.VisitedPLMNID, avp.Vbit|avp.Mbit, uint32(*vendorID), datatype.OctetString(*plmnID))
	log.Printf("\nSending ULR to %s\n%s\n", c.RemoteAddr(), m)
	_, err := m.WriteTo(c)
	return err
}

func handleAuthenticationInformationAnswer(done chan struct{}) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		log.Printf("Received Authentication-Information Answer from %s\n%s\n", c.RemoteAddr(), m)
		var aia AIA
		err := m.Unmarshal(&aia)
		if err != nil {
			log.Printf("AIA Unmarshal failed: %s", err)
		} else {
			log.Printf("Unmarshaled Authentication-Information Answer:\n%#+v\n", aia)
		}
		ok := struct{}{}
		done <- ok
	}
}

func handleUpdateLocationAnswer(done chan struct{}) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		log.Printf("Received Update-Location Answer from %s\n%s\n", c.RemoteAddr(), m)
		var ula ULA
		err := m.Unmarshal(&ula)
		if err != nil {
			log.Printf("ULA Unmarshal failed: %s", err)
		} else {
			log.Printf("Unmarshaled UL Answer:\n%#+v\n", ula)
		}
		ok := struct{}{}
		done <- ok
	}
}

func handleAll() diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		log.Printf("Received Meesage From %s\n%s\n", c.RemoteAddr(), m)
	}
}
