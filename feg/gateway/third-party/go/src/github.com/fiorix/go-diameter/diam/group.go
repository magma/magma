// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"fmt"

	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// GroupedAVPType is the identifier of the GroupedAVP data type.
// It must not conflict with other values from the datatype package.
const GroupedAVPType = 50

// GroupedAVP that is different from the dummy datatype.Grouped.
type GroupedAVP struct {
	AVP []*AVP
}

// DecodeGrouped decodes a Grouped AVP from a datatype.Grouped (byte array).
func DecodeGrouped(data datatype.Grouped, application uint32, dictionary *dict.Parser) (*GroupedAVP, error) {
	g := &GroupedAVP{}
	b := []byte(data)
	for n := 0; n < len(b); {
		avp, err := DecodeAVP(b[n:], application, dictionary)
		if err != nil {
			return nil, err
		}
		g.AVP = append(g.AVP, avp)
		n += avp.Len()
	}
	// TODO: handle nested groups?
	return g, nil
}

// Serialize implements the datatype.Type interface.
func (g *GroupedAVP) Serialize() []byte {
	b := make([]byte, g.Len())
	var n int
	for _, a := range g.AVP {
		a.SerializeTo(b[n:])
		n += a.Len()
	}
	return b
}

// Len implements the datatype.Type interface.
func (g *GroupedAVP) Len() int {
	var l int
	for _, a := range g.AVP {
		l += a.Len()
	}
	return l
}

// Padding implements the datatype.Type interface.
func (g *GroupedAVP) Padding() int {
	return 0
}

// Type implements the datatype.Type interface.
func (g *GroupedAVP) Type() datatype.TypeID {
	return GroupedAVPType
}

// String implements the datatype.Type interface.
func (g *GroupedAVP) String() string {
	var b bytes.Buffer
	for n, a := range g.AVP {
		if n > 0 {
			fmt.Fprint(&b, ",")
		}
		fmt.Fprint(&b, a)
	}
	return b.String()
}

// AddAVP adds the AVP to the GroupedAVP. It is not safe for concurrent calls.
func (g *GroupedAVP) AddAVP(a *AVP) {
	g.AVP = append(g.AVP, a)
}
