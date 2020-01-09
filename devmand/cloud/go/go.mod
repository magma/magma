// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.
//
module orc8r/devmand/cloud/go

replace magma/orc8r/cloud/go => ../../../orc8r/cloud/go

require (
	github.com/go-openapi/errors v0.18.0
	github.com/go-openapi/strfmt v0.18.0
	github.com/go-openapi/swag v0.18.0
	github.com/go-openapi/validate v0.18.0
	github.com/golang/protobuf v1.3.2
	github.com/labstack/echo v0.0.0-20181123063414-c54d9e8eed6c
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.4.0

	magma/orc8r/cloud/go v0.0.0
)

go 1.13
