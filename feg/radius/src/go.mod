module fbc/cwf/radius

replace (
	fbc/lib/go/http => ../lib/go/http
	fbc/lib/go/libgraphql => ../lib/go/libgraphql
	fbc/lib/go/log => ../lib/go/log
	fbc/lib/go/oc => ../lib/go/oc
	fbc/lib/go/oc/helpers => ../lib/go/oc/helpers
	fbc/lib/go/radius => ../lib/go/radius
)

require (
	fbc/lib/go/libgraphql v0.0.0-00010101000000-000000000000
	fbc/lib/go/log v0.0.0
	fbc/lib/go/oc v0.0.0
	fbc/lib/go/radius v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.3.1
	github.com/google/uuid v1.1.1
	github.com/mitchellh/mapstructure v1.1.2
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v0.9.3 // indirect
	github.com/stretchr/testify v1.3.0
	go.opencensus.io v0.21.0
	go.uber.org/zap v1.10.0
	golang.org/x/net v0.0.0-20190311183353-d8887717615a
	golang.org/x/xerrors v0.0.0-20190513163551-3ee3066db522
	google.golang.org/grpc v1.21.1
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	layeh.com/radius v0.0.0-20190322222518-890bc1058917
)
