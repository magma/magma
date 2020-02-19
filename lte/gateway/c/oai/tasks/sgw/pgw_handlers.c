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
 * Unless required by applicable law or agreed to in writing, software * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*! \file pgw_handlers.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#define PGW
#define S5_HANDLERS_C

#include <arpa/inet.h>
#include <netinet/in.h>
#include <stdint.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

#include "assertions.h"
#include "intertask_interface.h"
#include "log.h"
#include "spgw_config.h"
#include "pgw_pco.h"
#include "pgw_ue_ip_address_alloc.h"
#include "pgw_handlers.h"
#include "pcef_handlers.h"
#include "common_defs.h"
#include "3gpp_23.003.h"
#include "3gpp_23.401.h"
#include "3gpp_24.008.h"
#include "3gpp_29.274.h"
#include "common_types.h"
#include "hashtable.h"
#include "intertask_interface_types.h"
#include "ip_forward_messages_types.h"
#include "itti_types.h"
#include "pgw_config.h"
#include "s11_messages_types.h"
#include "service303.h"
#include "sgw_context_manager.h"
#include "sgw_ie_defs.h"

static void get_session_req_data(
  spgw_state_t *spgw_state,
  const itti_s11_create_session_request_t *saved_req,
  struct pcef_create_session_data *data);
static char convert_digit_to_char(char digit);
extern spgw_config_t spgw_config;
extern uint32_t sgw_get_new_s1u_teid(void);
extern void print_bearer_ids_helper(const ebi_t*, uint32_t);
//--------------------------------------------------------------------------------

