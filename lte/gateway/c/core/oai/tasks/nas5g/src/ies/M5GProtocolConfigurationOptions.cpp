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

#include <iostream>
#include <sstream>
#include <cstdint>
#include <cstring>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GProtocolConfigurationOptions.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"

namespace magma5g {
ProtocolConfigurationOptions::ProtocolConfigurationOptions() {}
ProtocolConfigurationOptions::~ProtocolConfigurationOptions() {}

int decode_bstring(bstring* bstr, const uint16_t pdulen,
                   const uint8_t* const buffer, const uint32_t buflen) {
  if (buflen < pdulen) {
    return TLV_BUFFER_TOO_SHORT;
  }

  if ((bstr) && (buffer)) {
    *bstr = blk2bstr(buffer, pdulen);
    return pdulen;
  } else {
    return TLV_BUFFER_TOO_SHORT;
  }
}

int encode_bstring(const_bstring const str, uint8_t* const buffer,
                   const uint32_t buflen) {
  if (str) {
    if (blength(str) > 0) {
      CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, blength(str), buflen);
      memcpy((void*)buffer, (void*)str->data, blength(str));
      return blength(str);
    } else {
      return 0;
    }
  } else {
    return 0;
  }
}

int ProtocolConfigurationOptions::m5g_decode_protocol_configuration_options(
    protocol_configuration_options_t* pco, const uint8_t* const buffer,
    const uint32_t len) {
  int decoded = 0;
  int decode_result = 0;

  if (((*(buffer + decoded) >> 7) & 0x1) != 1) {
    return TLV_VALUE_DOESNT_MATCH;
  }

  /*
   * Bits 7 to 4 of octet 3 are spare, read as 0
   */
  if (((*(buffer + decoded) & 0x78) >> 3) != 0) {
    return TLV_VALUE_DOESNT_MATCH;
  }

  pco->configuration_protocol = (*(buffer + decoded) >> 1) & 0x7;
  decoded++;
  pco->num_protocol_or_container_id = 0;

  while (3 <= ((int32_t)len - (int32_t)decoded)) {
    DECODE_U16(
        buffer + decoded,
        pco->protocol_or_container_ids[pco->num_protocol_or_container_id].id,
        decoded);
    DECODE_U8(buffer + decoded,
              pco->protocol_or_container_ids[pco->num_protocol_or_container_id]
                  .length,
              decoded);

    if (0 < pco->protocol_or_container_ids[pco->num_protocol_or_container_id]
                .length) {
      if ((decode_result = decode_bstring(
               &pco->protocol_or_container_ids
                    [pco->num_protocol_or_container_id]
                        .contents,
               pco->protocol_or_container_ids[pco->num_protocol_or_container_id]
                   .length,
               buffer + decoded, len - decoded)) < 0) {
        return decode_result;
      } else {
        decoded += decode_result;
      }
    } else {
      pco->protocol_or_container_ids[pco->num_protocol_or_container_id]
          .contents = NULL;
    }
    pco->num_protocol_or_container_id += 1;
  }

  return decoded;
}

// Decode Protocol Configuration Options IE
int ProtocolConfigurationOptions::DecodeProtocolConfigurationOptions(
    ProtocolConfigurationOptions* protocolconfigurationoptions, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  int decoded = 0;
  int decoded2 = 0;
  uint16_t ielen = 0;

  if (iei) {
    decoded++;
  }

  DECODE_U16(buffer + decoded, ielen, decoded);
  CHECK_LENGTH_DECODER(len - decoded, ielen);

  decoded2 = m5g_decode_protocol_configuration_options(
      &(protocolconfigurationoptions->pco), buffer + decoded, (uint32_t)ielen);
  if (decoded2 < 0) return decoded2;
  if (decoded2 != ielen) return -1;

  return decoded + decoded2;
};

int ProtocolConfigurationOptions::m5g_encode_protocol_configuration_options(
    const protocol_configuration_options_t* const pco, uint8_t* buffer,
    const uint32_t len) {
  uint8_t num_protocol_or_container_id = 0;
  uint32_t encoded = 0;
  int encode_result = 0;

  *(buffer + encoded) = 0x00 | (1 << 7) | (pco->configuration_protocol & 0x7);
  encoded++;

  while (num_protocol_or_container_id < pco->num_protocol_or_container_id) {
    ENCODE_U16(buffer + encoded,
               pco->protocol_or_container_ids[num_protocol_or_container_id].id,
               encoded);

    *(buffer + encoded) =
        pco->protocol_or_container_ids[num_protocol_or_container_id].length;
    encoded++;

    if ((encode_result = encode_bstring(
             pco->protocol_or_container_ids[num_protocol_or_container_id]
                 .contents,
             buffer + encoded, len - encoded)) < 0)
      return encode_result;
    else
      encoded += encode_result;

    num_protocol_or_container_id += 1;
  }
  return encoded;
}

int ProtocolConfigurationOptions::EncodeProtocolConfigurationOptions(
    ProtocolConfigurationOptions* protocolconfigurationoptions, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  uint8_t* lenPtr = NULL;
  uint16_t pco_len = 0;
  uint32_t encoded = 0;

  if (iei) {
    *buffer = REQUEST_EXTENDED_PROTOCOL_CONFIGURATION_OPTIONS_TYPE;
    encoded++;
  }

  lenPtr = (buffer + encoded);

  pco_len = m5g_encode_protocol_configuration_options(
      &(protocolconfigurationoptions->pco), buffer + encoded + sizeof(uint16_t),
      len - encoded);

  if (iei) {
    ENCODE_U16(lenPtr, pco_len, encoded);
    encoded += pco_len;
  } else {
    ENCODE_U16(lenPtr, pco_len - 1, encoded);
    encoded += pco_len - 1;
  }

  return encoded;
}

}  // namespace magma5g
