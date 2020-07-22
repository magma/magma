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

#include "AdditionalUpdateType.h"

//------------------------------------------------------------------------------
int decode_additional_update_type(
    additional_update_type_t* additionalupdatetype, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  int decoded = 0;

  *additionalupdatetype = *buffer & 0x1;
  decoded++;
  return decoded;
}

//------------------------------------------------------------------------------
int encode_additional_update_type(
    additional_update_type_t* additionalupdatetype, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  return 0;
}
