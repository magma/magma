/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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
#include <string.h>
#include <string>
#include <conversions.h>

#include "pcef_handlers.h"
#include "PCEFClient.h"
#include "MobilityClientAPI.h"
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "lte/protos/session_manager.pb.h"
#include "lte/protos/subscriberdb.pb.h"
#include "spgw_types.h"

extern "C" {
}

#define ULI_DATA_SIZE 13

static char _convert_digit_to_char(char digit);

static void create_session_response(
  spgw_state_t* state,
  const std::string& imsi,
  const std::string& apn,
  itti_sgi_create_end_point_response_t sgi_response,
  s5_create_session_request_t session_request,
  const grpc::Status& status,
  s_plus_p_gw_eps_bearer_context_information_t* ctx_p)
{
  s5_create_session_response_t s5_response = {0};
  s5_response.context_teid = session_request.context_teid;
  s5_response.eps_bearer_id = session_request.eps_bearer_id;
  s5_response.sgi_create_endpoint_resp = sgi_response;
  s5_response.failure_cause = S5_OK;

  if (!status.ok()) {
    //BUFFER_TO_IN_ADDR (sgi_response.paa.ipv4_address, addr);
    release_ipv4_address(imsi.c_str(), apn.c_str(),
                         &sgi_response.paa.ipv4_address);
    s5_response.failure_cause = PCEF_FAILURE;
  }
  handle_s5_create_session_response(state, ctx_p, s5_response);
}

static void pcef_fill_create_session_req(
  const struct pcef_create_session_data *session_data,
  magma::LocalCreateSessionRequest *sreq)
{
  sreq->set_apn(session_data->apn);
  sreq->set_msisdn(session_data->msisdn, session_data->msisdn_len);
  sreq->set_spgw_ipv4(session_data->sgw_ip);
  sreq->set_plmn_id(session_data->mcc_mnc, session_data->mcc_mnc_len);
  sreq->set_imsi_plmn_id(
    session_data->imsi_mcc_mnc, session_data->imsi_mcc_mnc_len);

  if (session_data->imeisv_exists) {
    sreq->set_imei(session_data->imeisv, IMEISV_DIGITS_MAX);
  }
  if (session_data->uli_exists) {
    sreq->set_user_location(session_data->uli, ULI_DATA_SIZE);
  }

  // QoS Info
  sreq->mutable_qos_info()->set_apn_ambr_dl(session_data->ambr_dl);
  sreq->mutable_qos_info()->set_apn_ambr_ul(session_data->ambr_ul);
  sreq->mutable_qos_info()->set_priority_level(session_data->pl);
  sreq->mutable_qos_info()->set_preemption_capability(session_data->pci);
  sreq->mutable_qos_info()->set_preemption_vulnerability(session_data->pvi);
  sreq->mutable_qos_info()->set_qos_class_id(session_data->qci);
}

void pcef_create_session(
    spgw_state_t* state, const char* imsi, const char* ip4, const char* ip6,
    const pcef_create_session_data* session_data,
    itti_sgi_create_end_point_response_t sgi_response,
    s5_create_session_request_t session_request,
    s_plus_p_gw_eps_bearer_context_information_t* ctx_p) {
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

  sreq.mutable_sid()->set_id("IMSI" + imsi_str);
  sreq.set_rat_type(magma::RATType::TGPP_LTE);
  if ((session_data->pdn_type == IPv4) && (!ip4_str.empty())) {
    sreq.set_ue_ipv4(ip4_str);
  } else if ((session_data->pdn_type == IPv6) && (!ip6_str.empty())) {
    sreq.set_ue_ipv6(ip6_str);
  } else if (
      (session_data->pdn_type == IPv4_AND_v6) && (!ip4_str.empty()) &&
      (!ip6_str.empty())) {
    sreq.set_ue_ipv4(ip4_str);
    sreq.set_ue_ipv6(ip6_str);
  }
  sreq.set_pdn_type(session_data->pdn_type);
  sreq.set_bearer_id(session_request.eps_bearer_id);
  pcef_fill_create_session_req(session_data, &sreq);

  auto apn = std::string(session_data->apn);
  // call the `CreateSession` gRPC method and execute the inline function
  magma::PCEFClient::create_session(
      sreq,
      [imsi_str, apn, sgi_response, session_request, ctx_p, state](
          grpc::Status status, magma::LocalCreateSessionResponse response) {
        create_session_response(
            state, imsi_str, apn, sgi_response, session_request, status, ctx_p);
      });
}

