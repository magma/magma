package s8_proxy_test

import (
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/s8_proxy"
	"magma/feg/gateway/services/s8_proxy/test_init"

	"github.com/stretchr/testify/assert"
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
	// Create Session Request message
	csReq := &protos.CreateSessionRequestPgw{
		Imsi:   IMSI1,
		Msisdn: "00111",
		Mei:    "111",
		ServingNetwork: &protos.ServingNetwork{
			Mcc: "222",
			Mnc: "333",
		},
		RatType: 0,
		BearerContext: &protos.BearerContext{
			Id: 5,
			AgwUserPlaneFteid: &protos.Fteid{
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

		Apn:            "internet.com",
		SelectionMode:  "",
		ApnRestriction: 0,
		Ambr: &protos.Ambr{
			BrUl: 999,
			BrDl: 888,
		},
		Uli: &protos.UserLocationInformation{
			Lac:    1,
			Ci:     2,
			Sac:    3,
			Rac:    4,
			Tac:    5,
			Eci:    6,
			MeNbi:  7,
			EMeNbi: 8,
		},
		IndicationFlag: nil,
	}

	csRes, err := s8_proxy.CreateSession(csReq)
	if err != nil {
		t.Fatalf("S8_proxy client Create Session Error: %v", err)
		return
	}


	assert.NoError(t, err)
	assert.NotEmpty(t, csRes)

	// check fteid was received properly
	assert.NotEqual(t, 0, csRes.PgwFteidU.Teid)
	assert.NotEmpty(t, csRes.PgwFteidU.Ipv4Address)
	assert.Empty(t, csRes.PgwFteidU.Ipv6Address)
	
	t.Logf("Create Session: %#+v", *csRes)
}
