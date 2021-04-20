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

#ifndef ESM_INFORMATION_TRANSFER_FLAG_SEEN
#define ESM_INFORMATION_TRANSFER_FLAG_SEEN

#include <stdint.h>

#define ESM_INFORMATION_TRANSFER_FLAG_MINIMUM_LENGTH 1
#define ESM_INFORMATION_TRANSFER_FLAG_MAXIMUM_LENGTH 1

typedef uint8_t esm_information_transfer_flag_t;

int encode_esm_information_transfer_flag(
    esm_information_transfer_flag_t* esminformationtransferflag, uint8_t iei,
    uint8_t* buffer, uint32_t len);

uint8_t encode_u8_esm_information_transfer_flag(
    esm_information_transfer_flag_t* esminformationtransferflag);

int decode_esm_information_transfer_flag(
    esm_information_transfer_flag_t* esminformationtransferflag, uint8_t iei,
    uint8_t* buffer, uint32_t len);

int decode_u8_esm_information_transfer_flag(
    esm_information_transfer_flag_t* esminformationtransferflag, uint8_t iei,
    uint8_t value, uint32_t len);

#endif /* ESM INFORMATION TRANSFER FLAG_SEEN */
