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

/*! \file TLVDecoder.h
  \brief
  \author Philippe MOREL, Sebastien ROUX, Lionel GAUTHIER
  \company Eurecom
*/
#ifndef FILE_TLV_DECODER_SEEN
#define FILE_TLV_DECODER_SEEN

#include <stddef.h>
#include <stdint.h>

#include "common_defs.h"
#include "log.h"
#include "bstrlib.h"

#define IES_DECODE_U8(bUFFER, dECODED, vALUE)                                  \
  DECODE_U8(bUFFER + dECODED, vALUE, dECODED)

#define IES_DECODE_U16(bUFFER, dECODED, vALUE)                                 \
  DECODE_U16(bUFFER + dECODED, vALUE, dECODED)

#define IES_DECODE_U24(bUFFER, dECODED, vALUE)                                 \
  DECODE_U24(bUFFER + dECODED, vALUE, dECODED)

#define IES_DECODE_U32(bUFFER, dECODED, vALUE)                                 \
  DECODE_U32(bUFFER + dECODED, vALUE, dECODED)

extern int errorCodeDecoder;

int decode_bstring(
    bstring* octetstring, const uint16_t pdulen, const uint8_t* const buffer,
    const uint32_t buflen);

bstring dump_bstring_xml(const bstring bstr);

void tlv_decode_perror(void);

#define CHECK_PDU_POINTER_AND_LENGTH_DECODER(bUFFER, mINIMUMlENGTH, lENGTH)    \
  if (bUFFER == NULL) {                                                        \
    OAILOG_WARNING(LOG_NAS, "Got NULL pointer for the payload\n");             \
    errorCodeDecoder = TLV_BUFFER_NULL;                                        \
    return TLV_BUFFER_NULL;                                                    \
  }                                                                            \
  if (lENGTH < mINIMUMlENGTH) {                                                \
    OAILOG_WARNING(                                                            \
        LOG_NAS, "Expecting at least %d bytes, got %d\n", mINIMUMlENGTH,       \
        lENGTH);                                                               \
    errorCodeDecoder = TLV_BUFFER_TOO_SHORT;                                   \
    return TLV_BUFFER_TOO_SHORT;                                               \
  }

#define CHECK_PDU_POINTER_AND_LENGTH_DECODER_FOR_MANDATORY_IES(                \
    bUFFER, mINIMUMlENGTH, lENGTH)                                             \
  if (bUFFER == NULL) {                                                        \
    OAILOG_WARNING(LOG_NAS, "Got NULL pointer for the payload\n");             \
    errorCodeDecoder = TLV_MANDATORY_FIELD_NOT_PRESENT;                        \
    return TLV_MANDATORY_FIELD_NOT_PRESENT;                                    \
  }                                                                            \
  if (lENGTH < mINIMUMlENGTH) {                                                \
    OAILOG_WARNING(                                                            \
        LOG_NAS, "Expecting at least %d bytes, got %d\n", mINIMUMlENGTH,       \
        lENGTH);                                                               \
    errorCodeDecoder = TLV_MANDATORY_FIELD_NOT_PRESENT;                        \
    return TLV_MANDATORY_FIELD_NOT_PRESENT;                                    \
  }

#define CHECK_LENGTH_DECODER(bUFFERlENGTH, lENGTH)                             \
  if (bUFFERlENGTH < lENGTH) {                                                 \
    errorCodeDecoder = TLV_BUFFER_TOO_SHORT;                                   \
    return TLV_BUFFER_TOO_SHORT;                                               \
  }

#define CHECK_MESSAGE_TYPE(mESSAGE_tYPE, bUFFER)                               \
  {                                                                            \
    if (mESSAGE_tYPE != bUFFER) {                                              \
      errorCodeDecoder = TLV_WRONG_MESSAGE_TYPE;                               \
      return errorCodeDecoder;                                                 \
    }                                                                          \
  }

#define CHECK_IEI_DECODER(iEI, bUFFER)                                         \
  if (iEI != bUFFER) {                                                         \
    OAILOG_WARNING(                                                            \
        LOG_NAS,                                                               \
        "IEI is different than the one expected."                              \
        "(Got: 0x%x, expecting: 0x%x\n",                                       \
        bUFFER, iEI);                                                          \
    errorCodeDecoder = TLV_UNEXPECTED_IEI;                                     \
    return TLV_UNEXPECTED_IEI;                                                 \
  }

#endif /* define (FILE_TLV_DECODER_SEEN) */
