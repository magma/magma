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
#define ONE_K 1024
#include <sstream>
#include <cstdint>

using namespace std;
namespace magma5g {
// New QOSRule Class
class NewQOSRulePktFilter {
 public:
  uint8_t spare : 2;
  uint8_t pkt_filter_dir : 2;
  uint8_t pkt_filter_id : 4;
  uint8_t len;
  uint8_t contents[4 * ONE_K];  // need to revisit if the QOS rules occupy more
                                // space than 4k.
  NewQOSRulePktFilter();
  ~NewQOSRulePktFilter();
};

// QOSRule Class
class QOSRule {
 public:
  uint8_t qos_rule_id;
  uint16_t len;
  uint8_t rule_oper_code : 3;
  uint8_t dqr_bit : 1;
  uint8_t no_of_pkt_filters : 4;
  NewQOSRulePktFilter new_qos_rule_pkt_filter[16];
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
#define QOSRULE_MIN_LEN 3
  uint8_t iei;
  uint8_t length;
  QOSRule qos_rule[32];  // need to revisit based on max num of QOS rules
                         // exchanged btw UE and core.
  QOSRulesMsg();
  ~QOSRulesMsg();
  int EncodeQOSRulesMsg(
      QOSRulesMsg* qos_rules, uint8_t iei, uint8_t* buffer, uint32_t len);
  int DecodeQOSRulesMsg(
      QOSRulesMsg* qos_rules, uint8_t iei, uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g
