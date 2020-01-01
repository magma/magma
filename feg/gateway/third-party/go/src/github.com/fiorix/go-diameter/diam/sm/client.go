// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sm

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

var (
	// ErrMissingStateMachine is returned by Dial or DialTLS when
	// the Client does not have a valid StateMachine set.
	ErrMissingStateMachine = errors.New("client state machine is nil")

	// ErrHandshakeTimeout is returned by Dial or DialTLS when the
	// client does not receive a handshake answer from the server.
	//
	// If the client is configured to retransmit messages, the
	// handshake timeout only occurs after all retransmits are
	// attempted and none has an aswer.
	ErrHandshakeTimeout = errors.New("handshake timeout (no response)")
)

// A Client is a diameter client that automatically performs a handshake
// with the server after the connection is established.
//
// It sends a Capabilities-Exchange-Request with the AVPs defined in it,
// and expects a Capabilities-Exchange-Answer with a success (2001) result
// code. If enabled, the client will send Device-Watchdog-Request messages
// in background until the connection is terminated.
//
// By default, retransmission and watchdog are disabled. Retransmission is
// enabled by setting MaxRetransmits to a number greater than zero, and
// watchdog is enabled by setting EnableWatchdog to true.
//
// A custom message handler for Device-Watchdog-Answer (DWA) can be registered.
// However, that will be overwritten if watchdog is enabled.
type Client struct {
	Dict                        *dict.Parser  // Dictionary parser (uses dict.Default if unset)
	Handler                     *StateMachine // Message handler
	MaxRetransmits              uint          // Max number of retransmissions before aborting
	RetransmitInterval          time.Duration // Interval between retransmissions (default 1s)
	EnableWatchdog              bool          // Enable automatic DWR
	WatchdogInterval            time.Duration // Interval between DWRs (default 5s)
	WatchdogStream              uint          // Stream to send DWR on (for multistreaming protocols), default is 0
	SupportedVendorID           []*diam.AVP   // Supported vendor ID
	AcctApplicationID           []*diam.AVP   // Acct applications
	AuthApplicationID           []*diam.AVP   // Auth applications
	VendorSpecificApplicationID []*diam.AVP   // Vendor specific applications
}

// Dial calls the address set as ip:port, performs a handshake and optionally
// start a watchdog goroutine in background.
func (cli *Client) Dial(addr string) (diam.Conn, error) {
	return cli.DialExt("tcp", addr, 0, nil)
}

// DialNetwork calls the network address set as ip:port, performs a handshake and optionally
// start a watchdog goroutine in background.
func (cli *Client) DialNetwork(network, addr string) (diam.Conn, error) {
	return cli.DialExt(network, addr, 0, nil)
}

// DialNetworkBind calls the network address set as ip:port, performs a handshake and optionally
// start a watchdog goroutine in background.
func (cli *Client) DialNetworkBind(network, laddr, raddr string) (diam.Conn, error) {
	return cli.dial(func() (diam.Conn, error) {
		return diam.DialNetworkBind(network, laddr, raddr, cli.Handler, cli.Dict)
	})
}

// DialTimeout is like Dial, but with timeout
func (cli *Client) DialTimeout(addr string, timeout time.Duration) (diam.Conn, error) {
	return cli.DialExt("tcp", addr, timeout, nil)
}

// DialTLS is like Dial, but using TLS.
func (cli *Client) DialTLS(addr, certFile, keyFile string) (diam.Conn, error) {
	return cli.DialTLSExt("tcp", addr, certFile, keyFile, 0, nil)
}

// DialTLSTimeout is like DialTimeout, but using TLS.
func (cli *Client) DialTLSTimeout(addr, certFile, keyFile string, timeout time.Duration) (diam.Conn, error) {
	return cli.DialTLSExt("tcp", addr, certFile, keyFile, timeout, nil)
}

// DialNetworkTLS calls the network address set as ip:port, performs a handshake and optionally
// start a watchdog goroutine in background.
func (cli *Client) DialNetworkTLS(network, addr, certFile, keyFile string, laddr net.Addr) (diam.Conn, error) {
	return cli.DialTLSExt(network, addr, certFile, keyFile, 0, nil)
}

