/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package receivers

import (
	"regexp"
	"testing"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/fsclient/mocks"

	"github.com/prometheus/alertmanager/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testNID              = "test"
	otherNID             = "other"
	testAlertmanagerFile = `global:
  resolve_timeout: 5m
  http_config: {}
  smtp_hello: localhost
  smtp_require_tls: true
  pagerduty_url: https://events.pagerduty.com/v2/enqueue
  hipchat_api_url: https://api.hipchat.com/
  opsgenie_api_url: https://api.opsgenie.com/
  wechat_api_url: https://qyapi.weixin.qq.com/cgi-bin/
  victorops_api_url: https://alert.victorops.com/integrations/generic/20131114/alert/
route:
  receiver: null_receiver
  group_by:
  - alertname
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  routes:
  - receiver: other_tenant_base_route
    match:
      tenantID: other
  - receiver: deprecated_network_base_route
    match:
      networkID: deprecated
receivers:
- name: null_receiver
- name: test_receiver
- name: receiver
- name: other_tenant_base_route
- name: deprecated_network_base_route
- name: test_slack
  slack_configs:
  - api_url: http://slack.com/12345
    channel: string
    username: string
- name: other_receiver
  slack_configs:
  - api_url: http://slack.com/54321
    channel: string
    username: string
- name: test_webhook
  webhook_configs:
  - url: http://webhook.com/12345
    send_resolved: true
- name: test_email
  email_configs:
  - to: test@mail.com
    from: testUser
    smarthost: http://mail-server.com
    headers:
      name: value
      foo: bar
templates: []`
)

type testClients struct {
	aClient  AlertmanagerClient
	fsClient *mocks.FSClient
}

func TestClient_CreateReceiver(t *testing.T) {
	mtClient, mtFSClient := newMultiTenantTestClient()
	stClient, stFSClient := newSingleTenantTestClient()
	for _, clients := range []testClients{{mtClient, mtFSClient}, {stClient, stFSClient}} {
		// Create Slack Receiver
		err := clients.aClient.CreateReceiver(testNID, sampleSlackReceiver)
		assert.NoError(t, err)
		clients.fsClient.AssertCalled(t, "WriteFile", "test/alertmanager.yml", mock.Anything, mock.Anything)

		// Create Webhook Receiver
		err = clients.aClient.CreateReceiver(testNID, sampleWebhookReceiver)
		assert.NoError(t, err)
		clients.fsClient.AssertCalled(t, "WriteFile", "test/alertmanager.yml", mock.Anything, mock.Anything)

		// Create Email receiver
		err = clients.aClient.CreateReceiver(testNID, sampleEmailReceiver)
		assert.NoError(t, err)
		clients.fsClient.AssertCalled(t, "WriteFile", "test/alertmanager.yml", mock.Anything, mock.Anything)

		// create duplicate receiver
		err = clients.aClient.CreateReceiver(testNID, Receiver{Name: "receiver"})
		assert.Regexp(t, regexp.MustCompile("notification config name \".*receiver\" is not unique"), err.Error())
	}
}

func TestClient_GetReceivers(t *testing.T) {
	mtClient, _ := newMultiTenantTestClient()
	recs, err := mtClient.GetReceivers(testNID)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(recs))
	assert.Equal(t, "receiver", recs[0].Name)
	assert.Equal(t, "slack", recs[1].Name)
	assert.Equal(t, "webhook", recs[2].Name)
	assert.Equal(t, "email", recs[3].Name)

	recs, err = mtClient.GetReceivers(otherNID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(recs))

	recs, err = mtClient.GetReceivers("bad_nid")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(recs))

	stClient, _ := newSingleTenantTestClient()
	recs, err = stClient.GetReceivers(testNID)
	assert.NoError(t, err)
	assert.Equal(t, 9, len(recs))

	recs, err = stClient.GetReceivers(otherNID)
	assert.NoError(t, err)
	assert.Equal(t, 9, len(recs))

	recs, err = stClient.GetReceivers("bad_nid")
	assert.NoError(t, err)
	assert.Equal(t, 9, len(recs))

}

