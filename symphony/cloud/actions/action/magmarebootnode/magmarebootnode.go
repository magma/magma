// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package magmarebootnode

import (
	"log"
	"net/http"

	"github.com/facebookincubator/symphony/cloud/actions/core"
)

type action struct {
	orc8rClient *http.Client
}

// New returns a new action
func New(orc8rClient *http.Client) core.Action {
	return &action{orc8rClient}
}

// ID returns the string identifier for this trigger
func (*action) ID() core.ActionID {
	return core.MagmaRebootNodeActionID
}

// Description is a description when rebooting magma
func (a *action) Description() string {
	return "reboot a magma node"
}

// Execute executes the action
func (a *action) Execute(ctx core.ActionContext) error {
	p := ctx.TriggerPayload
	rule := ctx.Rule
	networkID := p["networkID"]
	gatewayID := p["gatewayID"]

	log.Printf("running action:%v, networkID:%v, gatewayID:%v, ruleID: %v",
		core.MagmaRebootNodeActionID, networkID, gatewayID, rule.ID)
	return nil
}
