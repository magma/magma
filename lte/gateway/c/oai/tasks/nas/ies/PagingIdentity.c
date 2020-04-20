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
#include "PagingIdentity.h"

//------------------------------------------------------------------------------
int decode_paging_identity(
    paging_identity_t* pagingidentity, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  *pagingidentity = *buffer & 0x1;
  decoded++;
  return decoded;
}

//------------------------------------------------------------------------------
int encode_paging_identity(
    paging_identity_t* pagingidentity, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, PAGING_IDENTITY_MINIMUM_LENGTH, len);

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  *(buffer + encoded) = 0x00 | (*pagingidentity & 0x1);
  encoded++;
  return encoded;
}
