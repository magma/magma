// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Converts Wireshark diameter dictionaries to go-diameter format.
// Use: wireshark-dict-tool < wireshark-dict.xml > new-dict.xml
//
// Some wireshark dictionaries must be slightly fixed before they can
// be converted by this tool.

package main

// TODO: Improve the parser and fix AVP properties during conversion:
// <avp name=".." code=".." must="" may="" must-not="" may-encrypt="">

import (
	"encoding/xml"
	"log"
	"os"

	"github.com/fiorix/go-diameter/diam/dict"
)

func main() {
	wsd, err := load(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	var newDict = &dict.File{}
	for _, app := range wsd.App {
		newApp := &dict.App{
			ID:   app.ID,
			Type: app.Type,
			Name: app.Name,
		}
		copyVendors(wsd.Vendor, newApp)
		copyCommands(app.Cmd, newApp)
		copyAvps(app.AVP, newApp)
		newDict.App = append(newDict.App, newApp)
	}
	os.Stdout.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "\t")
	enc.Encode(newDict)
}

func copyVendors(src []*Vendor, dst *dict.App) {
	for _, vendor := range src {
		dst.Vendor = append(dst.Vendor, &dict.Vendor{
			ID:   vendor.ID,
			Name: vendor.Name,
		})
	}
}

func copyCommands(src []*Cmd, dst *dict.App) {
	for _, cmd := range src {
		newCmd := &dict.Command{
			Code:  cmd.Code,
			Name:  cmd.Name,
			Short: cmd.Name,
		}
		copyCmdRules(cmd.Request.Fixed.Rule, &newCmd.Request, false)
		copyCmdRules(cmd.Request.Required.Rule, &newCmd.Request, true)
		copyCmdRules(cmd.Request.Optional.Rule, &newCmd.Request, false)
		copyCmdRules(cmd.Answer.Fixed.Rule, &newCmd.Answer, false)
		copyCmdRules(cmd.Answer.Required.Rule, &newCmd.Answer, true)
		copyCmdRules(cmd.Answer.Optional.Rule, &newCmd.Answer, false)
		dst.Command = append(dst.Command, newCmd)
	}
}

func copyCmdRules(src []*Rule, dst *dict.CommandRule, required bool) {
	for _, req := range src {
		dst.Rule = append(dst.Rule, &dict.Rule{
			AVP:      req.Name,
			Required: required,
			Min:      req.Min,
			Max:      req.Max,
		})
	}
}

func copyAvps(src []*AVP, dst *dict.App) {
	for _, avp := range src {
		newAVP := &dict.AVP{
			Name: avp.Name,
			Code: avp.Code,
		}
		if avp.Type.Name == "" && avp.Grouped != nil {
			newAVP.Data = dict.Data{TypeName: "Grouped"}
		} else {
			newAVP.Data = dict.Data{TypeName: avp.Type.Name}
		}
		switch avp.MayEncrypt {
		case "yes":
			newAVP.MayEncrypt = "Y"
		case "no":
			newAVP.MayEncrypt = "N"
		default:
			newAVP.MayEncrypt = "-"
		}
		switch avp.Mandatory {
		case "must":
			newAVP.Must = "M"
		case "may":
			newAVP.May = "P"
		default:
			newAVP.Must = ""
		}
		if newAVP.May != "" {
			switch avp.Protected {
			case "may":
				newAVP.May = "P"
			default:
				newAVP.May = ""
			}
		}
		for _, p := range avp.Enum {
			newAVP.Data.Enum = append(newAVP.Data.Enum,
				&dict.Enum{
					Name: p.Name,
					Code: p.Code,
				})
		}
		for _, grp := range avp.Grouped {
			for _, p := range grp.GAVP {
				newAVP.Data.Rule = append(newAVP.Data.Rule,
					&dict.Rule{
						AVP: p.Name,
						Min: p.Min,
						Max: p.Max,
					})
			}
			for _, p := range grp.Required.Rule {
				newAVP.Data.Rule = append(newAVP.Data.Rule,
					&dict.Rule{
						AVP:      p.Name,
						Required: true,
						Min:      p.Min,
						Max:      p.Max,
					})
			}
			for _, p := range grp.Optional.Rule {
				newAVP.Data.Rule = append(newAVP.Data.Rule,
					&dict.Rule{
						AVP:      p.Name,
						Required: false,
						Min:      p.Min,
						Max:      p.Max,
					})
			}
		}
		dst.AVP = append(dst.AVP, newAVP)
	}
}
