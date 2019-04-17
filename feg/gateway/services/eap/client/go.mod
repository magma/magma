module magma/feg/gateway/services/eap/client

replace (
	magma/feg/cloud/go => ../../../../../feg/cloud/go
	magma/feg/cloud/go/protos => ../../../../../feg/cloud/go/protos
	magma/feg/gateway => ../../..

	magma/lte/cloud/go => ../../../../../lte/cloud/go
	magma/orc8r/cloud/go => ../../../../../orc8r/cloud/go
)

require (
	golang.org/x/net v0.0.0-20190311183353-d8887717615a
	google.golang.org/grpc v1.19.1
	magma/feg/cloud/go/protos v0.0.0
	magma/lte/cloud/go v0.0.0 // indirect
	magma/orc8r/cloud/go v0.0.0
)
