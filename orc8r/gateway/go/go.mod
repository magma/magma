module magma/orc8r/gateway

replace magma/orc8r/cloud/go => ../../../orc8r/cloud/go

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.2
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	google.golang.org/grpc v1.25.0

	magma/orc8r/cloud/go v0.0.0
)

go 1.13
