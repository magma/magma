// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package magmarebootnode

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/facebookincubator/symphony/pkg/actions/core"
	"github.com/facebookincubator/symphony/pkg/orc8r"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMagmaRebootNode(t *testing.T) {
	var handlerCalled bool
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/networks/network1/gateways/gateway1/command/reboot", r.URL.Path)
		_, err := io.WriteString(w, "ok")
		assert.NoError(t, err)
		handlerCalled = true
	}))
	defer srv.Close()

	uri, err := url.Parse(srv.URL)
	require.NoError(t, err)

	orc8rClient := srv.Client()
	orc8rClient.Transport = orc8r.Transport{
		Base: orc8rClient.Transport,
		Host: uri.Host,
	}

	action := New(orc8rClient)
	ac := core.ActionContext{
		TriggerPayload: map[string]interface{}{
			"networkID": "network1",
			"gatewayID": "gateway1",
		},
		Rule:       core.Rule{},
		RuleAction: &core.ActionsRuleAction{},
	}
	err = action.Execute(ac)
	assert.NoError(t, err)
	assert.True(t, handlerCalled)
}
