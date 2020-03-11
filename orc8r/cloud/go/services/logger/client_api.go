/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package logger

import (
	"math/rand"

	"magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

const ServiceName = "LOGGER"

// getLoggerClient is a utility function to get a RPC connection to the
// loggingService service
func getLoggerClient() (protos.LoggingServiceClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := errors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewLoggingServiceClient(conn), err

}

// determines if we should log for the specific instance given samplingRate
func shouldLog(samplingRate float64) bool {
	return rand.Float64() < samplingRate
}

/////////////////////////////////////
// User call this directly to log //
////////////////////////////////////
func LogEntriesToDest(entries []*protos.LogEntry, destination protos.LoggerDestination, samplingRate float64) error {
	lg, err := getLoggerClient()
	if err != nil {
		return err
	}
	if !shouldLog(samplingRate) {
		return nil
	}
	req := protos.LogRequest{Entries: entries, Destination: destination}
	_, err = lg.Log(context.Background(), &req)
	return err

}

// Log entries to Scribe with SamplingRate 1 (i.e. no sampling)
func LogToScribe(entries []*protos.LogEntry) error {
	return LogEntriesToDest(entries, protos.LoggerDestination_SCRIBE, 1)
}

// Log entries to Scribe with given samplingRate
func LogToScribeWithSamplingRate(entries []*protos.LogEntry, samplingRate float64) error {
	return LogEntriesToDest(entries, protos.LoggerDestination_SCRIBE, samplingRate)
}
