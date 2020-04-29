module magma/gateway

replace (
	magma/orc8r/lib/go => ../../../orc8r/lib/go
	magma/orc8r/lib/go/protos => ../../../orc8r/lib/go/protos
)

require (
	github.com/aeden/traceroute v0.0.0-20181124220833-147686d9cb0f
	github.com/emakeev/snowflake v0.0.0-20200206205012-767080b052fe
	github.com/go-openapi/strfmt v0.19.4
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.3
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/moriyoshi/routewrapper v0.0.0-20180228100351-e52d8d14cf39
	github.com/prometheus/client_model v0.2.0
	github.com/shirou/gopsutil v2.20.3+incompatible
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	google.golang.org/grpc v1.27.1
	magma/orc8r/lib/go v0.0.0-00010101000000-000000000000
	magma/orc8r/lib/go/protos v0.0.0-00010101000000-000000000000
)

go 1.13
