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
#include "AccessPointName.h"

int decode_access_point_name(
    AccessPointName* accesspointname, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
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

  if ((decode_result = decode_bstring(
           accesspointname, ielen, buffer + decoded, len - decoded)) < 0)
    return decode_result;
  else
    decoded += decode_result;

  return decoded;
}

int encode_access_point_name(
    AccessPointName accesspointname, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint8_t* lenPtr                                       = NULL;
  uint32_t encoded                                      = 0;
  int encode_result                                     = 0;
  uint32_t length_index                                 = 0;
  uint32_t index                                        = 0;
  uint32_t index_copy                                   = 0;
  uint8_t apn_encoded[ACCESS_POINT_NAME_MAXIMUM_LENGTH] = {0};

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, ACCESS_POINT_NAME_MINIMUM_LENGTH, len);

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;
  index        = 0;  // index on original APN string
  length_index = 0;  // marker where to write partial length
  index_copy   = 1;

  while ((accesspointname->data[index] != 0) &&
         (index < accesspointname->slen)) {
    if (accesspointname->data[index] == '.') {
      apn_encoded[length_index] = index_copy - length_index - 1;
      length_index              = index_copy;
      index_copy                = length_index + 1;
    } else {
      apn_encoded[index_copy] = accesspointname->data[index];
      index_copy++;
    }

    index++;
  }

  apn_encoded[length_index] = index_copy - length_index - 1;
  bstring bapn              = blk2bstr(apn_encoded, index_copy);

  if ((encode_result = encode_bstring(bapn, buffer + encoded, len - encoded)) <
      0) {
    bdestroy(bapn);
    return encode_result;
  } else {
    encoded += encode_result;
  }
  bdestroy(bapn);
  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);
  return encoded;
}

void dump_access_point_name_xml(AccessPointName accesspointname, uint8_t iei) {
  OAILOG_DEBUG(LOG_NAS, "<Access Point Name>\n");

  if (iei > 0)
    /*
     * Don't display IEI if = 0
     */
    OAILOG_DEBUG(LOG_NAS, "    <IEI>0x%X</IEI>\n", iei);
  bstring b = dump_bstring_xml(accesspointname);
  OAILOG_DEBUG(LOG_NAS, "%s</Access Point Name>\n", bdata(b));
  bdestroy(b);
}