bool pcef_end_session(char *imsi, char *apn)
{
  magma::LocalEndSessionRequest request;
  request.mutable_sid()->set_id("IMSI" + std::string(imsi));
  request.set_apn(apn);
  magma::PCEFClient::end_session(
    request, [&](grpc::Status status, magma::LocalEndSessionResponse response) {
      return; // For now, do nothing. TODO: handle errors asynchronously
    });
  return true;
}


/*
 * Converts ascii values in [0,9] to [48,57]=['0','9']
 * else if they are in [48,57] keep them the same
 * else log an error and return '0'=48 value
 */
static char _convert_digit_to_char(char digit)
{
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
  struct pcef_create_session_data* data)
{
  data->mcc_mnc[0] = _convert_digit_to_char(saved_req->serving_network.mcc[0]);
  data->mcc_mnc[1] = _convert_digit_to_char(saved_req->serving_network.mcc[1]);
  data->mcc_mnc[2] = _convert_digit_to_char(saved_req->serving_network.mcc[2]);
  data->mcc_mnc[3] = _convert_digit_to_char(saved_req->serving_network.mnc[0]);
  data->mcc_mnc[4] = _convert_digit_to_char(saved_req->serving_network.mnc[1]);
  data->mcc_mnc_len = 5;
  if ((saved_req->serving_network.mnc[2] & 0xf) != 0xf) {
    data->mcc_mnc[5] =
      _convert_digit_to_char(saved_req->serving_network.mnc[2]);
    data->mcc_mnc[6] = '\0';
    data->mcc_mnc_len += 1;
  } else {
    data->mcc_mnc[5] = '\0';
  }
}

static void get_imsi_plmn_from_session_req(
  const itti_s11_create_session_request_t* saved_req,
  struct pcef_create_session_data* data)
{
  data->imsi_mcc_mnc[0] = _convert_digit_to_char(saved_req->imsi.digit[0]);
  data->imsi_mcc_mnc[1] = _convert_digit_to_char(saved_req->imsi.digit[1]);
  data->imsi_mcc_mnc[2] = _convert_digit_to_char(saved_req->imsi.digit[2]);
  data->imsi_mcc_mnc[3] = _convert_digit_to_char(saved_req->imsi.digit[3]);
  data->imsi_mcc_mnc[4] = _convert_digit_to_char(saved_req->imsi.digit[4]);
  data->imsi_mcc_mnc_len = 5;
  // Check if 2 or 3 digit by verifying mnc[2] has a valid value
  if ((saved_req->serving_network.mnc[2] & 0xf) != 0xf) {
    data->imsi_mcc_mnc[5] = _convert_digit_to_char(saved_req->imsi.digit[5]);
    data->imsi_mcc_mnc[6] = '\0';
    data->imsi_mcc_mnc_len += 1;
  } else {
    data->imsi_mcc_mnc[5] = '\0';
  }
}

static int get_uli_from_session_req(
  const itti_s11_create_session_request_t *saved_req,
  char *uli)
{
  if (!saved_req->uli.present) {
    return 0;
  }

  uli[0] = 130; // TAI and ECGI - defined in 29.061

  // TAI as defined in 29.274 8.21.4
  uli[1] = ((saved_req->uli.s.tai.mcc[1] & 0xf) << 4) |
           ((saved_req->uli.s.tai.mcc[0] & 0xf));
  uli[2] = ((saved_req->uli.s.tai.mnc[2] & 0xf) << 4) |
           ((saved_req->uli.s.tai.mcc[2] & 0xf));
  uli[3] = ((saved_req->uli.s.tai.mnc[1] & 0xf) << 4) |
           ((saved_req->uli.s.tai.mnc[0] & 0xf));
  uli[4] = (saved_req->uli.s.tai.tac >> 8) & 0xff;
  uli[5] = saved_req->uli.s.tai.tac & 0xff;

  // ECGI as defined in 29.274 8.21.5
  uli[6] = ((saved_req->uli.s.ecgi.mcc[1] & 0xf) << 4) |
           ((saved_req->uli.s.ecgi.mcc[0] & 0xf));
  uli[7] = ((saved_req->uli.s.ecgi.mnc[2] & 0xf) << 4) |
           ((saved_req->uli.s.ecgi.mcc[2] & 0xf));
  uli[8] = ((saved_req->uli.s.ecgi.mnc[1] & 0xf) << 4) |
           ((saved_req->uli.s.ecgi.mnc[0] & 0xf));
  uli[9] = (saved_req->uli.s.ecgi.eci >> 24) & 0xf;
  uli[10] = (saved_req->uli.s.ecgi.eci >> 16) & 0xff;
  uli[11] = (saved_req->uli.s.ecgi.eci >> 8) & 0xff;
  uli[12] = saved_req->uli.s.ecgi.eci & 0xff;
  uli[13] = '\0';
  return 1;
}

