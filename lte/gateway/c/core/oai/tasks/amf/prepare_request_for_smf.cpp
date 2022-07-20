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
#include "lte/gateway/c/core/oai/include/nas/networkDef.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/common/common_defs.h"
#include <cstdint>
#include <cstring>
#include <string>
#include "lte/gateway/c/core/oai/lib/mobility_client/MobilityServiceClient.hpp"
#include <unistd.h>
#include <thread>
#include "lte/gateway/c/core/oai/lib/n11/SmfServiceClient.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_smfDefs.hpp"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/protos/session_manager.pb.h"
#include "lte/gateway/c/core/oai/lib/n11/M5GMobilityServiceClient.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_client_servicer.hpp"

#define VERSION_0 0

using grpc::Channel;
using grpc::ChannelCredentials;
using grpc::CreateChannel;
using grpc::InsecureChannelCredentials;
using grpc::Status;
using magma::lte::AllocateIPAddressResponse;
using magma::lte::IPAddress;
using magma::lte::MobilityServiceClient;
using magma::lte::QosPolicy;
using magma::lte::SetSMSessionContext;
using magma::lte::TeidSet;
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
status_code_e create_session_grpc_req_on_gnb_setup_rsp(
    amf_smf_establish_t* message, char* imsi, uint32_t version,
    std::shared_ptr<smf_context_t> smf_ctx) {
  status_code_e rc = RETURNerror;
  magma::lte::SetSMSessionContext req;

  auto imsi_str = std::string(imsi);
  auto* req_common = req.mutable_common_context();
  auto* req_rat_specific =
      req.mutable_rat_specific_context()->mutable_m5gsm_session_context();

  // IMSI retrieved from amf context
  req_common->mutable_sid()->set_type(
      magma::lte::SubscriberID_IDType::SubscriberID_IDType_IMSI);
  req_common->mutable_sid()->set_id("IMSI" + imsi_str);

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

#if 0 
  for (int i = 0; i < smf_ctx->qos_flow_list.maxNumOfQosFlows; i++) {
    QosPolicy* qosPolicy = req_rat_specific->add_qos_policy();
    qosPolicy->set_version(
        smf_ctx->qos_flow_list.item[i].qos_flow_req_item.qos_flow_version);
    qosPolicy->set_policy_state(QosPolicy::INSTALL);
    magma::lte::PolicyRule* rule = qosPolicy->mutable_qos();
    rule->set_id(
        (const char*)smf_ctx->qos_flow_list.item[i].qos_flow_req_item.rule_id);
  }
#endif

  OAILOG_DEBUG(LOG_AMF_APP, "Sending PDU session Setup Request to SMF");

  OAILOG_INFO(LOG_AMF_APP, "Sending msg(grpc) to :[sessiond] for ue: [%s]\n",
              imsi);

  AMFClientServicer::getInstance().set_smf_session(req);

  return rc;
}
/***************************************************************************
**                                                                        **
** Name:   amf_send_grpc_req_on_gnb_pdu_sess_mod_rsp()                    **
**                                                                        **
** Description:                                                           **
** Creating session grpc req upon receiving resource modify response      **
** from gNB sends qos flow info to UPF through SMF                        **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_send_grpc_req_on_gnb_pdu_sess_mod_rsp(
    amf_smf_establish_t* message, char* imsi, uint32_t version,
    std::shared_ptr<smf_context_t> smf_ctx) {
  int rc = RETURNerror;
  magma::lte::SetSMSessionContext req;

  auto* req_common = req.mutable_common_context();
  auto* req_rat_specific =
      req.mutable_rat_specific_context()->mutable_m5gsm_session_context();
  // IMSI retrieved from amf context
  req_common->mutable_sid()->mutable_id()->assign(
      imsi);  // string id
              // IMSI retrieved from amf context
  auto imsi_str = std::string(imsi);
  req_common->mutable_sid()->set_type(
      magma::lte::SubscriberID_IDType::SubscriberID_IDType_IMSI);
  req_common->mutable_sid()->set_id("IMSI" + imsi_str);

  req_common->set_rat_type(magma::lte::RATType::TGPP_NR);
  // PDU session state to ACTIVE_2
  req_common->set_sm_session_state(magma::lte::SMSessionFSMState::ACTIVE_2);
  // Version Added
  req_common->set_sm_session_version(version);
  req_rat_specific->set_pdu_session_id((message->pdu_session_id));
  req_rat_specific->set_request_type(
      magma::lte::RequestType::EXISTING_PDU_SESSION);
  TeidSet* gnode_endpoint = req_rat_specific->mutable_gnode_endpoint();
  gnode_endpoint->set_teid(message->gnb_gtp_teid);
  char ipv4_str[INET_ADDRSTRLEN] = {0};
  inet_ntop(AF_INET, message->gnb_gtp_teid_ip_addr, ipv4_str, INET_ADDRSTRLEN);
  req_rat_specific->mutable_gnode_endpoint()->set_end_ipv4_addr(ipv4_str);

  qos_flow_list_t* pti_flow_list = smf_ctx->get_proc_flow_list();

  for (int i = 0; i < pti_flow_list->maxNumOfQosFlows; i++) {
    QosPolicy* qosPolicy = req_rat_specific->add_qos_policy();
    qosPolicy->set_version(
        pti_flow_list->item[i].qos_flow_req_item.qos_flow_version);
    if (SMF_CAUSE_FAILURE == message->cause_value) {
      qosPolicy->set_policy_state(QosPolicy::REJECT);
    } else {
      qosPolicy->set_policy_state(QosPolicy::INSTALL);
    }
    magma::lte::PolicyRule* rule = qosPolicy->mutable_qos();
    rule->set_id((const char*)pti_flow_list->item[i].qos_flow_req_item.rule_id);
  }

  OAILOG_DEBUG(LOG_AMF_APP, "Sending PDU Session Modification Response to SMF");

  OAILOG_INFO(
      LOG_AMF_APP,
      "Sending msg(grpc) to :[sessiond] for ue: [%s] pdu session :[%u]\n", imsi,
      message->pdu_session_id);

  AMFClientServicer::getInstance().set_smf_session(req);

  OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
}

/***************************************************************************
**                                                                        **
** Name:    amf_smf_create_session_req()                                  **
**                                                                        **
** Description: Fill session establishment gRPC request to SMF            **
**                                                                        **
**                                                                        **
***************************************************************************/
status_code_e amf_smf_create_session_req(
    char* imsi, uint8_t* apn, uint32_t pdu_session_id,
    uint32_t pdu_session_type, uint32_t gnb_gtp_teid, uint8_t pti,
    uint8_t* gnb_gtp_teid_ip_addr, char* ue_ipv4_addr, char* ue_ipv6_addr,
    const ambr_t& state_ambr, const eps_subscribed_qos_profile_t& qos_profile) {
  imsi64_t imsi64 = INVALID_IMSI64;
  ue_m5gmm_context_s* ue_mm_context = NULL;
  amf_context_t* amf_ctxt_p = NULL;
  OAILOG_FUNC_IN(LOG_AMF_APP);
  OAILOG_INFO(
      LOG_AMF_APP,
      "Sending msg(grpc) to :[sessiond] for ue: [%s] pdu session: [%u]\n", imsi,
      pdu_session_id);

  IMSI_STRING_TO_IMSI64((char*)imsi, &imsi64);
  ue_mm_context = lookup_ue_ctxt_by_imsi(imsi64);

  if (ue_mm_context) {
    amf_ctxt_p = &ue_mm_context->amf_context;
  }

  if (!(amf_ctxt_p)) {
    OAILOG_ERROR(LOG_NAS_AMF, "IMSI is invalid\n");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  OAILOG_FUNC_RETURN(
      LOG_AMF_APP, AMFClientServicer::getInstance().amf_smf_create_pdu_session(
                       imsi, apn, pdu_session_id, pdu_session_type,
                       gnb_gtp_teid, pti, gnb_gtp_teid_ip_addr, ue_ipv4_addr,
                       ue_ipv6_addr, state_ambr, VERSION_0, qos_profile));
}

/***************************************************************************
 * **                                                                        **
 * ** Name:    amf_smf_initiate_pdu_session_creation()                       **
 * **                                                                        **
 * ** Description: Initiate PDU Session Creation process                     **
 * **                                                                        **
 * **                                                                        **
 * ***************************************************************************/
status_code_e amf_smf_initiate_pdu_session_creation(
    amf_smf_establish_t* message, char* imsi, uint32_t version) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  imsi64_t imsi64 = INVALID_IMSI64;
  amf_context_t* amf_ctxt_p = NULL;
  ue_m5gmm_context_s* ue_mm_context = NULL;
  std::shared_ptr<smf_context_t> smf_ctx;

  IMSI_STRING_TO_IMSI64((char*)imsi, &imsi64);
  ue_mm_context = lookup_ue_ctxt_by_imsi(imsi64);
  smf_ctx = amf_get_smf_context_by_pdu_session_id(ue_mm_context,
                                                  message->pdu_session_id);

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

  if (message->pdu_session_type == NET_PDN_TYPE_IPV4) {
    AMFClientServicer::getInstance().allocate_ipv4_address(
        imsi, smf_ctx->dnn.c_str(), message->pdu_session_id, message->pti,
        NET_PDN_TYPE_IPV4, message->gnb_gtp_teid, message->gnb_gtp_teid_ip_addr,
        4);
  } else if (message->pdu_session_type == NET_PDN_TYPE_IPV6) {
    AMFClientServicer::getInstance().allocate_ipv6_address(
        imsi, smf_ctx->dnn.c_str(), message->pdu_session_id, message->pti,
        NET_PDN_TYPE_IPV6, message->gnb_gtp_teid, message->gnb_gtp_teid_ip_addr,
        4);
  } else if (message->pdu_session_type == NET_PDN_TYPE_IPV4V6) {
    AMFClientServicer::getInstance().allocate_ipv4v6_address(
        imsi, smf_ctx->dnn.c_str(), message->pdu_session_id, message->pti,
        NET_PDN_TYPE_IPV4V6, message->gnb_gtp_teid,
        message->gnb_gtp_teid_ip_addr, 4);
  }

  OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNok);
}

