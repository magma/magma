/*
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

#include <string>

#ifdef __cplusplus
extern "C" {
#endif
#include "amf_service_handler.h"
#include "log.h"
#include "conversions.h"
#ifdef __cplusplus
}
#endif

#include "AmfServiceImpl.h"
#include "lte/protos/session_manager.pb.h"
#include "lte/protos/subscriberdb.pb.h"

namespace grpc {

class ServerContext;
}  // namespace grpc

using grpc::ServerContext;
using grpc::Status;
using magma::lte::SetSmNotificationContext;
using magma::lte::SetSMSessionContextAccess;
using magma::lte::SmContextVoid;
using magma::lte::SmfPduSessionSmContext;

namespace magma {
using namespace lte;

AmfServiceImpl::AmfServiceImpl() {}

Status AmfServiceImpl::SetAmfNotification(
    ServerContext* context, const SetSmNotificationContext* notif,
    SmContextVoid* response) {
  OAILOG_INFO(LOG_UTIL, "Received  GRPC SetSmNotificationContext request\n");
  // ToDo processing ITTI,ZMQ

  return Status::OK;
}
// Set message from SessionD received
Status AmfServiceImpl::SetSmfSessionContext(
    ServerContext* context, const SetSMSessionContextAccess* request,
    SmContextVoid* response) {
  struct in_addr ip_addr       = {0};
  char ip_str[INET_ADDRSTRLEN] = {0};
  uint32_t ip_int              = 0;
  OAILOG_INFO(
      LOG_UTIL, "Received GRPC SetSmfSessionContext request from SMF\n");

  itti_n11_create_pdu_session_response_t itti_msg;
  auto& req_common = request->common_context();
  auto& req_m5g    = request->rat_specific_context().m5g_session_context_rsp();

  // CommonSessionContext
  strcpy(itti_msg.imsi, req_common.sid().id().c_str());
  itti_msg.sm_session_fsm_state =
      (sm_session_fsm_state_t) req_common.sm_session_state();
  itti_msg.sm_session_version = req_common.sm_session_version();

  // RatSpecificContextAccess
  memcpy((&itti_msg.pdu_session_id), req_m5g.pdu_session_id().c_str(), 1);
  itti_msg.pdu_session_type  = (pdu_session_type_t) req_m5g.pdu_session_type();
  itti_msg.selected_ssc_mode = (ssc_mode_t) req_m5g.selected_ssc_mode();
  itti_msg.m5gsm_cause       = (m5g_sm_cause_t) req_m5g.m5gsm_cause();
  itti_msg.always_on_pdu_session_indication =
      req_m5g.always_on_pdu_session_indication();
  itti_msg.allowed_ssc_mode = (ssc_mode_t) req_m5g.allowed_ssc_mode();
  itti_msg.m5gsm_congetion_re_attempt_indicator =
      req_m5g.m5gsm_congetion_re_attempt_indicator();
  itti_msg.pdu_address.redirect_address_type =
      (redirect_address_type_t) req_m5g.pdu_address().redirect_address_type();
  // PDU IP address coming from SMF in human-readable format has to be packed
  // into 4 raw bytes in hex for NAS5G layer
  strcpy(ip_str, req_m5g.pdu_address().redirect_server_address().c_str());
  inet_pton(AF_INET, ip_str, &(ip_addr.s_addr));
  ip_int = ntohl(ip_addr.s_addr);
  INT32_TO_BUFFER(ip_int, itti_msg.pdu_address.redirect_server_address);
  send_n11_create_pdu_session_resp_itti(&itti_msg);
  return Status::OK;
}

}  // namespace magma
