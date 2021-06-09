module magma/orc8r/cloud/test

go 1.13

replace (
	// magma/feg/cloud/go => ../../../feg/cloud/go
	// magma/feg/cloud/go/protos => ../../../feg/cloud/go/protos
	magma/gateway => ../../gateway/go
	// magma/lte/cloud/go => ../../../lte/cloud/go
	magma/orc8r/cloud/go => ../go
	magma/orc8r/lib/go => ../../lib/go
	magma/orc8r/lib/go/protos => ../../lib/go/protos
)

require (
	github.com/go-openapi/runtime v0.19.28
	github.com/go-openapi/swag v0.19.9
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8
	magma/orc8r/cloud/go v0.0.0-00010101000000-000000000000
)
