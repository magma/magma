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
module magma/orc8r/cloud/go

replace (
	magma/gateway => ../../gateway/go
	magma/orc8r/lib/go => ../../lib/go
	magma/orc8r/lib/go/protos => ../../lib/go/protos
)

require (
	github.com/DATA-DOG/go-sqlmock v1.3.3
	github.com/Masterminds/squirrel v1.1.1-0.20190513200039-d13326f0be73
	github.com/emakeev/snowflake v0.0.0-20200206205012-767080b052fe
	github.com/facebookincubator/ent v0.0.0-20191128071424-29c7b0a0d805
	github.com/facebookincubator/prometheus-configmanager v0.0.0-20200717220759-a8282767b087
	github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/strfmt v0.19.4
	github.com/go-openapi/swag v0.19.5
	github.com/go-openapi/validate v0.19.3
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-swagger/go-swagger v0.21.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.3
	github.com/google/uuid v1.1.1
	github.com/labstack/echo v0.0.0-20181123063414-c54d9e8eed6c
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0
	github.com/lib/pq v1.2.0
	github.com/mattn/go-sqlite3 v1.11.0
	github.com/olivere/elastic/v7 v7.0.6
	github.com/pkg/errors v0.8.1
	github.com/prometheus/alertmanager v0.17.0
	github.com/prometheus/client_golang v1.2.1
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.7.0
	github.com/prometheus/procfs v0.0.5
	github.com/prometheus/prometheus v0.0.0-20190607092147-e23fa22233cf
	github.com/spf13/cobra v0.0.3
	github.com/stretchr/testify v1.4.0
	github.com/thoas/go-funk v0.7.0
	github.com/vektra/mockery v0.0.0-20181123154057-e78b021dcbb5
	golang.org/x/lint v0.0.0-20190409202823-959b441ac422
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	golang.org/x/tools v0.0.0-20191012152004-8de300cfc20a
	google.golang.org/grpc v1.27.1
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0
	gopkg.in/yaml.v2 v2.2.8
	magma/gateway v0.0.0
	magma/orc8r/lib/go v0.0.0-00010101000000-000000000000
	magma/orc8r/lib/go/protos v0.0.0
)

go 1.12
