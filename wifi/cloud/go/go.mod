module magma/wifi/cloud/go

replace (
	magma/gateway => ./../../../orc8r/gateway/go
	magma/orc8r/cloud/go => ./../../../orc8r/cloud/go
	magma/orc8r/lib/go => ./../../../orc8r/lib/go
	magma/orc8r/lib/go/protos => ./../../../orc8r/lib/go/protos
)

require (
	github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/strfmt v0.19.4
	github.com/go-openapi/swag v0.19.5
	github.com/go-openapi/validate v0.19.3
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.3
	github.com/labstack/echo v3.3.10+incompatible
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.5.1
	magma/orc8r/cloud/go v0.0.0
	magma/orc8r/lib/go v0.0.0
)

go 1.12
