// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package magmarebootnode

import (
	"fmt"
	"log"
	"net/http"

	"github.com/facebookincubator/symphony/pkg/actions/core"
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

// DataType is the expected type for this action
func (a *action) DataType() core.DataType {
	return core.DataTypeString
}

// Execute executes the action
func (a *action) Execute(ctx core.ActionContext) error {
	p := ctx.TriggerPayload
	rule := ctx.Rule
	networkID := p["networkID"]
	gatewayID := p["gatewayID"]

	log.Printf("running action:%v, networkID:%v, gatewayID:%v, ruleID: %v",
		core.MagmaRebootNodeActionID, networkID, gatewayID, rule.ID)

	url := fmt.Sprintf("/networks/%s/gateways/%s/command/reboot", networkID, gatewayID)
	res, err := a.orc8rClient.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("rebooting node: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("node reboot received status %d", res.StatusCode)
	}
	return nil
}
