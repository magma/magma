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

#include "sctpd_event_handler.h"

#include <lte/protos/sctpd.grpc.pb.h>

namespace magma {
namespace sctpd {

SctpdEventHandler::SctpdEventHandler(SctpdUplinkClient& client)
    : _client(client) {}

int SctpdEventHandler::HandleNewAssoc(
    uint32_t ppid, uint32_t assoc_id, uint32_t instreams, uint32_t outstreams,
    std::string& ran_cp_ipaddr) {
  NewAssocReq req;
  NewAssocRes res;

  req.set_ppid(ppid);
  req.set_assoc_id(assoc_id);
  req.set_instreams(instreams);
  req.set_outstreams(outstreams);
  req.set_ran_cp_ipaddr(ran_cp_ipaddr);

  return _client.newAssoc(req, &res);
}

void SctpdEventHandler::HandleCloseAssoc(
    uint32_t ppid, uint32_t assoc_id, bool reset) {
  CloseAssocReq req;
  CloseAssocRes res;

  req.set_ppid(ppid);
  req.set_assoc_id(assoc_id);
  req.set_is_reset(reset);

  _client.closeAssoc(req, &res);
}

void SctpdEventHandler::HandleRecv(
    uint32_t ppid, uint32_t assoc_id, uint32_t stream,
    const std::string& payload) {
  SendUlReq req;
  SendUlRes res;

  req.set_ppid(ppid);
  req.set_assoc_id(assoc_id);
  req.set_stream(stream);
  req.set_payload(payload);

  _client.sendUl(req, &res);
}

}  // namespace sctpd
}  // namespace magma
