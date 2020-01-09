// Copyright (c) Facebook, Inc. and its affiliates.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
module magma/feg/gateway

replace (
	magma/feg/cloud/go => ../../feg/cloud/go
	magma/feg/cloud/go/protos => ../../feg/cloud/go/protos

	magma/lte/cloud/go => ../../lte/cloud/go
	magma/orc8r/cloud/go => ../../orc8r/cloud/go
	magma/orc8r/gateway => ../../orc8r/gateway/go
)

require (
	github.com/fiorix/go-diameter/v4 v4.0.1-0.20200102140015-a6c006d17e34
	github.com/go-redis/redis v6.14.1+incompatible
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.2
	github.com/gorilla/mux v1.6.2
	github.com/ishidawataru/sctp v0.0.0-20190922091402-408ec287e38c
	github.com/prometheus/client_golang v0.9.3-0.20190127221311-3c4408c8b829
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4
	github.com/prometheus/common v0.2.0
	github.com/shirou/gopsutil v2.18.10+incompatible
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20191007182048-72f939374954
	google.golang.org/grpc v1.25.0

	magma/feg/cloud/go v0.0.0
	magma/feg/cloud/go/protos v0.0.0

	magma/lte/cloud/go v0.0.0
	magma/orc8r/cloud/go v0.0.0
	magma/orc8r/gateway v0.0.0
)

go 1.13
