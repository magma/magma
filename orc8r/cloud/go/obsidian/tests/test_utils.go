/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package tests

import (
	"sync"
	"testing"
	"time"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/access"
	access_tests "magma/orc8r/cloud/go/obsidian/access/tests"
	"magma/orc8r/cloud/go/obsidian/server"
	"magma/orc8r/lib/go/util"
)

const TEST_ADMIN_OPERATOR_ID = "Obsidian_Unit_Test_Admin_Operator"

var TestOperatorSerialNumber string

type Testcase struct {
	Name                      string
	Method                    string
	Url                       string
	Payload                   string
	Expected                  string
	Skip_payload_verification bool
	Expect_http_error_status  bool
}

func RunTest(t *testing.T, tst Testcase) (int, string, error) {
	t.Logf("\nTEST CASE: * %s *\n\t%s %s\n\tPayload: %s\n\tExpected: %s\n",
		tst.Name,
		tst.Method,
		tst.Url,
		tst.Payload,
		tst.Expected)
	status, response, err := SendHttpRequest(tst.Method, tst.Url, tst.Payload)
	if err != nil {
		t.Fatalf(
			"\n****** FAILED [%s]:: HTTP Request Error: %s, %d, %s  ******\n",
			tst.Name, err, status, response)
	} else {
		t.Logf("Result:\n\tStatus: %d\n\tReceived: %s\n", status, response)
		failedStatus := status >= 300 || status < 200
		if failedStatus != tst.Expect_http_error_status {
			t.Fatalf("\n\n****** %s::\n\tBad HTTP Status Code: %d\n******\n",
				tst.Name, status)
		} else if !(tst.Skip_payload_verification ||
			util.CompareJSON(response, tst.Expected)) {
			t.Fatalf(
				"\n\n****** %s::\nBAD RESPONSE:\n%s\nEXPECTED:\n%s\n******\n",
				tst.Name, response, tst.Expected)
		}
	}
	return status, response, err
}

var obsidianPort int
var lock sync.Mutex

func StartObsidian(t *testing.T) int {
	lock.Lock()
	defer lock.Unlock()

	if obsidianPort != 0 {
		return obsidianPort
	}
	// make sure to setup & enable access control before we start Obsidian & its
	// unit tests
	TestOperatorSerialNumber =
		access_tests.StartMockAccessControl(t, TEST_ADMIN_OPERATOR_ID)

	obsidian.Port = util.GetFreeTcpPort(obsidian.Port)
	if obsidian.Port == 0 {
		t.Fatalf("Failed to get a Free REST Server TCP Port")
	}
	obsidian.TLS = false

	// Start REST server
	go server.Start()
	time.Sleep(time.Millisecond * 100) // Some time for http server to start
	obsidianPort = obsidian.Port
	return obsidianPort
}

// SendHttpRequest is a wrapper to util.SendHttpRequest with extra test Admin
// Certificate HTTP Header which is needed to pass access control midleware
func SendHttpRequest(method, url, payload string) (int, string, error) {
	return util.SendHttpRequest(
		method,
		url,
		payload,
		[]string{access.CLIENT_CERT_SN_KEY, TestOperatorSerialNumber})
}