// DialExt - Optionally binds client to laddr, calls the network address set as ip:port,
// performs a handshake and optionally start a watchdog goroutine in background.
func (cli *Client) DialExt(network, addr string, timeout time.Duration, laddr net.Addr) (diam.Conn, error) {
	return cli.dial(func() (diam.Conn, error) {
		return diam.DialExt(network, addr, cli.Handler, cli.Dict, timeout, laddr)
	})
}

// DialTLSExt - Optionally binds client to laddr, calls the network address set as ip:port, performs a
// handshake and optionally start a watchdog goroutine in background.
func (cli *Client) DialTLSExt(
	network, addr, certFile, keyFile string, timeout time.Duration, laddr net.Addr) (diam.Conn, error) {

	return cli.dial(func() (diam.Conn, error) {
		return diam.DialTLSExt(network, addr, certFile, keyFile, cli.Handler, cli.Dict, timeout, laddr)
	})
}

// NewConn is like Dial, but using an already open net.Conn.
func (cli *Client) NewConn(rw net.Conn, addr string) (diam.Conn, error) {
	return cli.dial(func() (diam.Conn, error) {
		return diam.NewConn(rw, addr, cli.Handler, cli.Dict)
	})
}

type dialFunc func() (diam.Conn, error)

func (cli *Client) dial(f dialFunc) (diam.Conn, error) {
	if err := cli.validate(); err != nil {
		return nil, err
	}
	c, err := f()
	if err != nil {
		return c, err
	}
	c, err = cli.handshake(c)
	return c, err
}

func (cli *Client) validate() error {
	if cli.Handler == nil {
		return ErrMissingStateMachine
	}
	if cli.Dict == nil {
		cli.Dict = dict.Default
	}
	if cli.RetransmitInterval == 0 {
		// Set default RetransmitInterval.
		cli.RetransmitInterval = time.Second
	}
	if cli.WatchdogInterval == 0 {
		// Set default WatchdogInterval
		cli.WatchdogInterval = 5 * time.Second
	}
	// Make sure the applications supplied to Client are supported locally
	for _, submittedAcctApp := range cli.AcctApplicationID {
		acctAppID := uint32(submittedAcctApp.Data.(datatype.Unsigned32))
		isSupported := false
		for _, localApp := range cli.Handler.supportedApps {
			if localApp.AppType == "acct" && localApp.ID == acctAppID {
				isSupported = true
				break
			}
		}
		if isSupported == false {
			err := fmt.Errorf("Client attempts to advertise unsupported application - type: acct, id: %d", acctAppID)
			return err
		}

	}
	for _, submittedAuthApp := range cli.AuthApplicationID {
		authAppID := uint32(submittedAuthApp.Data.(datatype.Unsigned32))
		isSupported := false
		for _, localApp := range cli.Handler.supportedApps {
			if localApp.AppType == "auth" && localApp.ID == authAppID {
				isSupported = true
				break
			}
		}
		if isSupported == false {
			err := fmt.Errorf("Client attempts to advertise unsupported application - type: auth, id: %d", authAppID)
			return err
		}

	}
	return nil
}

func (cli *Client) handshake(c diam.Conn) (diam.Conn, error) {
	var (
		hostAddresses []datatype.Address
		err           error
	)
	if len(cli.Handler.cfg.HostIPAddresses) > 0 {
		hostAddresses = cli.Handler.cfg.HostIPAddresses
	} else {
		hostAddresses, err = getLocalAddresses(c)
		if err != nil {
			c.Close()
			return nil, fmt.Errorf("diameter handshake failure: %v", err)
		}
	}

	m := cli.makeCER(hostAddresses)
	// Ignore CER, but not DWR.
	cerClientHandler := func(c diam.Conn, m *diam.Message) {}
	// See sm.go for Base Diam Idx declarations
	cli.Handler.mux.HandleIdx(baseCERIdx, diam.HandlerFunc(cerClientHandler))
	cli.Handler.mux.HandleFunc("CER", cerClientHandler)
	// Handle CEA and DWA.
	errc := make(chan error)
	cli.Handler.mux.Handle("CEA", handleCEA(cli.Handler, errc))

	var dwac chan struct{}
	if cli.EnableWatchdog {
		dwac = make(chan struct{})
		cli.Handler.mux.Handle("DWA", handshakeOK(handleDWA(cli.Handler, dwac)))
	}
	for i := 0; i < (int(cli.MaxRetransmits) + 1); i++ {
		_, err := m.WriteTo(c)
		if err != nil {
			c.Close()
			return nil, err
		}
		select {
		case err, ok := <-errc: // Wait for CEA.
			if ok && err != nil {
				close(errc)
				c.Close()
				return nil, err
			}
			if cli.EnableWatchdog {
				go cli.watchdog(c, dwac)
			}
			return c, nil
		case <-time.After(cli.RetransmitInterval):
		}
	}
	c.Close()
	return nil, ErrHandshakeTimeout
}

