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

/* using this stub code we are going to test Encoding functionality of
 * PDU Session Est Accept Message */

#include <iostream>
#include <iomanip>
#include <cstring>
#include "M5GPDUSessionEstablishmentAccept.h"
#include "M5GCommonDefs.h"

using namespace std;
using namespace magma5g;
namespace magma5g {
int encode(void) {
  int ret = 0;
  uint8_t buffer[31];
  int len = 31;
  PDUSessionEstablishmentAcceptMsg msg;

  // Message to be Encoded
  msg.extended_protocol_discriminator.extended_proto_discriminator = 46;
  msg.pdu_session_identity.pdu_session_id                          = 1;
  msg.pti.pti                                                      = 1;
  msg.message_type.msg_type                                        = 0xC2;
  msg.pdu_session_type.type_val                                    = 1;
  msg.ssc_mode.mode_val                                            = 1;
  msg.qos_rules.length                                             = 18;
  msg.qos_rules.qos_rule[0].qos_rule_id                            = 1;
  msg.qos_rules.qos_rule[0].len                                    = 3;
  msg.qos_rules.qos_rule[0].rule_oper_code                         = 1;
  msg.qos_rules.qos_rule[0].dqr_bit                                = 1;
  msg.qos_rules.qos_rule[0].no_of_pkt_filters                      = 0;
  msg.qos_rules.qos_rule[0].qos_rule_precedence                    = 0;
  msg.qos_rules.qos_rule[0].spare                                  = 0;
  msg.qos_rules.qos_rule[0].segregation                            = 0;
  msg.qos_rules.qos_rule[0].qfi                                    = 1;
  msg.qos_rules.qos_rule[1].qos_rule_id                            = 1;
  msg.qos_rules.qos_rule[1].len                                    = 3;
  msg.qos_rules.qos_rule[1].rule_oper_code                         = 1;
  msg.qos_rules.qos_rule[1].dqr_bit                                = 0;
  msg.qos_rules.qos_rule[1].no_of_pkt_filters                      = 0;
  msg.qos_rules.qos_rule[1].qos_rule_precedence                    = 0;
  msg.qos_rules.qos_rule[1].spare                                  = 0;
  msg.qos_rules.qos_rule[1].segregation                            = 0;
  msg.qos_rules.qos_rule[1].qfi                                    = 1;
  msg.qos_rules.qos_rule[2].qos_rule_id                            = 2;
  msg.qos_rules.qos_rule[2].len                                    = 3;
  msg.qos_rules.qos_rule[2].rule_oper_code                         = 1;
  msg.qos_rules.qos_rule[2].dqr_bit                                = 0;
  msg.qos_rules.qos_rule[2].no_of_pkt_filters                      = 0;
  msg.qos_rules.qos_rule[2].qos_rule_precedence                    = 0;
  msg.qos_rules.qos_rule[2].spare                                  = 0;
  msg.qos_rules.qos_rule[2].segregation                            = 0;
  msg.qos_rules.qos_rule[2].qfi                                    = 2;
  msg.session_ambr.length                                          = 6;
  msg.session_ambr.dl_unit                                         = 0;
  msg.session_ambr.dl_session_ambr                                 = 0;
  msg.session_ambr.ul_unit                                         = 0;
  msg.session_ambr.ul_session_ambr                                 = 0;
  msg.dnn.iei                                                      = 0x25;
  msg.dnn.len                                                      = 12;
  msg.dnn.dnn = "carrier.com";

  MLOG(MDEBUG) << "\n\n---Encoding Message---\n\n";
  ret = msg.EncodePDUSessionEstablishmentAcceptMsg(&msg, buffer, len);

  MLOG(MDEBUG) << "---Encoded Message---";
  for (size_t i = 0; i <= sizeof(buffer); i++) {
    MLOG(MDEBUG) << setfill('0') << hex << setw(2) << int(buffer[i]);
  }

  return 0;
}
}  // namespace magma5g

// Main Function to call test decode function
int main(void) {
  int ret;
  ret = magma5g::encode();

  return 0;
}
