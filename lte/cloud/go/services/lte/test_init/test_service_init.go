/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package test_init

import (
	"strings"
	"testing"

	"magma/lte/cloud/go/lte"
	lte_service "magma/lte/cloud/go/services/lte"
	"magma/lte/cloud/go/services/lte/servicers"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/streamer/protos"
	"magma/orc8r/cloud/go/test_utils"
)

func StartTestService(t *testing.T) {
	streams := []string{
		lte.SubscriberStreamName,
		lte.PolicyStreamName,
		lte.BaseNameStreamName,
		lte.MappingsStreamName,
		lte.NetworkWideRulesStreamName,
		lte.RatingGroupStreamName,
	}
	labels := map[string]string{
		orc8r.StreamProviderLabel: "true",
	}
	annotations := map[string]string{
		orc8r.StreamProviderStreamsAnnotation: strings.Join(streams, orc8r.AnnotationListSeparator),
	}
	srv, lis := test_utils.NewTestOrchestratorService(t, lte.ModuleName, lte_service.ServiceName, labels, annotations)
	protos.RegisterStreamProviderServer(srv.GrpcServer, servicers.NewLTEStreamProviderServicer())
	go srv.RunTest(lis)
}
