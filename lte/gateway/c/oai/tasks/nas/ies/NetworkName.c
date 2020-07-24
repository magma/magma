/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#include "TLVEncoder.h"
#include "TLVDecoder.h"
#include "NetworkName.h"

int decode_network_name(
    NetworkName* networkname, uint8_t iei, uint8_t* buffer, uint32_t len) {
  int decoded   = 0;
  uint8_t ielen = 0;
  int decode_result;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);

  if (((*buffer >> 7) & 0x1) != 1) {
    errorCodeDecoder = TLV_VALUE_DOESNT_MATCH;
    return TLV_VALUE_DOESNT_MATCH;
  }

  networkname->codingscheme                 = (*(buffer + decoded) >> 5) & 0x7;
  networkname->addci                        = (*(buffer + decoded) >> 4) & 0x1;
  networkname->numberofsparebitsinlastoctet = (*(buffer + decoded) >> 1) & 0x7;

  if ((decode_result = decode_bstring(
           &networkname->textstring, ielen, buffer + decoded, len - decoded)) <
      0)
    return decode_result;
  else
    decoded += decode_result;

#if NAS_DEBUG
  dump_network_name_xml(networkname, iei);
#endif
  return decoded;
}

int encode_network_name(
    NetworkName* networkname, uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;
  int encode_result;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, NETWORK_NAME_MINIMUM_LENGTH, len);
#if NAS_DEBUG
  dump_network_name_xml(networkname, iei);
#endif

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;
  *(buffer + encoded) = 0x00 | (1 << 7) |
                        ((networkname->codingscheme & 0x7) << 4) |
                        ((networkname->addci & 0x1) << 3) |
                        (networkname->numberofsparebitsinlastoctet & 0x7);
  encoded++;

  if ((encode_result = encode_bstring(
           networkname->textstring, buffer + encoded, len - encoded)) < 0)
    return encode_result;
  else
    encoded += encode_result;

  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);
  return encoded;
}

void dump_network_name_xml(NetworkName* networkname, uint8_t iei) {
  OAILOG_DEBUG(LOG_NAS, "<Network Name>\n");

  if (iei > 0)
    /*
     * Don't display IEI if = 0
     */
    OAILOG_DEBUG(LOG_NAS, "    <IEI>0x%X</IEI>\n", iei);

  OAILOG_DEBUG(
      LOG_NAS, "    <Coding scheme>%u</Coding scheme>\n",
      networkname->codingscheme);
  OAILOG_DEBUG(LOG_NAS, "    <Add CI>%u</Add CI>\n", networkname->addci);
  OAILOG_DEBUG(
      LOG_NAS,
      "    <Number of spare bits in last octet>%u</Number of spare bits in "
      "last "
      "octet>\n",
      networkname->numberofsparebitsinlastoctet);
  bstring b = dump_bstring_xml(networkname->textstring);
  OAILOG_DEBUG(LOG_NAS, "%s", bdata(b));
  bdestroy(b);
  OAILOG_DEBUG(LOG_NAS, "</Network Name>\n");
}
