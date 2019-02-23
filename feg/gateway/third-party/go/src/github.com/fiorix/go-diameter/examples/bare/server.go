// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Experimental diameter server that currently does nothing but print
// incoming messages.
package main

import (
	"log"
	"net"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/dict"
)

func main() {
	dict.Default.LoadFile("diam_app.xml")
	srv, err := net.Listen("tcp", ":3868")
	if err != nil {
		panic(err)
	}
	for {
		if conn, err := srv.Accept(); err != nil {
			panic(err)
		} else {
			go handleClient(conn)
		}
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	for {
		m, err := diam.ReadMessage(conn, dict.Default)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(m)
	}
}
