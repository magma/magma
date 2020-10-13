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
 * PDU Session Est Reject Message */

#include <iostream>
#include <iomanip>
#include <cstring>
#include "PDUSessionEstablishmentReject.h"
#include "CommonDefs.h"

using namespace std;
using namespace magma5g;
namespace magma5g {
int encode(void) {
  int ret = 0;
  uint8_t buffer[5];
  int len = 5;
  PDUSessionEstablishmentRejectMsg msg;

  // Message to be Encoded
  msg.extendedprotocoldiscriminator.extendedprotodiscriminator = 46;
  msg.pdusessionidentity.pdusessionid                          = 1;
  msg.pti.pti                                                  = 1;
  msg.messagetype.msgtype                                      = 0xC3;
  msg.m5gsmcause.causevalue                                    = 7;

  MLOG(MDEBUG) << "\n\n---Encoding Message---\n\n";
  ret = msg.EncodePDUSessionEstablishmentRejectMsg(&msg, buffer, len);

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
