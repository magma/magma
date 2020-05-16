// +build all gx qos

/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package integration

import (
	"fmt"
	"io/ioutil"
	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/feg/cloud/go/protos"
	"magma/lte/cloud/go/plugin/models"
	"os"

	"math/rand"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

func configFileManager(fn string) (chan string, error) {
	info, err := os.Stat(fn)
	if err != nil {
		return nil, err
	}
	orig, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	ch := make(chan string)
	go func() {
		new := orig
		for mod := range ch {
			mod += "\n"
			new = append(orig, []byte(mod)...)
			ioutil.WriteFile(fn, new, info.Mode())
		}
		ioutil.WriteFile(fn, orig, info.Mode())
	}()
	return ch, nil
}

const (
	pipelinedCfgFn      = "pipelined.yml"
	cleanRestartYaml    = "clean_restart: true"
	nonCleanRestartYaml = "clean_restart: false"
)

//TestQosEnforcementRestart
// This test verifies the QOS configuration(uplink) present in the rules
// - Set an expectation for a  CCR-I to be sent up to PCRF, to which it will
//   respond with a rule install (static-ULQos) with QOS config setting with
//   maximum uplink bitrate.
// - Generate traffic and verify if the traffic observed bitrate matches the configured
// bitrate, restart pipelined and verify if Qos remains enforced
func testQosEnforcementRestart(t *testing.T, cfgCh chan string, restartCfg string) {
	tr := NewTestRunner(t)

	err := tr.RestartService("pipelined", true)
	fmt.Println("Restarting service")

	if err != nil {
		fmt.Printf("error restarting pipelined %v", err)
		assert.Fail(t, "failed restarting pipelined")
	}

	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)
	assert.NoError(t, usePCRFMockDriver())
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
		clearPCRFMockDriver()
	}()

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	imsi := ues[0].GetImsi()

	ki := rand.Intn(1000000)
	monitorKey := fmt.Sprintf("monitor-ULQos-%d", ki)
	ruleKey := fmt.Sprintf("static-ULQos-%d", ki)

	uplinkBwMax := uint32(1000000)
	qos := &models.FlowQos{MaxReqBwUl: &uplinkBwMax}
	rule := getStaticPassAll(ruleKey, monitorKey, 0, models.PolicyRuleTrackingTypeONLYPCRF, 3, qos)

	err = ruleManager.AddStaticRuleToDB(rule)
	assert.NoError(t, err)
	tr.WaitForPoliciesToSync()

	usageMonitorInfo := getUsageInformation(monitorKey, 15*MegaBytes)
	initRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGxCCAnswer(diam.Success).
		SetStaticRuleInstalls([]string{ruleKey}, []string{}).
		SetUsageMonitorInfo(usageMonitorInfo)
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	// On unexpected requests, just return the default update answer
	assert.NoError(t, setPCRFExpectations([]*protos.GxCreditControlExpectation{initExpectation},
		protos.NewGxCCAnswer(diam.Success)))

	tr.AuthenticateAndAssertSuccess(imsi)
	req := &cwfprotos.GenTrafficRequest{
		Imsi:   imsi,
		Volume: &wrappers.StringValue{Value: *swag.String("5M")},
	}
	verifyEgressRate(t, tr, req, float64(uplinkBwMax))

	// Assert that enforcement_stats rules are properly installed and the right
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+imsi][ruleKey]
	assert.NotNil(t, record, fmt.Sprintf("No policy usage record for imsi: %v", imsi))

	// restart the pipelined service and verify that Qos is still enforced
	oldCnt := tr.ScanContainerLogs("pipelined", "Starting pipelined")

	// modify pipelined yml to set clean_restart
	cfgCh <- restartCfg
	err = tr.RestartService("pipelined", true)
	if err != nil {
		fmt.Printf("error restarting pipelined %v", err)
		assert.Fail(t, "failed restarting pipelined")
	}
	waitForPipelinedRestart := func() bool {
		cnt := tr.ScanContainerLogs("pipelined", "Starting pipelined")
		fmt.Printf("curr restart count %d old count %d\n", cnt, oldCnt)
		return ((oldCnt + 1) == cnt)
	}
	assert.Eventually(t, waitForPipelinedRestart, time.Minute, 2*time.Second)

	// verify the egress rate after the restart of pipelined
	verifyEgressRate(t, tr, req, float64(uplinkBwMax))

	tr.DisconnectAndAssertSuccess(imsi)
	assert.NoError(t, err)
	time.Sleep(3 * time.Second)
}

func TestQosRestartMeter(t *testing.T) {
	t.Log("Running TestQosRestartMeter")
	cfgCh, err := configFileManager(pipelinedCfgFn)
	defer func() {
		close(cfgCh)
	}()

	if err != nil {
		t.Logf("failed modifying pipelined configs %v", err)
		t.Fail()
	}

	// clean restart test
	t.Log("Running TestQosRestartMeter - clean restart")
	testQosEnforcementRestart(t, cfgCh, cleanRestartYaml)

	// non clean restart test
	t.Log("Running TestQosRestartMeter - non clean restart")
	testQosEnforcementRestart(t, cfgCh, nonCleanRestartYaml)
}
