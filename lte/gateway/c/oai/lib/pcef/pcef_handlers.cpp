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

#include <grpcpp/impl/codegen/status.h>
#include <cstring>
#include <string>
#include <conversions.h>
#include <common_defs.h>

#ifdef __cplusplus
extern "C" {
#endif

#include "common_defs.h"
#include "log.h"

#ifdef __cplusplus
}
#endif

#include "pcef_handlers.h"
#include "PCEFClient.h"
#include "MobilityClientAPI.h"
#include "itti_types.h"
#include "lte/protos/session_manager.pb.h"
#include "spgw_types.h"

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
  lte_context->set_imsi_plmn_id(
      session_data->imsi_mcc_mnc, session_data->imsi_mcc_mnc_len);
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
int send_itti_pcef_create_session_response(
    const std::string& imsi, s5_create_session_request_t session_request,
    const grpc::Status& status) {
  MessageDef* message_p = nullptr;

  message_p =
      itti_alloc_new_message(TASK_GRPC_SERVICE, PCEF_CREATE_SESSION_RESPONSE);

  if (!message_p) {
    OAILOG_ERROR(
        LOG_UTIL, "Message PCEF Create Session Response allocation failed\n");
    return RETURNerror;
  }

  itti_pcef_create_session_response_t* pcef_create_session_resp_p = nullptr;
  pcef_create_session_resp_p = &message_p->ittiMsg.pcef_create_session_response;

  pcef_create_session_resp_p->rpc_status = PCEF_STATUS_OK;
  if (!status.ok()) {
    pcef_create_session_resp_p->rpc_status = PCEF_STATUS_FAILED;
  }

  pcef_create_session_resp_p->teid          = session_request.context_teid;
  pcef_create_session_resp_p->eps_bearer_id = session_request.eps_bearer_id;
  pcef_create_session_resp_p->sgi_status    = session_request.status;

  IMSI_STRING_TO_IMSI64(imsi.c_str(), &message_p->ittiMsgHeader.imsi);

  OAILOG_DEBUG(LOG_UTIL, "Sending PCEF create session response to SPGW task");
  return send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_SPGW_APP, message_p);
}

void pcef_create_session(
    const char* imsi, const char* ip4, const char* ip6,
    const pcef_create_session_data* session_data,
    s5_create_session_request_t session_request) {
  auto imsi_str = std::string(imsi);
  std::string ip4_str, ip6_str;

  if (ip4) {
    ip4_str = ip4;
  }
  if (ip6) {
    ip6_str = ip6;
  }
  // Change ip to spgw_ip. Get it from sgw_app_t sgw_app;
  magma::LocalCreateSessionRequest sreq;
  pcef_fill_create_session_req(
      imsi_str, ip4_str, ip6_str, session_request.eps_bearer_id, session_data,
      &sreq);

  auto apn = std::string(session_data->apn);
  // call the `CreateSession` gRPC method and execute the inline function
  magma::PCEFClient::create_session(
      sreq, [imsi_str, session_request](
                const grpc::Status& status,
                const magma::LocalCreateSessionResponse& response) {
        send_itti_pcef_create_session_response(
            imsi_str, session_request, status);
      });
}

bool pcef_end_session(char* imsi, char* apn) {
  magma::LocalEndSessionRequest request;
  request.mutable_sid()->set_id("IMSI" + std::string(imsi));
  request.set_apn(apn);
  magma::PCEFClient::end_session(
      request,
      [&](grpc::Status status, magma::LocalEndSessionResponse response) {
        return;  // For now, do nothing. TODO: handle errors asynchronously
      });
  return true;
}

