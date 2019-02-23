/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package http2 contains a minimal implementation of non-TLS http/2 server
// and client
package http2

import (
	"net"
	"net/http"

	"github.com/golang/glog"
	"golang.org/x/net/http2"
)

// H2CServer is a minimal http/2 server supports non-TLS only
type H2CServer struct {
	*http2.Server
}

func NewH2CServer() *H2CServer {
	return &H2CServer{&http2.Server{}}
}

func (server *H2CServer) Run(addr string, handler http.HandlerFunc) {
	tcpListener, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Fatalf("net.Listen err: %v\n", err)
	}
	server.Serve(tcpListener, handler)
}

func (server *H2CServer) Serve(listener net.Listener, handler http.HandlerFunc) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			glog.Fatalf("l.accept err: %v\n", err)
		}
		go server.ServeConn(conn, &http2.ServeConnOpts{
			Handler: handler,
		})
	}
}
