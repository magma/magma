// Copyright (c) Facebook, Inc. and its affiliates.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
module magma/feg/cloud/go

replace (
	magma/feg/cloud/go/protos => ../../../feg/cloud/go/protos
	magma/lte/cloud/go => ../../../lte/cloud/go
	magma/orc8r/cloud/go => ../../../orc8r/cloud/go
)

require (
	github.com/go-openapi/errors v0.18.0
	github.com/go-openapi/strfmt v0.18.0
	github.com/go-openapi/swag v0.18.0
	github.com/go-openapi/validate v0.18.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.0
	github.com/stretchr/testify v1.3.0
	golang.org/x/net v0.0.0-20190301231341-16b79f2e4e95
	google.golang.org/grpc v1.19.0

	magma/feg/cloud/go/protos v0.0.0
	magma/lte/cloud/go v0.0.0
	magma/orc8r/cloud/go v0.0.0
)
