module magma/feg/cloud/go/protos

require (
	github.com/golang/protobuf v1.3.3
	google.golang.org/grpc v1.27.1

	magma/lte/cloud/go v0.0.0
	magma/orc8r/lib/go/protos v0.0.0
)

replace (
	magma/gateway => ../../../../orc8r/gateway/go
	magma/lte/cloud/go => ../../../../lte/cloud/go
	magma/orc8r/cloud/go => ../../../../orc8r/cloud/go
	magma/orc8r/lib/go => ../../../../orc8r/lib/go
	magma/orc8r/lib/go/protos => ../../../../orc8r/lib/go/protos
)

go 1.12
