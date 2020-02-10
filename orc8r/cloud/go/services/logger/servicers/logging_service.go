/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"errors"
	"fmt"

	"magma/orc8r/cloud/go/services/logger/exporters"
	"magma/orc8r/lib/go/protos"

	"golang.org/x/net/context"
)

type LoggingService struct {
	exporters map[protos.LoggerDestination]exporters.Exporter
}

func NewLoggingService(exporters map[protos.LoggerDestination]exporters.Exporter) (*LoggingService, error) {
	if exporters == nil {
		return nil, errors.New("exporters cannot be nil")
	}
	return &LoggingService{exporters: exporters}, nil
}

// Invoke the right exporter to export logMessages based on the loggerDestination specified by the client.
// Input: LogRequest which specifies a loggerDestination and a slice of logEntries
// Output: error if any
func (srv *LoggingService) Log(ctx context.Context, request *protos.LogRequest) (*protos.Void, error) {
	if request == nil {
		return new(protos.Void), errors.New("Empty LogRequest")
	}
	exporter, ok := srv.exporters[request.Destination]
	if !ok {
		return new(protos.Void),
			fmt.Errorf("LoggerDestination %v not supported", request.Destination)
	}
	err := exporter.Submit(request.Entries)
	return new(protos.Void), err
}
