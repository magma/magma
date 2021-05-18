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
#pragma once

#include "Ngap_Cause.h"
#include "Ngap_NGAP-PDU.h"
#include "Ngap_ProtocolIE-Field.h"
#include "bstrlib.h"

#ifdef __cplusplus
extern "C" {
#endif
#include "common_defs.h"
#include "ngap_amf_encoder.h"
#include "ngap_amf_decoder.h"
#include "ngap_amf_nas_procedures.h"
#ifdef __cplusplus
}
#endif

#define NGAP_SETUP_FAILURE_FIND_PROTOCOLIE_BY_ID(                              \
    IE_TYPE, ie, container, IE_ID)                                             \
  do {                                                                         \
    IE_TYPE** ptr;                                                             \
    ie = NULL;                                                                 \
    for (ptr = container->protocolIEs.list.array;                              \
         ptr < &container->protocolIEs.list                                    \
                    .array[container->protocolIEs.list.count];                 \
         ptr++) {                                                              \
      if ((*ptr)->id == IE_ID) {                                               \
        ie = *ptr;                                                             \
        break;                                                                 \
      }                                                                        \
    }                                                                          \
  } while (0)

#define MACRO_GNB_ID_TO_CELL_IDENTITY(mACRO, cELL_iD, bITsTRING)               \
  do {                                                                         \
    (bITsTRING)->buf         = (uint8_t*) calloc(5, sizeof(uint8_t));          \
    (bITsTRING)->buf[0]      = ((mACRO) >> 20);                                \
    (bITsTRING)->buf[1]      = (mACRO) >> 12;                                  \
    (bITsTRING)->buf[2]      = (mACRO) >> 4;                                   \
    (bITsTRING)->buf[3]      = (((mACRO) &0x0f) << 4) | ((cELL_iD) >> 4);      \
    (bITsTRING)->buf[4]      = ((cELL_iD) &0x0f) << 4;                         \
    (bITsTRING)->size        = 5;                                              \
    (bITsTRING)->bits_unused = 4;                                              \
  } while (0)

#define MCC_HUNDREDS(vALUE) ((vALUE) / 100)
/* When MNC is only composed of 2 digits, set the hundreds unit to 0xf */
#define MNC_HUNDREDS(vALUE, mNCdIGITlENGTH)                                    \
  (mNCdIGITlENGTH == 2 ? 15 : (vALUE) / 100)
#define MCC_MNC_DECIMAL(vALUE) (((vALUE) / 10) % 10)
#define MCC_MNC_DIGIT(vALUE) ((vALUE) % 10)

/* Convert an integer on 24 bits to the given bUFFER */
#define INT24_TO_BUFFER(x, buf)                                                \
  do {                                                                         \
    (buf)[0] = (x) >> 16;                                                      \
    (buf)[1] = (x) >> 8;                                                       \
    (buf)[2] = (x);                                                            \
  } while (0)

#define MCC_MNC_TO_TBCD(mCC, mNC, mNCdIGITlENGTH, tBCDsTRING)                  \
  do {                                                                         \
    char _buf[3];                                                              \
    _buf[0] = (MCC_MNC_DECIMAL(mCC) << 4) | MCC_HUNDREDS(mCC);                 \
    _buf[1] = (MNC_HUNDREDS(mNC, mNCdIGITlENGTH) << 4) | MCC_MNC_DIGIT(mCC);   \
    _buf[2] = (MCC_MNC_DIGIT(mNC) << 4) | MCC_MNC_DECIMAL(mNC);                \
    OCTET_STRING_fromBuf(tBCDsTRING, _buf, 3);                                 \
  } while (0)

#define INT24_TO_OCTET_STRING(x, aSN)                                          \
  do {                                                                         \
    (aSN)->buf = (uint8_t*) calloc(3, sizeof(uint8_t));                        \
    INT24_TO_BUFFER(x, ((aSN)->buf));                                          \
    (aSN)->size = 3;                                                           \
  } while (0)

#define MCC_MNC_TO_PLMNID(mCC, mNC, mNCdIGITlENGTH, oCTETsTRING)               \
  do {                                                                         \
    (oCTETsTRING)->buf    = (uint8_t*) calloc(3, sizeof(uint8_t));             \
    (oCTETsTRING)->buf[0] = (MCC_MNC_DECIMAL(mCC) << 4) | MCC_HUNDREDS(mCC);   \
    (oCTETsTRING)->buf[1] =                                                    \
        (MNC_HUNDREDS(mNC, mNCdIGITlENGTH) << 4) | MCC_MNC_DIGIT(mCC);         \
    (oCTETsTRING)->buf[2] = (MCC_MNC_DIGIT(mNC) << 4) | MCC_MNC_DECIMAL(mNC);  \
    (oCTETsTRING)->size   = 3;                                                 \
  } while (0)

// Base test function
int ngap_ng_setup_failure_stream(
    const Ngap_Cause_PR cause_type, const long cause_value, bstring& stream);

int ngap_ng_setup_failure_pdu(
    const Ngap_Cause_PR cause_type, const long cause_value,
    Ngap_NGAP_PDU_t& encode_pdu);

bool ng_setup_failure_decode(const_bstring const raw, Ngap_NGAP_PDU_t* pdu);

bool ngap_initiate_ue_message(bstring& stream);

bool generator_ngap_pdusession_resource_setup_req(bstring& stream);

bool generator_itti_ngap_pdusession_resource_setup_req(bstring& stream);
