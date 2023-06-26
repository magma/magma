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

#include "lte/gateway/c/core/oai/lib/pcef/pcef_handlers.hpp"

#include <cstring>
#include <string>

#include <grpcpp/impl/codegen/status.h>
#include "lte/protos/session_manager.pb.h"

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/itti/itti_types.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/lib/pcef/PCEFClient.hpp"
#include "lte/gateway/c/core/oai/lib/mobility_client/MobilityClientAPI.hpp"

extern task_zmq_ctx_t grpc_service_task_zmq_ctx;

#define ULI_DATA_SIZE 13

// TODO Clean up pcef_create_session_data structure to include
// imsi/ip/bearer_id etc.
static void pcef_fill_create_session_req(
    std::string& imsi, std::string& ip4, std::string& ip6, ebi_t eps_bearer_id,
    const struct pcef_create_session_data* session_data,
    magma::LocalCreateSessionRequest* sreq) {
  // Common Context
  auto common_context = sreq->mutable_common_context();
  common_context->mutable_sid()->set_id("IMSI" + imsi);
  if (!ip4.empty()) {
    common_context->set_ue_ipv4(ip4);
  }
  if (!ip6.empty()) {
    common_context->set_ue_ipv6(ip6);
  }
  common_context->set_apn(session_data->apn);
  common_context->set_msisdn(session_data->msisdn, session_data->msisdn_len);
  common_context->set_rat_type(magma::RATType::TGPP_LTE);

  // LTE Context
  auto lte_context =
      sreq->mutable_rat_specific_context()->mutable_lte_context();
  lte_context->set_bearer_id(eps_bearer_id);
  lte_context->set_spgw_ipv4(session_data->sgw_ip);
  lte_context->set_plmn_id(session_data->mcc_mnc, session_data->mcc_mnc_len);
  lte_context->set_imsi_plmn_id(session_data->imsi_mcc_mnc,
                                session_data->imsi_mcc_mnc_len);
  auto cc = session_data->charging_characteristics;
  if (cc.length > 0) {
    OAILOG_DEBUG(LOG_SPGW_APP, "Sending Charging Characteristics to PCEF");
    lte_context->set_charging_characteristics(cc.value, cc.length);
  }
  if (session_data->imeisv_exists) {
    lte_context->set_imei(session_data->imeisv, IMEISV_DIGITS_MAX);
  }
  if (session_data->uli_exists) {
    OAILOG_DEBUG(LOG_SPGW_APP, "Sending ULI to PCEF");
    lte_context->set_user_location(session_data->uli, ULI_DATA_SIZE);
  }

  // QoS Info
  magma::QosInformationRequest qos_info;
  qos_info.set_apn_ambr_dl(session_data->ambr_dl);
  qos_info.set_apn_ambr_ul(session_data->ambr_ul);
  qos_info.set_priority_level(session_data->pl);
  qos_info.set_preemption_capability(session_data->pci);
  qos_info.set_preemption_vulnerability(session_data->pvi);
  qos_info.set_qos_class_id(session_data->qci);
  lte_context->mutable_qos_info()->CopyFrom(qos_info);
}

/**
 * Send an ITTI message from GRPC service task to SPGW task when the
 * PCEFClient receives the response for the aynchronous Create Sesssion RPC
 */
status_code_e send_itti_pcef_create_session_response(
    const std::string& imsi, s5_create_session_request_t session_request,
    const grpc::Status& status) {
#if MME_UNIT_TEST
  return RETURNok;
#endif
  MessageDef* message_p = nullptr;

  message_p =
      itti_alloc_new_message(TASK_GRPC_SERVICE, PCEF_CREATE_SESSION_RESPONSE);

  if (!message_p) {
    OAILOG_ERROR(LOG_UTIL,
                 "Message PCEF Create Session Response allocation failed\n");
    return RETURNerror;
  }

  itti_pcef_create_session_response_t* pcef_create_session_resp_p = nullptr;
  pcef_create_session_resp_p = &message_p->ittiMsg.pcef_create_session_response;

  pcef_create_session_resp_p->rpc_status = PCEF_STATUS_OK;
  if (!status.ok()) {
    pcef_create_session_resp_p->rpc_status = PCEF_STATUS_FAILED;
  }

  pcef_create_session_resp_p->teid = session_request.context_teid;
  pcef_create_session_resp_p->eps_bearer_id = session_request.eps_bearer_id;
  pcef_create_session_resp_p->sgi_status = session_request.status;

  IMSI_STRING_TO_IMSI64(imsi.c_str(), &message_p->ittiMsgHeader.imsi);

  OAILOG_DEBUG(LOG_UTIL, "Sending PCEF create session response to SPGW task");
  return send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_SPGW_APP, message_p);
}

