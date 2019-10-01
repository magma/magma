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
	magma/lte/cloud/go => ../../lte/cloud/go
	magma/orc8r/cloud/go => ../../orc8r/cloud/go
	magma/orc8r/gateway => ../../orc8r/gateway/go
)

require (
	fbc/cwf/radius v0.0.0
	fbc/lib/go/radius v0.0.0-00010101000000-000000000000

	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.1
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/onsi/ginkgo v1.7.0 // indirect
	github.com/onsi/gomega v1.4.3 // indirect
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.3.0
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	google.golang.org/grpc v1.21.1
	magma/cwf/cloud/go v0.0.0-00010101000000-000000000000
	magma/feg/cloud/go v0.0.0
	magma/feg/cloud/go/protos v0.0.0
	magma/feg/gateway v0.0.0-00010101000000-000000000000
	magma/lte/cloud/go v0.0.0
	magma/orc8r/cloud/go v0.0.0
)
