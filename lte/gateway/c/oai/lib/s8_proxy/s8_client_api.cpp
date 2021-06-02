/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#include <grpcpp/impl/codegen/status.h>
#include "feg/protos/s8_proxy.grpc.pb.h"
#include "orc8r/protos/common.pb.h"
#include "s8_client_api.h"
#include "S8Client.h"
#include "pcef_handlers.h"
extern "C" {
#include "intertask_interface.h"
#include "log.h"
#include "s8_messages_types.h"
#include "common_defs.h"
#include "common_types.h"
extern task_zmq_ctx_t grpc_service_task_zmq_ctx;
}

static void convert_proto_msg_to_itti_csr(
    magma::feg::CreateSessionResponsePgw& response,
    s8_create_session_response_t* s5_response, bearer_qos_t dflt_bearer_qos);

static void get_qos_from_proto_msg(
    const magma::feg::QosInformation& proto_qos, bearer_qos_t* bearer_qos) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  bearer_qos->pci       = proto_qos.pci();
  bearer_qos->pl        = proto_qos.priority_level();
  bearer_qos->pvi       = proto_qos.preemption_vulnerability();
  bearer_qos->qci       = proto_qos.qci();
  bearer_qos->gbr.br_ul = proto_qos.gbr().br_ul();
  bearer_qos->gbr.br_dl = proto_qos.gbr().br_dl();
  bearer_qos->mbr.br_ul = proto_qos.mbr().br_ul();
  bearer_qos->mbr.br_dl = proto_qos.mbr().br_dl();
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void get_fteid_from_proto_msg(
    const magma::feg::Fteid& proto_fteid, fteid_t* pgw_fteid) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  pgw_fteid->teid = proto_fteid.teid();
  if (proto_fteid.ipv4_address().c_str()) {
    pgw_fteid->ipv4 = true;
    inet_pton(
        AF_INET, proto_fteid.ipv4_address().c_str(),
        &(pgw_fteid->ipv4_address));
  }
  if (proto_fteid.ipv6_address().c_str()) {
    pgw_fteid->ipv6 = true;
    inet_pton(
        AF_INET6, proto_fteid.ipv6_address().c_str(),
        &(pgw_fteid->ipv6_address));
  }
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void get_paa_from_proto_msg(
    const magma::feg::PDNType& proto_pdn_type,
    const magma::feg::PdnAddressAllocation& proto_paa, paa_t* paa) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  switch (proto_pdn_type) {
    case magma::feg::PDNType::IPV4: {
      paa->pdn_type = IPv4;
      auto ip       = proto_paa.ipv4_address();
      inet_pton(AF_INET, ip.c_str(), &(paa->ipv4_address));
      break;
    }
    case magma::feg::PDNType::IPV6: {
      paa->pdn_type = IPv6;
      auto ip       = proto_paa.ipv6_address();
      inet_pton(AF_INET6, ip.c_str(), &(paa->ipv6_address));
      paa->ipv6_prefix_length = IPV6_PREFIX_LEN;
      break;
    }
    case magma::feg::PDNType::IPV4V6: {
      paa->pdn_type = IPv4_AND_v6;
      auto ip       = proto_paa.ipv4_address();
      inet_pton(AF_INET, ip.c_str(), &(paa->ipv4_address));
      auto ipv6 = proto_paa.ipv6_address();
      inet_pton(AF_INET6, ipv6.c_str(), &(paa->ipv6_address));
      paa->ipv6_prefix_length = IPV6_PREFIX_LEN;
      break;
    }
    case magma::feg::PDNType::NonIP: {
      OAILOG_ERROR(LOG_SGW_S8, " pdn_type NonIP is not supported \n");
      break;
    }
    default:
      OAILOG_ERROR(
          LOG_SGW_S8,
          "Received invalid pdn_type in create session response \n");
      break;
  }
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void recv_s8_delete_session_response(
    imsi64_t imsi64, teid_t context_teid, const grpc::Status& status,
    magma::feg::DeleteSessionResponsePgw& response) {
  OAILOG_FUNC_IN(LOG_SGW_S8);

  s8_delete_session_response_t* s8_delete_session_rsp = NULL;
  MessageDef* message_p                               = NULL;
  message_p = itti_alloc_new_message(TASK_GRPC_SERVICE, S8_DELETE_SESSION_RSP);
  if (!message_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to allocate memory for S8_DELETE_SESSION_RSP for "
        "context_teid" TEID_FMT "\n",
        context_teid);
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  s8_delete_session_rsp         = &message_p->ittiMsg.s8_delete_session_rsp;
  message_p->ittiMsgHeader.imsi = imsi64;
  s8_delete_session_rsp->context_teid = context_teid;

  if (status.ok()) {
    if (response.has_gtp_error()) {
      s8_delete_session_rsp->cause = response.mutable_gtp_error()->cause();
    } else {
      s8_delete_session_rsp->cause = REQUEST_ACCEPTED;
    }
  } else {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Received gRPC error for delete session response for "
        "context_teid " TEID_FMT "\n",
        context_teid);
    s8_delete_session_rsp->cause = REMOTE_PEER_NOT_RESPONDING;
  }
  OAILOG_INFO_UE(
      LOG_UTIL, imsi64,
      "Sending delete session response to sgw_s8 task for "
      "context_teid " TEID_FMT "\n",
      context_teid);
  if ((send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_SGW_S8, message_p)) !=
      RETURNok) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to send delete session response to sgw_s8 task for"
        "context_teid " TEID_FMT "\n",
        context_teid);
  }
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void recv_s8_create_session_response(
    imsi64_t imsi64, teid_t context_teid, bearer_qos_t dflt_bearer_qos,
    const grpc::Status& status,
    magma::feg::CreateSessionResponsePgw& response) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  s8_create_session_response_t* s5_response = NULL;
  MessageDef* message_p                     = NULL;
  message_p = itti_alloc_new_message(TASK_GRPC_SERVICE, S8_CREATE_SESSION_RSP);
  if (!message_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to allocate memory for S8_CREATE_SESSION_RSP for "
        "context_teid" TEID_FMT "\n",
        context_teid);
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  s5_response                   = &message_p->ittiMsg.s8_create_session_rsp;
  message_p->ittiMsgHeader.imsi = imsi64;
  s5_response->context_teid     = context_teid;
  if (status.ok()) {
    convert_proto_msg_to_itti_csr(response, s5_response, dflt_bearer_qos);
  } else {
    OAILOG_ERROR(
        LOG_SGW_S8,
        "Received gRPC error for create session response for "
        "context_teid " TEID_FMT "\n",
        context_teid);
    s5_response->cause = REMOTE_PEER_NOT_RESPONDING;
  }
  OAILOG_DEBUG(LOG_UTIL, "Sending create session response to sgw_s8 task");
  if ((send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_SGW_S8, message_p)) !=
      RETURNok) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to send S8 CREATE SESSION RESPONSE message to sgw_s8 task "
        "for context_teid " TEID_FMT "\n",
        context_teid);
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void convert_serving_network_to_proto_msg(
    magma::feg::ServingNetwork* serving_network, ServingNetwork_t itti_msg_sn) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  char mcc[3]     = {0};
  char mnc[3]     = {0};
  uint8_t mnc_len = 0;

  mcc[0] = convert_digit_to_char(itti_msg_sn.mcc[0]);
  mcc[1] = convert_digit_to_char(itti_msg_sn.mcc[1]);
  mcc[2] = convert_digit_to_char(itti_msg_sn.mcc[2]);
  mnc[0] = convert_digit_to_char(itti_msg_sn.mnc[0]);
  mnc[1] = convert_digit_to_char(itti_msg_sn.mnc[1]);
  if ((itti_msg_sn.mnc[2] & 0xf) != 0xf) {
    mnc[2]  = convert_digit_to_char(itti_msg_sn.mnc[2]);
    mnc_len = 3;
  } else {
    mnc_len = 2;
  }
  serving_network->set_mcc(mcc, 3);
  serving_network->set_mnc(mnc, mnc_len);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void convert_uli_to_proto_msg(
    magma::feg::UserLocationInformation* uli, Uli_t msg_uli) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  uli->set_lac(msg_uli.s.lai.lac);
  uli->set_ci(msg_uli.s.cgi.ci);
  uli->set_sac(msg_uli.s.sai.sac);
  uli->set_rac(msg_uli.s.rai.rac);
  uli->set_tac(msg_uli.s.tai.tac);
  uli->set_eci(msg_uli.s.ecgi.cell_identity.enb_id);
  uli->set_menbi(0);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void convert_paa_to_proto_msg(
    const itti_s11_create_session_request_t* msg,
    magma::feg::CreateSessionRequestPgw* csr) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  switch (msg->pdn_type) {
    case IPv4:
      csr->set_pdn_type(magma::feg::PDNType::IPV4);
      if (msg->paa.ipv4_address.s_addr) {
        char ip4_str[INET_ADDRSTRLEN];
        inet_ntop(
            AF_INET, &(msg->paa.ipv4_address.s_addr), ip4_str, INET_ADDRSTRLEN);
        csr->mutable_paa()->set_ipv4_address(ip4_str);
      }
      break;
    case IPv6: {
      csr->set_pdn_type(magma::feg::PDNType::IPV6);
      char ip6_str[INET6_ADDRSTRLEN];
      inet_ntop(AF_INET6, &(msg->paa.ipv6_address), ip6_str, INET6_ADDRSTRLEN);
      csr->mutable_paa()->set_ipv6_address(ip6_str);
      csr->mutable_paa()->set_ipv6_prefix(msg->paa.ipv6_prefix_length);
      break;
    }
    case IPv4_AND_v6: {
      csr->set_pdn_type(magma::feg::PDNType::IPV4V6);
      if (msg->paa.ipv4_address.s_addr) {
        char ip4_str[INET_ADDRSTRLEN];
        inet_ntop(
            AF_INET, &(msg->paa.ipv4_address.s_addr), ip4_str, INET_ADDRSTRLEN);
        csr->mutable_paa()->set_ipv4_address(ip4_str);
      }
      char ip6_str[INET6_ADDRSTRLEN];
      inet_ntop(AF_INET6, &(msg->paa.ipv6_address), ip6_str, INET6_ADDRSTRLEN);
      csr->mutable_paa()->set_ipv6_address(ip6_str);
      csr->mutable_paa()->set_ipv6_prefix(msg->paa.ipv6_prefix_length);
      break;
    }
    default:
      std::cout << "[ERROR] Invalid pdn type " << std::endl;
      break;
  }
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void convert_indication_flag_to_proto_msg(
    const itti_s11_create_session_request_t* msg,
    magma::feg::CreateSessionRequestPgw* csr) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
