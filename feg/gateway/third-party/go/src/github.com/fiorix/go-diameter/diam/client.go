// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter client.

package diam

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// DialNetwork connects to the peer pointed to by network & addr and returns the Conn that
// can be used to send diameter messages. Incoming messages are handled
// by the handler, which is typically nil and DefaultServeMux is used.
// If dict is nil, dict.Default is used.
func DialNetwork(network, addr string, handler Handler, dp *dict.Parser) (Conn, error) {
	return DialExt(network, addr, handler, dp, 0, nil)
}

func DialNetworkBind(network, laddr, raddr string, handler Handler, dp *dict.Parser) (Conn, error) {
	var (
		err     error
		netAddr net.Addr
	)

	if laddr != "" {
		netAddr, err = resolveAddress(network, laddr)
		if err != nil {
			return nil, err
		}
	}
	return DialExt(network, raddr, handler, dp, 0, netAddr)
}

func DialNetworkTimeout(network, addr string, handler Handler, dp *dict.Parser, timeout time.Duration) (Conn, error) {
	return DialExt(network, addr, handler, dp, timeout, nil)
}

// Dial connects to the peer pointed to by addr and returns the Conn that
// can be used to send diameter messages. Incoming messages are handled
// by the handler, which is typically nil and DefaultServeMux is used.
// If dict is nil, dict.Default is used.
func Dial(addr string, handler Handler, dp *dict.Parser) (Conn, error) {
	return DialNetwork("tcp", addr, handler, dp)
}

func DialTimeout(addr string, handler Handler, dp *dict.Parser, timeout time.Duration) (Conn, error) {
	return DialNetworkTimeout("tcp", addr, handler, dp, timeout)
}

// DialExt - extended dial API connects to the peer pointed to by network &
// addr and returns the Conn that can be used to send diameter messages.
// Incoming messages are handled by the handler, which is typically nil and
// DefaultServeMux is used. Allows binding dailer socket to given laddr.
// If dict is nil, dict.Default is used.
func DialExt(
	network, addr string, handler Handler, dp *dict.Parser, timeout time.Duration, laddr net.Addr) (Conn, error) {

	srv := &Server{Network: network, Addr: addr, Handler: handler, Dict: dp, LocalAddr: laddr}
	return dial(srv, timeout)
}

// dial network wrapper
func dial(srv *Server, timeout time.Duration) (Conn, error) {
	network := srv.Network
	if len(network) == 0 {
		network = "tcp"
	}
	addr := srv.Addr
	if len(addr) == 0 {
		addr = ":3868"
	}
	var rw net.Conn
	var err error
	dialer := getMultistreamDialer(network, timeout, srv.LocalAddr)
	rw, err = dialer.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	c, err := srv.newConn(rw)
	if err != nil {
		return nil, err
	}
	go c.serve()
	return c.writer, nil
}

// DialTLS is the same as Dial, but for TLS.
func DialTLS(addr, certFile, keyFile string, handler Handler, dp *dict.Parser) (Conn, error) {
	return DialTLSExt("tcp", addr, certFile, keyFile, handler, dp, 0, nil)
}

// DialTLSTimeout is the same as DialTimeout, but for TLS.
func DialTLSTimeout(addr, certFile, keyFile string, handler Handler, dp *dict.Parser, timeout time.Duration) (Conn, error) {
	return DialTLSExt("tcp", addr, certFile, keyFile, handler, dp, timeout, nil)
}

// DialNetworkTLS is the same as DialNetwork, but for TLS.
func DialNetworkTLS(network, addr, certFile, keyFile string, handler Handler, dp *dict.Parser) (Conn, error) {
	return DialTLSExt(network, addr, certFile, keyFile, handler, dp, 0, nil)
}

// DialTLSExt is the same as DialExt, but for TLS.
func DialTLSExt(
	network,
	addr,
	certFile,
	keyFile string,
	handler Handler,
	dp *dict.Parser,
	timeout time.Duration,
	laddr net.Addr) (Conn, error) {

	srv := &Server{Network: network, Addr: addr, Handler: handler, Dict: dp, LocalAddr: laddr}
	return dialTLS(srv, certFile, keyFile, timeout)
}

// dialTLS net TCP wrapper
func dialTLS(srv *Server, certFile, keyFile string, timeout time.Duration) (Conn, error) {
	var err error
	network := srv.Network
	if len(network) == 0 {
		network = "tcp"
	}
	addr := srv.Addr
	if len(addr) == 0 {
		addr = ":3868"
	}
	var config *tls.Config
	if srv.TLSConfig == nil {
		config = &tls.Config{InsecureSkipVerify: true}
	} else {
		config = TLSConfigClone(srv.TLSConfig)
	}
	if len(certFile) != 0 {
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}
	}

	var rw net.Conn
	dialer := getDialer(network, timeout, srv.LocalAddr)
	rw, err = dialer.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	c, err := srv.newConn(tls.Client(rw, config))
	if err != nil {
		return nil, err
	}
	go c.serve()
	return c.writer, nil
}

// NewConn is the same as Dial, but using an already open net.Conn.
func NewConn(rw net.Conn, addr string, handler Handler, dp *dict.Parser) (Conn, error) {
	srv := &Server{Addr: addr, Handler: handler, Dict: dp}

	c, err := srv.newConn(rw)
	if err != nil {
		return nil, err
	}
	go c.serve()
	return c.writer, nil
}
