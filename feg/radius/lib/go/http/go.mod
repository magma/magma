module fbc/lib/go/http

go 1.12

replace (
	fbc/lib/go/log => ../log
	fbc/lib/go/oc/helpers => ../oc/helpers
)

require (
	fbc/lib/go/log v0.0.0
	github.com/google/uuid v1.1.1
	github.com/justinas/alice v0.0.0-20171023064455-03f45bd4b7da
	github.com/pkg/errors v0.8.1
	github.com/stretchr/objx v0.1.1 // indirect
	github.com/stretchr/testify v1.3.0
	go.opencensus.io v0.21.0
	go.uber.org/multierr v1.1.0
	go.uber.org/zap v1.10.0
)
