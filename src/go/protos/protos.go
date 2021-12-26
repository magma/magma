// Copyright 2021 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protos

//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

//go:generate protoc --go_out=. --go_opt=Msctpd.proto=magma/sctpd;sctpd -I ../../../lte/protos sctpd.proto
//go:generate protoc --go-grpc_out=. --go-grpc_opt=Msctpd.proto=magma/sctpd;sctpd -I ../../../lte/protos sctpd.proto

//go:generate go run github.com/golang/mock/mockgen -source magma/sctpd/sctpd_grpc.pb.go -destination magma/sctpd/mock_sctpd/mock_sctpd_grpc.pb.go

//go:generate protoc --go_out=. --go_opt=Mcommon.proto=magma/orc8r;orc8r --go_opt=Mdigest.proto=magma/orc8r;orc8r -I ../../../orc8r/protos common.proto digest.proto
//go:generate protoc --go_out=. --go_opt=Mapn.proto=magma/apn;apn -I ../../../lte/protos apn.proto
//go:generate protoc --go_out=. --go_opt=Msubscriberdb.proto=magma/subscriberdb;subscriberdb --go_opt=Morc8r/protos/common.proto=github.com/magma/magma/src/go/protos/magma/orc8r --go_opt=Morc8r/protos/digest.proto=github.com/magma/magma/src/go/protos/magma/orc8r  --go_opt=Mlte/protos/apn.proto=github.com/magma/magma/src/go/protos/magma/apn;apn -I ../../../lte/protos -I ../../../ subscriberdb.proto
//go:generate protoc --go_out=. --go_opt=Mmobilityd.proto=magma/mobilityd;mobilityd --go_opt=Morc8r/protos/common.proto=github.com/magma/magma/src/go/protos/magma/orc8r --go_opt=Mlte/protos/subscriberdb.proto=github.com/magma/magma/src/go/protos/magma/subscriberdb -I ../../../lte/protos -I ../../../ mobilityd.proto
//go:generate protoc --go_out=. --go_opt=Mpolicydb.proto=magma/policydb;policydb --go_opt=Morc8r/protos/common.proto=github.com/magma/magma/src/go/protos/magma/orc8r --go_opt=Mlte/protos/mobilityd.proto=github.com/magma/magma/src/go/protos/magma/mobilityd -I ../../../lte/protos -I ../../../ policydb.proto
//go:generate protoc --go_out=. --go_opt=Msession_manager.proto=magma/session_manager;session_manager --go_opt=Morc8r/protos/common.proto=github.com/magma/magma/src/go/protos/magma/orc8r --go_opt=Mlte/protos/policydb.proto=github.com/magma/magma/src/go/protos/magma/policydb --go_opt=Mlte/protos/apn.proto=github.com/magma/magma/src/go/protos/magma/apn --go_opt=Mlte/protos/subscriberdb.proto=github.com/magma/magma/src/go/protos/magma/subscriberdb  -I ../../../lte/protos -I ../../../ session_manager.proto
//go:generate protoc --go_out=. --go_opt=Mpipelined.proto=magma/pipelined;pipelined --go_opt=Morc8r/protos/common.proto=github.com/magma/magma/src/go/protos/magma/orc8r --go_opt=Mlte/protos/policydb.proto=github.com/magma/magma/src/go/protos/magma/policydb --go_opt=Mlte/protos/apn.proto=github.com/magma/magma/src/go/protos/magma/apn --go_opt=Mlte/protos/subscriberdb.proto=github.com/magma/magma/src/go/protos/magma/subscriberdb  --go_opt=Mlte/protos/session_manager.proto=github.com/magma/magma/src/go/protos/magma/session_manager --go_opt=Mlte/protos/mobilityd.proto=github.com/magma/magma/src/go/protos/magma/mobilityd -I ../../../lte/protos -I ../../../ pipelined.proto