static int get_msisdn_from_session_req(
  const itti_s11_create_session_request_t *saved_req,
  char *msisdn)
{
  int len = saved_req->msisdn.length;
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

static int get_imeisv_from_session_req(
  const itti_s11_create_session_request_t *saved_req,
  char *imeisv)
{
  if (saved_req->mei.present & MEI_IMEISV) {
    // IMEISV as defined in 3GPP TS 23.003 MEI_IMEISV
    imeisv[0] = saved_req->mei.choice.imeisv.u.num.tac8;
    imeisv[1] = saved_req->mei.choice.imeisv.u.num.tac7;
    imeisv[2] = saved_req->mei.choice.imeisv.u.num.tac6;
    imeisv[3] = saved_req->mei.choice.imeisv.u.num.tac5;
    imeisv[4] = saved_req->mei.choice.imeisv.u.num.tac4;
    imeisv[5] = saved_req->mei.choice.imeisv.u.num.tac3;
    imeisv[6] = saved_req->mei.choice.imeisv.u.num.tac2;
    imeisv[7] = saved_req->mei.choice.imeisv.u.num.tac1;
    imeisv[8] = saved_req->mei.choice.imeisv.u.num.snr6;
    imeisv[9] = saved_req->mei.choice.imeisv.u.num.snr5;
    imeisv[10] = saved_req->mei.choice.imeisv.u.num.snr4;
    imeisv[11] = saved_req->mei.choice.imeisv.u.num.snr3;
    imeisv[12] = saved_req->mei.choice.imeisv.u.num.snr2;
    imeisv[13] = saved_req->mei.choice.imeisv.u.num.snr1;
    imeisv[14] = saved_req->mei.choice.imeisv.u.num.svn2;
    imeisv[15] = saved_req->mei.choice.imeisv.u.num.svn1;
    imeisv[IMEISV_DIGITS_MAX] = '\0';

    return 1;
  }
  return 0;
}

void get_session_req_data(
  spgw_state_t* spgw_state,
  const itti_s11_create_session_request_t* saved_req,
  struct pcef_create_session_data* data)
{
  const bearer_qos_t *qos;

  data->msisdn_len = get_msisdn_from_session_req(saved_req, data->msisdn);

  data->imeisv_exists = get_imeisv_from_session_req(saved_req, data->imeisv);
  data->uli_exists = get_uli_from_session_req(saved_req, data->uli);
  get_plmn_from_session_req(saved_req, data);
  get_imsi_plmn_from_session_req(saved_req, data);

  memcpy(data->apn, saved_req->apn, APN_MAX_LENGTH + 1);
  data->pdn_type = saved_req->pdn_type;

  inet_ntop(
    AF_INET,
    &spgw_state->sgw_ip_address_S1u_S12_S4_up,
    data->sgw_ip,
    INET_ADDRSTRLEN);

  // QoS Info
  data->ambr_dl = saved_req->ambr.br_dl;
  data->ambr_ul = saved_req->ambr.br_ul;
  qos = &saved_req->bearer_contexts_to_be_created.bearer_contexts[0]
    .bearer_level_qos;
  data->pl = qos->pl;
  data->pci = qos->pci;
  data->pvi = qos->pvi;
  data->qci = qos->qci;
}
