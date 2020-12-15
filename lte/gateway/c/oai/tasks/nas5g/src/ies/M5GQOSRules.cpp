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

#include <sstream>
#include <cstdint>
#include <cstring>
#include "M5GQOSRules.h"
#include "M5GCommonDefs.h"

using namespace std;
namespace magma5g {
NewQOSRulePktFilter::NewQOSRulePktFilter(){};
NewQOSRulePktFilter::~NewQOSRulePktFilter(){};
QOSRule::QOSRule(){};
QOSRulesMsg::QOSRulesMsg(){};
QOSRule::~QOSRule(){};
QOSRulesMsg::~QOSRulesMsg(){};

// Decode QOSRules IE
int QOSRulesMsg::DecodeQOSRulesMsg(
    QOSRulesMsg* qos_rules, uint8_t iei, uint8_t* buffer, uint32_t len) {
  // Not yet Implemented, will be suppported POST MVC
  return (0);
};

// Encode QOSRules IE
int QOSRulesMsg::EncodeQOSRulesMsg(
    QOSRulesMsg* qos_rules, uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint16_t encoded = 0;
  uint8_t i        = 0;
  uint8_t j        = 0;

  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, QOSRULE_MIN_LEN, len);

  if (iei > 0) {
    CHECK_IEI_ENCODER((unsigned char) iei, qos_rules->iei);
    *buffer = iei;
    MLOG(MDEBUG) << "In EncodeQOSRulesMsg: iei" << hex << int(*buffer);
    encoded++;
  }

  IES_ENCODE_U16(buffer, encoded, qos_rules->length);
  MLOG(MDEBUG) << "Length : " << hex << int(qos_rules->length);
  while (encoded < (qos_rules->length) && i <= 255) {
    *(buffer + encoded) = qos_rules->qos_rule[i].qos_rule_id;
    MLOG(MDEBUG) << "qos_rule_id: " << hex << int(*(buffer + encoded));
    encoded++;
    IES_ENCODE_U16(buffer, encoded, qos_rules->qos_rule[i].len);
    *(buffer + encoded) =
        0x00 | ((qos_rules->qos_rule[i].rule_oper_code & 0x07) << 5) |
        ((qos_rules->qos_rule[i].dqr_bit & 0x01) << 4) |
        (qos_rules->qos_rule[i].no_of_pkt_filters & 0x0f);
    MLOG(MDEBUG) << "rule_oper_code, dqr_bit, no_of_pkt_filters: " << hex
                 << int(*(buffer + encoded));
    encoded++;
    for (j = 0; j < qos_rules->qos_rule[i].no_of_pkt_filters; j++) {
      *(buffer + encoded) =
          0x00 |
          ((qos_rules->qos_rule[i].new_qos_rule_pkt_filter[j].spare & 0x03)
           << 6) |
          ((qos_rules->qos_rule[i].new_qos_rule_pkt_filter[j].pkt_filter_dir &
            0x03)
           << 4) |
          (qos_rules->qos_rule[i].new_qos_rule_pkt_filter[j].pkt_filter_id &
           0x0f);
      MLOG(MDEBUG) << "pkt_filter_dir, pkt_filter_id: " << hex
                   << int(*(buffer + encoded));
      encoded++;
      *(buffer + encoded) =
          qos_rules->qos_rule[i].new_qos_rule_pkt_filter[j].len;
      MLOG(MDEBUG) << "len: " << hex << int(*(buffer + encoded));
      encoded++;
      memcpy(
          buffer + encoded,
          qos_rules->qos_rule[i].new_qos_rule_pkt_filter[j].contents,
          qos_rules->qos_rule[i].new_qos_rule_pkt_filter[j].len);
      BUFFER_PRINT_LOG(
          buffer + encoded,
          qos_rules->qos_rule[i].new_qos_rule_pkt_filter[j].len);
      encoded = encoded + qos_rules->qos_rule[i].new_qos_rule_pkt_filter[j].len;
      encoded++;
    }

    *(buffer + encoded) = qos_rules->qos_rule[i].qos_rule_precedence;
    MLOG(MDEBUG) << "qos_rule_precedence: " << hex << int(*(buffer + encoded));
    encoded++;
    *(buffer + encoded) = 0x00 | ((qos_rules->qos_rule[i].spare & 0x01) << 7) |
                          ((qos_rules->qos_rule[i].segregation & 0x01) << 6) |
                          (qos_rules->qos_rule[i].qfi & 0x3f);
    MLOG(MDEBUG) << "segregation, qfi: " << hex << int(*(buffer + encoded));
    encoded++;
    i++;
  }

  return (encoded);
};
}  // namespace magma5g
