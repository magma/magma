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

#ifndef EPS_UPDATE_RESULT_SEEN
#define EPS_UPDATE_RESULT_SEEN

#include <stdint.h>

#define EPS_UPDATE_RESULT_MINIMUM_LENGTH 1
#define EPS_UPDATE_RESULT_MAXIMUM_LENGTH 1

#define EPS_UPDATE_RESULT_TA_UPDATED 0
#define EPS_UPDATE_RESULT_COMBINED_TA_LA_UPDATED 1
#define EPS_UPDATE_RESULT_TA_UPDATED_AND_ISR_ACTIVATED 4
#define EPS_UPDATE_RESULT_COMBINED_TA_LA_UPDATED_AND_ISR_ACTIVATED 5

typedef uint8_t eps_update_result_t;

int encode_eps_update_result(
    eps_update_result_t* epsupdateresult, uint8_t iei, uint8_t* buffer,
    uint32_t len);

uint8_t encode_u8_eps_update_result(eps_update_result_t* epsupdateresult);

int decode_eps_update_result(
    eps_update_result_t* epsupdateresult, uint8_t iei, uint8_t* buffer,
    uint32_t len);

int decode_u8_eps_update_result(
    eps_update_result_t* epsupdateresult, uint8_t iei, uint8_t value,
    uint32_t len);

#endif /* EPS UPDATE RESULT_SEEN */
