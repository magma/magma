module magma/feg/gateway/services/eap/client

replace (
	magma/feg/cloud/go => ../../../../../feg/cloud/go
	magma/feg/cloud/go/protos => ../../../../../feg/cloud/go/protos
	magma/feg/gateway => ../../..

	magma/lte/cloud/go => ../../../../../lte/cloud/go
	magma/orc8r/cloud/go => ../../../../../orc8r/cloud/go
)

require (
	magma/feg/cloud/go v0.0.0
	magma/feg/cloud/go/protos v0.0.0
	magma/lte/cloud/go v0.0.0
	magma/orc8r/cloud/go v0.0.0
)
