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

//go:generate protoc --go_out=. --go_opt=Mcommon.proto=magma/orc8r;orc8r -I ../../../orc8r/protos common.proto
//go:generate protoc --go_out=. --go_opt=Mmconfigs.proto=magma/mconfig;mconfig --go_opt=Morc8r/protos/common.proto=github.com/magma/magma/src/go/protos/magma/orc8r -I ../../../lte/protos/mconfig -I ../../../ mconfigs.proto
//go:generate protoc --go-grpc_out=. --go-grpc_opt=Mmconfigs.proto=magma/mconfig;mconfig -I ../../../lte/protos/mconfig -I ../../../ mconfigs.proto

//go:generate protoc --go_out=. --go_opt=Mcapture.proto=magma/capture;capture -I magma/capture capture.proto
