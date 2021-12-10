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

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_38.401.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include <cstdint>
#include <cstring>
#include <string>
#include "lte/gateway/c/core/oai/lib/mobility_client/MobilityServiceClient.h"
#include <unistd.h>
#include <thread>
#include "lte/gateway/c/core/oai/lib/n11/SmfServiceClient.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_smfDefs.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/protos/session_manager.pb.h"
#include "lte/gateway/c/core/oai/lib/n11/M5GMobilityServiceClient.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"
#include "include/amf_client_servicer.h"

#define VERSION_0 0

using grpc::Channel;
using grpc::ChannelCredentials;
using grpc::CreateChannel;
using grpc::InsecureChannelCredentials;
using grpc::Status;
using magma::lte::AllocateIPAddressResponse;
using magma::lte::IPAddress;
using magma::lte::MobilityServiceClient;
using magma::lte::SetSMSessionContext;
using magma::lte::TeidSet;
using magma5g::AsyncM5GMobilityServiceClient;
using magma5g::AsyncSmfServiceClient;

namespace magma5g {

/***************************************************************************
**                                                                        **
** Name:   create_session_grpc_req_on_gnb_setup_rsp()                     **
**                                                                        **
** Description:                                                           **
** Creating session grpc req upon receiving resource setup response       **
** from gNB which have gNB IP and TEID, and sends to UPF through SMF      **
**                                                                        **
**                                                                        **
***************************************************************************/
int create_session_grpc_req_on_gnb_setup_rsp(
    amf_smf_establish_t* message, char* imsi, uint32_t version) {
  int rc = RETURNerror;
  magma::lte::SetSMSessionContext req;

  auto* req_common = req.mutable_common_context();
  auto* req_rat_specific =
      req.mutable_rat_specific_context()->mutable_m5gsm_session_context();
  // IMSI retrieved from amf context
  req_common->mutable_sid()->mutable_id()->assign(imsi);  // string id
  req_common->mutable_sid()->set_type(
      magma::lte::SubscriberID_IDType::SubscriberID_IDType_IMSI);
  req_common->set_rat_type(magma::lte::RATType::TGPP_NR);
  // PDU session state to CREATING
  req_common->set_sm_session_state(magma::lte::SMSessionFSMState::CREATING_0);
  // Version Added
  req_common->set_sm_session_version(version);
  req_rat_specific->set_pdu_session_id((message->pdu_session_id));
  req_rat_specific->set_request_type(magma::lte::RequestType::INITIAL_REQUEST);

  TeidSet* gnode_endpoint = req_rat_specific->mutable_gnode_endpoint();

  gnode_endpoint->set_teid(message->gnb_gtp_teid);

  char ipv4_str[INET_ADDRSTRLEN] = {0};
  inet_ntop(AF_INET, message->gnb_gtp_teid_ip_addr, ipv4_str, INET_ADDRSTRLEN);
  req_rat_specific->mutable_gnode_endpoint()->set_end_ipv4_addr(ipv4_str);

  OAILOG_INFO(
      LOG_AMF_APP,
      "Sending msg(grpc) to :[sessiond] for ue: [%s] pdu session :[%u]\n", imsi,
      message->pdu_session_id);

  AMFClientServicer::getInstance().set_smf_session(req);

  return rc;
}

/***************************************************************************
**                                                                        **
** Name:    amf_smf_create_ipv4_session_grpc_req()                        **
**                                                                        **
** Description: Fill session establishment gRPC request to SMF            **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_smf_create_ipv4_session_grpc_req(
    char* imsi, uint8_t* apn, uint32_t pdu_session_id,
    uint32_t pdu_session_type, uint32_t gnb_gtp_teid, uint8_t pti,
    uint8_t* gnb_gtp_teid_ip_addr, char* ipv4_addr, const ambr_t& state_ambr) {
  imsi64_t imsi64                   = INVALID_IMSI64;
  ue_m5gmm_context_s* ue_mm_context = NULL;
  amf_context_t* amf_ctxt_p         = NULL;

  OAILOG_INFO(
      LOG_AMF_APP,
      "Sending msg(grpc) to :[sessiond] for ue: [%s] pdu session: [%u]\n", imsi,
      pdu_session_id);

  IMSI_STRING_TO_IMSI64((char*) imsi, &imsi64);
  ue_mm_context = lookup_ue_ctxt_by_imsi(imsi64);

  if (ue_mm_context) {
    amf_ctxt_p = &ue_mm_context->amf_context;
  }

  if (!(amf_ctxt_p)) {
    OAILOG_ERROR(LOG_NAS_AMF, "IMSI is invalid\n");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  return AMFClientServicer::getInstance().amf_smf_create_pdu_session_ipv4(
      imsi, apn, pdu_session_id, pdu_session_type, gnb_gtp_teid, pti,
      gnb_gtp_teid_ip_addr, ipv4_addr, VERSION_0, state_ambr);
}

/***************************************************************************
 * **                                                                        **
 * ** Name:    amf_smf_create_pdu_session()                                  **
 * **                                                                        **
 * ** Description: Trigger PDU Session Creation in SMF                       **
 * **                                                                        **
 * **                                                                        **
 * ***************************************************************************/
int amf_smf_create_pdu_session(
    amf_smf_establish_t* message, char* imsi, uint32_t version) {
  imsi64_t imsi64                   = INVALID_IMSI64;
  amf_context_t* amf_ctxt_p         = NULL;
  ue_m5gmm_context_s* ue_mm_context = NULL;
  std::shared_ptr<smf_context_t> smf_ctx;

  IMSI_STRING_TO_IMSI64((char*) imsi, &imsi64);
  ue_mm_context = lookup_ue_ctxt_by_imsi(imsi64);
  smf_ctx       = amf_get_smf_context_by_pdu_session_id(
      ue_mm_context, message->pdu_session_id);

  if (ue_mm_context) {
    amf_ctxt_p = &ue_mm_context->amf_context;
  }

  if (!(amf_ctxt_p)) {
    OAILOG_ERROR(LOG_NAS_AMF, "IMSI is invalid\n");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  OAILOG_INFO(
      LOG_AMF_APP,
      "Sending msg(grpc) to :[mobilityd] for ue: [%s] ip-addr pdu session: "
      "[%u]\n",
      imsi, message->pdu_session_id);
  AMFClientServicer::getInstance().allocate_ipv4_address(
      imsi, smf_ctx->dnn.c_str(), message->pdu_session_id, message->pti,
      AF_INET, message->gnb_gtp_teid, message->gnb_gtp_teid_ip_addr, 4);

  return (RETURNok);
}

/***************************************************************************
**                                                                        **
** Name:    release_session_gprc_req()                                    **
**                                                                        **
** Description: Sends session release request to SMF                      **
**                                                                        **
**                                                                        **
***************************************************************************/
int release_session_gprc_req(amf_smf_release_t* message, char* imsi) {
  magma::lte::SetSMSessionContext req;
  auto* req_common = req.mutable_common_context();
  req_common->mutable_sid()->mutable_id()->assign(imsi);
  req_common->mutable_sid()->set_type(
      magma::lte::SubscriberID_IDType::SubscriberID_IDType_IMSI);
  req_common->set_sm_session_state(magma::lte::SMSessionFSMState::RELEASED_4);
  req_common->set_sm_session_version(1);  // uint32
  auto* req_rat_specific =
      req.mutable_rat_specific_context()->mutable_m5gsm_session_context();
  req_rat_specific->set_pdu_session_id(message->pdu_session_id);
  req_rat_specific->set_procedure_trans_identity(
      (const char*) (&(message->pti)));
  req_common->set_rat_type(magma::lte::RATType::TGPP_NR);

  OAILOG_INFO(
      LOG_AMF_APP,
      "Sending msg(grpc) to :[sessiond] for ue: [%s] release,  pdu session: "
      "[%d]\n",
      imsi, message->pdu_session_id);

  AMFClientServicer::getInstance().set_smf_session(req);

  return RETURNok;
}
}  // namespace magma5g
