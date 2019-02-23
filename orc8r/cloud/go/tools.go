// +build tools

/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

// Put all binary tool dependencies in here so they can be tracked by the go
// module.

import _ "magma/orc8r/cloud/go/tools/combine_swagger"
import _ "golang.org/x/lint/golint"
import _ "github.com/golang/protobuf/protoc-gen-go"
import _ "github.com/go-swagger/go-swagger/cmd/swagger"
