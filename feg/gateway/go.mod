// Copyright (c) Facebook, Inc. and its affiliates.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
module magma/feg/gateway

replace (
	github.com/fiorix/go-diameter/v4 => github.com/emakeev/go-diameter/v4 v4.0.4

	magma/feg/cloud/go => ../../feg/cloud/go
	magma/feg/cloud/go/protos => ../../feg/cloud/go/protos
	magma/gateway => ../../orc8r/gateway/go
	magma/lte/cloud/go => ../../lte/cloud/go
	magma/orc8r/cloud/go => ../../orc8r/cloud/go
	magma/orc8r/lib/go => ../../orc8r/lib/go
	magma/orc8r/lib/go/protos => ../../orc8r/lib/go/protos
)

require (
	github.com/fiorix/go-diameter/v4 v4.0.1-0.20200120193412-55a1c21738f9
	github.com/go-openapi/swag v0.18.0
	github.com/go-redis/redis v6.14.1+incompatible
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.3
	github.com/ishidawataru/sctp v0.0.0-20190922091402-408ec287e38c
	github.com/prometheus/client_golang v0.9.3-0.20190127221311-3c4408c8b829
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.2.0
	github.com/shirou/gopsutil v2.18.10+incompatible
	github.com/stretchr/testify v1.4.0
	github.com/thoas/go-funk v0.4.0
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	google.golang.org/grpc v1.27.1

	magma/feg/cloud/go v0.0.0
	magma/feg/cloud/go/protos v0.0.0
	magma/gateway v0.0.0

	magma/lte/cloud/go v0.0.0
	magma/orc8r/cloud/go v0.0.0
	magma/orc8r/lib/go v0.0.0
	magma/orc8r/lib/go/protos v0.0.0
)

go 1.13
