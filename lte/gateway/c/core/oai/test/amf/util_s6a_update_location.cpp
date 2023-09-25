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

#include <iostream>

#include "lte/gateway/c/core/oai/test/amf/util_nas5g_pkt.hpp"
#include "lte/gateway/c/core/oai/include/s6a_messages_types.hpp"
#include "lte/gateway/c/core/oai/test/amf/util_s6a_update_location.hpp"

namespace magma5g {
// api to mock handling of s6a_update_location_ans

s6a_update_location_ans_t util_amf_send_s6a_ula(const std::string& imsi) {
  s6a_update_location_ans_t itti_msg = {};

  // building s6a_update_location_ans_t
  strncpy(itti_msg.imsi, imsi.c_str(), imsi.size());
  itti_msg.imsi_length = imsi.size();
  itti_msg.result.present = S6A_RESULT_BASE;
  itti_msg.result.choice.base = DIAMETER_SUCCESS;

  // ambr
  memset(&itti_msg.subscription_data.subscribed_ambr.br_unit, 0,
         sizeof(apn_ambr_bitrate_unit_t));
  memset(&itti_msg.subscription_data.subscribed_ambr.br_ul, 100000000,
         sizeof(bitrate_t));
  memset(&itti_msg.subscription_data.subscribed_ambr.br_dl, 200000000,
         sizeof(bitrate_t));

  // apnconfig
  itti_msg.subscription_data.apn_config_profile.context_identifier = 0;
  itti_msg.subscription_data.apn_config_profile.apn_configuration[0]
      .context_identifier = 0;

  itti_msg.subscription_data.apn_config_profile.nb_apns = 1;

  memset(&itti_msg.subscription_data.apn_config_profile.apn_configuration[0]
              .ambr.br_unit,
         0, sizeof(apn_ambr_bitrate_unit_t));
  memset(&itti_msg.subscription_data.apn_config_profile.apn_configuration[0]
              .ambr.br_ul,
         100000000, sizeof(bitrate_t));
  memset(&itti_msg.subscription_data.apn_config_profile.apn_configuration[0]
              .ambr.br_dl,
         200000000, sizeof(bitrate_t));

  memcpy(&itti_msg.subscription_data.apn_config_profile.apn_configuration[0]
              .service_selection,
         "internet", strlen("internet") + 1);
  memset(
      &itti_msg.subscription_data.apn_config_profile.apn_configuration[0]
           .service_selection_length,
      strlen("internet") + 1,
      sizeof(itti_msg.subscription_data.apn_config_profile.apn_configuration[0]
                 .service_selection_length));

  itti_msg.subscription_data.apn_config_profile.apn_configuration[0]
      .subscribed_qos.qci = 8;
  itti_msg.subscription_data.apn_config_profile.apn_configuration[0]
      .subscribed_qos.allocation_retention_priority.priority_level = 8;
  itti_msg.subscription_data.apn_config_profile.apn_configuration[0]
      .subscribed_qos.allocation_retention_priority.pre_emp_vulnerability =
      static_cast<pre_emption_vulnerability_t>(
          PRE_EMPTION_VULNERABILITY_ENABLED);
  itti_msg.subscription_data.apn_config_profile.apn_configuration[0]
      .subscribed_qos.allocation_retention_priority.pre_emp_capability =
      static_cast<pre_emption_capability_t>(PRE_EMPTION_CAPABILITY_ENABLED);

  return itti_msg;
}

}  // namespace magma5g
