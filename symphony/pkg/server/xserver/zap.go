// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xserver

import (
	"go.uber.org/zap"
	"gocloud.dev/server/requestlog"
)

// A ZapLogger writes log entries to *zap.Logger.
type ZapLogger struct {
	logger *zap.Logger
}

// NewZapLogger returns a new logger that writes to logger.
func NewZapLogger(logger *zap.Logger) *ZapLogger {
	return &ZapLogger{logger}
}

// Log implements Logger interface.
func (l *ZapLogger) Log(ent *requestlog.Entry) {
	l.logger.Info("http request",
		zap.String("method", ent.RequestMethod),
		zap.String("url", ent.RequestURL),
		zap.Int("status", ent.Status),
		zap.String("user_agent", ent.UserAgent),
		zap.String("remote_ip", ent.RemoteIP),
		zap.String("server_ip", ent.ServerIP),
		zap.String("referer", ent.Referer),
		zap.Stringer("trace_id", ent.TraceID),
		zap.Stringer("span_id", ent.SpanID),
		zap.Duration("latency", ent.Latency),
		zap.Int64("bytes_in", ent.RequestHeaderSize+ent.RequestBodySize),
		zap.Int64("bytes_out", ent.ResponseHeaderSize+ent.ResponseBodySize),
	)
}