void handle_s5_create_session_request(
  spgw_state_t* spgw_state,
  teid_t context_teid,
  ebi_t eps_bearer_id)
{
  OAILOG_FUNC_IN(LOG_PGW_APP);
  s_plus_p_gw_eps_bearer_context_information_t *new_bearer_ctxt_info_p = NULL;
  hashtable_rc_t hash_rc = HASH_TABLE_OK;
  itti_sgi_create_end_point_response_t sgi_create_endpoint_resp = {0};
  s5_create_session_response_t s5_response = {0};
  struct in_addr inaddr;
  char *imsi = NULL;
  char *apn = NULL;

  OAILOG_DEBUG(
    LOG_PGW_APP,
    "Handle s5_create_session_request, for context sgw s11 teid, " TEID_FMT
    "EPS bearer id %u\n",
    context_teid,
    eps_bearer_id);
  hash_rc = hashtable_ts_get(
    spgw_state->sgw_state.s11_bearer_context_information,
    context_teid,
    (void**) &new_bearer_ctxt_info_p);

  if (HASH_TABLE_OK != hash_rc) {
    OAILOG_ERROR(
      LOG_PGW_APP,
      "Failed to fetch sgw bearer context from the received context "
      "teid" TEID_FMT "\n",
      context_teid);
    sgi_create_endpoint_resp.status = SGI_STATUS_ERROR_CONTEXT_NOT_FOUND;
    goto err;
  }

  // PCO processing
  protocol_configuration_options_t* pco_req =
    &new_bearer_ctxt_info_p->sgw_eps_bearer_context_information.saved_message
       .pco;
  protocol_configuration_options_t pco_resp = {0};
  protocol_configuration_options_ids_t pco_ids;
  memset(&pco_ids, 0, sizeof pco_ids);

  if (pgw_process_pco_request(pco_req, &pco_resp, &pco_ids) != RETURNok) {
    OAILOG_ERROR(
      LOG_PGW_APP,
      "Error in processing PCO in create session request for "
      "context_id: " TEID_FMT "\n",
      context_teid);
    sgi_create_endpoint_resp.status = SGI_STATUS_ERROR_FAILED_TO_PROCESS_PCO;
    goto err;
  }
  copy_protocol_configuration_options(&sgi_create_endpoint_resp.pco, &pco_resp);
  clear_protocol_configuration_options(&pco_resp);

  // IP forward will forward packets to this teid
  sgi_create_endpoint_resp.context_teid = context_teid;
  sgi_create_endpoint_resp.eps_bearer_id = eps_bearer_id;
  sgi_create_endpoint_resp.paa.pdn_type =
    new_bearer_ctxt_info_p->sgw_eps_bearer_context_information.saved_message
      .pdn_type;

  imsi =
    (char*)
      new_bearer_ctxt_info_p->sgw_eps_bearer_context_information.imsi.digit;

  apn = (char*) new_bearer_ctxt_info_p->sgw_eps_bearer_context_information
          .pdn_connection.apn_in_use;

  switch (sgi_create_endpoint_resp.paa.pdn_type) {
    case IPv4:
      // Use NAS by default if no preference is set.
      //
      // For context, the protocol configuration options (PCO) section of
      // packet from the UE is optional, which means that it is perfectly
      // valid UE to send no PCO preferences at all. The previous logic only
      // allocates an IPv4 address if the UE has explicitly set the PCO
      // parameter for allocating IPv4 via NAS signaling (as opposed to via
      // DHCPv4). This means that, in the absence of either parameter being,
      // set the does not know what to do, so we need a default option as well.
      //
      // Since we only support the NAS signaling option right now, we will
      // default to using NAS signaling UNLESS we see a preference for DHCPv4.
      // This means that all IPv4 addresses are now allocated via NAS signaling
      // unless specified otherwise.
      //
      // In the long run, we will want to evolve the logic to use whatever
      // information we have to choose the ``best" allocation method. This means
      // adding new bitfields to pco_ids in pgw_pco.h, setting them in pgw_pco.c
      // and using them here in conditional logic. We will also want to
      // implement different logic between the PDN types.
      if (!pco_ids.ci_ipv4_address_allocation_via_dhcpv4) {
        if (0 == allocate_ue_ipv4_address(imsi, apn, &inaddr)) {
          increment_counter(
            "ue_pdn_connection", 1, 2, "pdn_type", "ipv4", "result", "success");
          sgi_create_endpoint_resp.paa.ipv4_address = inaddr;
          OAILOG_DEBUG(
            LOG_PGW_APP,
            "Allocated IPv4 address for imsi <%s>, apn <%s>\n",
            imsi,
            apn);
          sgi_create_endpoint_resp.status = SGI_STATUS_OK;
        } else {
          increment_counter(
            "ue_pdn_connection", 1, 2, "pdn_type", "ipv4", "result", "failure");
          OAILOG_ERROR(
            LOG_PGW_APP,
            "Failed to allocate IPv4 PAA for PDN type IPv4 for "
            "imsi <%s> and apn <%s>\n",
            imsi,
            apn);
          sgi_create_endpoint_resp.status =
            SGI_STATUS_ERROR_ALL_DYNAMIC_ADDRESSES_OCCUPIED;
        }
      }
      break;

    case IPv6:
      increment_counter(
        "ue_pdn_connection", 1, 2, "pdn_type", "ipv4v6", "result", "failure");
      OAILOG_ERROR(LOG_PGW_APP, "IPV6 PDN type NOT Supported\n");
      sgi_create_endpoint_resp.status = SGI_STATUS_ERROR_SERVICE_NOT_SUPPORTED;

      break;

    case IPv4_AND_v6:
      if (0 == allocate_ue_ipv4_address(imsi, apn, &inaddr)) {
        increment_counter(
          "ue_pdn_connection", 1, 2, "pdn_type", "ipv4v6", "result", "success");
        sgi_create_endpoint_resp.paa.ipv4_address = inaddr;
        OAILOG_DEBUG(LOG_PGW_APP, "Allocated IPv4 address\n");
        sgi_create_endpoint_resp.status = SGI_STATUS_OK;
        sgi_create_endpoint_resp.paa.pdn_type = IPv4;
      } else {
        increment_counter(
          "ue_pdn_connection", 1, 2, "pdn_type", "ipv4v6", "result", "failure");
        OAILOG_ERROR(
          LOG_PGW_APP,
          "Failed to allocate IPv4 PAA for PDN type IPv4_AND_v6\n");
        sgi_create_endpoint_resp.status =
          SGI_STATUS_ERROR_ALL_DYNAMIC_ADDRESSES_OCCUPIED;
      }
      break;

    default:
      AssertFatal(
        0, "BAD paa.pdn_type %d", sgi_create_endpoint_resp.paa.pdn_type);
      break;
  }
  if (sgi_create_endpoint_resp.status == SGI_STATUS_OK) {
    // create session in PCEF and return
    s5_create_session_request_t session_req = {0};
    session_req.context_teid = context_teid;
    session_req.eps_bearer_id = eps_bearer_id;
    char ip_str[INET_ADDRSTRLEN];
    inet_ntop(AF_INET, &(inaddr.s_addr), ip_str, INET_ADDRSTRLEN);
    struct pcef_create_session_data session_data;
    get_session_req_data(
      spgw_state,
      &new_bearer_ctxt_info_p->sgw_eps_bearer_context_information.saved_message,
      &session_data);
    pcef_create_session(
      imsi, ip_str, &session_data, sgi_create_endpoint_resp, session_req);
    OAILOG_FUNC_OUT(LOG_PGW_APP);
  }
err:
  s5_response.context_teid = context_teid;
  s5_response.eps_bearer_id = eps_bearer_id;
  s5_response.sgi_create_endpoint_resp = sgi_create_endpoint_resp;
  s5_response.failure_cause = S5_OK;

  OAILOG_DEBUG(
    LOG_PGW_APP,
    "Sending S5 Create Session Response to SGW: with context teid, " TEID_FMT
    "EPS Bearer Id = %u\n",
    s5_response.context_teid,
    s5_response.eps_bearer_id);
  handle_s5_create_session_response(s5_response);
  OAILOG_FUNC_OUT(LOG_PGW_APP);
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

/*
 * Converts ascii values in [0,9] to [48,57]=['0','9']
 * else if they are in [48,57] keep them the same
 * else log an error and return '0'=48 value
 */
static char convert_digit_to_char(char digit)
{
  if ((digit >= 0) && (digit <= 9)) {
    return (digit + '0');
  } else if ((digit >= '0') && (digit <= '9')){
    return digit;
  } else {
    OAILOG_ERROR(
      LOG_PGW_APP,
      "The input value for digit is not in a valid range: "
      "Session request would likely be rejected on Gx or Gy interface\n");
    return '0';
  }
}

static void get_plmn_from_session_req(
  const itti_s11_create_session_request_t* saved_req,
  struct pcef_create_session_data* data)
{
  data->mcc_mnc[0] = convert_digit_to_char(saved_req->serving_network.mcc[0]);
  data->mcc_mnc[1] = convert_digit_to_char(saved_req->serving_network.mcc[1]);
  data->mcc_mnc[2] = convert_digit_to_char(saved_req->serving_network.mcc[2]);
  data->mcc_mnc[3] = convert_digit_to_char(saved_req->serving_network.mnc[0]);
  data->mcc_mnc[4] = convert_digit_to_char(saved_req->serving_network.mnc[1]);
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
  struct pcef_create_session_data* data)
{
  data->imsi_mcc_mnc[0] = convert_digit_to_char(saved_req->imsi.digit[0]);
  data->imsi_mcc_mnc[1] = convert_digit_to_char(saved_req->imsi.digit[1]);
  data->imsi_mcc_mnc[2] = convert_digit_to_char(saved_req->imsi.digit[2]);
  data->imsi_mcc_mnc[3] = convert_digit_to_char(saved_req->imsi.digit[3]);
  data->imsi_mcc_mnc[4] = convert_digit_to_char(saved_req->imsi.digit[4]);
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

static void get_session_req_data(
  spgw_state_t *spgw_state,
  const itti_s11_create_session_request_t *saved_req,
  struct pcef_create_session_data *data)
{
  const bearer_qos_t *qos;

  data->msisdn_len = get_msisdn_from_session_req(saved_req, data->msisdn);

  data->imeisv_exists = get_imeisv_from_session_req(saved_req, data->imeisv);
  data->uli_exists = get_uli_from_session_req(saved_req, data->uli);
  get_plmn_from_session_req(saved_req, data);
  get_imsi_plmn_from_session_req(saved_req, data);

  memcpy(data->apn, saved_req->apn, APN_MAX_LENGTH + 1);

  inet_ntop(
    AF_INET,
    &spgw_state->sgw_state.sgw_ip_address_S1u_S12_S4_up,
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

//------------------------------------------------------------------------------

int spgw_send_nw_init_activate_bearer_rsp(
  gtpv2c_cause_value_t cause,
  Imsi_t imsi,
  uint8_t eps_bearer_id)
{
  uint32_t rc = RETURNok;
  OAILOG_FUNC_IN(LOG_PGW_APP);

  OAILOG_INFO(
    LOG_PGW_APP,
    "Sending Create Bearer Rsp to PCRF with EBI %d\n",
    eps_bearer_id);
  // Send Create Bearer Rsp to PCRF
  // TODO-Uncomment once implemented at PCRF
  /* rc = send_dedicated_bearer_actv_rsp(act_ded_bearer_rsp->ebi,
       act_ded_bearer_rsp->cause);*/
  OAILOG_FUNC_RETURN(LOG_PGW_APP, rc);
}

//------------------------------------------------------------------------------

uint32_t spgw_handle_nw_init_deactivate_bearer_rsp(
  gtpv2c_cause_t cause,
  ebi_t lbi)
{
  uint32_t rc = RETURNok;
  OAILOG_FUNC_IN(LOG_PGW_APP);

  OAILOG_INFO(
    LOG_PGW_APP,
    "Sending Delete Bearer Rsp to PCRF with LBI %u with cause :%d\n",
    lbi,
    cause.cause_value);
  // Send Delete Bearer Rsp to PCRF
  // TODO-Uncomment once implemented at PCRF
  // rc = send_dedicated_bearer_deactv_rsp(lbi, cause);
  OAILOG_FUNC_RETURN(LOG_PGW_APP, rc);
}