void pcef_create_session(std::string imsi_str, const char* ip4, const char* ip6,
                         const pcef_create_session_data* session_data,
                         s5_create_session_request_t session_request) {
  std::string ip4_str, ip6_str;

  if (ip4) {
    ip4_str = ip4;
  }
  if (ip6) {
    ip6_str = ip6;
  }
  // Change ip to spgw_ip. Get it from sgw_app_t sgw_app;
  magma::LocalCreateSessionRequest sreq;
  pcef_fill_create_session_req(imsi_str, ip4_str, ip6_str,
                               session_request.eps_bearer_id, session_data,
                               &sreq);

  auto apn = std::string(session_data->apn);
  // call the `CreateSession` gRPC method and execute the inline function
  magma::PCEFClient::create_session(
      sreq, [imsi_str, session_request](
                const grpc::Status& status,
                const magma::LocalCreateSessionResponse& response) {
        send_itti_pcef_create_session_response(imsi_str, session_request,
                                               status);
      });
}

bool pcef_end_session(const std::string imsi, const std::string apn) {
  magma::LocalEndSessionRequest request;
  request.mutable_sid()->set_id("IMSI" + imsi);
  request.set_apn(apn);
  magma::PCEFClient::end_session(
      request,
      [&](grpc::Status status, magma::LocalEndSessionResponse response) {
        return;  // For now, do nothing. TODO: handle errors asynchronously
      });
  return true;
}

void pcef_send_policy2bearer_binding(const char* imsi,
                                     const uint8_t default_bearer_id,
                                     const char* policy_rule_name,
                                     const uint8_t eps_bearer_id,
                                     const uint32_t eps_bearer_agw_teid,
                                     const uint32_t eps_bearer_enb_teid) {
  magma::PolicyBearerBindingRequest request;
  request.mutable_sid()->set_id("IMSI" + std::string(imsi));
  request.set_linked_bearer_id(default_bearer_id);
  request.set_policy_rule_id(policy_rule_name);
  request.set_bearer_id(eps_bearer_id);
  request.mutable_teids()->set_enb_teid(eps_bearer_enb_teid);
  request.mutable_teids()->set_agw_teid(eps_bearer_agw_teid);
  magma::PCEFClient::bind_policy2bearer(
      request,
      [&](grpc::Status status, magma::PolicyBearerBindingResponse response) {
        return;  // For now, do nothing. TODO: handle errors asynchronously
      });
}

void pcef_update_teids(const char* imsi, uint8_t default_bearer_id,
                       uint32_t enb_teid, uint32_t agw_teid) {
#if MME_UNIT_TEST
  return;
#endif
  magma::UpdateTunnelIdsRequest request;
  request.mutable_sid()->set_id("IMSI" + std::string(imsi));
  request.set_bearer_id(default_bearer_id);
  request.set_enb_teid(enb_teid);
  request.set_agw_teid(agw_teid);

  magma::PCEFClient::update_teids(
      request,
      [request](grpc::Status status, magma::UpdateTunnelIdsResponse response) {
        if (status.error_code() == grpc::ABORTED) {
#if MME_UNIT_TEST
          return;
#endif
          MessageDef* message_p = DEPRECATEDitti_alloc_new_message_fatal(
              TASK_GRPC_SERVICE, GX_NW_INITIATED_DEACTIVATE_BEARER_REQ);
          itti_gx_nw_init_deactv_bearer_request_t* itti_msg =
              &message_p->ittiMsg.gx_nw_init_deactv_bearer_request;
          std::string imsi = request.sid().id();
          OAILOG_INFO(
              LOG_UTIL,
              "Received grpc::ABORTED for update_teids RPC for %s with error "
              "msg: %s. Deactivating bearer %u.",
              imsi.c_str(), status.error_message().c_str(),
              request.bearer_id());
          // strip off "IMSI" prefix
          imsi = imsi.substr(4, std::string::npos);
          itti_msg->imsi_length = imsi.size();
          strcpy(itti_msg->imsi, imsi.c_str());
          itti_msg->lbi = request.bearer_id();
          itti_msg->no_of_bearers = 1;
          itti_msg->ebi[0] = request.bearer_id();
          send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_SPGW_APP,
                           message_p);
        }

        return;
      });
}

