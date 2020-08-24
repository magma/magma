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
#include <future>
#include <memory>
#include <utility>

#include <glog/logging.h>
#include <gtest/gtest.h>

#include "SessiondMocks.h"

using ::testing::Test;

namespace magma {

MATCHER_P(CheckCount, count, "") {
  int arg_count = arg.size();
  return arg_count == count;
}

MATCHER_P2(CheckUpdateRequestCount, monitorCount, chargingCount, "") {
  auto req = static_cast<const UpdateSessionRequest>(arg);
  return req.updates().size() == chargingCount &&
         req.usage_monitors().size() == monitorCount;
}

MATCHER_P3(CheckTerminateRequestCount, imsi, monitorCount, chargingCount, "") {
  auto req = static_cast<const SessionTerminateRequest>(arg);
  return req.sid() == imsi && req.credit_usages().size() == chargingCount &&
         req.monitor_usages().size() == monitorCount;
}

MATCHER_P2(CheckActivateFlows, imsi, rule_count, "") {
  auto request = static_cast<const ActivateFlowsRequest*>(arg);
  return request->sid().id() == imsi && request->rule_ids_size() == rule_count;
}

MATCHER_P5(
    CheckSessionInfos, imsi_list, ip_address_list, cfg, static_rule_lists,
    dynamic_rule_ids_lists, "") {
  auto infos = static_cast<const std::vector<SessionState::SessionInfo>>(arg);

  if (infos.size() != imsi_list.size()) return false;

  for (size_t i = 0; i < infos.size(); i++) {
    if (infos[i].imsi != imsi_list[i]) return false;
    if (infos[i].ip_addr != ip_address_list[i]) return false;
    if (infos[i].static_rules.size() != static_rule_lists[i].size())
      return false;
    if (infos[i].dynamic_rules.size() != dynamic_rule_ids_lists[i].size())
      return false;
    for (size_t r_index = 0; i < infos[i].static_rules.size(); i++) {
      if (infos[i].static_rules[r_index] != static_rule_lists[i][r_index])
        return false;
    }
    for (size_t r_index = 0; i < infos[i].dynamic_rules.size(); i++) {
      if (infos[i].dynamic_rules[r_index].id() !=
          dynamic_rule_ids_lists[i][r_index])
        return false;
    }
    // check ambr field if config has qos_info
    if (cfg.rat_specific_context.has_lte_context() &&
        cfg.rat_specific_context.lte_context().has_qos_info()) {
      const auto& qos_info = cfg.rat_specific_context.lte_context().qos_info();
      if (!infos[i].ambr) {
        return false;
      } else if (infos[i].ambr->max_bandwidth_ul() != qos_info.apn_ambr_ul()) {
        return false;
      } else if (infos[i].ambr->max_bandwidth_dl() != qos_info.apn_ambr_dl()) {
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

};  // namespace magma