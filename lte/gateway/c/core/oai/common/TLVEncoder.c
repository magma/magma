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

/*! \file TLVEncoder.c
  \brief
  \author Philippe MOREL, Sebastien ROUX, Lionel GAUTHIER
  \company Eurecom
*/

#include <stdint.h>
#include <string.h>

#include "lte/gateway/c/core/oai/common/TLVEncoder.h"

int errorCodeEncoder = 0;

int encode_bstring(const_bstring const str, uint8_t* const buffer,
                   const uint32_t buflen) {
  if (!str) return 0;
  if (blength(str) <= 0) return 0;
  if (buffer == NULL) return TLV_BUFFER_NULL;
  if (buflen < blength(str)) return TLV_BUFFER_TOO_SHORT;
  memcpy((void*)buffer, (void*)str->data, blength(str));
  return blength(str);
}
