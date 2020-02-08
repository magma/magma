module magma/gateway

replace (
	magma/orc8r/lib/go => ../../../orc8r/lib/go
	magma/orc8r/lib/go/protos => ../../../orc8r/lib/go/protos
)

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.3
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	google.golang.org/grpc v1.27.1
	magma/orc8r/lib/go v0.0.0-00010101000000-000000000000
	magma/orc8r/lib/go/protos v0.0.0-00010101000000-000000000000
)

go 1.13