/*
 * Converts ascii values in [0,9] to [48,57]=['0','9']
 * else if they are in [48,57] keep them the same
 * else log an error and return '0'=48 value
 */
char convert_digit_to_char(char digit) {
  if ((digit >= 0) && (digit <= 9)) {
    return (digit + '0');
  } else if ((digit >= '0') && (digit <= '9')) {
    return digit;
  } else {
    OAILOG_ERROR(
        LOG_SPGW_APP,
        "The input value for digit is not in a valid range: "
        "Session request would likely be rejected on Gx or Gy interface\n");
    return '0';
  }
}

void get_plmn_from_session_req(
    const itti_s11_create_session_request_t* saved_req, char* mcc_mnc) {
  mcc_mnc[0] = convert_digit_to_char(saved_req->serving_network.mcc[0]);
  mcc_mnc[1] = convert_digit_to_char(saved_req->serving_network.mcc[1]);
  mcc_mnc[2] = convert_digit_to_char(saved_req->serving_network.mcc[2]);
  mcc_mnc[3] = convert_digit_to_char(saved_req->serving_network.mnc[0]);
  mcc_mnc[4] = convert_digit_to_char(saved_req->serving_network.mnc[1]);
  int mcc_mnc_len = 5;
  if ((saved_req->serving_network.mnc[2] & 0xf) != 0xf) {
    mcc_mnc[5] = convert_digit_to_char(saved_req->serving_network.mnc[2]);
    mcc_mnc[6] = '\0';
    mcc_mnc_len += 1;
  } else {
    mcc_mnc[5] = '\0';
  }
}

void get_imsi_plmn_from_session_req(const std::string imsi,
                                    struct pcef_create_session_data* data) {
  const char* imsi_digit = imsi.c_str();
  data->imsi_mcc_mnc[0] = convert_digit_to_char(imsi_digit[0]);
  data->imsi_mcc_mnc[1] = convert_digit_to_char(imsi_digit[1]);
  data->imsi_mcc_mnc[2] = convert_digit_to_char(imsi_digit[2]);
  data->imsi_mcc_mnc[3] = convert_digit_to_char(imsi_digit[3]);
  data->imsi_mcc_mnc[4] = convert_digit_to_char(imsi_digit[4]);
  data->imsi_mcc_mnc_len = 5;
  // Check the mcc_mnc_len which is deocoded from serving network mnc value
  if (data->mcc_mnc_len == 6) {
    data->imsi_mcc_mnc[5] = convert_digit_to_char(imsi_digit[5]);
    data->imsi_mcc_mnc[6] = '\0';
    data->imsi_mcc_mnc_len += 1;
  } else {
    data->imsi_mcc_mnc[5] = '\0';
  }
}