/***************************************************************************
**                                                                        **
** Name:    release_session_gprc_req()                                    **
**                                                                        **
** Description: Sends session release request to SMF                      **
**                                                                        **
**                                                                        **
***************************************************************************/
status_code_e release_session_gprc_req(amf_smf_release_t* message, char* imsi) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  magma::lte::SetSMSessionContext req;
  auto imsi_str = std::string(imsi);
  auto* req_common = req.mutable_common_context();

  // Encode subscriber as IMSI
  req_common->mutable_sid()->set_type(
      magma::lte::SubscriberID_IDType::SubscriberID_IDType_IMSI);
  req_common->mutable_sid()->set_id("IMSI" + imsi_str);

  req_common->set_sm_session_state(magma::lte::SMSessionFSMState::RELEASED_4);
  req_common->set_sm_session_version(1);  // uint32
  auto* req_rat_specific =
      req.mutable_rat_specific_context()->mutable_m5gsm_session_context();
  req_rat_specific->set_pdu_session_id(message->pdu_session_id);
  req_rat_specific->set_procedure_trans_identity(message->pti);
  req_common->set_rat_type(magma::lte::RATType::TGPP_NR);

  OAILOG_INFO(
      LOG_AMF_APP,
      "Sending msg(grpc) to :[sessiond] for ue: [%s] release,  pdu session: "
      "[%d]\n",
      imsi, message->pdu_session_id);

  AMFClientServicer::getInstance().set_smf_session(req);

  OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNok);
}

