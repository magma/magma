module magma/orc8r/cloud/test

go 1.13

replace (
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191016111102-bec269661e48
	magma/gateway => ../../gateway/go
	magma/orc8r/cloud/api/v1/go => ../api/v1/go
	magma/orc8r/cloud/go => ../go
	magma/orc8r/lib/go => ../../lib/go
	magma/orc8r/lib/go/protos => ../../lib/go/protos
)

require (
	github.com/go-openapi/runtime v0.19.5
	github.com/go-openapi/strfmt v0.19.4 // indirect
	github.com/go-openapi/swag v0.19.15
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.3 // indirect
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/googleapis/gnostic v0.2.0 // indirect
	github.com/mitchellh/mapstructure v1.3.2 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.5.1
	github.com/prometheus/common v0.9.1
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.0.0-20191016110408-35e52d86657a
	k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	magma/orc8r/cloud/api/v1/go v0.0.0-00010101000000-000000000000
)
