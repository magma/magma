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
#include "common_defs.h"
extern task_zmq_ctx_t grpc_service_task_zmq_ctx;
}

static void recv_s8_create_session_response(
    imsi64_t imsi64, teid_t context_teid, const grpc::Status& status,
    magma::feg::CreateSessionResponsePgw& response) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  /*TODO send create session response to sgw_s8 task */
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
  uli->set_eci(msg_uli.s.ecgi.cell_identity.cell_id);
  uli->set_menbi(msg_uli.s.ecgi.cell_identity.enb_id);
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
  char sgw_s5s8_up_ip[INET_ADDRSTRLEN];
  inet_ntop(
      AF_INET, &msg_bc->s5_s8_u_pgw_fteid.ipv4_address, sgw_s5s8_up_ip,
      INET_ADDRSTRLEN);
  bc->mutable_user_plane_fteid()->set_ipv4_address(sgw_s5s8_up_ip);
  bc->mutable_user_plane_fteid()->set_teid(msg_bc->s5_s8_u_sgw_fteid.teid);
  convert_qos_to_proto_msg(&msg_bc->bearer_level_qos, bc->mutable_qos());
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void fill_s8_create_session_req(
    const itti_s11_create_session_request_t* msg,
    magma::feg::CreateSessionRequestPgw* csr) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  csr->Clear();
  char msisdn[MSISDN_LENGTH + 1];
  int msisdn_len = get_msisdn_from_session_req(msg, msisdn);
  csr->set_imsi((char*) msg->imsi.digit, msg->imsi.length);
  csr->set_msisdn((char*) msisdn, msisdn_len);
  char mcc[3];
  mcc[0] = convert_digit_to_char(msg->serving_network.mcc[0]);
  mcc[1] = convert_digit_to_char(msg->serving_network.mcc[1]);
  mcc[2] = convert_digit_to_char(msg->serving_network.mcc[2]);
  char mnc[3];
  uint8_t mnc_len = 0;
  mnc[0]          = convert_digit_to_char(msg->serving_network.mnc[0]);
  mnc[1]          = convert_digit_to_char(msg->serving_network.mnc[1]);
  if ((msg->serving_network.mnc[2] & 0xf) != 0xf) {
    mnc[2]  = convert_digit_to_char(msg->serving_network.mnc[2]);
    mnc_len = 3;
  } else {
    mnc[2]  = '\0';
    mnc_len = 2;
  }

  magma::feg::ServingNetwork* serving_network = csr->mutable_serving_network();
  serving_network->set_mcc(mcc, 3);
  serving_network->set_mnc(mnc, mnc_len);
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
  csr->set_charging_characteristics(
      msg->charging_characteristics.value,
      msg->charging_characteristics.length);

  convert_indication_flag_to_proto_msg(msg, csr);
  convert_time_zone_to_proto_msg(&msg->ue_time_zone, csr->mutable_time_zone());
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

void send_s8_create_session_request(
    teid_t sgw_s11_teid, const itti_s11_create_session_request_t* msg,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  magma::feg::CreateSessionRequestPgw csr_req;

  OAILOG_INFO_UE(
      LOG_SGW_S8, imsi64,
      "Sending create session request for context_tied " TEID_FMT "\n",
      sgw_s11_teid);
  fill_s8_create_session_req(msg, &csr_req);

  magma::S8Client::s8_create_session_request(
      csr_req,
      [imsi64, sgw_s11_teid](
          grpc::Status status, magma::feg::CreateSessionResponsePgw response) {
        recv_s8_create_session_response(imsi64, sgw_s11_teid, status, response);
      });
}
