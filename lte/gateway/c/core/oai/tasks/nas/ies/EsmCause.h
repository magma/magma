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

#ifndef ESM_CAUSE_SEEN
#define ESM_CAUSE_SEEN

#include <stdint.h>

#define ESM_CAUSE_MINIMUM_LENGTH 1
#define ESM_CAUSE_MAXIMUM_LENGTH 1

// warning coding flaws in ESM, do not use uint8_t yet.
typedef int esm_cause_t;

int encode_esm_cause(
    esm_cause_t* esmcause, uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_esm_cause(
    esm_cause_t* esmcause, uint8_t iei, uint8_t* buffer, uint32_t len);

#endif /* ESM CAUSE_SEEN */
