// +build all qos

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
	"math/rand"
	"os"
	"testing"
	"time"

	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/feg/cloud/go/protos"
	"magma/lte/cloud/go/plugin/models"

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

	// do not use restartPipeline functon. Otherwise we are not testing the case where attach
	// comes while pipelined is still rebooting.
	err := tr.RestartService("pipelined")
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

	uplinkBwMax := uint32(100000)
	qos := &models.FlowQos{MaxReqBwUl: &uplinkBwMax}
	rule := getStaticPassAll(ruleKey, monitorKey, 0, models.PolicyRuleTrackingTypeONLYPCRF, 3, qos)

	err = ruleManager.AddStaticRuleToDB(rule)
	assert.NoError(t, err)
	tr.WaitForPoliciesToSync()

	usageMonitorInfo := getUsageInformation(monitorKey, 1.5*MegaBytes)
	initRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGxCCAnswer(diam.Success).
		SetStaticRuleInstalls([]string{ruleKey}, []string{}).
		SetUsageMonitorInfo(usageMonitorInfo)
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	// On unexpected requests, just return the default update answer
	assert.NoError(t, setPCRFExpectations([]*protos.GxCreditControlExpectation{initExpectation},
		protos.NewGxCCAnswer(diam.Success)))

	tr.AuthenticateAndAssertSuccessWithRetries(imsi, 5)
	req := &cwfprotos.GenTrafficRequest{
		Imsi:   imsi,
		Volume: &wrappers.StringValue{Value: *swag.String("500k")},
	}
	verifyEgressRate(t, tr, req, float64(uplinkBwMax))

	// Assert that enforcement_stats rules are properly installed and the right
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+imsi][ruleKey]
	assert.NotNil(t, record, fmt.Sprintf("No policy usage record for imsi: %v", imsi))

	// modify pipelined yml to set clean_restart
	cfgCh <- restartCfg

	restartPipelined(t, tr)

	// verify the egress rate after the restart of pipelined
	verifyEgressRate(t, tr, req, float64(uplinkBwMax))

	tr.DisconnectAndAssertSuccess(imsi)
	assert.NoError(t, err)
	time.Sleep(3 * time.Second)
}

func restartPipelined(t *testing.T, tr *TestRunner) {
	oldCount := tr.ScanContainerLogs("pipelined", "Starting pipelined")
	err := tr.RestartService("pipelined")
	if err != nil {
		fmt.Printf("error restarting pipelined %v", err)
		assert.Fail(t, "failed restarting pipelined")
	}
	waitForPipelinedRestart := func() bool {
		cnt := tr.ScanContainerLogs("pipelined", "Starting pipelined")
		fmt.Printf("curr restart count %d old count %d\n", cnt, oldCount)
		return ((oldCount + 1) == cnt)
	}
	assert.Eventually(t, waitForPipelinedRestart, time.Minute, 2*time.Second)
}

func TestQosRestartMeterClean(t *testing.T) {
	fmt.Println("\nRunning TestQosRestartMeterClean...")
	cfgCh, err := configFileManager(pipelinedCfgFn)
	defer func() {
		close(cfgCh)
		// add a additional second for original file to be syncd
		time.Sleep(time.Second)
	}()
	if err != nil {
		t.Logf("failed modifying pipelined configs %v", err)
		t.Fail()
	}
	// clean restart test
	testQosEnforcementRestart(t, cfgCh, cleanRestartYaml)
}

func TestQosRestartMeterNonClean(t *testing.T) {
	fmt.Println("\nRunning TestQosRestartMeterNonClean...")
	cfgCh, err := configFileManager(pipelinedCfgFn)
	defer func() {
		close(cfgCh)
		// add a additional second for original file to be syncd
		time.Sleep(time.Second)
	}()
	if err != nil {
		t.Logf("failed modifying pipelined configs %v", err)
		t.Fail()
	}
	// non clean restart test
	testQosEnforcementRestart(t, cfgCh, nonCleanRestartYaml)
}
