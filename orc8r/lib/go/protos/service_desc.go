/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package protos

import "google.golang.org/grpc"

func GetLegacyBootstrapperDesc() *grpc.ServiceDesc {
	return &grpc.ServiceDesc{
		ServiceName: "magma.Bootstrapper",
		HandlerType: _Bootstrapper_serviceDesc.HandlerType,
		Methods:     _Bootstrapper_serviceDesc.Methods,
		Streams:     _Bootstrapper_serviceDesc.Streams,
		Metadata:    _Bootstrapper_serviceDesc.Metadata,
	}
}

func GetLegacyDispatcherDesc() *grpc.ServiceDesc {
	return &grpc.ServiceDesc{
		ServiceName: "magma.SyncRPCService",
		HandlerType: _SyncRPCService_serviceDesc.HandlerType,
		Methods:     _SyncRPCService_serviceDesc.Methods,
		Streams:     _SyncRPCService_serviceDesc.Streams,
		Metadata:    _SyncRPCService_serviceDesc.Metadata,
	}
}

func GetLegacyDirectorydDesc() *grpc.ServiceDesc {
	return &grpc.ServiceDesc{
		ServiceName: "magma.DirectoryService",
		HandlerType: _DirectoryService_serviceDesc.HandlerType,
		Methods:     _DirectoryService_serviceDesc.Methods,
		Streams:     _DirectoryService_serviceDesc.Streams,
		Metadata:    _DirectoryService_serviceDesc.Metadata,
	}
}

func GetLegacyLoggerDesc() *grpc.ServiceDesc {
	return &grpc.ServiceDesc{
		ServiceName: "magma.LoggingService",
		HandlerType: _LoggingService_serviceDesc.HandlerType,
		Methods:     _LoggingService_serviceDesc.Methods,
		Streams:     _LoggingService_serviceDesc.Streams,
		Metadata:    _LoggingService_serviceDesc.Metadata,
	}
}

func GetLegacyMetricsdDesc() *grpc.ServiceDesc {
	return &grpc.ServiceDesc{
		ServiceName: "magma.MetricsController",
		HandlerType: _MetricsController_serviceDesc.HandlerType,
		Methods:     _MetricsController_serviceDesc.Methods,
		Streams:     _MetricsController_serviceDesc.Streams,
		Metadata:    _MetricsController_serviceDesc.Metadata,
	}
}

func GetLegacyStreamerDesc() *grpc.ServiceDesc {
	return &grpc.ServiceDesc{
		ServiceName: "magma.Streamer",
		HandlerType: _Streamer_serviceDesc.HandlerType,
		Methods:     _Streamer_serviceDesc.Methods,
		Streams:     _Streamer_serviceDesc.Streams,
		Metadata:    _Streamer_serviceDesc.Metadata,
	}
}
