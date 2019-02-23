/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package registry

import (
	"sync"
	"time"

	"magma/orc8r/cloud/go/services/materializer"

	"github.com/golang/glog"
)

type materializerRegistry struct {
	sync.RWMutex
	applications map[string]materializer.Application
}

var registry = materializerRegistry{
	applications: make(map[string]materializer.Application),
}

// RegisterApplications registers a collection of applications with the
// materializer app registry. This function is thread safe.
func RegisterApplications(apps ...materializer.Application) {
	registry.Lock()
	defer registry.Unlock()
	for _, app := range apps {
		registerUnsafe(app)
	}
}

// RegisterApplication registers a MaterializerApplication with the service's
// application registry.
func RegisterApplication(app materializer.Application) {
	registry.Lock()
	defer registry.Unlock()
	registerUnsafe(app)
}

func registerUnsafe(app materializer.Application) {
	registry.applications[app.Name] = app
}

// RunAll is the function which runs all registered stream processors. It does not return unless an error
// is encountered
func RunAll() error {
	errorChan := make(chan error)
	timeout := 10 * time.Second
	registry.Lock()
	for _, application := range registry.applications {
		for _, processor := range application.Processors {
			go runProcessor(processor, errorChan, timeout)
		}
	}
	registry.Unlock()
	err := <-errorChan
	return err
}

// Run runs a single application specified by applicationName. It does not return unless an error is encountered.
func Run(applicationName string) error {
	errorChan := make(chan error)
	timeout := 10 * time.Second
	registry.Lock()
	for _, processor := range registry.applications[applicationName].Processors {
		go runProcessor(processor, errorChan, timeout)
	}
	registry.Unlock()
	err := <-errorChan
	return err
}

func runProcessor(processor materializer.StreamProcessor, errorChan chan error, timeout time.Duration) {
	err := processor.Run()
	if err != nil {
		select {
		case errorChan <- err:
		case <-time.After(timeout):
			glog.Warningf("Write to error channel timed out: %s", err)
		}
	}
}

// StopAll stops all registered stream processors
func StopAll() {
	registry.Lock()
	defer registry.Unlock()
	for _, application := range registry.applications {
		for _, processor := range application.Processors {
			processor.Stop()
		}
	}
}

// Stop stops the registered processors for a specific application
func Stop(applicationName string) {
	registry.Lock()
	defer registry.Unlock()
	for _, processor := range registry.applications[applicationName].Processors {
		processor.Stop()
	}
}
