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
#include "DLNASTransport.h"
#include "CommonDefs.h"

using namespace std;
using namespace magma5g;
namespace magma5g {
int encode(void) {
  int ret = 0;
  uint8_t buffer[37];
  int len = 37;
  DLNASTransportMsg msg;

  // Message to be Encoded
  msg.extendedprotocoldiscriminator.extendedprotodiscriminator     = 126;
  msg.securityheadertype.securityhdr                               = 0;
  msg.sparehalfoctet.spare                                         = 0;
  msg.messagetype.msgtype                                          = 0x68;
  msg.payloadcontainertype.typeval                                 = 1;
  msg.payloadcontainer.len                                         = 32;
  msg.payloadcontainer.smfmsg.header.extendedprotocoldiscriminator = 46;
  msg.payloadcontainer.smfmsg.header.pdusessionid                  = 1;
  msg.payloadcontainer.smfmsg.header.proceduretractionid           = 1;
  msg.payloadcontainer.smfmsg.header.messagetype                   = 0xc2;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.pdusessiontype
      .typeval                                                              = 1;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.sscmode.modeval = 1;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.length =
      18;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[0]
      .qosruleid = 1;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[0]
      .len = 3;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[0]
      .ruleopercode = 1;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[0]
      .dqrbit = 1;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[0]
      .noofpktfilters = 0;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[0]
      .qosruleprecedence = 0;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[0]
      .spare = 0;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[0]
      .segregation = 0;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[0]
      .qfi = 1;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[1]
      .qosruleid = 1;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[1]
      .len = 3;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[1]
      .ruleopercode = 1;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[1]
      .dqrbit = 0;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[1]
      .noofpktfilters = 0;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[1]
      .qosruleprecedence = 0;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[1]
      .spare = 0;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[1]
      .segregation = 0;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[1]
      .qfi = 1;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[2]
      .qosruleid = 2;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[2]
      .len = 3;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[2]
      .ruleopercode = 1;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[2]
      .dqrbit = 0;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[2]
      .noofpktfilters = 0;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[2]
      .qosruleprecedence = 0;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[2]
      .spare = 0;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[2]
      .segregation = 0;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.qosrules.qosrule[2]
      .qfi = 2;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.sessionambr.length =
      6;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.sessionambr.dlunit =
      0;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.sessionambr
      .dlsessionambr = 0;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.sessionambr.ulunit =
      0;
  msg.payloadcontainer.smfmsg.pdusessionestablishmentaccept.sessionambr
      .ulsessionambr = 0;

  MLOG(MDEBUG) << "\n\n---Encoding Message---\n\n";
  ret = msg.EncodeDLNASTransportMsg(&msg, buffer, len);

  MLOG(MDEBUG) << "---Encoded Message---";
  for(size_t i=0; i <= sizeof(buffer); i++)
  {
    MLOG(MDEBUG) << setfill('0') << hex << setw(2) <<int(buffer[i]);
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
