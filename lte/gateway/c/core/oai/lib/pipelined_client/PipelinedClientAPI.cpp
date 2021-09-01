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

#include "PipelinedClientAPI.h"

#include <grpcpp/security/credentials.h>
#include <cstdint>
#include <cstring>
#include <string>

#include "conversions.h"
#include "common_defs.h"
#include "service303.h"
#include "spgw_types.h"
#include "common_types.h"

#include "PipelinedServiceClient.h"

using grpc::Channel;
using grpc::ChannelCredentials;
using grpc::CreateChannel;
using grpc::InsecureChannelCredentials;
using grpc::Status;
using magma::lte::CauseIE;
using magma::lte::PipelinedServiceClient;
using magma::lte::UESessionContextResponse;

void handle_upf_classifier_rpc_call_done(
    const grpc::Status& status, UESessionContextResponse response);

int upf_classifier_add_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, int vlan, struct in_addr enb,
    uint32_t i_tei, uint32_t o_tei, const char* imsi,
    struct ip_flow_dl* flow_dl, uint32_t flow_precedence_dl, const char* apn) {
  std::string imsi_str = std::string(imsi);
  std::string apn_str  = std::string(apn);

  if (ue_ipv6 == nullptr) {
    if (flow_dl == nullptr) {
      PipelinedServiceClient::UpdateUEIPv4SessionSet(
          ue, vlan, enb, i_tei, o_tei, imsi_str, flow_precedence_dl, apn_str,
          UE_SESSION_ACTIVE_STATE, handle_upf_classifier_rpc_call_done);
    } else {
      PipelinedServiceClient::UpdateUEIPv4SessionSetWithFlowdl(
          ue, vlan, enb, i_tei, o_tei, imsi_str, *flow_dl, flow_precedence_dl,
          apn_str, UE_SESSION_ACTIVE_STATE,
          handle_upf_classifier_rpc_call_done);
    }
  } else {
    if (flow_dl == nullptr) {
      PipelinedServiceClient::UpdateUEIPv4v6SessionSet(
          ue, *ue_ipv6, vlan, enb, i_tei, o_tei, imsi_str, flow_precedence_dl,
          apn_str, UE_SESSION_ACTIVE_STATE,
          handle_upf_classifier_rpc_call_done);
    } else {
      PipelinedServiceClient::UpdateUEIPv4v6SessionSetWithFlowdl(
          ue, *ue_ipv6, vlan, enb, i_tei, o_tei, imsi_str, *flow_dl,
          flow_precedence_dl, apn, UE_SESSION_ACTIVE_STATE,
          handle_upf_classifier_rpc_call_done);
    }
  }
  return RETURNok;
}

int upf_classifier_del_tunnel(
    struct in_addr enb, struct in_addr ue, struct in6_addr* ue_ipv6,
    uint32_t i_tei, uint32_t o_tei, struct ip_flow_dl* flow_dl) {
  if (ue_ipv6 == nullptr) {
    if (flow_dl == nullptr) {
      PipelinedServiceClient::UpdateUEIPv4SessionSet(
          enb, ue, i_tei, o_tei, UE_SESSION_UNREGISTERED_STATE,
          handle_upf_classifier_rpc_call_done);
    } else {
      PipelinedServiceClient::UpdateUEIPv4SessionSetWithFlowdl(
          enb, ue, i_tei, o_tei, *flow_dl, UE_SESSION_UNREGISTERED_STATE,
          handle_upf_classifier_rpc_call_done);
    }
  } else {
    if (flow_dl == nullptr) {
      PipelinedServiceClient::UpdateUEIPv4v6SessionSet(
          enb, ue, *ue_ipv6, i_tei, o_tei, UE_SESSION_UNREGISTERED_STATE,
          handle_upf_classifier_rpc_call_done);
    } else {
      PipelinedServiceClient::UpdateUEIPv4v6SessionSetWithFlowdl(
          enb, ue, *ue_ipv6, i_tei, o_tei, *flow_dl,
          UE_SESSION_UNREGISTERED_STATE, handle_upf_classifier_rpc_call_done);
    }
  }
  return RETURNok;
}

