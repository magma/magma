module fbc/lib/go/oc

go 1.12

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
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v0.9.3-0.20190127221311-3c4408c8b829
	github.com/smartystreets/goconvey v0.0.0-20190330032615-68dc04aab96a // indirect
	github.com/stretchr/testify v1.3.0
	go.opencensus.io v0.21.0
	go.uber.org/zap v1.10.0
)
