// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Skeleton of the dictionary file.  Part of go-diameter.

package dict

import (
	"encoding/xml"
	"fmt"

	"github.com/fiorix/go-diameter/v4/diam/datatype"
)

// File is the dictionary root element of a XML file.  See diam_base.xml.
type File struct {
	XMLName xml.Name `xml:"diameter"`
	App     []*App   `xml:"application"` // Support for multiple applications
}

// App defines a diameter application in XML and its multiple AVPs.
type App struct {
	ID      uint32     `xml:"id,attr"`   // Application Id
	Type    string     `xml:"type,attr"` // Application type
	Name    string     `xml:"name,attr"` // Application name
	Vendor  []*Vendor  `xml:"vendor"`    // Support for multiple vendors
	Command []*Command `xml:"command"`   // Diameter commands
	AVP     []*AVP     `xml:"avp"`       // Each application support multiple AVPs
}

// Vendor defines diameter vendors in XML, that can be used to translate
// the VendorId AVP of incoming messages.
type Vendor struct {
	ID   uint32 `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

// Command defines a diameter command (CE, CC, etc)
type Command struct {
	Code    uint32      `xml:"code,attr"`
	Name    string      `xml:"name,attr"`
	Short   string      `xml:"short,attr"`
	Request CommandRule `xml:"request"`
	Answer  CommandRule `xml:"answer"`
}

func (cmd *Command) String() string {
	if cmd == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s (%s) CODE: %d", cmd.Name, cmd.Short, cmd.Code)
}

// CommandRule contains rules for a given command.
type CommandRule struct {
	Rule []*Rule `xml:"rule"`
}

// AVP represents a dictionary AVP that is loaded from XML.
type AVP struct {
	Name       string `xml:"name,attr"`
	Code       uint32 `xml:"code,attr"`
	Must       string `xml:"must,attr"`
	May        string `xml:"may,attr"`
	MustNot    string `xml:"must-not,attr"`
	MayEncrypt string `xml:"may-encrypt,attr"`
	VendorID   uint32 `xml:"vendor-id,attr"`
	Data       Data   `xml:"data"`
	App        *App   `xml:"none"` // Link back to diameter application
}

// Data of an AVP can be EnumItem or a Parser of multiple AVPs.
type Data struct {
	Type     datatype.TypeID `xml:"-"`
	TypeName string          `xml:"type,attr"`
	Enum     []*Enum         `xml:"item"` // In case of Enumerated AVP data
	Rule     []*Rule         `xml:"rule"` // In case of Grouped AVPs
}

// Enum contains the code and name of Enumerated items.
type Enum struct {
	// rfc6733 (section 4.3.1):
	// The Enumerated format is derived from the Integer32 Basic AVP Format.
	Code int32  `xml:"code,attr"`
	Name string `xml:"name,attr"`
}

// Rule defines the usage rules of an AVP.
type Rule struct {
	AVP      string `xml:"avp,attr"` // AVP Name
	Required bool   `xml:"required,attr"`
	Min      int    `xml:"min,attr"`
	Max      int    `xml:"max,attr"`
}
