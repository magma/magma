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

/* using this stub code we are going to test Decoding functionality of
 * DL NAS Treansport with PDU Session Est Request Message */

#include <iostream>
#include <iomanip>
#include <cstring>
#include "M5GDLNASTransport.h"
#include "M5GCommonDefs.h"

using namespace std;
using namespace magma5g;
namespace magma5g {
int encode(void) {
  int ret = 0;
  uint8_t buffer[37];
  int len = 37;
  DLNASTransportMsg msg;

  // Message to be Encoded
  msg.extended_protocol_discriminator.extended_proto_discriminator     = 126;
  msg.sec_header_type.sec_hdr                                          = 0;
  msg.spare_half_octet.spare                                           = 0;
  msg.message_type.msg_type                                            = 0x68;
  msg.payload_container_type.type_val                                  = 1;
  msg.payload_container.len                                            = 32;
  msg.payload_container.smf_msg.header.extended_protocol_discriminator = 46;
  msg.payload_container.smf_msg.header.pdu_session_id                  = 1;
  msg.payload_container.smf_msg.header.procedure_transaction_id        = 1;
  msg.payload_container.smf_msg.header.message_type                    = 0xc2;
  msg.payload_container.smf_msg.pdu_session_estab_accept.pdu_session_type
      .type_val                                                            = 1;
  msg.payload_container.smf_msg.pdu_session_estab_accept.ssc_mode.mode_val = 1;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.length  = 18;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[0]
      .qos_rule_id = 1;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[0]
      .len = 3;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[0]
      .rule_oper_code = 1;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[0]
      .dqr_bit = 1;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[0]
      .no_of_pkt_filters = 0;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[0]
      .qos_rule_precedence = 0;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[0]
      .spare = 0;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[0]
      .segregation = 0;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[0]
      .qfi = 1;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[1]
      .qos_rule_id = 1;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[1]
      .len = 3;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[1]
      .rule_oper_code = 1;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[1]
      .dqr_bit = 0;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[1]
      .no_of_pkt_filters = 0;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[1]
      .qos_rule_precedence = 0;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[1]
      .spare = 0;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[1]
      .segregation = 0;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[1]
      .qfi = 1;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[2]
      .qos_rule_id = 2;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[2]
      .len = 3;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[2]
      .rule_oper_code = 1;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[2]
      .dqr_bit = 0;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[2]
      .no_of_pkt_filters = 0;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[2]
      .qos_rule_precedence = 0;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[2]
      .spare = 0;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[2]
      .segregation = 0;
  msg.payload_container.smf_msg.pdu_session_estab_accept.qos_rules.qos_rule[2]
      .qfi = 2;
  msg.payload_container.smf_msg.pdu_session_estab_accept.session_ambr.length =
      6;
  msg.payload_container.smf_msg.pdu_session_estab_accept.session_ambr.dl_unit =
      0;
  msg.payload_container.smf_msg.pdu_session_estab_accept.session_ambr
      .dl_session_ambr = 0;
  msg.payload_container.smf_msg.pdu_session_estab_accept.session_ambr.ul_unit =
      0;
  msg.payload_container.smf_msg.pdu_session_estab_accept.session_ambr
      .ul_session_ambr = 0;

  MLOG(MDEBUG) << "\n\n---Encoding Message---\n\n";
  ret = msg.EncodeDLNASTransportMsg(&msg, buffer, len);

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
