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
#include "M5GEAPMessage.h"
#include "M5GCommonDefs.h"

using namespace std;
namespace magma5g {
EAPMessageMsg::EAPMessageMsg(){};
EAPMessageMsg::~EAPMessageMsg(){};

// Decode EAP Message
int EAPMessageMsg::DecodeEAPMessageMsg(
    EAPMessageMsg* eap_message, uint8_t iei, uint8_t* buffer, uint32_t len) {
  int decoded    = 0;
  uint32_t ielen = 0;

  /*** Will be supported POST MVC ***/
  if (iei > 0) {
    CHECK_IEI_DECODER(iei, (unsigned char) *buffer);
    IES_DECODE_U16(buffer, decoded, ielen);
  }
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  return decoded;
};

// Encode EAP Message
int EAPMessageMsg::EncodeEAPMessageMsg(
    EAPMessageMsg* eap_message, uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, EAP_MIN_LENGTH, len);

  if (iei > 0) {
    CHECK_IEI_ENCODER(iei, (unsigned char) eap_message->iei);
    *buffer = iei;
    encoded++;
  }

  MLOG(MDEBUG) << "EncodeEAPMessage : ";
  IES_ENCODE_U16(buffer, encoded, eap_message->len);
  MLOG(MDEBUG) << "Length = " << hex << int(eap_message->len);
  std::copy(eap_message->eap.begin(), eap_message->eap.end(), buffer + encoded);
  BUFFER_PRINT_LOG(buffer + encoded, eap_message->eap.length());
  encoded = encoded + eap_message->eap.length();

  return encoded;
};
}  // namespace magma5g
