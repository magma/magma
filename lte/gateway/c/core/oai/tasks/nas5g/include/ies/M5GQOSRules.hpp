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

#pragma once
#include <sstream>
#include <cstdint>

#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"

namespace magma5g {
// New QOSRule Class
class NewQOSRulePktFilter {
 public:
  uint8_t spare : 2;
  uint8_t pkt_filter_dir : 2;
  uint8_t pkt_filter_id : 4;
  uint8_t len;

  // Size of the content can be max of packet_filter_contents_t
  uint8_t contents[2 * TRAFFIC_FLOW_TEMPLATE_IPV6_ADDR_SIZE + sizeof(uint8_t)];

  NewQOSRulePktFilter();
  ~NewQOSRulePktFilter();
};

// QOSRule Class
class QOSRule {
 public:
#define QOS_ADD_RULE_MIN_LEN 3
#define QOS_DEL_RULE_MIN_LEN 1
#define QOS_RULE_DQR_BIT_SET 0x1
  uint8_t qos_rule_id;
  uint16_t len;
  uint8_t rule_oper_code : 3;
  uint8_t dqr_bit : 1;
  uint8_t no_of_pkt_filters : 4;

  // Max packet filter supported is 4
  NewQOSRulePktFilter new_qos_rule_pkt_filter[4];
  uint8_t qos_rule_precedence;
  uint8_t spare : 1;
  uint8_t segregation : 1;
  uint8_t qfi : 6;

  QOSRule();
  ~QOSRule();
};

// QOSRules IE Class
class QOSRulesMsg {
 public:
#define QOS_RULES_MSG_MIN_LEN 3
#define QOS_RULE_ENTRY_MAX 4
#define QOS_RULES_MSG_BUF_LEN_MAX 4096
  uint8_t iei;
  uint16_t length;
  QOSRule qos_rule[QOS_RULE_ENTRY_MAX];
  QOSRulesMsg();
  ~QOSRulesMsg();

  uint16_t EncodeQOSRulesMsgData(QOSRulesMsg* qos_rules, uint8_t* buffer,
                                 uint32_t len);
  int EncodeQOSRulesMsg(QOSRulesMsg* qos_rules, uint8_t iei, uint8_t* buffer,
                        uint32_t len);
  uint8_t DecodeQOSRulesMsgData(QOSRulesMsg* qos_rules, uint8_t* buffer,
                                uint32_t len);
  int DecodeQOSRulesMsg(QOSRulesMsg* qos_rules, uint8_t iei, uint8_t* buffer,
                        uint32_t len);
};
}  // namespace magma5g
