module fbc/lib/go/http

go 1.20

replace (
	fbc/lib/go/log => ../log
	fbc/lib/go/oc/helpers => ../oc/helpers
)

require (
	fbc/lib/go/log v0.0.0
	github.com/google/uuid v1.1.1
	github.com/justinas/alice v0.0.0-20171023064455-03f45bd4b7da
	github.com/stretchr/testify v1.3.0
	go.opencensus.io v0.21.0
	go.uber.org/multierr v1.1.0
	go.uber.org/zap v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/hashicorp/golang-lru v0.5.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.1.1 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
)
