/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package nghttpxlogger_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/logger/nghttpxlogger"
	"magma/orc8r/cloud/go/services/logger/test_init"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockParser struct {
}

func (parser *mockParser) Parse(str string) (*nghttpxlogger.NghttpxMessage, error) {
	res := strings.Split(str, " ")
	status := strings.Trim(res[0], "\x00")
	time, err := strconv.Atoi(res[1])
	if err != nil {
		time = 123456
	}
	return &nghttpxlogger.NghttpxMessage{
			Time:   int64(time),
			Normal: map[string]string{"status": status}},
		nil
}

func TestNghttpxLogger_Run(t *testing.T) {

	mockExporter := test_init.StartTestServiceWithMockExporterExposed(t)

	mockExporter.On("Submit", mock.AnythingOfType("[]*protos.LogEntry")).Return(nil)
	parser := mockParser{}
	logger, err := nghttpxlogger.NewNghttpLogger(time.Second, &parser)
	assert.NoError(t, err)
	// create temp file
	f, err := ioutil.TempFile("", "nghttpxlogger-test-")
	assert.NoError(t, err)
	fileName := f.Name()
	// start tailing logfile
	logger.Run(fileName)
	// logrotate the logfile after 5 seconds
	go SimulateLogRotation(fileName, t)
	// write lines to file
	for i := 0; i < 10; i++ {
		// this line has to be in format "<string> <int>" for testing purpose. See Parse() on mockParser.
		f.Write([]byte(fmt.Sprintf("testLine %v\n", i+1)))
		time.Sleep(time.Second)
	}
	// assert
	for i := 0; i < 10; i++ {
		msg := nghttpxlogger.NghttpxMessage{Time: int64(i + 1), Normal: map[string]string{"status": "testLine"}}
		logEntries := []*protos.LogEntry{
			{
				Category:  "perfpipe_magma_rest_api_stats",
				NormalMap: msg.Normal,
				Time:      msg.Time,
			},
		}
		mockExporter.AssertCalled(t, "Submit", logEntries)
	}
	mockExporter.AssertNumberOfCalls(t, "Submit", 10)
	f.Close()
	os.Remove(fileName)
	mockExporter.AssertExpectations(t)
}

func SimulateLogRotation(fileName string, t *testing.T) {
	time.Sleep(5 * time.Second)
	//copytruncate is used for logrotatino for nghttpx.log
	err := os.Truncate(fileName, 0)
	assert.NoError(t, err)
}
