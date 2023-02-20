module fbc/lib/go/oc

go 1.19

replace (
	fbc/lib/go/http => ../http
	fbc/lib/go/log => ../log
	fbc/lib/go/oc/helpers => ../oc/helpers

)

require (
	contrib.go.opencensus.io/exporter/aws v0.0.0-20181029163544-2befc13012d0
	contrib.go.opencensus.io/exporter/jaeger v0.1.0
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	fbc/lib/go/http v0.0.0
	fbc/lib/go/oc/helpers v0.0.0
	github.com/jessevdk/go-flags v1.4.1-0.20181221193153-c0795c8afcf4
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/prometheus/client_golang v0.9.3-0.20190127221311-3c4408c8b829
	github.com/stretchr/testify v1.5.1
	go.opencensus.io v0.21.0
	go.uber.org/zap v1.10.0
)

require (
	fbc/lib/go/log v0.0.0 // indirect
	github.com/apache/thrift v0.12.0 // indirect
	github.com/aws/aws-sdk-go v1.34.0 // indirect
	github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.2.0 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/hashicorp/golang-lru v0.5.0 // indirect
	github.com/jmespath/go-jmespath v0.3.0 // indirect
	github.com/justinas/alice v0.0.0-20171023064455-03f45bd4b7da // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.0.0-20190115171406-56726106282f // indirect
	github.com/prometheus/common v0.2.0 // indirect
	github.com/prometheus/procfs v0.0.0-20190117184657-bf6a532e95b1 // indirect
	github.com/stretchr/objx v0.1.1 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	golang.org/x/sync v0.0.0-20190227155943-e225da77a7e6 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	google.golang.org/api v0.3.2 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
)
