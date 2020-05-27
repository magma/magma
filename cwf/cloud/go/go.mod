// Copyright (c) Facebook, Inc. and its affiliates.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
module magma/cwf/cloud/go

replace (
	magma/feg/cloud/go => ../../../feg/cloud/go
	magma/feg/cloud/go/protos => ../../../feg/cloud/go/protos
	magma/gateway => ../../../orc8r/gateway/go
	magma/lte/cloud/go => ../../../lte/cloud/go
	magma/orc8r/cloud/go => ../../../orc8r/cloud/go
	magma/orc8r/lib/go => ../../../orc8r/lib/go
	magma/orc8r/lib/go/protos => ../../../orc8r/lib/go/protos
)

require (
	github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/strfmt v0.19.4
	github.com/go-openapi/swag v0.18.0
	github.com/go-openapi/validate v0.18.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.3
	github.com/labstack/echo v0.0.0-20181123063414-c54d9e8eed6c
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.2.1
	github.com/prometheus/common v0.7.0
	github.com/robfig/cron/v3 v3.0.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55
	google.golang.org/grpc v1.27.1

	magma/feg/cloud/go v0.0.0
	magma/feg/cloud/go/protos v0.0.0
	magma/lte/cloud/go v0.0.0
	magma/orc8r/cloud/go v0.0.0
	magma/orc8r/lib/go v0.0.0
	magma/orc8r/lib/go/protos v0.0.0
)

go 1.13
