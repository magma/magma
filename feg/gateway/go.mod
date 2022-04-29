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

// TODO remove golang.org/x/net line once Go Upgrade (https://github.com/magma/magma/pull/12151) is merged
replace (
	golang.org/x/net => golang.org/x/net v0.0.0-20210520170846-37e1c6afe023
	magma/feg/cloud/go => ../../feg/cloud/go
	magma/feg/cloud/go/protos => ../../feg/cloud/go/protos
	magma/gateway => ../../orc8r/gateway/go
	magma/lte/cloud/go => ../../lte/cloud/go
	magma/orc8r/cloud/go => ../../orc8r/cloud/go
	magma/orc8r/lib/go => ../../orc8r/lib/go
	magma/orc8r/lib/go/protos => ../../orc8r/lib/go/protos
)

require (
	github.com/deepmap/oapi-codegen v1.9.0
	github.com/emakeev/milenage v1.0.0
	github.com/emakeev/snowflake v0.0.0-20200206205012-767080b052fe
	github.com/envoyproxy/go-control-plane v0.10.1
	github.com/fiorix/go-diameter/v4 v4.0.4
	github.com/go-openapi/swag v0.19.15
	github.com/go-redis/redis v6.14.1+incompatible
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.1.2
	github.com/ishidawataru/sctp v0.0.0-20191218070446-00ab2ac2db07
	github.com/labstack/echo/v4 v4.2.1
	github.com/mennanov/fieldmask-utils v0.5.0
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.26.0
	github.com/shirou/gopsutil/v3 v3.21.5
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/thoas/go-funk v0.7.0
	github.com/wmnsk/go-gtp v0.8.1
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/genproto v0.0.0-20210524171403-669157292da3
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.28.0
	layeh.com/radius v0.0.0-20201203135236-838e26d0c9be
	magma/feg/cloud/go v0.0.0
	magma/feg/cloud/go/protos v0.0.0
	magma/gateway v0.0.0
	magma/lte/cloud/go v0.0.0
	magma/orc8r/cloud/go v0.0.0
	magma/orc8r/lib/go v0.0.0
	magma/orc8r/lib/go/protos v0.0.0
)

go 1.13