bool get_uli_from_session_req(
    const itti_s11_create_session_request_t* saved_req, char* uli) {
  if (!saved_req->uli.present) {
    return 0;
  }

  uli[0] = 130;  // TAI and ECGI - defined in 29.061

  // TAI as defined in 29.274 8.21.4
  uli[1] = ((saved_req->uli.s.tai.plmn.mcc_digit2 & 0xf) << 4) |
           ((saved_req->uli.s.tai.plmn.mcc_digit1 & 0xf));
  uli[2] = ((saved_req->uli.s.tai.plmn.mnc_digit3 & 0xf) << 4) |
           ((saved_req->uli.s.tai.plmn.mcc_digit3 & 0xf));
  uli[3] = ((saved_req->uli.s.tai.plmn.mnc_digit2 & 0xf) << 4) |
           ((saved_req->uli.s.tai.plmn.mnc_digit1 & 0xf));
  uli[4] = (saved_req->uli.s.tai.tac >> 8) & 0xff;
  uli[5] = saved_req->uli.s.tai.tac & 0xff;

  // ECGI as defined in 29.274 8.21.5
  uli[6] = ((saved_req->uli.s.ecgi.plmn.mcc_digit2 & 0xf) << 4) |
           ((saved_req->uli.s.ecgi.plmn.mcc_digit1 & 0xf));
  uli[7] = ((saved_req->uli.s.ecgi.plmn.mnc_digit3 & 0xf) << 4) |
           ((saved_req->uli.s.ecgi.plmn.mcc_digit3 & 0xf));
  uli[8] = ((saved_req->uli.s.ecgi.plmn.mnc_digit2 & 0xf) << 4) |
           ((saved_req->uli.s.ecgi.plmn.mnc_digit1 & 0xf));
  uli[9] = (saved_req->uli.s.ecgi.cell_identity.enb_id >> 16) & 0xf;
  uli[10] = (saved_req->uli.s.ecgi.cell_identity.enb_id >> 8) & 0xff;
  uli[11] = saved_req->uli.s.ecgi.cell_identity.enb_id & 0xff;
  uli[12] = saved_req->uli.s.ecgi.cell_identity.cell_id & 0xff;
  uli[13] = '\0';

  char hex_uli[3 * ULI_DATA_SIZE + 1];
  OAILOG_DEBUG(LOG_SPGW_APP, "Session request ULI %s",
               bytes_to_hex(uli, ULI_DATA_SIZE, hex_uli));
  return 1;
}

int get_msisdn_from_session_req(
    const itti_s11_create_session_request_t* saved_req, char* msisdn) {
  int len = saved_req->msisdn.length;
  if (len == 0) {
    return len;
  }

  int i, j;

  for (i = 0; i < len; ++i) {
    j = i << 1;
    msisdn[j] = (saved_req->msisdn.digit[i] & 0xf) + '0';
    msisdn[j + 1] = ((saved_req->msisdn.digit[i] >> 4) & 0xf) + '0';
  }
  if ((saved_req->msisdn.digit[len - 1] & 0xf0) == 0xf0) {
    len = (len << 1) - 1;
  } else {
    len = len << 1;
  }
  return len;
}

void convert_imeisv_to_string(char* imeisv) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  uint8_t idx = 0;
  for (; idx < IMEISV_DIGITS_MAX; idx++) {
    imeisv[idx] = convert_digit_to_char(imeisv[idx]);
  }
  imeisv[idx] = '\0';

  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

int get_imeisv_from_session_req(
    const itti_s11_create_session_request_t* saved_req, char* imeisv) {
  if (saved_req->mei.present & MEI_IMEISV) {
    // IMEISV as defined in 3GPP TS 23.003 MEI_IMEISV
    imeisv[0] = saved_req->mei.choice.imeisv.u.num.tac1;
    imeisv[1] = saved_req->mei.choice.imeisv.u.num.tac2;
    imeisv[2] = saved_req->mei.choice.imeisv.u.num.tac3;
    imeisv[3] = saved_req->mei.choice.imeisv.u.num.tac4;
    imeisv[4] = saved_req->mei.choice.imeisv.u.num.tac5;
    imeisv[5] = saved_req->mei.choice.imeisv.u.num.tac6;
    imeisv[6] = saved_req->mei.choice.imeisv.u.num.tac7;
    imeisv[7] = saved_req->mei.choice.imeisv.u.num.tac8;
    imeisv[8] = saved_req->mei.choice.imeisv.u.num.snr1;
    imeisv[9] = saved_req->mei.choice.imeisv.u.num.snr2;
    imeisv[10] = saved_req->mei.choice.imeisv.u.num.snr3;
    imeisv[11] = saved_req->mei.choice.imeisv.u.num.snr4;
    imeisv[12] = saved_req->mei.choice.imeisv.u.num.snr5;
    imeisv[13] = saved_req->mei.choice.imeisv.u.num.snr6;
    imeisv[14] = saved_req->mei.choice.imeisv.u.num.svn1;
    imeisv[15] = saved_req->mei.choice.imeisv.u.num.svn2;
    imeisv[IMEISV_DIGITS_MAX] = '\0';
    return 1;
  }
  return 0;
}

bool pcef_delete_dedicated_bearer(const char* imsi, const ebi_list_t ebi_list) {
  auto imsi_str = std::string(imsi);

  // TODO(pruthvihebbani) : Send grpc message to session manager to delete
  // dedicated bearer
  return true;
}
