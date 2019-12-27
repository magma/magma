// package servce implements S6a GRPC proxy service which sends AIR, ULR messages over diameter connection,
// waits (blocks) for diameter's AIAs, ULAs & returns their RPC representation
package service

import (
	"net"
	"sync"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/fiorix/go-diameter/v4/diam/sm"
	"github.com/fiorix/go-diameter/v4/examples/s6a_proxy/protos"
	"golang.org/x/net/context"
)

const (
	ULR_RAT_TYPE     = 1004
	ULR_FLAGS        = 1<<1 | 1<<5
	TIMEOUT_SECONDS  = 10
	MAX_DIAM_RETRIES = 1
	PRODUCT_NAME     = "s6a_proxy"
)

type s6aProxy struct {
	mu         sync.RWMutex
	cfg        *S6aProxyConfig
	smClient   *sm.Client
	conn       diam.Conn
	sessionsMu sync.Mutex
	sessions   map[string]chan interface{}

	// test related fields
	airSendLocks [diam.MaxOutboundSCTPStreams]sync.Mutex
}

func NewS6aProxy(cfg *S6aProxyConfig) (*s6aProxy, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	cfg = cfg.CloneWithDefaults()

	mux := sm.New(&sm.Settings{
		OriginHost:       datatype.DiameterIdentity(cfg.Host),
		OriginRealm:      datatype.DiameterIdentity(cfg.Realm),
		VendorID:         datatype.Unsigned32(VENDOR_3GPP),
		ProductName:      PRODUCT_NAME,
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
		HostIPAddresses: []datatype.Address{
			datatype.Address(net.ParseIP("127.0.0.1")),
		},
	})

	mux.HandleFunc("ALL", func(diam.Conn, *diam.Message) {}) // Catch all.

	proxy := &s6aProxy{
		cfg: cfg,
		smClient: &sm.Client{
			Dict:               dict.Default,
			Handler:            mux,
			MaxRetransmits:     cfg.Retransmits,
			RetransmitInterval: time.Second * 3,
			EnableWatchdog:     true,
			// WatchdogInterval:   time.Second * time.Duration(cfg.WatchdogInterval),
			WatchdogInterval: time.Millisecond * 20,
			WatchdogStream:   diam.MaxOutboundSCTPStreams - 1,
			SupportedVendorID: []*diam.AVP{
				diam.NewAVP(avp.SupportedVendorID, avp.Mbit, 0, datatype.Unsigned32(VENDOR_3GPP)),
			},
			VendorSpecificApplicationID: []*diam.AVP{
				diam.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(diam.TGPP_S6A_APP_ID)),
						diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(diam.TGPP_S6A_APP_ID)),
					},
				}),
			},
		},
		conn:     nil,
		sessions: make(map[string]chan interface{}),
	}
	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.AuthenticationInformation, Request: false},
		handleAIA(proxy))

	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_S6A_APP_ID, Code: diam.UpdateLocation, Request: false},
		handleULA(proxy))

	return proxy, nil
}

// S6AProxyServer implementation
//
// AuthenticationInformation sends AIR over diameter connection,
// waits (blocks) for AIA & returns its RPC representation
func (s *s6aProxy) AuthenticationInformation(
	ctx context.Context, req *protos.AuthenticationInformationRequest) (*protos.AuthenticationInformationAnswer, error,
) {
	return s.AuthenticationInformationImpl(req)
}

// UpdateLocation sends ULR (Code 316) over diameter connection,
// waits (blocks) for ULA & returns its RPC representation
func (s *s6aProxy) UpdateLocation(
	ctx context.Context, req *protos.UpdateLocationRequest) (*protos.UpdateLocationAnswer, error,
) {
	return s.UpdateLocationImpl(req)
}

// CancelLocation sends CLR (Code 317) over diameter connection,
// Not implemented for now
func (s *s6aProxy) CancelLocation(
	ctx context.Context, req *protos.CancelLocationRequest) (*protos.CancelLocationAnswer, error,
) {
	panic("Not implemented")
}

// PurgeUE sends PUR (Code 321) over diameter connection,
// Not implemented
func (s *s6aProxy) PurgeUE(ctx context.Context, req *protos.PurgeUERequest) (*protos.PurgeUEAnswer, error) {
	panic("Not implemented")
}

// Reset (Code 322) over diameter connection,
// Not implemented
func (s *s6aProxy) Reset(ctx context.Context, in *protos.ResetRequest) (*protos.ResetAnswer, error) {
	panic("Not implemented")
}
