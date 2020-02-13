/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package nghttpxlogger

import (
	"time"

	"magma/orc8r/cloud/go/services/logger"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/hpcloud/tail"
)

const (
	NGHTTPX_SCRIBE_CATEGORY = "perfpipe_magma_rest_api_stats"
	SAMPLING_RATE           = 1
)

type NghttpxLogger struct {
	readInterval time.Duration
	parser       NghttpxParser
}

// Nghttpx log message is in the following format:
//${time_iso8601}@|@${remote_addr}@|@${http_host}@|@${server_port}@|@${request}@|@${status}@|@${body_bytes_sent}bytes@|@${request_time}ms
// Scribe requires int and string(normal) messages to be separated
type NghttpxMessage struct {
	Int    map[string]int64  `json:"int"`
	Normal map[string]string `json:"normal"`
	Time   int64             `json:"time"`
}

// Returns a NghttpxLogger
// Potential loss of data:
// When logger restarts, the logger will drop all log entries which happened while it was down
func NewNghttpLogger(readInterval time.Duration, parser NghttpxParser) (*NghttpxLogger, error) {
	return &NghttpxLogger{readInterval: readInterval, parser: parser}, nil
}

func (nghttpxlogger *NghttpxLogger) Run(filepath string) {
	go nghttpxlogger.tail(filepath)
}

func (nghttpxlogger *NghttpxLogger) tail(filepath string) {
	t, err := tail.TailFile(filepath, tail.Config{Poll: true, Follow: true})
	if err != nil {
		glog.Errorf("Error opening file %v for tailing: %v\n", filepath, err)
		return
	}
	for line := range t.Lines {
		msg, err := nghttpxlogger.parser.Parse(line.Text)
		if err != nil {
			glog.Errorf("err parsing %s in nghttpx.log: %v\n", line.Text, err)
			continue
		}

		entries := []*protos.LogEntry{{
			Category:  NGHTTPX_SCRIBE_CATEGORY,
			NormalMap: msg.Normal,
			IntMap:    msg.Int,
			Time:      msg.Time,
		}}
		// TODO: change to a lower samplingRate as we see fit
		err = logger.LogToScribeWithSamplingRate(entries, SAMPLING_RATE)
		if err != nil {
			glog.Errorf("err sending nghttpx log: %v\n", err)
		}
	}
}
