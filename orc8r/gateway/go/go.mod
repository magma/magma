module magma/gateway

// TODO remove golang.org/x/net line once Go Upgrade (https://github.com/magma/magma/pull/12151) is merged
replace (
	golang.org/x/net => golang.org/x/net v0.0.0-20210520170846-37e1c6afe023
	magma/orc8r/lib/go => ../../../orc8r/lib/go
	magma/orc8r/lib/go/protos => ../../../orc8r/lib/go/protos
)

require (
	github.com/aeden/traceroute v0.0.0-20181124220833-147686d9cb0f
	github.com/alicebob/gopher-json v0.0.0-20200520072559-a9ecdc9d1d3a // indirect
	github.com/alicebob/miniredis v2.5.0+incompatible
	github.com/emakeev/snowflake v0.0.0-20200206205012-767080b052fe
	github.com/go-redis/redis v6.14.1+incompatible
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.5.2
	github.com/gomodule/redigo v2.0.0+incompatible // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/moriyoshi/routewrapper v0.0.0-20180228100351-e52d8d14cf39
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.19.0 // indirect
	github.com/prometheus/client_model v0.2.0
	github.com/shirou/gopsutil/v3 v3.21.5
	github.com/stretchr/testify v1.7.0
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da // indirect
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f
	google.golang.org/grpc v1.43.0
	magma/orc8r/lib/go v0.0.0-00010101000000-000000000000
	magma/orc8r/lib/go/protos v0.0.0-00010101000000-000000000000
)

go 1.13
