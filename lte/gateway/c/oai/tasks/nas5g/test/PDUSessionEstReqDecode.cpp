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
 * PDU Session Est Request Message */

#include <iostream>
#include <M5GPDUSessionEstablishmentRequest.h>

using namespace std;
using namespace magma5g;

namespace magma5g {
// Testing Decoding functionality
int decode(void) {
  int ret = 0;

  // Message to be decoded
  uint8_t buffer[] = {0x2e, 0x01, 0x01, 0xC1, 0xFF, 0xFF, 0x91, 0xA1,
                      0x12, 0x01, 0x81, 0x22, 0x04, 0x09, 0x00, 0x00,
                      0x05, 0X25, 0X0C, 0x07, 0x63, 0x61, 0x72, 0x72,
                      0x69, 0x65, 0x72, 0x03, 0x63, 0x6F, 0x6D};
  int len          = 31;
  PDUSessionEstablishmentRequestMsg Req;

  // Decoding Message
  ret = Req.DecodePDUSessionEstablishmentRequestMsg(&Req, buffer, len);
  return 0;
}
}  // namespace magma5g

// Main function to call test decode function
int main(void) {
  int ret;
  ret = magma5g::decode();
  return 0;
}
