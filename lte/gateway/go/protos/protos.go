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

//go:generate protoc --go_out=. --go_opt=Msctpd.proto=magma/sctpd;sctpd -I ../../../protos sctpd.proto
//go:generate protoc --go-grpc_out=. --go-grpc_opt=Msctpd.proto=magma/sctpd;sctpd -I ../../../protos sctpd.proto
