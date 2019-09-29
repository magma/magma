/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package ods

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"fbc/cwf/radius/monitoring/counters"

	"github.com/stretchr/testify/require"
	"go.opencensus.io/tag"
	"go.uber.org/zap"
)

// TestAnalyticsModulesAuthenticate tests the Analytics module handling of the Authenticate RADIUS packet
func TestSendOdsCounters(t *testing.T) {
	// Arrange
	logger, _ := zap.NewDevelopment()
	Init(&Config{
		ReportingPeriod: time.Duration(time.Second),
		GraphURL:        "http://127.0.0.1:1234/ods",
		Entity:          "entity",
		Category:        "123",
		Prefix:          "lalala",
	}, logger)

	var gotMetrics = make(chan bool, 1)

	http.HandleFunc("/ods", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			gotMetrics <- false
			return
		}
		t.Logf("got event: %s", body)
		gotMetrics <- true
	})

	go func() {
		if err := http.ListenAndServe(":1234", nil); err != nil {
			panic(err)
		}
	}()

	// Act
	tg, _ := tag.NewKey("test")
	op := counters.NewOperation("test", tg)
	op.Start()
	op.Success()
	time.Sleep(time.Duration(time.Second))

	// Assert
	timeout := time.NewTimer(5 * time.Second)
	select {
	case success := <-gotMetrics:
		require.Equal(t, true, success)
	case <-timeout.C:
		require.Fail(t, "timed out waiting for metrics to propagate")
	}
}
