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
#include <glog/logging.h>
#include <google/protobuf/util/message_differencer.h>
#include <gtest/gtest.h>
#include <stdio.h>

#include <future>
#include <memory>
#include <string>
#include <utility>
#include <vector>

#include "SessiondMocks.h"

using ::testing::Test;

namespace magma {

MATCHER_P(CheckCount, count, "") {
  int arg_count = arg.size();
  return arg_count == count;
}

MATCHER_P(CheckRuleCount, count, "") {
  int arg_count = arg.size();
  return arg_count == count;
}

MATCHER_P(CheckRuleNames, list_static_rules, "") {
  std::vector<RuleToProcess> to_process = arg;
  if (to_process.size() != list_static_rules.size()) {
    return false;
  }
  for (RuleToProcess val : to_process) {
    bool found = false;
    for (const std::string rule_to_check : list_static_rules) {
      if (val.rule.id() == rule_to_check) {
        found = true;
        break;
      }
    }
    if (!found) {
      return false;
    }
  }
  return true;
}

MATCHER_P(CheckRulesToProcess, expected, "") {
  std::vector<RuleToProcess> to_process = arg;
  // basic size check
  if (to_process.size() != expected.size()) {
    return false;
  }
  for (RuleToProcess val : to_process) {
    for (uint32_t i = 0; i < expected.size(); i++) {
      if (val.rule.id() == expected[i].rule.id()) {
        // check teids
        return val.teids == expected[i].teids;
      }
    }
  }
  return false;
}

MATCHER_P(CheckTeids, configured_teids, "") {
  Teids pipelined_req_teids = static_cast<const Teids>(arg);

  if ((pipelined_req_teids.agw_teid() == configured_teids.agw_teid()) &&
      (pipelined_req_teids.enb_teid() == configured_teids.enb_teid())) {
    return true;
  }

  return false;
}

MATCHER_P(CheckTeidVector, expected, "") {
  const std::vector<Teids> req_teids =
      static_cast<const std::vector<Teids>>(arg);
  if (expected.size() != req_teids.size()) {
    return false;
  }
  for (uint32_t i = 0; i < req_teids.size(); i++) {
    const Teids& expected_teids = expected[i];
    const Teids& actual_teids = req_teids[i];
    if ((expected_teids.agw_teid() != actual_teids.agw_teid()) ||
        (expected_teids.enb_teid() != actual_teids.enb_teid())) {
      return false;
    }
  }

  return true;
}

MATCHER_P2(CheckUpdateRequestCount, monitorCount, chargingCount, "") {
  auto req = static_cast<const UpdateSessionRequest>(arg);
  return req.updates().size() == chargingCount &&
         req.usage_monitors().size() == monitorCount;
}

MATCHER_P(CheckUpdateRequestNumber, request_number, "") {
  auto request = static_cast<const UpdateSessionRequest&>(arg);
  for (const auto& credit_usage_update : request.updates()) {
    int req_number = credit_usage_update.request_number();
    return req_number == request_number;
  }
  return false;
}

MATCHER_P(CheckCoreRequest, expected_request, "") {
  auto req = static_cast<const CreateSessionRequest&>(arg);
  auto ex_req = static_cast<const CreateSessionRequest&>(expected_request);
  if (!google::protobuf::util::MessageDifferencer::Equals(
          ex_req.requested_units(), req.requested_units())) {
    return false;
  }

  // Add other check for the request
  return true;
}

MATCHER_P3(CheckTerminateRequestCount, imsi, monitorCount, chargingCount, "") {
  auto req = static_cast<const SessionTerminateRequest>(arg);
  return req.common_context().sid().id() == imsi &&
         req.credit_usages().size() == chargingCount &&
         req.monitor_usages().size() == monitorCount;
}

MATCHER_P6(CheckSessionInfos, imsi_list, ip_address_list, ipv6_address_list,
           cfg, rule_ids_lists, versions_lists, "") {
  auto infos = static_cast<const std::vector<SessionState::SessionInfo>>(arg);

  if (infos.size() != imsi_list.size()) {
    return false;
  }

  for (size_t i = 0; i < infos.size(); i++) {
    SessionState::SessionInfo info = infos[i];
    if (info.imsi != imsi_list[i]) {
      return false;
    }
    if (info.ip_addr != ip_address_list[i]) {
      return false;
    }
    if (info.ipv6_addr != ipv6_address_list[i]) {
      return false;
    }

    std::vector<std::string> expected_gx_rules = rule_ids_lists[i];
    if (info.gx_rules.size() != expected_gx_rules.size()) {
      return false;
    }
    for (size_t r_index = 0; i < info.gx_rules.size(); i++) {
      if (info.gx_rules[r_index].rule.id() != expected_gx_rules[r_index])
        return false;
    }

    std::vector<uint32_t> expected_versions = versions_lists[i];
    for (size_t r_index = 0; i < info.gx_rules.size(); i++) {
      if (info.gx_rules[r_index].version != expected_versions[r_index])
        return false;
    }

    // check ambr field if config has qos_info
    if (cfg.rat_specific_context.has_lte_context() &&
        cfg.rat_specific_context.lte_context().has_qos_info()) {
      const auto& qos_info = cfg.rat_specific_context.lte_context().qos_info();
      if (!info.ambr) {
        return false;
      } else if (info.ambr->max_bandwidth_ul() != qos_info.apn_ambr_ul()) {
        return false;
      } else if (info.ambr->max_bandwidth_dl() != qos_info.apn_ambr_dl()) {
        return false;
      } else if (static_cast<int>(info.ambr->br_unit()) !=
                 static_cast<int>(qos_info.br_unit())) {
        return false;
      }
    }
  }
  return true;
}

MATCHER_P(CheckEventType, expectedEventType, "") {
  return (arg.event_type() == expectedEventType);
}

MATCHER_P2(CheckCreateBearerReq, imsi, rule_count, "") {
  auto request = static_cast<const CreateBearerRequest>(arg);
  return request.sid().id() == imsi &&
         request.policy_rules().size() == rule_count;
}

MATCHER_P3(CheckDeleteOneBearerReq, imsi, link_bearer_id, eps_bearer_id, "") {
  auto request = static_cast<const DeleteBearerRequest>(arg);

  return request.sid().id() == imsi &&
         request.link_bearer_id() == uint32_t(link_bearer_id) &&
         request.eps_bearer_ids_size() == 1 &&
         request.eps_bearer_ids(0) == uint32_t(eps_bearer_id);
}

MATCHER_P(CheckSubset, ids, "") {
  auto request = static_cast<const RulesToProcess>(arg);
  for (size_t i = 0; i < request.size(); i++) {
    if (ids.find(request[i].rule.id()) != ids.end()) {
      return true;
    }
  }
  return false;
}

MATCHER_P(CheckPolicyID, id, "") {
  auto request = static_cast<const RulesToProcess>(arg);
  for (size_t i = 0; i < request.size(); i++) {
    if (request[i].rule.id() == id) {
      return true;
    }
  }
  return false;
}

MATCHER_P2(CheckPolicyIDs, count, ids, "") {
  auto request = static_cast<const RulesToProcess>(arg);
  if (request.size() != (unsigned int)count) {
    return false;
  }
  for (size_t i = 0; i < request.size(); i++) {
    if (ids.find(request[i].rule.id()) != ids.end()) {
      return true;
    }
  }
  return false;
}

MATCHER_P(CheckSubscriberQuotaUpdate, quota, "") {
  auto update = static_cast<std::vector<SubscriberQuotaUpdate>>(arg);
  if (update.size() != 1) {
    return false;
  }
  std::cerr << "\n\n" << update[0].update_type() << " \n\n";
  return update[0].update_type() == quota;
}

MATCHER_P2(CheckCreateSession, imsi, promise_p, "") {
  auto req = static_cast<const CreateSessionRequest*>(arg);
  promise_p->set_value(req->session_id());
  auto res = req->common_context().sid().id() == imsi;
  return res;
}

MATCHER_P(CheckSingleUpdate, expected_update, "") {
  auto request = static_cast<const UpdateSessionRequest*>(arg);
  if (request->updates_size() != 1) {
    return false;
  }

  auto& update = request->updates(0);
  bool val =
      update.usage().type() == expected_update.usage().type() &&
      update.usage().bytes_tx() == expected_update.usage().bytes_tx() &&
      update.usage().bytes_rx() == expected_update.usage().bytes_rx() &&
      update.common_context().sid().id() ==
          expected_update.common_context().sid().id() &&
      update.usage().charging_key() == expected_update.usage().charging_key();
  return val;
}

MATCHER_P(CheckTerminate, imsi, "") {
  auto request = static_cast<const SessionTerminateRequest*>(arg);
  return request->common_context().sid().id() == imsi;
}

MATCHER_P6(CheckActivateFlowsForTunnIds, imsi, ipv4, ipv6, enb_teid, agw_teid,
           rule_count, "") {
  auto request = static_cast<const ActivateFlowsRequest*>(arg);
  std::cerr << "Got " << request->policies_size() << " rules" << std::endl;
  auto res = request->sid().id() == imsi && request->ip_addr() == ipv4 &&
             request->ipv6_addr() == ipv6 &&
             request->uplink_tunnel() == agw_teid &&
             request->downlink_tunnel() == enb_teid &&
             request->policies_size() == rule_count;
  return res;
}

MATCHER_P(CheckDeactivateFlows, imsi, "") {
  auto request = static_cast<const DeactivateFlowsRequest*>(arg);
  return request->sid().id() == imsi;
}

MATCHER_P(CheckSrvResponse, expected_response, "") {
  auto actual_response = static_cast<const SetSMSessionContextAccess>(arg);

  auto expected_session_ambr = &(expected_response->rat_specific_context()
                                     .m5g_session_context_rsp()
                                     .session_ambr());

  auto actual_response_ambr = &(actual_response.rat_specific_context()
                                    .m5g_session_context_rsp()
                                    .session_ambr());

  auto unit_res =
      (expected_session_ambr->br_unit() == actual_response_ambr->br_unit());

  auto ul_ambr_res = (expected_session_ambr->max_bandwidth_ul() ==
                      actual_response_ambr->max_bandwidth_ul());

  auto dl_ambr_res = (expected_session_ambr->max_bandwidth_dl() ==
                      actual_response_ambr->max_bandwidth_dl());

  return (unit_res && ul_ambr_res && dl_ambr_res);
}

MATCHER_P(CheckSendRequest, expected_request, "") {
  auto req = static_cast<const CreateSessionRequest>(arg);
  auto imsi = req.common_context().sid().id();

  auto apn = req.common_context().apn();
  auto rat_type = req.common_context().rat_type();

  auto imsi_exp = expected_request.common_context().sid().id();

  auto apn_exp = expected_request.common_context().apn();
  auto rat_type_exp = expected_request.common_context().rat_type();

  return (imsi == imsi_exp && apn == apn_exp && rat_type == rat_type_exp);
}

MATCHER_P(SessionCheck, request, "") {
  auto req = static_cast<const SessionState::SessionInfo>(arg);

  auto imsi_req = req.subscriber_id;
  uint32_t teid = req.local_f_teid;

  bool res = request.subscriber_id == imsi_req && request.local_f_teid == teid;

  return res;
}

};  // namespace magma
