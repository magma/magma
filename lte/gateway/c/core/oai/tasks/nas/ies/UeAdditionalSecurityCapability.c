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

#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

#include "bstrlib.h"

#include "TLVDecoder.h"
#include "TLVEncoder.h"
#include "UeAdditionalSecurityCapability.h"

//------------------------------------------------------------------------------
int decode_ue_additional_security_capability(
    ue_additional_security_capability_t* uasc, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded   = 0;
  uint8_t ielen = 0;
  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }
  DECODE_U8(buffer + decoded, ielen, decoded);
  memset(uasc, 0, sizeof(ue_additional_security_capability_t));
  OAILOG_TRACE(
      LOG_NAS_EMM, "decode_ue_additional_security_capability len = %d\n",
      ielen);
  CHECK_LENGTH_DECODER(len - decoded, ielen);

  uasc->_5g_ea = (*(buffer + decoded++));
  uasc->_5g_ea = uasc->_5g_ea << 8;
  uasc->_5g_ea |= (*(buffer + decoded++));

  uasc->_5g_ia = (*(buffer + decoded++));
  uasc->_5g_ia = uasc->_5g_ia << 8;
  uasc->_5g_ia |= (*(buffer + decoded++));

  OAILOG_TRACE(
      LOG_NAS_EMM, "ue_additional_security_capability decoded=%u\n", decoded);

  if ((ielen + 2) != decoded) {
    decoded = ielen + 1 + (iei > 0 ? 1 : 0) /* Size of header for this IE */;
    OAILOG_TRACE(
        LOG_NAS_EMM, "ue_additional_security_capability then decoded=%u\n",
        decoded);
  }
  return decoded;
}

//------------------------------------------------------------------------------
int encode_ue_additional_security_capability(
    ue_additional_security_capability_t* uasc, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, UE_ADDITIONAL_SECURITY_CAPABILITY_MINIMUM_LENGTH, len);

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;
  *(buffer + encoded) = uasc->_5g_ea >> 8;
  encoded++;
  *(buffer + encoded) = (uasc->_5g_ea & 0x00FF);
  encoded++;
  *(buffer + encoded) = uasc->_5g_ia >> 8;
  encoded++;
  *(buffer + encoded) = (uasc->_5g_ia & 0x00FF);
  encoded++;

  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);
  return encoded;
}