#define INDICATION_FLAG_SIZE 3
  char indication_flag[INDICATION_FLAG_SIZE] = {0};
  indication_flag[0] = (msg->indication_flags.daf << DAF_FLAG_BIT_POS) |
                       (msg->indication_flags.dtf << DTF_FLAG_BIT_POS) |
                       (msg->indication_flags.hi << HI_FLAG_BIT_POS) |
                       (msg->indication_flags.dfi << DFI_FLAG_BIT_POS) |
                       (msg->indication_flags.oi << OI_FLAG_BIT_POS) |
                       (msg->indication_flags.isrsi << ISRSI_FLAG_BIT_POS) |
                       (msg->indication_flags.israi << ISRAI_FLAG_BIT_POS) |
                       (msg->indication_flags.sgwci << SGWCI_FLAG_BIT_POS);

  indication_flag[1] = (msg->indication_flags.sqci << SQSI_FLAG_BIT_POS) |
                       (msg->indication_flags.uimsi << UIMSI_FLAG_BIT_POS) |
                       (msg->indication_flags.cfsi << CFSI_FLAG_BIT_POS) |
                       (msg->indication_flags.crsi << CRSI_FLAG_BIT_POS) |
                       (msg->indication_flags.p << P_FLAG_BIT_POS) |
                       (msg->indication_flags.pt << PT_FLAG_BIT_POS) |
                       (msg->indication_flags.si << SI_FLAG_BIT_POS) |
                       (msg->indication_flags.msv << MSV_FLAG_BIT_POS);

  indication_flag[2] = (msg->indication_flags.s6af << S6AF_FLAG_BIT_POS) |
                       (msg->indication_flags.s4af << S4AF_FLAG_BIT_POS) |
                       (msg->indication_flags.mbmdt << MBMDT_FLAG_BIT_POS) |
                       (msg->indication_flags.israu << ISRAU_FLAG_BIT_POS) |
                       (msg->indication_flags.ccrsi << CCRSI_FLAG_BIT_POS);

  csr->set_indication_flag(indication_flag, INDICATION_FLAG_SIZE);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void convert_time_zone_to_proto_msg(
    const UETimeZone_t* msg_tz, magma::feg::TimeZone* csr_tz) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  csr_tz->set_delta_seconds(msg_tz->time_zone);
  csr_tz->set_daylight_saving_time(msg_tz->daylight_saving_time);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void convert_qos_to_proto_msg(
    const bearer_qos_t* msg_qos, magma::feg::QosInformation* proto_qos) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  proto_qos->set_pci(msg_qos->pci);
  proto_qos->set_priority_level(msg_qos->pl);
  proto_qos->set_preemption_capability(msg_qos->pci);
  proto_qos->set_preemption_vulnerability(msg_qos->pvi);
  proto_qos->set_qci(msg_qos->qci);
  proto_qos->mutable_gbr()->set_br_ul(msg_qos->gbr.br_ul);
  proto_qos->mutable_gbr()->set_br_dl(msg_qos->gbr.br_dl);
  proto_qos->mutable_mbr()->set_br_ul(msg_qos->mbr.br_ul);
  proto_qos->mutable_mbr()->set_br_dl(msg_qos->mbr.br_dl);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void convert_bearer_context_to_proto(
    const bearer_context_to_be_created_t* msg_bc,
    magma::feg::BearerContext* bc) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  bc->set_id(msg_bc->eps_bearer_id);
  char sgw_s8_up_ip[INET_ADDRSTRLEN];
  inet_ntop(
      AF_INET, &msg_bc->s5_s8_u_sgw_fteid.ipv4_address.s_addr, sgw_s8_up_ip,
      INET_ADDRSTRLEN);
  bc->mutable_user_plane_fteid()->set_ipv4_address(sgw_s8_up_ip);
  bc->mutable_user_plane_fteid()->set_teid(msg_bc->s5_s8_u_sgw_fteid.teid);
  convert_qos_to_proto_msg(&msg_bc->bearer_level_qos, bc->mutable_qos());
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void convert_imeisv_to_string(char* imeisv) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  uint8_t idx = 0;
  for (; idx < IMEISV_DIGITS_MAX; idx++) {
    imeisv[idx] = convert_digit_to_char(imeisv[idx]);
  }
  imeisv[idx] = '\0';

  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void convert_pco_to_proto_msg(
    protocol_configuration_options_t pco,
    magma::feg::CreateSessionRequestPgw* csr) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  magma::feg::PcoProtocolOrContainerId* proto_pco = NULL;
  pco_protocol_or_container_id_t csr_pco          = {0};
  csr->mutable_protocol_configuration_options()->set_config_protocol(
      pco.configuration_protocol);
  for (uint8_t idx = 0; idx < pco.num_protocol_or_container_id; idx++) {
    proto_pco = csr->mutable_protocol_configuration_options()
                    ->add_proto_or_container_id();
    csr_pco = pco.protocol_or_container_ids[idx];
    proto_pco->set_id(csr_pco.id);
    proto_pco->set_contents(
        std::string(bdata(csr_pco.contents), blength(csr_pco.contents)));
  }
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void fill_s8_create_session_req(
    const itti_s11_create_session_request_t* msg,
    magma::feg::CreateSessionRequestPgw* csr, teid_t sgw_s8_teid) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  csr->Clear();
  char msisdn[MSISDN_LENGTH + 1];
  int msisdn_len = get_msisdn_from_session_req(msg, msisdn);
  csr->set_imsi((char*) msg->imsi.digit, msg->imsi.length);
  csr->set_msisdn((char*) msisdn, msisdn_len);
  char imeisv[IMEISV_DIGITS_MAX + 1];
  get_imeisv_from_session_req(msg, imeisv);
  convert_imeisv_to_string(imeisv);
  csr->set_mei(imeisv, IMEISV_DIGITS_MAX);
  convert_serving_network_to_proto_msg(
      csr->mutable_serving_network(), msg->serving_network);
  convert_uli_to_proto_msg(csr->mutable_uli(), msg->uli);
  csr->set_rat_type(magma::feg::RATType::EUTRAN);
  convert_paa_to_proto_msg(msg, csr);
  csr->set_apn(msg->apn);
  csr->mutable_ambr()->set_br_ul(msg->ambr.br_ul);
  csr->mutable_ambr()->set_br_dl(msg->ambr.br_dl);

  if (msg->bearer_contexts_to_be_created.num_bearer_context) {
    magma::feg::BearerContext* bc = csr->mutable_bearer_context();
    convert_bearer_context_to_proto(
        &msg->bearer_contexts_to_be_created.bearer_contexts[0], bc);
  }
  csr->set_c_agw_teid(sgw_s8_teid);
  csr->set_charging_characteristics(
      msg->charging_characteristics.value,
      msg->charging_characteristics.length);

  convert_indication_flag_to_proto_msg(msg, csr);
  convert_time_zone_to_proto_msg(&msg->ue_time_zone, csr->mutable_time_zone());
  convert_pco_to_proto_msg(msg->pco, csr);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

