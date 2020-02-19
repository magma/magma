// Copyright (c) Facebook, Inc. and its affiliates.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
module magma/cwf/gateway

replace (
	fbc/cwf/radius => ../../feg/radius/src/
	fbc/lib/go => ../../feg/radius/lib/go
	fbc/lib/go/http => ../../feg/radius/lib/go/http
	fbc/lib/go/libgraphql => ../../feg/radius/lib/go/libgraphql
	fbc/lib/go/log => ../../feg/radius/lib/go/log
	fbc/lib/go/machine => ../../feg/radius/lib/go/machine
	fbc/lib/go/oc => ../../feg/radius/lib/go/oc
	fbc/lib/go/radius => ../../feg/radius/lib/go/radius

	magma/cwf/cloud/go => ../../cwf/cloud/go
	magma/feg/cloud/go => ../../feg/cloud/go
	magma/feg/cloud/go/protos => ../../feg/cloud/go/protos
	magma/feg/gateway => ../../feg/gateway
	magma/gateway => ../../orc8r/gateway/go
	magma/lte/cloud/go => ../../lte/cloud/go
	magma/orc8r/cloud/go => ../../orc8r/cloud/go
	magma/orc8r/lib/go => ../../orc8r/lib/go
	magma/orc8r/lib/go/protos => ../../orc8r/lib/go/protos
)

require (
	fbc/cwf/radius v0.0.0
	fbc/lib/go/radius v0.0.0-00010101000000-000000000000
	github.com/go-openapi/swag v0.18.0
	github.com/go-redis/redis v6.15.5+incompatible
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.3
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550 // indirect
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/sys v0.0.0-20191025090151-53bf42e6b339 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	google.golang.org/grpc v1.27.1
	magma/cwf/cloud/go v0.0.0-00010101000000-000000000000
	magma/feg/cloud/go/protos v0.0.0
	magma/feg/gateway v0.0.0-00010101000000-000000000000
	magma/lte/cloud/go v0.0.0
	magma/orc8r/cloud/go v0.0.0
	magma/orc8r/lib/go v0.0.0
	magma/orc8r/lib/go/protos v0.0.0
)

go 1.13
