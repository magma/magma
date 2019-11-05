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
	github.com/creack/pty v1.1.9 // indirect
	github.com/go-openapi/swag v0.18.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/groupcache v0.0.0-20191025150517-4a4ac3fbac33 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/google/pprof v0.0.0-20191025152101-a8b9f9d2d3ce // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/jstemmer/go-junit-report v0.9.1 // indirect
	github.com/kr/pty v1.1.8 // indirect
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.3.0
	go.opencensus.io v0.22.1 // indirect
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/exp v0.0.0-20191024150812-c286b889502e // indirect
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8 // indirect
	golang.org/x/lint v0.0.0-20190930215403-16217165b5de // indirect
	golang.org/x/mobile v0.0.0-20191025110607-73ccc5ba0426 // indirect
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/sys v0.0.0-20191025090151-53bf42e6b339 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	golang.org/x/tools v0.0.0-20191025023517-2077df36852e // indirect
	golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898 // indirect
	google.golang.org/api v0.3.1 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	google.golang.org/grpc v1.21.1
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	magma/cwf/cloud/go v0.0.0-00010101000000-000000000000
	magma/feg/cloud/go v0.0.0
	magma/feg/cloud/go/protos v0.0.0
	magma/feg/gateway v0.0.0-00010101000000-000000000000
	magma/lte/cloud/go v0.0.0
	magma/orc8r/cloud/go v0.0.0
)

go 1.13
