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
#include "NrUeSecurityCapability.h"

//------------------------------------------------------------------------------
int decode_nr_ue_security_capability(
    nr_ue_security_capability_t* nruesecuritycapability, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  int decoded   = 0;
  uint8_t ielen = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  memset(nruesecuritycapability, 0, sizeof(nr_ue_security_capability_t));
  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  nruesecuritycapability->nea = *(buffer + decoded);
  decoded++;
  nruesecuritycapability->nia = *(buffer + decoded);
  decoded++;

  if (len >= (decoded + 2)) {
    nruesecuritycapability->nr_present = 1;
    nruesecuritycapability->nea        = *(buffer + decoded);
    decoded++;
    nruesecuritycapability->nia = *(buffer + decoded) & 0x7f;
    decoded++;
  }
  return decoded;
}

//------------------------------------------------------------------------------
int encode_nr_ue_security_capability(
    nr_ue_security_capability_t* nruesecuritycapability, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, NR_UE_SECURITY_CAPABILITY_MAXIMUM_LENGTH, len);

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;
  *(buffer + encoded) = nruesecuritycapability->nea;
  encoded++;
  *(buffer + encoded) = nruesecuritycapability->nia;
  encoded++;

  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);
  return encoded;
}
