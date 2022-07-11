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
module magma/lte/cloud/go

replace (
	magma/feg/cloud/go => ../../../feg/cloud/go
	magma/feg/cloud/go/protos => ../../../feg/cloud/go/protos
	magma/gateway => ../../../orc8r/gateway/go
	magma/lte/cloud/go => ../../../lte/cloud/go
	magma/orc8r/cloud/go => ../../../orc8r/cloud/go
	magma/orc8r/lib/go => ../../../orc8r/lib/go
	magma/orc8r/lib/go/protos => ../../../orc8r/lib/go/protos
)

require (
	github.com/Masterminds/squirrel v1.1.1-0.20190513200039-d13326f0be73
	github.com/go-openapi/errors v0.20.2
	github.com/go-openapi/strfmt v0.21.1
	github.com/go-openapi/swag v0.19.15
	github.com/go-openapi/validate v0.20.3
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/gogf/gf v1.16.6
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.6
	github.com/google/uuid v1.1.2
	github.com/hashicorp/go-multierror v1.1.1
	github.com/influxdata/tdigest v0.0.1
	github.com/labstack/echo/v4 v4.7.2
	github.com/lib/pq v1.2.0
	github.com/magma/milenage v1.0.2
	github.com/olivere/elastic/v7 v7.0.6
	github.com/prometheus/client_golang v1.5.1
	github.com/prometheus/common v0.9.1
	github.com/stretchr/testify v1.7.1
	github.com/thoas/go-funk v0.7.0
	github.com/warthog618/sms v0.3.0
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f
	google.golang.org/genproto v0.0.0-20211208223120-3a66f561d7aa
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0
	magma/feg/cloud/go v0.0.0
	magma/feg/cloud/go/protos v0.0.0
	magma/orc8r/cloud/go v0.0.0
	magma/orc8r/lib/go v0.0.0
	magma/orc8r/lib/go/protos v0.0.0
)

require (
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bmatcuk/doublestar v1.3.4 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-openapi/analysis v0.21.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/loads v0.21.0 // indirect
	github.com/go-openapi/runtime v0.21.1 // indirect
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/jxskiss/base62 v1.0.0 // indirect
	github.com/labstack/gommon v0.3.1 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-sqlite3 v1.14.9 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/procfs v0.0.8 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	go.mongodb.org/mongo-driver v1.8.2 // indirect
	go.opentelemetry.io/otel v1.0.0-RC2 // indirect
	go.opentelemetry.io/otel/trace v1.0.0-RC2 // indirect
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce // indirect
	golang.org/x/sys v0.0.0-20220319134239-a9b59b0215f8 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.8 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	magma/gateway v0.0.0 // indirect
)

go 1.18
