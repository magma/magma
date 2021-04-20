/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include "includes/GRPCReceiver.h"
#include <ostream>          // for operator<<, char_traits
#include "magma_logging.h"  // for MLOG

namespace magma {

void GRPCReceiver::rpc_response_loop() {
  running_ = true;
  void* tag;
  bool ok = false;
  while (running_) {
    if (!queue_.Next(&tag, &ok)) {
      return;
    }
    if (!ok) {
      MLOG(MINFO) << "gRPC receiver encountered error while processing request";
      continue;
    }
    static_cast<AsyncResponse*>(tag)->handle_response();
  }
}

void GRPCReceiver::stop() {
  running_ = false;
  queue_.Shutdown();
  // Pop all items in the queue until it is empty
  // https://github.com/grpc/grpc/issues/8610
  void* tag;
  bool ok;
  while (queue_.Next(&tag, &ok)) {
  }
}

}  // namespace magma
