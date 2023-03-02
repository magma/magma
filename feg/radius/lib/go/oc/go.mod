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
	github.com/prometheus/client_golang v1.11.1
	github.com/stretchr/testify v1.5.1
	go.opencensus.io v0.21.0
	go.uber.org/zap v1.10.0
)

require (
	fbc/lib/go/log v0.0.0 // indirect
	github.com/apache/thrift v0.13.0 // indirect
	github.com/aws/aws-sdk-go v1.34.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/hashicorp/golang-lru v0.5.0 // indirect
	github.com/jmespath/go-jmespath v0.3.0 // indirect
	github.com/justinas/alice v0.0.0-20171023064455-03f45bd4b7da // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.26.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/stretchr/objx v0.1.1 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a // indirect
	golang.org/x/sys v0.0.0-20210603081109-ebe580a85c40 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	google.golang.org/api v0.3.2 // indirect
	google.golang.org/protobuf v1.26.0-rc.1 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)
