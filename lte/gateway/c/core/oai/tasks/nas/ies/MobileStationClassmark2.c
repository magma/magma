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
#include "MobileStationClassmark2.h"

int decode_mobile_station_classmark_2(
    MobileStationClassmark2* mobilestationclassmark2, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  int decoded   = 0;
  uint8_t ielen = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  mobilestationclassmark2->revisionlevel     = (*(buffer + decoded) >> 5) & 0x3;
  mobilestationclassmark2->esind             = (*(buffer + decoded) >> 4) & 0x1;
  mobilestationclassmark2->a51               = (*(buffer + decoded) >> 3) & 0x1;
  mobilestationclassmark2->rfpowercapability = *(buffer + decoded) & 0x7;
  decoded++;
  mobilestationclassmark2->pscapability      = (*(buffer + decoded) >> 6) & 0x1;
  mobilestationclassmark2->ssscreenindicator = (*(buffer + decoded) >> 4) & 0x3;
  mobilestationclassmark2->smcapability      = (*(buffer + decoded) >> 3) & 0x1;
  mobilestationclassmark2->vbs               = (*(buffer + decoded) >> 2) & 0x1;
  mobilestationclassmark2->vgcs              = (*(buffer + decoded) >> 1) & 0x1;
  mobilestationclassmark2->fc                = *(buffer + decoded) & 0x1;
  decoded++;
  mobilestationclassmark2->cm3      = (*(buffer + decoded) >> 7) & 0x1;
  mobilestationclassmark2->lcsvacap = (*(buffer + decoded) >> 5) & 0x1;
  mobilestationclassmark2->ucs2     = (*(buffer + decoded) >> 4) & 0x1;
  mobilestationclassmark2->solsa    = (*(buffer + decoded) >> 3) & 0x1;
  mobilestationclassmark2->cmsp     = (*(buffer + decoded) >> 2) & 0x1;
  mobilestationclassmark2->a53      = (*(buffer + decoded) >> 1) & 0x1;
  mobilestationclassmark2->a52      = *(buffer + decoded) & 0x1;
  decoded++;
#if NAS_DEBUG
  dump_mobile_station_classmark_2_xml(mobilestationclassmark2, iei);
#endif
  return decoded;
}

int encode_mobile_station_classmark_2(
    MobileStationClassmark2* mobilestationclassmark2, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, MOBILE_STATION_CLASSMARK_2_MINIMUM_LENGTH, len);
#if NAS_DEBUG
  dump_mobile_station_classmark_2_xml(mobilestationclassmark2, iei);
#endif

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;
  *(buffer + encoded) = 0x00 |
                        ((mobilestationclassmark2->revisionlevel & 0x3) << 5) |
                        ((mobilestationclassmark2->esind & 0x1) << 4) |
                        ((mobilestationclassmark2->a51 & 0x1) << 3) |
                        (mobilestationclassmark2->rfpowercapability & 0x7);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((mobilestationclassmark2->pscapability & 0x1) << 6) |
      ((mobilestationclassmark2->ssscreenindicator & 0x3) << 4) |
      ((mobilestationclassmark2->smcapability & 0x1) << 3) |
      ((mobilestationclassmark2->vbs & 0x1) << 2) |
      ((mobilestationclassmark2->vgcs & 0x1) << 1) |
      (mobilestationclassmark2->fc & 0x1);
  encoded++;
  *(buffer + encoded) = 0x00 | ((mobilestationclassmark2->cm3 & 0x1) << 7) |
                        ((mobilestationclassmark2->lcsvacap & 0x1) << 5) |
                        ((mobilestationclassmark2->ucs2 & 0x1) << 4) |
                        ((mobilestationclassmark2->solsa & 0x1) << 3) |
                        ((mobilestationclassmark2->cmsp & 0x1) << 2) |
                        ((mobilestationclassmark2->a53 & 0x1) << 1) |
                        (mobilestationclassmark2->a52 & 0x1);
  encoded++;
  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);
  return encoded;
}

void dump_mobile_station_classmark_2_xml(
    MobileStationClassmark2* mobilestationclassmark2, uint8_t iei) {
  OAILOG_DEBUG(LOG_NAS, "<Mobile Station Classmark 2>\n");

  if (iei > 0)
    /*
     * Don't display IEI if = 0
     */
    OAILOG_DEBUG(LOG_NAS, "    <IEI>0x%X</IEI>\n", iei);

  OAILOG_DEBUG(
      LOG_NAS, "    <Revision level>%u</Revision level>\n",
      mobilestationclassmark2->revisionlevel);
  OAILOG_DEBUG(
      LOG_NAS, "    <ES IND>%u</ES IND>\n", mobilestationclassmark2->esind);
  OAILOG_DEBUG(LOG_NAS, "    <A51>%u</A51>\n", mobilestationclassmark2->a51);
  OAILOG_DEBUG(
      LOG_NAS, "    <RF power capability>%u</RF power capability>\n",
      mobilestationclassmark2->rfpowercapability);
  OAILOG_DEBUG(
      LOG_NAS, "    <PS capability>%u</PS capability>\n",
      mobilestationclassmark2->pscapability);
  OAILOG_DEBUG(
      LOG_NAS, "    <SS Screen indicator>%u</SS Screen indicator>\n",
      mobilestationclassmark2->ssscreenindicator);
  OAILOG_DEBUG(
      LOG_NAS, "    <SM capability>%u</SM capability>\n",
      mobilestationclassmark2->smcapability);
  OAILOG_DEBUG(LOG_NAS, "    <VBS>%u</VBS>\n", mobilestationclassmark2->vbs);
  OAILOG_DEBUG(LOG_NAS, "    <VGCS>%u</VGCS>\n", mobilestationclassmark2->vgcs);
  OAILOG_DEBUG(LOG_NAS, "    <FC>%u</FC>\n", mobilestationclassmark2->fc);
  OAILOG_DEBUG(LOG_NAS, "    <CM3>%u</CM3>\n", mobilestationclassmark2->cm3);
  OAILOG_DEBUG(
      LOG_NAS, "    <LCSVA CAP>%u</LCSVA CAP>\n",
      mobilestationclassmark2->lcsvacap);
  OAILOG_DEBUG(LOG_NAS, "    <UCS2>%u</UCS2>\n", mobilestationclassmark2->ucs2);
  OAILOG_DEBUG(
      LOG_NAS, "    <SoLSA>%u</SoLSA>\n", mobilestationclassmark2->solsa);
  OAILOG_DEBUG(LOG_NAS, "    <CMSP>%u</CMSP>\n", mobilestationclassmark2->cmsp);
  OAILOG_DEBUG(LOG_NAS, "    <A53>%u</A53>\n", mobilestationclassmark2->a53);
  OAILOG_DEBUG(LOG_NAS, "    <A52>%u</A52>\n", mobilestationclassmark2->a52);
  OAILOG_DEBUG(LOG_NAS, "</Mobile Station Classmark 2>\n");
}
