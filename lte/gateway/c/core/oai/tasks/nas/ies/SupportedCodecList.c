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

#include <stdint.h>

#include "TLVEncoder.h"
#include "TLVDecoder.h"
#include "SupportedCodecList.h"
#include "log.h"

int decode_supported_codec_list(
    SupportedCodecList* supportedcodeclist, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded   = 0;
  uint8_t ielen = 0, indx = 0;
  uint8_t total_coded_len = 0, decoded_len;

  if (iei > 0) {
    CHECK_IEI_DECODER(SUPPORTED_CODED_LIST_IE, *buffer);
    decoded++;
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);

  for (indx = 0;
       ((indx < SUPPORTED_CODEC_LIST_NUMBER_OF_SYSTEM_INDICATION) &&
        (total_coded_len < (ielen - 2)));
       indx++) {
    decoded_len                                    = decoded;
    supportedcodeclist[indx]->systemidentification = *(buffer + decoded);
    decoded++;
    supportedcodeclist[indx]->lengthofbitmap = *(buffer + decoded);
    decoded++;
    if (supportedcodeclist[indx]->lengthofbitmap == sizeof(uint8_t)) {
      IES_DECODE_U8(buffer, decoded, supportedcodeclist[indx]->codecbitmap);
    } else if (supportedcodeclist[indx]->lengthofbitmap == sizeof(uint16_t)) {
      IES_DECODE_U16(buffer, decoded, supportedcodeclist[indx]->codecbitmap);
    } else {
      OAILOG_ERROR(
          LOG_NAS, "Supported codec list: Invalid Bitmap length :%d\n",
          supportedcodeclist[indx]->lengthofbitmap);
    }
    total_coded_len += decoded - decoded_len;
  }
#if NAS_DEBUG
  dump_supported_codec_list_xml(supportedcodeclist, iei);
#endif
  return decoded;
}

int encode_supported_codec_list(
    SupportedCodecList* supportedcodeclist, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;
  uint8_t indx     = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, SUPPORTED_CODEC_LIST_MINIMUM_LENGTH, len);
#if NAS_DEBUG
  dump_supported_codec_list_xml(supportedcodeclist, iei);
#endif

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;
  for (indx = 0; indx < SUPPORTED_CODEC_LIST_NUMBER_OF_SYSTEM_INDICATION;
       indx++) {
    *(buffer + encoded) = supportedcodeclist[indx]->systemidentification;
    encoded++;
    *(buffer + encoded) = supportedcodeclist[indx]->lengthofbitmap;
    encoded++;
    if (supportedcodeclist[indx]->lengthofbitmap == sizeof(uint8_t)) {
      IES_ENCODE_U8(buffer, encoded, supportedcodeclist[indx]->codecbitmap);
    } else if (supportedcodeclist[indx]->lengthofbitmap == sizeof(uint16_t)) {
      IES_ENCODE_U16(buffer, encoded, supportedcodeclist[indx]->codecbitmap);
    } else {
      OAILOG_ERROR(
          LOG_NAS, "Encode supported codec list: Invalid Bitmal length: %d \n",
          supportedcodeclist[indx]->lengthofbitmap);
    }
  }

  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);
  return encoded;
}

void dump_supported_codec_list_xml(
    SupportedCodecList* supportedcodeclist, uint8_t iei) {
  OAILOG_DEBUG(LOG_NAS, "<Supported Codec List>\n");
  uint8_t indx = 0;

  if (iei > 0)
    /*
     * Don't display IEI if = 0
     */
    OAILOG_DEBUG(LOG_NAS, "    <IEI>0x%X</IEI>\n", iei);

  for (indx = 0; indx < SUPPORTED_CODEC_LIST_NUMBER_OF_SYSTEM_INDICATION;
       indx++) {
    OAILOG_DEBUG(
        LOG_NAS, "<System identification>%u</System identification>\n",
        supportedcodeclist[indx]->systemidentification);
    OAILOG_DEBUG(
        LOG_NAS, "<Length of bitmap>%u</Length of bitmap>\n",
        supportedcodeclist[indx]->lengthofbitmap);
    OAILOG_DEBUG(
        LOG_NAS, "<Codec bitmap>%u</Codec bitmap>\n",
        supportedcodeclist[indx]->codecbitmap);
  }
  OAILOG_DEBUG(LOG_NAS, "</Supported Codec List>\n");
}