// UE_SESSION_SUSPENDED_DATA_STATE -> DISCARDING
int upf_classifier_discard_data_on_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, uint32_t i_tei,
    struct ip_flow_dl* flow_dl) {
  if (ue_ipv6 == nullptr) {
    if (flow_dl == nullptr) {
      PipelinedServiceClient::UpdateUEIPv4SessionSet(
          ue, i_tei, UE_SESSION_SUSPENDED_DATA_STATE,
          handle_upf_classifier_rpc_call_done);
    } else {
      PipelinedServiceClient::UpdateUEIPv4SessionSetWithFlowdl(
          ue, i_tei, *flow_dl, UE_SESSION_SUSPENDED_DATA_STATE,
          handle_upf_classifier_rpc_call_done);
    }
  } else {
    if (flow_dl == nullptr) {
      PipelinedServiceClient::UpdateUEIPv4v6SessionSet(
          ue, *ue_ipv6, i_tei, UE_SESSION_SUSPENDED_DATA_STATE,
          handle_upf_classifier_rpc_call_done);
    } else {
      PipelinedServiceClient::UpdateUEIPv4v6SessionSetWithFlowdl(
          ue, *ue_ipv6, i_tei, *flow_dl, UE_SESSION_SUSPENDED_DATA_STATE,
          handle_upf_classifier_rpc_call_done);
    }
  }
  return RETURNok;
}

// No enb + ACTIVATE -> FORWARDING
int upf_classifier_forward_data_on_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, uint32_t i_tei,
    struct ip_flow_dl* flow_dl, uint32_t flow_precedence_dl) {
  if (ue_ipv6 == nullptr) {
    if (flow_dl == nullptr) {
      PipelinedServiceClient::UpdateUEIPv4SessionSet(
          ue, i_tei, flow_precedence_dl, UE_SESSION_RESUME_DATA_STATE,
          handle_upf_classifier_rpc_call_done);
    } else {
      PipelinedServiceClient::UpdateUEIPv4SessionSetWithFlowdl(
          ue, i_tei, *flow_dl, flow_precedence_dl, UE_SESSION_RESUME_DATA_STATE,
          handle_upf_classifier_rpc_call_done);
    }
  } else {
    if (flow_dl == nullptr) {
      PipelinedServiceClient::UpdateUEIPv4v6SessionSet(
          ue, *ue_ipv6, i_tei, flow_precedence_dl, UE_SESSION_RESUME_DATA_STATE,
          handle_upf_classifier_rpc_call_done);
    } else {
      PipelinedServiceClient::UpdateUEIPv4v6SessionSetWithFlowdl(
          ue, *ue_ipv6, i_tei, *flow_dl, flow_precedence_dl,
          UE_SESSION_RESUME_DATA_STATE, handle_upf_classifier_rpc_call_done);
    }
  }
  return RETURNok;
}

// UE + IDLE -> ADD PAGING
int upf_classifier_add_paging_rule(struct in_addr ue) {
  PipelinedServiceClient::UpdateUEIPv4SessionSet(
      ue, UE_SESSION_INSTALL_IDLE_STATE, handle_upf_classifier_rpc_call_done);
  return RETURNok;
}

// UE + ACTIVATE -> DELETE PAGING
int upf_classifier_delete_paging_rule(struct in_addr ue) {
  PipelinedServiceClient::UpdateUEIPv4SessionSet(
      ue, UE_SESSION_UNINSTALL_IDLE_STATE, handle_upf_classifier_rpc_call_done);
  return RETURNok;
}

void handle_upf_classifier_rpc_call_done(
    const grpc::Status& status, UESessionContextResponse response) {
  if (!status.ok()) {
    OAILOG_ERROR(
        LOG_UTIL, "Error Code=%d, Error Message=%s", status.error_code(),
        status.error_message().c_str());
  }

  if (response.cause_info().cause_ie() != CauseIE::REQUEST_ACCEPTED) {
    std::string ipv4_addr_str = "";
    std::string ipv6_addr_str = "";
    uint32_t oper_type        = response.operation_type();

    ipv4_addr_str = response.ue_ipv4_address().address();
    ipv6_addr_str = response.ue_ipv6_address().address();

    OAILOG_ERROR(
        LOG_UTIL,
        "Failed Message from pipelined UEv4:%s, UEv6:%s, OperationType=%d"
        "causeie=%d",
        ipv4_addr_str.c_str(), ipv6_addr_str.c_str(), oper_type,
        response.cause_info().cause_ie());
  }
}
