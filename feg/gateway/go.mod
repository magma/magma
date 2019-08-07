// Copyright (c) Facebook, Inc. and its affiliates.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
module magma/feg/gateway

replace (
	github.com/fiorix/go-diameter => ./third-party/go/src/github.com/fiorix/go-diameter

	magma/feg/cloud/go => ../../feg/cloud/go
	magma/feg/cloud/go/protos => ../../feg/cloud/go/protos

	magma/lte/cloud/go => ../../lte/cloud/go
	magma/orc8r/cloud/go => ../../orc8r/cloud/go
	magma/orc8r/gateway => ../../orc8r/gateway/go
)

require (
	github.com/fiorix/go-diameter v3.0.3-0.20180924121357-70410bd9fce3+incompatible
	github.com/go-redis/redis v6.14.1+incompatible
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.2.0
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.6.2
	github.com/ishidawataru/sctp v0.0.0-20180918013207-6e2cb1366111
	github.com/prometheus/client_golang v0.9.2
	github.com/prometheus/client_model v0.0.0-20190109181635-f287a105a20e
	github.com/prometheus/common v0.2.0
	github.com/shirou/gopsutil v2.18.10+incompatible
	github.com/stretchr/testify v1.3.0
	golang.org/x/net v0.0.0-20190110200230-915654e7eabc
	google.golang.org/grpc v1.17.0

	magma/feg/cloud/go v0.0.0
	magma/feg/cloud/go/protos v0.0.0

	magma/lte/cloud/go v0.0.0
	magma/orc8r/cloud/go v0.0.0
	magma/orc8r/gateway v0.0.0
)
