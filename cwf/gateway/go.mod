// Copyright (c) Facebook, Inc. and its affiliates.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
module magma/cwf/gateway

replace (
	magma/cwf/cloud/go => ../../cwf/cloud/go
	magma/feg/cloud/go => ../../feg/cloud/go
	magma/feg/gateway => ../../feg/gateway
	magma/lte/cloud/go => ../../lte/cloud/go
	magma/orc8r/cloud/go => ../../orc8r/cloud/go
)

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.3.0
	golang.org/x/net v0.0.0-20190110200230-915654e7eabc
	google.golang.org/grpc v1.17.0
	magma/cwf/cloud/go v0.0.0-00010101000000-000000000000
	magma/orc8r/cloud/go v0.0.0
)
