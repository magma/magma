package s8_proxy_test

import (
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/s8_proxy"
	"magma/feg/gateway/services/s8_proxy/test_init"
	lteprotos "magma/lte/cloud/go/protos"
)

const (
	IMSI1 = "001010000000055"
)

var (
	gtpServerAddr = "127.0.0.1:0"
	gtpClientAddr = "127.0.0.1:0"
)

func TestS8ProxyClient(t *testing.T) {
	// run both s8 and pgw
	err := test_init.StartS8AndPGWService(t, gtpClientAddr, gtpServerAddr)
	if err != nil {
		t.Fatal(err)
		return
	}

	// test create session through client
	csReq := &protos.CreateSessionRequestPgw{
		Sid: &lteprotos.SubscriberID{
			Id:   IMSI1,
			Type: lteprotos.SubscriberID_IMSI,
		},
		MSISDN:               "00111",
		MEI:                  "111",
		MCC:                  "222",
		MNC:                  "333",
		RatType:              0,
		IndicationFlag:       nil,
		BearerId:             5,
		UserPlaneTeid:        0,
		S5S8Ip4UserPane:      "127.0.0.10",
		S5S8Ip6UserPane:      "",
		Apn:                  "internet.com",
		SelectionMode:        "",
		PdnType:              0,
		PdnAddressAllocation: "",
		ApnRestriction:       0,
		AmbrUp:               0,
		AmbrDown:             0,
		Uli:                  "",
	}

	r, err := s8_proxy.CreateSession(csReq)
	if err != nil {
		t.Fatalf("S8_proxy client Create Session Error: %v", err)
		return
	}
	t.Logf("Create Session: %#+v", *r)
}
