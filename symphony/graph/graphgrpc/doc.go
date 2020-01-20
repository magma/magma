// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate protoc --go_out=plugins=grpc,paths=source_relative:. rpc.proto

package graphgrpc
