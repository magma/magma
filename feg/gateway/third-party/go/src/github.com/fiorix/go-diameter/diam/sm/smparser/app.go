// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smparser

import (
	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/dict"
)

// Role stores information whether SM is initialized as a Client or a Server
type Role uint8

// ServerRole and ClientRole enums are passed to smparser for proper CER/CEA verification
const (
	Server Role = iota + 1
	Client
)

// Application validates accounting, auth, and vendor specific application IDs.
type Application struct {
	AcctApplicationID           []*diam.AVP
	AuthApplicationID           []*diam.AVP
	VendorSpecificApplicationID []*diam.AVP
	id                          []uint32 // List of supported application IDs.
}

// Parse ensures all acct or auth applications in the CE
// exist in this server's dictionary.
func (app *Application) Parse(d *dict.Parser, localRole Role) (failedAVP *diam.AVP, err error) {
	failedAVP, err = app.validateAll(d, avp.AcctApplicationID, app.AcctApplicationID)
	if err != nil {
		return failedAVP, err
	}
	failedAVP, err = app.validateAll(d, avp.AuthApplicationID, app.AuthApplicationID)
	if err != nil {
		return failedAVP, err
	}
	if app.VendorSpecificApplicationID != nil {
		for _, vs := range app.VendorSpecificApplicationID {
			failedAVP, err := app.handleGroup(d, vs)
			if err != nil {
				return failedAVP, err
			}
		}
	}
	if app.ID() == nil {
		if localRole == Client {
			return nil, ErrMissingApplication
		}
		return nil, ErrNoCommonApplication

	}
	return nil, nil
}

// handleGroup handles the VendorSpecificApplicationID grouped AVP and
// validates accounting or auth applications.
func (app *Application) handleGroup(d *dict.Parser, gavp *diam.AVP) (failedAVP *diam.AVP, err error) {
	group, ok := gavp.Data.(*diam.GroupedAVP)
	if !ok {
		return gavp, &ErrUnexpectedAVP{gavp}
	}
	for _, a := range group.AVP {
		switch a.Code {
		case avp.AcctApplicationID:
			failedAVP, err = app.validate(d, a.Code, a)
		case avp.AuthApplicationID:
			failedAVP, err = app.validate(d, a.Code, a)
		}
	}
	return failedAVP, err
}

// validateAll is a convenience method to test a slice of application IDs.
func (app *Application) validateAll(d *dict.Parser, appType uint32, appAVPs []*diam.AVP) (failedAVP *diam.AVP, err error) {
	if appAVPs != nil {
		for _, a := range appAVPs {
			failedAVP, err = app.validate(d, appType, a)
			if err != nil {
				return failedAVP, err
			}
		}
	}
	return nil, nil
}

// validate ensures the given acct or auth application ID exists in
// the given dictionary.
func (app *Application) validate(d *dict.Parser, appType uint32, appAVP *diam.AVP) (failedAVP *diam.AVP, err error) {
	if appAVP == nil {
		return nil, nil
	}
	var typ string
	switch appType {
	case avp.AcctApplicationID:
		typ = "acct"
	case avp.AuthApplicationID:
		typ = "auth"
	}
	if appAVP.Code != appType {
		return appAVP, &ErrUnexpectedAVP{appAVP}
	}
	appID, ok := appAVP.Data.(datatype.Unsigned32)
	if !ok {
		return appAVP, &ErrUnexpectedAVP{appAVP}
	}
	id := uint32(appID)
	if id == 0xffffffff { // relay application id
		app.id = append(app.id, id)
		return nil, nil
	}
	avp, err := d.App(id)
	if err != nil {
		//TODO Log informational message to console?
	} else if len(avp.Type) > 0 && avp.Type != typ {
		return nil, ErrNoCommonApplication
	} else {
		app.id = append(app.id, id)
	}
	return nil, nil
}

// ID returns a list of supported application IDs.
// Must be called after Parse, otherwise it returns an empty array.
func (app *Application) ID() []uint32 {
	return app.id
}
