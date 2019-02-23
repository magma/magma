// Copyright (c) Facebook, Inc. and its affiliates.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
module magma/feg/cloud/go

replace (
	magma/lte/cloud/go => ../../../lte/cloud/go
	magma/orc8r/cloud/go => ../../../orc8r/cloud/go
)

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/fiorix/go-diameter v3.0.2+incompatible
	github.com/go-openapi/errors v0.18.0
	github.com/go-openapi/strfmt v0.18.0
	github.com/go-openapi/swag v0.18.0
	github.com/go-openapi/validate v0.18.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/mock v1.1.1
	github.com/golang/protobuf v1.2.0
	github.com/ishidawataru/sctp v0.0.0-20180918013207-6e2cb1366111
	github.com/pmezard/go-difflib v1.0.0
	github.com/stretchr/objx v0.1.1
	github.com/stretchr/testify v1.3.0
	golang.org/x/crypto v0.0.0-20190103213133-ff983b9c42bc
	golang.org/x/net v0.0.0-20190110200230-915654e7eabc
	golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be
	golang.org/x/text v0.3.0
	golang.org/x/tools v0.0.0-20181023010539-40a48ad93fbe
	google.golang.org/genproto v0.0.0-20190111180523-db91494dd46c
	google.golang.org/grpc v1.17.0

	magma/lte/cloud/go v0.0.0
	magma/orc8r/cloud/go v0.0.0
)
