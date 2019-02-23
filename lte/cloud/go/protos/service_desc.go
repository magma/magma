/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package protos

import "google.golang.org/grpc"

func GetLegacySubscriberdbDesc() *grpc.ServiceDesc {
	return &grpc.ServiceDesc{
		ServiceName: "magma.SubscriberDBController",
		HandlerType: _SubscriberDBController_serviceDesc.HandlerType,
		Methods:     _SubscriberDBController_serviceDesc.Methods,
		Streams:     _SubscriberDBController_serviceDesc.Streams,
		Metadata:    _SubscriberDBController_serviceDesc.Metadata,
	}
}

func GetLegacyMeteringDesc() *grpc.ServiceDesc {
	return &grpc.ServiceDesc{
		ServiceName: "magma.MeteringdRecordsController",
		HandlerType: _MeteringdRecordsController_serviceDesc.HandlerType,
		Methods:     _MeteringdRecordsController_serviceDesc.Methods,
		Streams:     _MeteringdRecordsController_serviceDesc.Streams,
		Metadata:    _MeteringdRecordsController_serviceDesc.Metadata,
	}
}
