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
#include <string.h>

#include "TLVEncoder.h"
#include "TLVDecoder.h"
#include "3gpp_24.301.h"
#include "UeSecurityCapability.h"

//------------------------------------------------------------------------------
int decode_ue_security_capability(
    ue_security_capability_t* uesecuritycapability, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  int decoded   = 0;
  uint8_t ielen = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  memset(uesecuritycapability, 0, sizeof(ue_security_capability_t));
  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  uesecuritycapability->eea = *(buffer + decoded);
  decoded++;
  uesecuritycapability->eia = *(buffer + decoded);
  decoded++;

  if (len >= (decoded + 2)) {
    uesecuritycapability->umts_present = 1;
    uesecuritycapability->uea          = *(buffer + decoded);
    decoded++;
    uesecuritycapability->uia = *(buffer + decoded) & 0x7f;
    decoded++;

    if (len >= (decoded + 1)) {
      uesecuritycapability->gprs_present = 1;
      uesecuritycapability->gea          = *(buffer + decoded) & 0x7f;
      decoded++;
    }
  }
  return decoded;
}

//------------------------------------------------------------------------------
int encode_ue_security_capability(
    ue_security_capability_t* uesecuritycapability, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, UE_SECURITY_CAPABILITY_MAXIMUM_LENGTH, len);

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;
  *(buffer + encoded) = uesecuritycapability->eea;
  encoded++;
  *(buffer + encoded) = uesecuritycapability->eia;
  encoded++;

  // From ETSI TS 124 301 V10.15.0 (2014-10) 9.9.3.36 Security capability:
  // Octets 5, 6, and 7 are optional. If octet 5 is included, then also octet 6
  // shall be included and octet 7 may be included. If a UE did not indicate
  // support of any security algorithm for Gb mode, octet 7 shall not be
  // included. If the UE did not indicate support of any security algorithm for
  // Iu mode and Gb mode, octets 5, 6, and 7 shall not be included. If the UE
  // did not indicate support of any security algorithm for Iu mode but
  // indicated support of a security algorithm for Gb mode, octets 5, 6, and 7
  // shall be included. In this case octets 5 and 6 are filled with the value of
  // zeroes.
  if (uesecuritycapability->umts_present) {
    *(buffer + encoded) = uesecuritycapability->uea;
    encoded++;
    *(buffer + encoded) = 0x00 | (uesecuritycapability->uia & 0x7f);
    encoded++;

    if (uesecuritycapability->gprs_present) {
      *(buffer + encoded) = 0x00 | (uesecuritycapability->gea & 0x7f);
      encoded++;
    }
  } else {
    if (uesecuritycapability->gprs_present) {
      *(buffer + encoded) = 0x00;
      encoded++;
      *(buffer + encoded) = 0x00;
      encoded++;
      *(buffer + encoded) = 0x00 | (uesecuritycapability->gea & 0x7f);
      encoded++;
    }
  }

  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);
  return encoded;
}
