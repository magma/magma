// Copyright 2020 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
module magma/feg/gateway

replace (
	magma/feg/cloud/go => ../../feg/cloud/go
	magma/feg/cloud/go/protos => ../../feg/cloud/go/protos
	magma/gateway => ../../orc8r/gateway/go
	magma/lte/cloud/go => ../../lte/cloud/go
	magma/orc8r/cloud/go => ../../orc8r/cloud/go
	magma/orc8r/lib/go => ../../orc8r/lib/go
	magma/orc8r/lib/go/protos => ../../orc8r/lib/go/protos
)

require (
	github.com/emakeev/milenage v1.0.0
	github.com/emakeev/snowflake v0.0.0-20200206205012-767080b052fe
	github.com/envoyproxy/go-control-plane v0.9.4
	github.com/fiorix/go-diameter/v4 v4.0.4
	github.com/go-openapi/swag v0.19.5
	github.com/go-redis/redis v6.14.1+incompatible
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.5.2
	github.com/ishidawataru/sctp v0.0.0-20191218070446-00ab2ac2db07
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.10.0
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.18.0
	github.com/shirou/gopsutil v2.20.3+incompatible
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/thoas/go-funk v0.7.0
	github.com/wmnsk/go-gtp v0.7.25
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/sys v0.0.0-20210426080607-c94f62235c83 // indirect
	google.golang.org/grpc v1.33.2
	google.golang.org/protobuf v1.26.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	layeh.com/radius v0.0.0-20200615152116-663b41c3bf86
	magma/feg/cloud/go v0.0.0
	magma/feg/cloud/go/protos v0.0.0
	magma/gateway v0.0.0
	magma/lte/cloud/go v0.0.0
	magma/orc8r/cloud/go v0.0.0
	magma/orc8r/lib/go v0.0.0
	magma/orc8r/lib/go/protos v0.0.0
)

go 1.13
