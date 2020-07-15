/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

// File registry.go provides a stream provider registry by forwarding calls to
// the service registry.

package providers

import (
	"fmt"
	"strings"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

// GetStreamProvider gets the stream provider for a stream name.
// Returns an error if no provider has been registered for the stream.
func GetStreamProvider(streamName string) (StreamProvider, error) {
	if len(streamName) == 0 {
		return nil, errors.New("stream name cannot be empty string")
	}

	services := getServicesForStream(streamName)

	n := len(services)
	if n == 0 {
		return nil, fmt.Errorf("no stream providers found for stream name %s", streamName)
	}
	if n != 1 {
		glog.Warningf("Found %d stream providers for stream name %s", n, streamName)
	}

	return NewRemoteProvider(services[0], streamName), nil
}

func getServicesForStream(streamName string) []string {
	services := registry.FindServices(orc8r.StreamProviderLabel)

	var ret []string
	for _, s := range services {
		streamsVal, err := registry.GetAnnotation(s, orc8r.StreamProviderStreamsAnnotation)
		// Ignore annotation errors, since they indicate either
		//	- service registry contents were recently updated
		//	- this service has incorrect annotations given its label
		if err != nil {
			glog.Warningf("Received error getting annotation %s for service %s: %v", orc8r.StreamProviderStreamsAnnotation, s, err)
			continue
		}
		streams := strings.Split(streamsVal, orc8r.AnnotationListSeparator)
		if funk.Contains(streams, streamName) {
			ret = append(ret, s)
		}
	}

	return ret
}
