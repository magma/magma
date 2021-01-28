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
#include "PlmnList.h"

int decode_plmn_list(
    PlmnList* plmnlist, uint8_t iei, uint8_t* buffer, uint32_t len) {
  int decoded   = 0;
  uint8_t ielen = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  plmnlist->mccdigit2 = (*(buffer + decoded) >> 4) & 0xf;
  plmnlist->mccdigit1 = *(buffer + decoded) & 0xf;
  decoded++;
  plmnlist->mncdigit3 = (*(buffer + decoded) >> 4) & 0xf;
  plmnlist->mccdigit3 = *(buffer + decoded) & 0xf;
  decoded++;
  plmnlist->mncdigit2 = (*(buffer + decoded) >> 4) & 0xf;
  plmnlist->mncdigit1 = *(buffer + decoded) & 0xf;
  decoded++;
#if NAS_DEBUG
  dump_plmn_list_xml(plmnlist, iei);
#endif
  return decoded;
}

int encode_plmn_list(
    PlmnList* plmnlist, uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, PLMN_LIST_MINIMUM_LENGTH, len);
#if NAS_DEBUG
  dump_plmn_list_xml(plmnlist, iei);
#endif

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((plmnlist->mccdigit2 & 0xf) << 4) | (plmnlist->mccdigit1 & 0xf);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((plmnlist->mncdigit3 & 0xf) << 4) | (plmnlist->mccdigit3 & 0xf);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((plmnlist->mncdigit2 & 0xf) << 4) | (plmnlist->mncdigit1 & 0xf);
  encoded++;
  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);
  return encoded;
}

void dump_plmn_list_xml(PlmnList* plmnlist, uint8_t iei) {
  OAILOG_DEBUG(LOG_NAS, "<Plmn List>\n");

  if (iei > 0)
    /*
     * Don't display IEI if = 0
     */
    OAILOG_DEBUG(LOG_NAS, "    <IEI>0x%X</IEI>\n", iei);

  OAILOG_DEBUG(
      LOG_NAS, "    <MCC digit 2>%u</MCC digit 2>\n", plmnlist->mccdigit2);
  OAILOG_DEBUG(
      LOG_NAS, "    <MCC digit 1>%u</MCC digit 1>\n", plmnlist->mccdigit1);
  OAILOG_DEBUG(
      LOG_NAS, "    <MNC digit 3>%u</MNC digit 3>\n", plmnlist->mncdigit3);
  OAILOG_DEBUG(
      LOG_NAS, "    <MCC digit 3>%u</MCC digit 3>\n", plmnlist->mccdigit3);
  OAILOG_DEBUG(
      LOG_NAS, "    <MNC digit 2>%u</MNC digit 2>\n", plmnlist->mncdigit2);
  OAILOG_DEBUG(
      LOG_NAS, "    <MNC digit 1>%u</MNC digit 1>\n", plmnlist->mncdigit1);
  OAILOG_DEBUG(LOG_NAS, "</Plmn List>\n");
}
