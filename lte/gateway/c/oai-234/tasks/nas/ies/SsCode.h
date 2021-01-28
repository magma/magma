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

#ifndef SS_CODE_SEEN
#define SS_CODE_SEEN
#include <stdint.h>

#define SS_CODE_MINIMUM_LENGTH 2
#define SS_CODE_MAXIMUM_LENGTH 2

typedef uint8_t ss_code_t;

int encode_ss_code(
    ss_code_t* sscode, uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_ss_code(
    ss_code_t* sscode, uint8_t iei, uint8_t* buffer, uint32_t len);

#endif /* SS CODE_SEEN */