func (cli *Client) makeCER(hostIPAddresses []datatype.Address) *diam.Message {
	m := diam.NewRequest(diam.CapabilitiesExchange, 0, cli.Dict)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, cli.Handler.cfg.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, cli.Handler.cfg.OriginRealm)
	for _, hostIPAddress := range hostIPAddresses {
		m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, hostIPAddress)
	}
	m.NewAVP(avp.VendorID, avp.Mbit, 0, cli.Handler.cfg.VendorID)
	m.NewAVP(avp.ProductName, 0, 0, cli.Handler.cfg.ProductName)
	if cli.Handler.cfg.OriginStateID != 0 {
		stateid := datatype.Unsigned32(cli.Handler.cfg.OriginStateID)
		m.NewAVP(avp.OriginStateID, avp.Mbit, 0, stateid)
	}
	if cli.SupportedVendorID != nil {
		for _, a := range cli.SupportedVendorID {
			m.AddAVP(a)
		}
	}
	if cli.AuthApplicationID != nil {
		for _, a := range cli.AuthApplicationID {
			m.AddAVP(a)
		}
	}
	m.NewAVP(avp.InbandSecurityID, avp.Mbit, 0, datatype.Unsigned32(0))
	if cli.AcctApplicationID != nil {
		for _, a := range cli.AcctApplicationID {
			m.AddAVP(a)
		}
	}
	if cli.VendorSpecificApplicationID != nil {
		for _, a := range cli.VendorSpecificApplicationID {
			m.AddAVP(a)
		}
	}
	if cli.Handler.cfg.FirmwareRevision != 0 {
		m.NewAVP(avp.FirmwareRevision, 0, 0, cli.Handler.cfg.FirmwareRevision)
	}
	return m
}

func (cli *Client) watchdog(c diam.Conn, dwac chan struct{}) {
	disconnect := c.(diam.CloseNotifier).CloseNotify()
	var osid = uint32(cli.Handler.cfg.OriginStateID)
	for {
		select {
		case <-disconnect:
			return
		case <-time.After(cli.WatchdogInterval):
			cli.dwr(c, osid, dwac)
		}
	}
}

func (cli *Client) dwr(c diam.Conn, osid uint32, dwac chan struct{}) {
	m := cli.makeDWR(osid)
	for i := 0; i < (int(cli.MaxRetransmits) + 1); i++ {
		_, err := m.WriteToStream(c, cli.WatchdogStream)
		if err != nil {
			return
		}
		select {
		case <-dwac:
			return
		case <-time.After(cli.RetransmitInterval):
		}
	}
	// Watchdog failed, disconnect.
	c.Close()
}

func (cli *Client) makeDWR(osid uint32) *diam.Message {
	m := diam.NewRequest(diam.DeviceWatchdog, 0, cli.Dict)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, cli.Handler.cfg.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, cli.Handler.cfg.OriginRealm)
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(osid))
	return m
}

func getLocalAddresses(c diam.Conn) ([]datatype.Address, error) {
	var addrStr string
	if c.LocalAddr() != nil {
		addrStr = c.LocalAddr().String()
	}
	addr, _, err := net.SplitHostPort(addrStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse local ip %s [%q]: %s", addrStr, c.LocalAddr(), err)
	}
	hostIPs := strings.Split(addr, "/")
	addresses := make([]datatype.Address, 0, len(hostIPs))
	for _, ipStr := range hostIPs {
		ip := net.ParseIP(ipStr)
		if ip != nil {
			addresses = append(addresses, datatype.Address(ip))
		}
	}
	return addresses, nil
}