func TestClient_UpdateReceiver(t *testing.T) {
	mtClient, mtFSClient := newMultiTenantTestClient()
	err := mtClient.UpdateReceiver(testNID, &Receiver{Name: "slack"})
	mtFSClient.AssertCalled(t, "WriteFile", "test/alertmanager.yml", mock.Anything, mock.Anything)
	assert.NoError(t, err)

	err = mtClient.UpdateReceiver(testNID, &Receiver{Name: "nonexistent"})
	mtFSClient.AssertNumberOfCalls(t, "WriteFile", 1)
	assert.Error(t, err)

	stClient, stFSClient := newSingleTenantTestClient()
	err = stClient.UpdateReceiver(testNID, &Receiver{Name: "receiver"})
	stFSClient.AssertCalled(t, "WriteFile", "test/alertmanager.yml", mock.Anything, mock.Anything)
	assert.NoError(t, err)

	err = stClient.UpdateReceiver(testNID, &Receiver{Name: "nonexistent"})
	stFSClient.AssertNumberOfCalls(t, "WriteFile", 1)
	assert.Error(t, err)
}

func TestClient_DeleteReceiver(t *testing.T) {
	mtClient, mtFSClient := newMultiTenantTestClient()
	err := mtClient.DeleteReceiver(testNID, "slack")
	mtFSClient.AssertCalled(t, "WriteFile", "test/alertmanager.yml", mock.Anything, mock.Anything)
	assert.NoError(t, err)

	err = mtClient.DeleteReceiver(testNID, "nonexistent")
	assert.Error(t, err)
	mtFSClient.AssertNumberOfCalls(t, "WriteFile", 1)

	stClient, stFSClient := newSingleTenantTestClient()
	err = stClient.DeleteReceiver("", "receiver")
	stFSClient.AssertCalled(t, "WriteFile", "test/alertmanager.yml", mock.Anything, mock.Anything)
	assert.NoError(t, err)

	err = stClient.DeleteReceiver("", "nonexistent")
	assert.Error(t, err)
	stFSClient.AssertNumberOfCalls(t, "WriteFile", 1)
}

func TestClient_ModifyTenantRoute(t *testing.T) {
	mtClient, mtFSClient := newMultiTenantTestClient()
	err := mtClient.ModifyTenantRoute(testNID, &config.Route{
		Receiver: "slack",
	})
	assert.NoError(t, err)
	mtFSClient.AssertCalled(t, "WriteFile", "test/alertmanager.yml", mock.Anything, mock.Anything)

	err = mtClient.ModifyTenantRoute(testNID, &config.Route{
		Receiver: "test",
		Routes: []*config.Route{{
			Receiver: "nonexistent",
		}},
	})
	assert.Error(t, err)
	mtFSClient.AssertNumberOfCalls(t, "WriteFile", 1)

	stClient, stFSClient := newSingleTenantTestClient()
	err = stClient.ModifyTenantRoute(testNID, &config.Route{
		Receiver: "receiver",
	})
	assert.NoError(t, err)

	stFSClient.AssertCalled(t, "WriteFile", "test/alertmanager.yml", mock.Anything, mock.Anything)
	err = stClient.ModifyTenantRoute("", &config.Route{
		Receiver: "test",
		Routes: []*config.Route{{
			Receiver: "nonexistent",
		}},
	})
	assert.Error(t, err)
	stFSClient.AssertNumberOfCalls(t, "WriteFile", 1)
}

func TestClient_GetRoute(t *testing.T) {
	mtClient, _ := newMultiTenantTestClient()

	route, err := mtClient.GetRoute(otherNID)
	assert.NoError(t, err)
	assert.Equal(t, config.Route{Receiver: "tenant_base_route", Match: map[string]string{"tenantID": "other"}}, *route)

	route, err = mtClient.GetRoute("deprecated")
	assert.NoError(t, err)
	assert.Equal(t, config.Route{Receiver: "tenant_base_route", Match: map[string]string{"networkID": "deprecated"}}, *route)

	route, err = mtClient.GetRoute("no-network")
	assert.Error(t, err)

	stClient, _ := newSingleTenantTestClient()
	route, err = stClient.GetRoute("")
	assert.NoError(t, err)
	assert.Equal(t, "null_receiver", route.Receiver)
}

func newMultiTenantTestClient() (AlertmanagerClient, *mocks.FSClient) {
	fsClient := &mocks.FSClient{}
	fsClient.On("ReadFile", mock.Anything).Return([]byte(testAlertmanagerFile), nil)
	fsClient.On("WriteFile", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	return NewClient("test/alertmanager.yml", "alertmanager-host:9093", "networkID", fsClient), fsClient
}

func newSingleTenantTestClient() (AlertmanagerClient, *mocks.FSClient) {
	fsClient := &mocks.FSClient{}
	fsClient.On("ReadFile", mock.Anything).Return([]byte(testAlertmanagerFile), nil)
	fsClient.On("WriteFile", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	return NewClient("test/alertmanager.yml", "alertmanager-host:9093", "", fsClient), fsClient
}
