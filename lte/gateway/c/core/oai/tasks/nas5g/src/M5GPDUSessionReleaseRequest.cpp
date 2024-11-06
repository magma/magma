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
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GPDUSessionReleaseRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
PDUSessionReleaseRequestMsg::PDUSessionReleaseRequestMsg(){};
PDUSessionReleaseRequestMsg::~PDUSessionReleaseRequestMsg(){};

// Decode PDUSessionReleaseRequest Message and its IEs
int PDUSessionReleaseRequestMsg::DecodePDUSessionReleaseRequestMsg(
    uint8_t* buffer, uint32_t len) {
  uint32_t decoded = 0;
  int decoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(buffer, PDU_SESSION_RELEASE_REQ_MIN_LEN,
                                       len);

  if ((decoded_result = this->extended_protocol_discriminator
                            .DecodeExtendedProtocolDiscriminatorMsg(
                                0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = this->pdu_session_identity.DecodePDUSessionIdentityMsg(
           0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           this->pti.DecodePTIMsg(0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = this->message_type.DecodeMessageTypeMsg(
           0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  return decoded;
}

// Encode PDUSessionReleaseRequest Message and its IEs
int PDUSessionReleaseRequestMsg::EncodePDUSessionReleaseRequestMsg(
    uint8_t* buffer, uint32_t len) {
  return 0;
}
}  // namespace magma5g