void send_s8_create_session_request(
    teid_t sgw_s11_teid, const itti_s11_create_session_request_t* msg,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  magma::feg::CreateSessionRequestPgw csr_req;
  bearer_qos_t dflt_bearer_qos = {0};

  // teid shall remain same for both sgw's s11 interface and s8 interface as
  // teid is allocated per PDN
  OAILOG_INFO_UE(
      LOG_SGW_S8, imsi64,
      "Sending create session request for context_tied " TEID_FMT "\n",
      sgw_s11_teid);

  fill_s8_create_session_req(msg, &csr_req, sgw_s11_teid);
  dflt_bearer_qos =
      msg->bearer_contexts_to_be_created.bearer_contexts[0].bearer_level_qos;

  magma::S8Client::s8_create_session_request(
      csr_req,
      [imsi64, sgw_s11_teid, dflt_bearer_qos](
          grpc::Status status, magma::feg::CreateSessionResponsePgw response) {
        recv_s8_create_session_response(
            imsi64, sgw_s11_teid, dflt_bearer_qos, status, response);
      });
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void get_pco_from_proto_msg(
    const magma::feg::ProtocolConfigurationOptions& proto_pco,
    protocol_configuration_options_t* s8_pco) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  uint8_t idx                          = 0;
  s8_pco->configuration_protocol       = proto_pco.config_protocol();
  s8_pco->num_protocol_or_container_id = proto_pco.proto_or_container_id_size();
  auto proto_pco_ids                   = proto_pco.proto_or_container_id();
  for (auto ptr = proto_pco_ids.begin(); ptr < proto_pco_ids.end(); ptr++) {
    s8_pco->protocol_or_container_ids[idx].id = ptr->id();
    if (ptr->contents().length()) {
      s8_pco->protocol_or_container_ids[idx].length = ptr->contents().length();
      s8_pco->protocol_or_container_ids[idx].contents = bfromcstr_with_str_len(
          ptr->contents().c_str(), ptr->contents().length());
    }
    ++idx;
  }
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void convert_proto_msg_to_itti_csr(
    magma::feg::CreateSessionResponsePgw& response,
    s8_create_session_response_t* s5_response, bearer_qos_t dflt_bearer_qos) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  s8_bearer_context_t* s8_bc = &(s5_response->bearer_context[0]);
  s8_bc->eps_bearer_id       = response.bearer_context().id();
  s5_response->eps_bearer_id = s8_bc->eps_bearer_id;

  if (response.has_gtp_error()) {
    s5_response->cause = response.mutable_gtp_error()->cause();
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  } else {
    s5_response->cause = REQUEST_ACCEPTED;
  }

  s8_bc->charging_id = response.bearer_context().charging_id();
  if (response.bearer_context().has_qos()) {
    get_qos_from_proto_msg(response.bearer_context().qos(), &s8_bc->qos);
  } else {
    // If qos is not received from PGW, set the qos that was sent in CS Req
    s8_bc->qos = dflt_bearer_qos;
  }

  if (response.has_protocol_configuration_options()) {
    get_pco_from_proto_msg(
        response.protocol_configuration_options(), &s5_response->pco);
  }
  s5_response->apn_restriction_value = response.apn_restriction();
  get_fteid_from_proto_msg(
      response.c_pgw_fteid(), &s5_response->pgw_s8_cp_teid);
  get_paa_from_proto_msg(
      response.pdn_type(), response.paa(), &s5_response->paa);
  get_fteid_from_proto_msg(
      response.bearer_context().user_plane_fteid(), &s8_bc->pgw_s8_up);

  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

void send_s8_delete_session_request(
    imsi64_t imsi64, Imsi_t imsi, teid_t sgw_s11_teid, teid_t pgw_s5_teid,
    ebi_t bearer_id,
    const itti_s11_delete_session_request_t* const delete_session_req_p) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  OAILOG_INFO_UE(
      LOG_SGW_S8, imsi64,
      "Sending delete session request for context_teid:" TEID_FMT "\n",
      sgw_s11_teid);

  magma::feg::DeleteSessionRequestPgw dsr_req;

  dsr_req.Clear();
  dsr_req.set_imsi((char*) imsi.digit, imsi.length);
  dsr_req.set_bearer_id(bearer_id);
  dsr_req.set_c_pgw_teid(pgw_s5_teid);
  dsr_req.set_c_agw_teid(sgw_s11_teid);
  convert_uli_to_proto_msg(dsr_req.mutable_uli(), delete_session_req_p->uli);
  convert_serving_network_to_proto_msg(
      dsr_req.mutable_serving_network(), delete_session_req_p->serving_network);
  magma::S8Client::s8_delete_session_request(
      dsr_req,
      [imsi64, sgw_s11_teid](
          grpc::Status status, magma::feg::DeleteSessionResponsePgw response) {
        recv_s8_delete_session_response(imsi64, sgw_s11_teid, status, response);
      });

  OAILOG_FUNC_OUT(LOG_SGW_S8);
}
