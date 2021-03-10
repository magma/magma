/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

extern "C" {
#include "sctpd_uplink_server.h"

// #include "assertions.h"
#include "bstrlib.h"
#include "log.h"

#include "sctp_defs.h"
#include "sctp_itti_messaging.h"
}

#include <memory>

#include <grpcpp/grpcpp.h>

#include <lte/protos/sctpd.grpc.pb.h>

namespace magma {
namespace mme {

using grpc::ServerContext;
using grpc::Status;

using magma::sctpd::CloseAssocReq;
using magma::sctpd::CloseAssocRes;
using magma::sctpd::NewAssocReq;
using magma::sctpd::NewAssocRes;
using magma::sctpd::SctpdUplink;
using magma::sctpd::SendUlReq;
using magma::sctpd::SendUlRes;

class SctpdUplinkImpl final : public SctpdUplink::Service {
 public:
  SctpdUplinkImpl();

  Status SendUl(
      ServerContext* context, const SendUlReq* req, SendUlRes* res) override;
  Status NewAssoc(
      ServerContext* context, const NewAssocReq* req,
      NewAssocRes* res) override;
  Status CloseAssoc(
      ServerContext* context, const CloseAssocReq* req,
      CloseAssocRes* res) override;
};

SctpdUplinkImpl::SctpdUplinkImpl() {}

Status SctpdUplinkImpl::SendUl(
    ServerContext* context, const SendUlReq* req, SendUlRes* res) {
  bstring payload;
  uint32_t ppid;
  uint32_t assoc_id;
  uint16_t stream;

  payload = blk2bstr(req->payload().c_str(), req->payload().size());
  if (payload == NULL) {
    OAILOG_ERROR(LOG_SCTP, "failed to allocate bstr for SendUl\n");
    return Status::OK;
  }

  ppid     = req->ppid();
  assoc_id = req->assoc_id();
  stream   = req->stream();

  if (sctp_itti_send_new_message_ind(&payload, ppid, assoc_id, stream) < 0) {
    OAILOG_ERROR(LOG_SCTP, "failed to send new_message_ind for SendUl\n");
    return Status::OK;
  }

  return Status::OK;
}

#include <assert.h>

Status SctpdUplinkImpl::NewAssoc(
    ServerContext* context, const NewAssocReq* req, NewAssocRes* res) {
  uint32_t ppid       = req->ppid();
  uint32_t assoc_id   = req->assoc_id();
  uint16_t instreams  = req->instreams();
  uint16_t outstreams = req->outstreams();
  bstring ran_cp_ipaddr =
      blk2bstr(req->ran_cp_ipaddr().c_str(), req->ran_cp_ipaddr().size());

  if (sctp_itti_send_new_association(
          ppid, assoc_id, instreams, outstreams, &ran_cp_ipaddr) < 0) {
    OAILOG_ERROR(LOG_SCTP, "failed to send new_association for NewAssoc\n");
    return Status::OK;
  }

  return Status::OK;
}

Status SctpdUplinkImpl::CloseAssoc(
    ServerContext* context, const CloseAssocReq* req, CloseAssocRes* res) {
  uint32_t ppid;
  uint32_t assoc_id;
  bool reset;

  ppid     = req->ppid();
  assoc_id = req->assoc_id();
  reset    = req->is_reset();

  if (sctp_itti_send_com_down_ind(ppid, assoc_id, reset) < 0) {
    OAILOG_ERROR(LOG_SCTP, "failed to send com_down_ind for CloseAssoc\n");
    return Status::OK;
  }

  return Status::OK;
}

}  // namespace mme
}  // namespace magma

using grpc::Server;
using grpc::ServerBuilder;

using magma::mme::SctpdUplinkImpl;

std::shared_ptr<SctpdUplinkImpl> _service = nullptr;
std::unique_ptr<Server> _server           = nullptr;

int start_sctpd_uplink_server(void) {
  _service = std::make_shared<SctpdUplinkImpl>();

  ServerBuilder builder;
  builder.AddListeningPort(UPSTREAM_SOCK, grpc::InsecureServerCredentials());
  builder.RegisterService(_service.get());

  _server = builder.BuildAndStart();

  return 0;
}

void stop_sctpd_uplink_server(void) {
  if (_server != nullptr) {
    _server->Shutdown();
    _server->Wait();
    _server = nullptr;
  }
  _service = nullptr;
}
