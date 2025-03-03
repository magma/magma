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
#include <string.h>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GNSSAI.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
namespace magma5g {
NSSAIMsg::NSSAIMsg() {};

NSSAIMsg::~NSSAIMsg() {};

int NSSAIMsg::EncodeNSSAIMsg(NSSAIMsg* NSSAI, uint8_t iei, uint8_t* buffer,
                             uint32_t len) {
  uint8_t encoded = 0;

  if (iei > 0) {
    CHECK_IEI_ENCODER(iei, (unsigned char)NSSAI->iei);
    ENCODE_U8(buffer, iei, encoded);
  }

  ENCODE_U8(buffer + encoded, NSSAI->len, encoded);

  switch (NSSAI->len) {
    case 0b00000001:  // SST
      ENCODE_U8(buffer + encoded, NSSAI->sst, encoded);
      break;
    case 0b00000010:  // SST and mapped HPLMN SST
      ENCODE_U8(buffer + encoded, NSSAI->sst, encoded);

      ENCODE_U8(buffer + encoded, NSSAI->hplmn_sst, encoded);
      break;
    case 0b00000100:  // SST and SD
      ENCODE_U8(buffer + encoded, NSSAI->sst, encoded);

      ENCODE_U8(buffer + encoded, NSSAI->sd[0], encoded);
      ENCODE_U8(buffer + encoded, NSSAI->sd[1], encoded);
      ENCODE_U8(buffer + encoded, NSSAI->sd[2], encoded);
      break;
    case 0b00000101:  // SST, SD and mapped HPLMN SST
      ENCODE_U8(buffer + encoded, NSSAI->sst, encoded);

      ENCODE_U8(buffer + encoded, NSSAI->sd[0], encoded);
      ENCODE_U8(buffer + encoded, NSSAI->sd[1], encoded);
      ENCODE_U8(buffer + encoded, NSSAI->sd[2], encoded);

      ENCODE_U8(buffer + encoded, NSSAI->hplmn_sst, encoded);
      break;
    case 0b00001000:  // SST, SD, mapped HPLMN SST and mapped HPLMN SD
      ENCODE_U8(buffer + encoded, NSSAI->sst, encoded);

      ENCODE_U8(buffer + encoded, NSSAI->sd[0], encoded);
      ENCODE_U8(buffer + encoded, NSSAI->sd[1], encoded);
      ENCODE_U8(buffer + encoded, NSSAI->sd[2], encoded);

      ENCODE_U8(buffer + encoded, NSSAI->hplmn_sst, encoded);

      ENCODE_U8(buffer + encoded, NSSAI->hplmn_sd[0], encoded);
      ENCODE_U8(buffer + encoded, NSSAI->hplmn_sd[1], encoded);
      ENCODE_U8(buffer + encoded, NSSAI->hplmn_sd[2], encoded);
      break;
    default:  // All other values are reserved
      break;
  }
  return (encoded);
};

int NSSAIMsg::DecodeNSSAIMsg(NSSAIMsg* NSSAI, uint8_t iei, uint8_t* buffer,
                             uint32_t len) {
  int decoded = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, (unsigned char)*buffer);
    NSSAI->iei = *(buffer + decoded);
    decoded++;
  }
  DECODE_U8(buffer + decoded, NSSAI->len, decoded);
  CHECK_LENGTH_DECODER(len - decoded, NSSAI->len);

  switch (NSSAI->len) {
    case 0b00000001:  // SST
      DECODE_U8(buffer + decoded, NSSAI->sst, decoded);
      break;
    case 0b00000010:  // SST and mapped HPLMN SST
      DECODE_U8(buffer + decoded, NSSAI->sst, decoded);
      DECODE_U8(buffer + decoded, NSSAI->hplmn_sst, decoded);
      break;
    case 0b00000100:  // SST and SD
      DECODE_U8(buffer + decoded, NSSAI->sst, decoded);

      DECODE_U8(buffer + decoded, NSSAI->sd[0], decoded);
      DECODE_U8(buffer + decoded, NSSAI->sd[1], decoded);
      DECODE_U8(buffer + decoded, NSSAI->sd[2], decoded);
      break;
    case 0b00000101:  // SST, SD and mapped HPLMN SST
      DECODE_U8(buffer + decoded, NSSAI->sst, decoded);

      DECODE_U8(buffer + decoded, NSSAI->sd[0], decoded);
      DECODE_U8(buffer + decoded, NSSAI->sd[1], decoded);
      DECODE_U8(buffer + decoded, NSSAI->sd[2], decoded);

      DECODE_U8(buffer + decoded, NSSAI->hplmn_sst, decoded);
      break;
    case 0b00001000:  // SST, SD, mapped HPLMN SST and mapped HPLMN SD
      DECODE_U8(buffer + decoded, NSSAI->sst, decoded);

      DECODE_U8(buffer + decoded, NSSAI->sd[0], decoded);
      DECODE_U8(buffer + decoded, NSSAI->sd[1], decoded);
      DECODE_U8(buffer + decoded, NSSAI->sd[2], decoded);

      DECODE_U8(buffer + decoded, NSSAI->hplmn_sst, decoded);

      DECODE_U8(buffer + decoded, NSSAI->hplmn_sd[0], decoded);
      DECODE_U8(buffer + decoded, NSSAI->hplmn_sd[1], decoded);
      DECODE_U8(buffer + decoded, NSSAI->hplmn_sd[2], decoded);
      break;
    default:  // All other values are reserved
      break;
  }

  return decoded;
}

NSSAIMsgList::NSSAIMsgList() {}

NSSAIMsgList::~NSSAIMsgList() {}

int NSSAIMsgList::EncodeNSSAIMsgList(NSSAIMsgList* NSSAI_list, uint8_t iei,
                                     uint8_t* buffer, uint32_t len) {
  uint8_t encoded = 0;

  if (iei > 0) {
    CHECK_IEI_ENCODER(iei, (unsigned char)NSSAI_list->iei);
    ENCODE_U8(buffer, iei, encoded);
  }

  ENCODE_U8(buffer + encoded, NSSAI_list->len, encoded);

  encoded +=
      nssai.EncodeNSSAIMsg(&nssai, 0, (buffer + encoded), (len - encoded));

  return (encoded);
};
}  // namespace magma5g
