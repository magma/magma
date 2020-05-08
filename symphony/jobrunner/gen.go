// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jobrunner

//go:generate go run github.com/google/addlicense -c Facebook -y 2004-present -l bsd ./
//go:generate protoc --go_out=plugins=grpc,paths=source_relative:. rpc.proto
