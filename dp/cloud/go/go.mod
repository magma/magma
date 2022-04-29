// Copyright 2022 The Magma Authors.
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
module magma/dp/cloud/go

go 1.13

replace (
	golang.org/x/net => golang.org/x/net v0.0.0-20210520170846-37e1c6afe023
	magma/dp/cloud/go => ../../../dp/cloud/go
	magma/gateway => ../../../orc8r/gateway/go
	magma/orc8r/cloud/go => ../../../orc8r/cloud/go
	magma/orc8r/lib/go => ../../../orc8r/lib/go
	magma/orc8r/lib/go/protos => ../../../orc8r/lib/go/protos
)

require (
	github.com/Masterminds/squirrel v1.1.1-0.20190513200039-d13326f0be73
	github.com/go-openapi/errors v0.20.1
	github.com/go-openapi/strfmt v0.21.1
	github.com/go-openapi/swag v0.19.15
	github.com/go-openapi/validate v0.20.3
	github.com/golang/glog v1.0.0
	github.com/golang/protobuf v1.5.2
	github.com/labstack/echo v3.3.10+incompatible
	github.com/olivere/elastic/v7 v7.0.6
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.28.0
	magma/orc8r/cloud/go v0.0.0-00010101000000-000000000000
	magma/orc8r/lib/go v0.0.0-00010101000000-000000000000
	magma/orc8r/lib/go/protos v0.0.0
)