/***************************************************************************
 * **                                                                        **
 * ** Name:    amf_app_pdu_session_modification_complete()                   **
 * **                                                                        **
 * ** Description: Handle PDU Session Modification Complete Msg              **
 * **                                                                        **
 * **                                                                        **
 * ***************************************************************************/
int amf_app_pdu_session_modification_complete(amf_smf_establish_t* message,
                                              char* imsi, uint32_t version) {
  imsi64_t imsi64 = INVALID_IMSI64;
  amf_context_t* amf_ctxt_p = NULL;
  ue_m5gmm_context_s* ue_mm_context = NULL;
  std::shared_ptr<smf_context_t> smf_ctx;
  amf_smf_establish_t amf_smf_grpc_ies;

  IMSI_STRING_TO_IMSI64((char*)imsi, &imsi64);
  ue_mm_context = lookup_ue_ctxt_by_imsi(imsi64);
  smf_ctx = amf_get_smf_context_by_pdu_session_id(ue_mm_context,
                                                  message->pdu_session_id);

  if (ue_mm_context) {
    amf_ctxt_p = &ue_mm_context->amf_context;
  }

  if (!(amf_ctxt_p)) {
    OAILOG_ERROR(LOG_NAS_AMF, "IMSI is invalid\n");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  if (!smf_ctx) {
    OAILOG_ERROR(LOG_NAS_AMF,
                 "session context not found using session id [%u]\n",
                 message->pdu_session_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  OAILOG_INFO(LOG_AMF_APP, "Received PDU Session Modification Complete [%s] \n",
              imsi);

  // Incrementing the  pdu session version
  smf_ctx->pdu_session_version++;

  amf_smf_grpc_ies.pdu_session_id = message->pdu_session_id;
  amf_smf_grpc_ies.cause_value = SMF_CAUSE_SUCCESS;

  // gnb tunnel info
  memset(&amf_smf_grpc_ies.gnb_gtp_teid_ip_addr, '\0',
         sizeof(amf_smf_grpc_ies.gnb_gtp_teid_ip_addr));
  memset(&amf_smf_grpc_ies.gnb_gtp_teid, '\0',
         sizeof(amf_smf_grpc_ies.gnb_gtp_teid));
  memcpy(&amf_smf_grpc_ies.gnb_gtp_teid_ip_addr,
         &smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr, 4);
  memcpy(&amf_smf_grpc_ies.gnb_gtp_teid, &smf_ctx->gtp_tunnel_id.gnb_gtp_teid,
         4);

  IMSI64_TO_STRING(ue_mm_context->amf_context.imsi64, imsi, 15);
  // Prepare and send modify setup response message to SMF through gRPC
  amf_send_grpc_req_on_gnb_pdu_sess_mod_rsp(
      &amf_smf_grpc_ies, imsi, smf_ctx->pdu_session_version, smf_ctx);

  return (RETURNok);
}

/***************************************************************************
 * **                                                                        **
 * ** Name:    amf_app_pdu_session_modification_command_reject()             **
 * **                                                                        **
 * ** Description: Handle PDU Session Modification Command reject Msg        **
 * **                                                                        **
 * **                                                                        **
 * ***************************************************************************/
int amf_app_pdu_session_modification_command_reject(
    amf_smf_establish_t* message, char* imsi, uint32_t version) {
  imsi64_t imsi64 = INVALID_IMSI64;
  amf_context_t* amf_ctxt_p = NULL;
  ue_m5gmm_context_s* ue_mm_context = NULL;
  std::shared_ptr<smf_context_t> smf_ctx;
  amf_smf_establish_t amf_smf_grpc_ies;

  IMSI_STRING_TO_IMSI64((char*)imsi, &imsi64);
  ue_mm_context = lookup_ue_ctxt_by_imsi(imsi64);
  smf_ctx = amf_get_smf_context_by_pdu_session_id(ue_mm_context,
                                                  message->pdu_session_id);

  if (ue_mm_context) {
    amf_ctxt_p = &ue_mm_context->amf_context;
  }

  if (!(amf_ctxt_p)) {
    OAILOG_ERROR(LOG_NAS_AMF, "IMSI is invalid\n");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  if (!smf_ctx) {
    OAILOG_ERROR(LOG_NAS_AMF,
                 "session context not found using session id [%u]\n",
                 message->pdu_session_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  OAILOG_INFO(LOG_AMF_APP,
              "Received PDU Session Modification Command Reject [%s] \n", imsi);
  // Incrementing the  pdu session version
  smf_ctx->pdu_session_version++;

  amf_smf_grpc_ies.pdu_session_id = message->pdu_session_id;
  amf_smf_grpc_ies.cause_value = SMF_CAUSE_FAILURE;

  // gnb tunnel info
  memset(&amf_smf_grpc_ies.gnb_gtp_teid_ip_addr, '\0',
         sizeof(amf_smf_grpc_ies.gnb_gtp_teid_ip_addr));
  memset(&amf_smf_grpc_ies.gnb_gtp_teid, '\0',
         sizeof(amf_smf_grpc_ies.gnb_gtp_teid));
  memcpy(&amf_smf_grpc_ies.gnb_gtp_teid_ip_addr,
         &smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr, 4);
  memcpy(&amf_smf_grpc_ies.gnb_gtp_teid, &smf_ctx->gtp_tunnel_id.gnb_gtp_teid,
         4);

  IMSI64_TO_STRING(ue_mm_context->amf_context.imsi64, imsi, 15);
  // Prepare and send modify setup response message to SMF through gRPC
  amf_send_grpc_req_on_gnb_pdu_sess_mod_rsp(
      &amf_smf_grpc_ies, imsi, smf_ctx->pdu_session_version, smf_ctx);

  return (RETURNok);
}
}  // namespace magma5g
