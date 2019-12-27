// Copyright 2013-2015 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Example of nested grouped AVPs.

package main

import (
	"bytes"
	"log"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

func main() {
	err := dict.Default.Load(bytes.NewReader(customApp))
	if err != nil {
		log.Fatal(err)
	}

	m := diam.NewMessage(1111, diam.RequestFlag, 999, 1, 2, dict.Default)
	m.NewAVP(avp.ProductName, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(7070, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(10)),
			diam.NewAVP(8080, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(20)),
					diam.NewAVP(9090, avp.Mbit, 0, &diam.GroupedAVP{
						AVP: []*diam.AVP{
							diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(30)),
							diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(40)),
						},
					}),
				},
			}),
		},
	})

	log.Printf("m:\n%s\n", m)

	var b bytes.Buffer
	_, err = m.WriteTo(&b)
	if err != nil {
		log.Fatal(err)
	}

	z, err := diam.ReadMessage(&b, dict.Default)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("z:\n%s\n", z)

	var y customMsg
	err = z.Unmarshal(&y)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("y:\n%+v\n", y)
}

type customMsg struct {
	ProductName  string  `avp:"Product-Name"`
	NestedGroupA nestedA `avp:"Nested-Group-A"`
}

type nestedA struct {
	VendorID     int     `avp:"Vendor-Id"`
	NestedGroupB nestedB `avp:"Nested-Group-B"`
}

type nestedB struct {
	VendorID     int     `avp:"Vendor-Id"`
	NestedGroupC nestedC `avp:"Nested-Group-C"`
}

type nestedC struct {
	VendorID []int `avp:"Vendor-Id"`
}

var customApp = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<diameter>
	<application id="999">
		<command code="1111" short="TC" name="Test-Command">
			<request>
				<rule avp="Product-Name" required="true" max="1"/>
				<rule avp="Nested-Group-A" required="true" max="1"/>
			</request>
			<answer>
				<rule avp="Result-Code" required="true" max="1"/>
			</answer>
		</command>

		<avp name="Nested-Group-A" code="7070">
			<data type="Grouped">
				<rule avp="Vendor-Id" required="true" max="1"/>
				<rule avp="Nested-Group-B" required="true" max="1"/>
			</data>
		</avp>

		<avp name="Nested-Group-B" code="8080">
			<data type="Grouped">
				<rule avp="Vendor-Id" required="true" max="1"/>
				<rule avp="Nested-Group-C" required="true" max="1"/>
			</data>
		</avp>

		<avp name="Nested-Group-C" code="9090">
			<data type="Grouped">
				<rule avp="Vendor-Id" required="true" max="2"/>
			</data>
		</avp>
	</application>
</diameter>
`)