//go:generate protoc --go-grpc_out=. --go-grpc_opt=Mcommon.proto=magma/orc8r;orc8r --go-grpc_opt=Mdigest.proto=magma/orc8r;orc8r -I ../../../orc8r/protos common.proto digest.proto
//go:generate protoc --go-grpc_out=. --go-grpc_opt=Mapn.proto=magma/apn;apn -I ../../../lte/protos apn.proto
//go:generate protoc --go-grpc_out=. --go-grpc_opt=Msubscriberdb.proto=magma/subscriberdb;subscriberdb --go-grpc_opt=Morc8r/protos/common.proto=github.com/magma/magma/src/go/protos/magma/orc8r --go-grpc_opt=Morc8r/protos/digest.proto=github.com/magma/magma/src/go/protos/magma/orc8r  --go-grpc_opt=Mlte/protos/apn.proto=github.com/magma/magma/src/go/protos/magma/apn;apn -I ../../../lte/protos -I ../../../ subscriberdb.proto
//go:generate protoc --go-grpc_out=. --go-grpc_opt=Mmobilityd.proto=magma/mobilityd;mobilityd --go-grpc_opt=Morc8r/protos/common.proto=github.com/magma/magma/src/go/protos/magma/orc8r --go-grpc_opt=Mlte/protos/subscriberdb.proto=github.com/magma/magma/src/go/protos/magma/subscriberdb -I ../../../lte/protos -I ../../../ mobilityd.proto
//go:generate protoc --go-grpc_out=. --go-grpc_opt=Mpolicydb.proto=magma/policydb;policydb --go-grpc_opt=Morc8r/protos/common.proto=github.com/magma/magma/src/go/protos/magma/orc8r --go-grpc_opt=Mlte/protos/mobilityd.proto=github.com/magma/magma/src/go/protos/magma/mobilityd -I ../../../lte/protos -I ../../../ policydb.proto
//go:generate protoc --go-grpc_out=. --go-grpc_opt=Msession_manager.proto=magma/session_manager;session_manager --go-grpc_opt=Morc8r/protos/common.proto=github.com/magma/magma/src/go/protos/magma/orc8r --go-grpc_opt=Mlte/protos/policydb.proto=github.com/magma/magma/src/go/protos/magma/policydb --go-grpc_opt=Mlte/protos/apn.proto=github.com/magma/magma/src/go/protos/magma/apn --go-grpc_opt=Mlte/protos/subscriberdb.proto=github.com/magma/magma/src/go/protos/magma/subscriberdb  -I ../../../lte/protos -I ../../../ session_manager.proto
//go:generate protoc --go-grpc_out=. --go-grpc_opt=Mpipelined.proto=magma/pipelined;pipelined --go-grpc_opt=Morc8r/protos/common.proto=github.com/magma/magma/src/go/protos/magma/orc8r --go-grpc_opt=Mlte/protos/policydb.proto=github.com/magma/magma/src/go/protos/magma/policydb --go-grpc_opt=Mlte/protos/apn.proto=github.com/magma/magma/src/go/protos/magma/apn --go-grpc_opt=Mlte/protos/subscriberdb.proto=github.com/magma/magma/src/go/protos/magma/subscriberdb  --go-grpc_opt=Mlte/protos/session_manager.proto=github.com/magma/magma/src/go/protos/magma/session_manager --go-grpc_opt=Mlte/protos/mobilityd.proto=github.com/magma/magma/src/go/protos/magma/mobilityd -I ../../../lte/protos -I ../../../ pipelined.proto

//go:generate go run github.com/golang/mock/mockgen -source magma/pipelined/pipelined_grpc.pb.go -destination magma/pipelined/mock_pipelined/mock_pipelined_grpc.pb.go

//go:generate protoc --go_out=. --go_opt=Mcapture.proto=magma/capture;capture -I magma/capture capture.proto
//go:generate protoc --go-grpc_out=. --go-grpc_opt=Mcapture.proto=magma/capture;capture -I magma/capture capture.proto

//go:generate protoc --go_out=. --go_opt=Mconfig.proto=magma/config;config -I magma/config config.proto
//go:generate protoc --go-grpc_out=. --go-grpc_opt=Mconfig.proto=magma/config;config -I magma/config config.proto
//go:generate go run github.com/golang/mock/mockgen -source magma/config/config_grpc.pb.go -destination magma/config/mock_config/mock_config_grpc.pb.go
