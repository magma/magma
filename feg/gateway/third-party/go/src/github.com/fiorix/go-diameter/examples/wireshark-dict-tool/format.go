// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Dictionary file structure.  Part of go-diameter.

package main

import "encoding/xml"

// File is the dictionary root element of a XML file.  See diam_base.xml.
type File struct {
	XMLName xml.Name  `xml:"dictionary"`
	Vendor  []*Vendor `xml:"vendor"`
	Base    Base      `xml:"base"`
	App     []*App    `xml:"application"` // Support for multiple applications
}

// Vendor defines diameter vendors in XML, that can be used to translate
// the VendorId AVP of incoming messages.
type Vendor struct {
	ID   uint32 `xml:"vendor-id,attr"`
	Name string `xml:"name,attr"`
}

// Base defined AVPs that can be used by all subsequent Applications
type Base struct {
	URI      string      `xml:"uri,attr"`
	Typedefn []*Typedefn `xml:"typedefn"`
	AVP      []*AVP      `xml:"avp"`
}

// Typedefn defines diameter data types for the base protocol
type Typedefn struct {
	Name   string `xml:"type-name,attr"`
	Parent string `xml:"type-parent,attr"`
}

// AVP represents a dictionary AVP that is loaded from XML
type AVP struct {
	Name       string     `xml:"name,attr"`
	Code       uint32     `xml:"code,attr"`
	Mandatory  string     `xml:"mandatory,attr"`
	MayEncrypt string     `xml:"may-encrypt,attr"`
	Protected  string     `xml:"protected,attr"`
	Type       DataType   `xml:"type"`
	Enum       []*Enum    `xml:"enum"`    // In case of Enumerated AVP
	Grouped    []*Grouped `xml:"grouped"` // In case of Grouped AVP
}

// DataType represents a dictionary data type.
type DataType struct {
	Name string `xml:"type-name,attr"`
}

// Enum contains the code and name of Enumerated items.
type Enum struct {
	Name string `xml:"name,attr"`
	Code int32  `xml:"code,attr"`
}

// Grouped represents a grouped AVP definition.
type Grouped struct {
	GAVP     []*Rule     `xml:"gavp"`
	Required GroupedRule `xml:"required"`
	Optional GroupedRule `xml:"optional"`
}

// GroupedRule defines the usage rules of a Grouped AVP.
type GroupedRule struct {
	Rule []*Rule `xml:"avprule"`
}

// Rule defines the usage rules of an AVP.
type Rule struct {
	Name string `xml:"name,attr"`
	Min  int    `xml:"minimum,attr"`
	Max  int    `xml:"maximum,attr"`
}

// App defines a diameter application in XML and its multiple AVPs.
type App struct {
	ID   uint32 `xml:"id,attr"`
	Type string `xml:"type,attr"`
	Name string `xml:"name,attr"`
	URI  string `xml:"uri,attr"`
	Cmd  []*Cmd `xml:"command"` // Diameter commands
	AVP  []*AVP `xml:"avp"`     // Each application support multiple AVPs
}

// Cmd defines a diameter command (CE, CC, etc)
type Cmd struct {
	Name    string  `xml:"name,attr"`
	Code    uint32  `xml:"code,attr"`
	Request CmdRule `xml:"requestrules"`
	Answer  CmdRule `xml:"answerrules"`
}

// CmdRule defines the rule for a Diameter command.
type CmdRule struct {
	Fixed    Rules `xml:"fixed"`
	Required Rules `xml:"required"`
	Optional Rules `xml:"optional"`
}

// Rules contains a list of AVP rules.
type Rules struct {
	Rule []*Rule `xml:"avprule"`
}
