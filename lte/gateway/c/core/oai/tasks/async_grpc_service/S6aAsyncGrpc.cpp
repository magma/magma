/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#include "S6aGatewayImpl.h"
extern "C" {
#include "S6aAsyncGrpc.h"
}

void init_async_grpc_server(void) {
  magma::init_async_s6a_grpc_server();
}

void stop_async_grpc_service(void) {
  magma::stop_async_s6a_grpc_server();
}
