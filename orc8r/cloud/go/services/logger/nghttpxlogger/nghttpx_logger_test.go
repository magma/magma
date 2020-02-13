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
	"sync"
	"testing"
	"time"

	"magma/orc8r/cloud/go/services/logger/nghttpxlogger"
	"magma/orc8r/cloud/go/services/logger/test_init"
	"magma/orc8r/lib/go/protos"

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
	defer func() {
		r := recover()
		_ = f.Close()
		_ = os.Remove(fileName)
		if r != nil {
			panic(r)
		}
	}()

	// start tailing logfile, logrotate after 3 seconds
	// a bit of a hack to prevent timing races with this test - we'll use a
	// mutex to prevent truncation from happening at the same time as a write
	// in real life, copytruncate log rotation could probably result in the
	// tailer losing a line, but we are ok with this.
	l := &sync.Mutex{}
	logger.Run(fileName)
	go SimulateLogRotation(t, fileName, l)

	// write lines to file
	for i := 0; i < 6; i++ {
		// this line has to be in format "<string> <int>" for testing purpose. See Parse() on mockParser.
		l.Lock()
		_, err := f.Write([]byte(fmt.Sprintf("testLine %v\n", i+1)))
		l.Unlock()
		assert.NoError(t, err)
		time.Sleep(time.Second)
	}

	// assert
	for i := 0; i < 6; i++ {
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
	mockExporter.AssertNumberOfCalls(t, "Submit", 6)
	mockExporter.AssertExpectations(t)
}

func SimulateLogRotation(t *testing.T, fileName string, lock *sync.Mutex) {
	lock.Lock()
	defer lock.Unlock()

	time.Sleep(3 * time.Second)
	//copytruncate is used for logrotation for nghttpx.log
	err := os.Truncate(fileName, 0)
	assert.NoError(t, err)
}
