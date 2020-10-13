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
#include "PDUSessionEstablishmentAccept.h"
#include "CommonDefs.h"

using namespace std;
using namespace magma5g;
namespace magma5g {
int encode(void) {
  int ret = 0;
  uint8_t buffer[31];
  int len = 31;
  PDUSessionEstablishmentAcceptMsg msg;

  // Message to be Encoded
  msg.extendedprotocoldiscriminator.extendedprotodiscriminator = 46;
  msg.pdusessionidentity.pdusessionid                          = 1;
  msg.pti.pti                                                  = 1;
  msg.messagetype.msgtype                                      = 0xC2;
  msg.pdusessiontype.typeval                                   = 1;
  msg.sscmode.modeval                                          = 1;
  msg.qosrules.length                                          = 18;
  msg.qosrules.qosrule[0].qosruleid                            = 1;
  msg.qosrules.qosrule[0].len                                  = 3;
  msg.qosrules.qosrule[0].ruleopercode                         = 1;
  msg.qosrules.qosrule[0].dqrbit                               = 1;
  msg.qosrules.qosrule[0].noofpktfilters                       = 0;
  msg.qosrules.qosrule[0].qosruleprecedence                    = 0;
  msg.qosrules.qosrule[0].spare                                = 0;
  msg.qosrules.qosrule[0].segregation                          = 0;
  msg.qosrules.qosrule[0].qfi                                  = 1;
  msg.qosrules.qosrule[1].qosruleid                            = 1;
  msg.qosrules.qosrule[1].len                                  = 3;
  msg.qosrules.qosrule[1].ruleopercode                         = 1;
  msg.qosrules.qosrule[1].dqrbit                               = 0;
  msg.qosrules.qosrule[1].noofpktfilters                       = 0;
  msg.qosrules.qosrule[1].qosruleprecedence                    = 0;
  msg.qosrules.qosrule[1].spare                                = 0;
  msg.qosrules.qosrule[1].segregation                          = 0;
  msg.qosrules.qosrule[1].qfi                                  = 1;
  msg.qosrules.qosrule[2].qosruleid                            = 2;
  msg.qosrules.qosrule[2].len                                  = 3;
  msg.qosrules.qosrule[2].ruleopercode                         = 1;
  msg.qosrules.qosrule[2].dqrbit                               = 0;
  msg.qosrules.qosrule[2].noofpktfilters                       = 0;
  msg.qosrules.qosrule[2].qosruleprecedence                    = 0;
  msg.qosrules.qosrule[2].spare                                = 0;
  msg.qosrules.qosrule[2].segregation                          = 0;
  msg.qosrules.qosrule[2].qfi                                  = 2;
  msg.sessionambr.length                                       = 6;
  msg.sessionambr.dlunit                                       = 0;
  msg.sessionambr.dlsessionambr                                = 0;
  msg.sessionambr.ulunit                                       = 0;
  msg.sessionambr.ulsessionambr                                = 0;

  MLOG(MDEBUG) << "\n\n---Encoding Message---\n\n";
  ret = msg.EncodePDUSessionEstablishmentAcceptMsg(&msg, buffer, len);

  MLOG(MDEBUG) << "---Encoded Message---";
  for (size_t i = 0; i < sizeof(buffer); i++) {
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
