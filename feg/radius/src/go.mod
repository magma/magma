module fbc/cwf/radius

replace (
	fbc/lib/go/machine => ../lib/go/machine
	fbc/lib/go/radius => ../lib/go/radius
)

require (
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	fbc/lib/go/machine v0.0.0-00010101000000-000000000000
	fbc/lib/go/radius v0.0.0-00010101000000-000000000000
	github.com/alicebob/gopher-json v0.0.0-20180125190556-5a6b3ba71ee6 // indirect
	github.com/alicebob/miniredis v2.5.0+incompatible
	github.com/cosiner/argv v0.0.1 // indirect
	github.com/donovanhide/eventsource v0.0.0-20171031113327-3ed64d21fb0b
	github.com/go-delve/delve v1.3.1 // indirect
	github.com/go-redis/redis v6.15.5+incompatible
	github.com/golang/protobuf v1.3.1
	github.com/gomodule/redigo v2.0.0+incompatible // indirect
	github.com/google/uuid v1.1.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.9 // indirect
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/peterh/liner v1.1.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v0.9.3
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/spf13/cobra v0.0.5 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/yuin/gopher-lua v0.0.0-20190514113301-1cd887cd7036 // indirect
	go.opencensus.io v0.21.0
	go.starlark.net v0.0.0-20190919145610-979af19b165c // indirect
	go.uber.org/atomic v1.4.0
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0
	golang.org/x/arch v0.0.0-20190927153633-4e8777c89be4 // indirect
	golang.org/x/crypto v0.0.0-20191002192127-34f69633bfdc // indirect
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	golang.org/x/sys v0.0.0-20191002091554-b397fe3ad8ed // indirect
	google.golang.org/grpc v1.21.1
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

go 1.13
