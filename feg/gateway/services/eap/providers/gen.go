/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package providers encapsulates supported EAP Authenticator Providers
//
//go:generate protoc -I ../protos -I . --go_out=plugins=grpc,paths=source_relative:. protos/eap_provider.proto
package providers
