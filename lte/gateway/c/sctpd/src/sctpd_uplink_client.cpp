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
#include <chrono>

#include "lte/gateway/c/sctpd/src/sctpd_uplink_client.h"
#include "lte/gateway/c/sctpd/src/util.h"

namespace magma {
namespace sctpd {

using grpc::ClientContext;

SctpdUplinkClient::SctpdUplinkClient(std::shared_ptr<Channel> channel) {
  _stub = SctpdUplink::NewStub(channel);
}

int SctpdUplinkClient::sendUl(const SendUlReq& req, SendUlRes* res) {
  assert(res != nullptr);

  ClientContext context;
  auto deadline = std::chrono::system_clock::now() +
                  std::chrono::milliseconds(1000 * RESPONSE_TIMEOUT);
  context.set_deadline(deadline);

  auto status = _stub->SendUl(&context, req, res);

  if (!status.ok()) {
    MLOG(MERROR) << "sctpul.sendul error";
    MLOG_grpcerr(status);
  }

  return status.ok() ? 0 : -1;
}

int SctpdUplinkClient::newAssoc(const NewAssocReq& req, NewAssocRes* res) {
  assert(res != nullptr);

  ClientContext context;
  auto deadline = std::chrono::system_clock::now() +
                  std::chrono::milliseconds(1000 * RESPONSE_TIMEOUT);
  context.set_deadline(deadline);

  auto status = _stub->NewAssoc(&context, req, res);

  if (!status.ok()) {
    MLOG(MERROR) << "sctpul.newassoc error";
    MLOG_grpcerr(status);
  }

  return status.ok() ? 0 : -1;
}

int SctpdUplinkClient::closeAssoc(const CloseAssocReq& req,
                                  CloseAssocRes* res) {
  assert(res != nullptr);

  ClientContext context;
  // Not putting a timeout event for closeAssoc events
  auto status = _stub->CloseAssoc(&context, req, res);

  if (!status.ok()) {
    MLOG(MERROR) << "sctpul.closeassoc error";
    MLOG_grpcerr(status);
  }

  return status.ok() ? 0 : -1;
}

}  // namespace sctpd
}  // namespace magma
