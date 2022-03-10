module fbc/cwf/radius

go 1.17

replace (
	fbc/lib/go/machine => ../lib/go/machine
	fbc/lib/go/radius => ../lib/go/radius
	magma/orc8r/lib/go => github.com/facebookincubator/magma/orc8r/lib/go v0.0.0-20200625002311-c12331e4dadc
	magma/orc8r/lib/go/protos => github.com/facebookincubator/magma/orc8r/lib/go/protos v0.0.0-20200625002311-c12331e4dadc
)

require (
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	fbc/lib/go/machine v0.0.0-00010101000000-000000000000
	fbc/lib/go/radius v0.0.0-00010101000000-000000000000
	github.com/alicebob/miniredis v2.5.0+incompatible
	github.com/donovanhide/eventsource v0.0.0-20171031113327-3ed64d21fb0b
	github.com/go-redis/redis v6.15.5+incompatible
	github.com/google/uuid v1.1.2
	github.com/mitchellh/mapstructure v1.1.2
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v0.9.3
	github.com/stretchr/testify v1.7.0
	go.opencensus.io v0.21.0
	go.uber.org/atomic v1.7.0
	go.uber.org/zap v1.18.1
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	google.golang.org/grpc v1.33.2
	google.golang.org/protobuf v1.26.0
	magma/orc8r/lib/go v0.0.0
)

require (
	github.com/alicebob/gopher-json v0.0.0-20180125190556-5a6b3ba71ee6 // indirect
	github.com/beorn7/perks v1.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gomodule/redigo v2.0.0+incompatible // indirect
	github.com/hashicorp/golang-lru v0.5.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.4.0 // indirect
	github.com/prometheus/procfs v0.0.0-20190507164030-5867b95ac084 // indirect
	github.com/stretchr/objx v0.1.1 // indirect
	github.com/yuin/gopher-lua v0.0.0-20190514113301-1cd887cd7036 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/sys v0.0.0-20191002091554-b397fe3ad8ed // indirect
	golang.org/x/text v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
