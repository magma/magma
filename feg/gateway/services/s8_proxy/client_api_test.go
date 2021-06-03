package s8_proxy_test

import (
	"testing"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/s8_proxy"
	"magma/feg/gateway/services/s8_proxy/test_init"

	"github.com/stretchr/testify/assert"
)

const (
	IMSI1          = "001010000000055"
	BEARER         = 5
	PGW_ADDRS      = "127.0.0.1:0"
	S8_PROXY_ADDRS = ":0"
	AGWTeidC       = 10
)

func TestS8ProxyClient(t *testing.T) {
	// run both s8 and pgw
	mockPgw, err := test_init.StartS8AndPGWService(t, S8_PROXY_ADDRS, PGW_ADDRS)
	if err != nil {
		t.Fatal(err)
		return
	}

	// in case pgwAddres has a 0 port, mock_pgw will chose the port. With this variable we make
	// sure we use the right address (this only happens in testing)
	actualPgwAddress := mockPgw.LocalAddr().String()

	//------------------------
	//---- Create Session ----
	csReq := getCreateSessionRequest(actualPgwAddress, AGWTeidC)

	csRes, err := s8_proxy.CreateSession(csReq)
	if err != nil {
		t.Fatalf("S8_proxy client Create Session Error: %v", err)
		return
	}

	assert.NoError(t, err)
	assert.NotEmpty(t, csRes)

	// check fteid was received properly
	assert.Equal(t, mockPgw.LastTEIDu, csRes.BearerContext.UserPlaneFteid.Teid)
	assert.NotEmpty(t, csRes.BearerContext.UserPlaneFteid.Ipv4Address)
	assert.Empty(t, csRes.BearerContext.UserPlaneFteid.Ipv6Address)

	t.Logf("Create Session: %#+v", *csRes)

	//------------------------
	//---- Delete session ----
	dsReq := &protos.DeleteSessionRequestPgw{
		PgwAddrs: actualPgwAddress,
		Imsi:     IMSI1,
		BearerId: BEARER,
		CAgwTeid: AGWTeidC,
		CPgwTeid: csRes.CPgwFteid.Teid,
		ServingNetwork: &protos.ServingNetwork{
			Mcc: "222",
			Mnc: "333",
		},
		Uli: &protos.UserLocationInformation{
			Tac: 5,
			Eci: 6,
		},
	}
	_, err = s8_proxy.DeleteSession(dsReq)
	assert.NoError(t, err)

	//------------------------
	//---- Echo Request ----
	eReq := &protos.EchoRequest{PgwAddrs: actualPgwAddress}
	_, err = s8_proxy.SendEcho(eReq)
	assert.NoError(t, err)
}

func getCreateSessionRequest(pgwAddrs string, cPgwTeid uint32) *protos.CreateSessionRequestPgw {
	_, offset := time.Now().Zone()
	return &protos.CreateSessionRequestPgw{
		PgwAddrs: pgwAddrs,
		Imsi:     IMSI1,
		Msisdn:   "00111",
		Mei:      "111",
		CAgwTeid: cPgwTeid,
		ServingNetwork: &protos.ServingNetwork{
			Mcc: "222",
			Mnc: "333",
		},
		RatType: protos.RATType_EUTRAN,
		BearerContext: &protos.BearerContext{
			Id: BEARER,
			UserPlaneFteid: &protos.Fteid{
				Ipv4Address: "127.0.0.10",
				Ipv6Address: "",
				Teid:        11,
			},
			Qos: &protos.QosInformation{
				Pci:                     0,
				PriorityLevel:           0,
				PreemptionCapability:    0,
				PreemptionVulnerability: 0,
				Qci:                     0,
				Gbr: &protos.Ambr{
					BrUl: 123,
					BrDl: 234,
				},
				Mbr: &protos.Ambr{
					BrUl: 567,
					BrDl: 890,
				},
			},
		},
		PdnType: protos.PDNType_IPV4,
		Paa: &protos.PdnAddressAllocation{
			Ipv4Address: "10.0.0.10",
			Ipv6Address: "",
			Ipv6Prefix:  0,
		},

		Apn:           "internet.com",
		SelectionMode: protos.SelectionModeType_APN_provided_subscription_verified,
		Ambr: &protos.Ambr{
			BrUl: 999,
			BrDl: 888,
		},
		Uli: &protos.UserLocationInformation{
			Tac: 5,
			Eci: 6,
		},
		IndicationFlag: nil,
		TimeZone: &protos.TimeZone{
			DeltaSeconds:       int32(offset),
			DaylightSavingTime: 0,
		},
	}

}