void pcef_send_policy2bearer_binding(
    const char* imsi, const uint8_t default_bearer_id,
    const char* policy_rule_name, const uint8_t eps_bearer_id,
    const uint32_t eps_bearer_agw_teid, const uint32_t eps_bearer_enb_teid) {
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

void pcef_update_teids(
    const char* imsi, uint8_t default_bearer_id, uint32_t enb_teid,
    uint32_t agw_teid) {
  magma::UpdateTunnelIdsRequest request;
  request.mutable_sid()->set_id("IMSI" + std::string(imsi));
  request.set_bearer_id(default_bearer_id);
  request.set_enb_teid(enb_teid);
  request.set_agw_teid(agw_teid);

  magma::PCEFClient::update_teids(
      request, [&](grpc::Status status,
                   magma::UpdateTunnelIdsResponse response) { return; });
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

static void get_plmn_from_session_req(
    const itti_s11_create_session_request_t* saved_req,
    struct pcef_create_session_data* data) {
  data->mcc_mnc[0]  = convert_digit_to_char(saved_req->serving_network.mcc[0]);
  data->mcc_mnc[1]  = convert_digit_to_char(saved_req->serving_network.mcc[1]);
  data->mcc_mnc[2]  = convert_digit_to_char(saved_req->serving_network.mcc[2]);
  data->mcc_mnc[3]  = convert_digit_to_char(saved_req->serving_network.mnc[0]);
  data->mcc_mnc[4]  = convert_digit_to_char(saved_req->serving_network.mnc[1]);
  data->mcc_mnc_len = 5;
  if ((saved_req->serving_network.mnc[2] & 0xf) != 0xf) {
    data->mcc_mnc[5] = convert_digit_to_char(saved_req->serving_network.mnc[2]);
    data->mcc_mnc[6] = '\0';
    data->mcc_mnc_len += 1;
  } else {
    data->mcc_mnc[5] = '\0';
  }
}

static void get_imsi_plmn_from_session_req(
    const itti_s11_create_session_request_t* saved_req,
    struct pcef_create_session_data* data) {
  data->imsi_mcc_mnc[0]  = convert_digit_to_char(saved_req->imsi.digit[0]);
  data->imsi_mcc_mnc[1]  = convert_digit_to_char(saved_req->imsi.digit[1]);
  data->imsi_mcc_mnc[2]  = convert_digit_to_char(saved_req->imsi.digit[2]);
  data->imsi_mcc_mnc[3]  = convert_digit_to_char(saved_req->imsi.digit[3]);
  data->imsi_mcc_mnc[4]  = convert_digit_to_char(saved_req->imsi.digit[4]);
  data->imsi_mcc_mnc_len = 5;
  // Check if 2 or 3 digit by verifying mnc[2] has a valid value
  if ((saved_req->serving_network.mnc[2] & 0xf) != 0xf) {
    data->imsi_mcc_mnc[5] = convert_digit_to_char(saved_req->imsi.digit[5]);
    data->imsi_mcc_mnc[6] = '\0';
    data->imsi_mcc_mnc_len += 1;
  } else {
    data->imsi_mcc_mnc[5] = '\0';
  }
}

static int get_uli_from_session_req(
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
  uli[9]  = (saved_req->uli.s.ecgi.cell_identity.enb_id >> 16) & 0xf;
  uli[10] = (saved_req->uli.s.ecgi.cell_identity.enb_id >> 8) & 0xff;
  uli[11] = saved_req->uli.s.ecgi.cell_identity.enb_id & 0xff;
  uli[12] = saved_req->uli.s.ecgi.cell_identity.cell_id & 0xff;
  uli[13] = '\0';

  char hex_uli[3 * ULI_DATA_SIZE + 1];
  OAILOG_DEBUG(
      LOG_SPGW_APP, "Session request ULI %s",
      bytes_to_hex(uli, ULI_DATA_SIZE, hex_uli));
  return 1;
}

int get_msisdn_from_session_req(
    const itti_s11_create_session_request_t* saved_req, char* msisdn) {
  int len = saved_req->msisdn.length;
  int i, j;

  for (i = 0; i < len; ++i) {
    j             = i << 1;
    msisdn[j]     = (saved_req->msisdn.digit[i] & 0xf) + '0';
    msisdn[j + 1] = ((saved_req->msisdn.digit[i] >> 4) & 0xf) + '0';
  }
  if ((saved_req->msisdn.digit[len - 1] & 0xf0) == 0xf0) {
    len = (len << 1) - 1;
  } else {
    len = len << 1;
  }
  return len;
}

int get_imeisv_from_session_req(
    const itti_s11_create_session_request_t* saved_req, char* imeisv) {
  if (saved_req->mei.present & MEI_IMEISV) {
    // IMEISV as defined in 3GPP TS 23.003 MEI_IMEISV
    imeisv[0]                 = saved_req->mei.choice.imeisv.u.num.tac8;
    imeisv[1]                 = saved_req->mei.choice.imeisv.u.num.tac7;
    imeisv[2]                 = saved_req->mei.choice.imeisv.u.num.tac6;
    imeisv[3]                 = saved_req->mei.choice.imeisv.u.num.tac5;
    imeisv[4]                 = saved_req->mei.choice.imeisv.u.num.tac4;
    imeisv[5]                 = saved_req->mei.choice.imeisv.u.num.tac3;
    imeisv[6]                 = saved_req->mei.choice.imeisv.u.num.tac2;
    imeisv[7]                 = saved_req->mei.choice.imeisv.u.num.tac1;
    imeisv[8]                 = saved_req->mei.choice.imeisv.u.num.snr6;
    imeisv[9]                 = saved_req->mei.choice.imeisv.u.num.snr5;
    imeisv[10]                = saved_req->mei.choice.imeisv.u.num.snr4;
    imeisv[11]                = saved_req->mei.choice.imeisv.u.num.snr3;
    imeisv[12]                = saved_req->mei.choice.imeisv.u.num.snr2;
    imeisv[13]                = saved_req->mei.choice.imeisv.u.num.snr1;
    imeisv[14]                = saved_req->mei.choice.imeisv.u.num.svn2;
    imeisv[15]                = saved_req->mei.choice.imeisv.u.num.svn1;
    imeisv[IMEISV_DIGITS_MAX] = '\0';

    return 1;
  }
  return 0;
}

void get_session_req_data(
    spgw_state_t* spgw_state,
    const itti_s11_create_session_request_t* saved_req,
    struct pcef_create_session_data* data) {
  const bearer_qos_t* qos;

  data->msisdn_len = get_msisdn_from_session_req(saved_req, data->msisdn);

  data->imeisv_exists = get_imeisv_from_session_req(saved_req, data->imeisv);
  data->uli_exists    = get_uli_from_session_req(saved_req, data->uli);
  get_plmn_from_session_req(saved_req, data);
  get_imsi_plmn_from_session_req(saved_req, data);
  memcpy(
      &data->charging_characteristics, &saved_req->charging_characteristics,
      sizeof(charging_characteristics_t));

  memcpy(data->apn, saved_req->apn, APN_MAX_LENGTH + 1);
  data->pdn_type = saved_req->pdn_type;

  inet_ntop(
      AF_INET, &spgw_state->sgw_ip_address_S1u_S12_S4_up, data->sgw_ip,
      INET_ADDRSTRLEN);

  // QoS Info
  data->ambr_dl = saved_req->ambr.br_dl;
  data->ambr_ul = saved_req->ambr.br_ul;
  qos           = &saved_req->bearer_contexts_to_be_created.bearer_contexts[0]
             .bearer_level_qos;
  data->pl  = qos->pl;
  data->pci = qos->pci;
  data->pvi = qos->pvi;
  data->qci = qos->qci;
}
