module magma/orc8r/gateway

replace magma/orc8r/cloud/go => ../../../orc8r/cloud/go

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.2.0
	golang.org/x/net v0.0.0-20190311183353-d8887717615a
	google.golang.org/grpc v1.19.0

	magma/orc8r/cloud/go v0.0.0
)

go 1.13
