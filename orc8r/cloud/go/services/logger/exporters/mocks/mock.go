/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package mocks

import (
	"encoding/json"
	"fmt"

	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/mock"
)

type mockExporter struct {
}

func NewMockExporter() *mockExporter {
	return &mockExporter{}
}

// prints out the marshaled result of logEntries
func (exporter *mockExporter) Submit(logEntries []*protos.LogEntry) error {
	logJson, err := json.Marshal(logEntries)
	if err != nil {
		return err
	}
	fmt.Printf("entries to Export in json: %v\n", string(logJson))
	return nil
}

// can assert methods called on this exporter
type ExposedMockExporter struct {
	mock.Mock
}

func NewExposedMockExporter() *ExposedMockExporter {
	return &ExposedMockExporter{}
}

func (exporter *ExposedMockExporter) Submit(logEntries []*protos.LogEntry) error {
	args := exporter.Called(logEntries)
	fmt.Printf("\n\nSUBMIT: %+v\n", logEntries)
	logJson, err := json.Marshal(logEntries)
	if err != nil {
		return args.Error(0)
	}
	fmt.Printf("entries to Export in json: %v\n", string(logJson))
	return args.Error(0)
}
