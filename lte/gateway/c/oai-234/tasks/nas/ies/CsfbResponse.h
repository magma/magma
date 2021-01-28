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

#ifndef CSFB_RESPONSE_SEEN
#define CSFB_RESPONSE_SEEN
#include <stdint.h>

#define CSFB_RESPONSE_MINIMUM_LENGTH 1
#define CSFB_RESPONSE_MAXIMUM_LENGTH 1

typedef uint8_t csfb_response_t;

int encode_csfb_response(
    csfb_response_t* csfbresponse, uint8_t iei, uint8_t* buffer, uint32_t len);

uint8_t encode_u8_csfb_response(csfb_response_t* csfbresponse);

int decode_csfb_response(
    csfb_response_t* csfbresponse, uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_u8_csfb_response(
    csfb_response_t* csfbresponse, uint8_t iei, uint8_t value, uint32_t len);

/*
 *  CSFB response value:reference 24301-e40:Table 9.9.3.5
 * Bits
 * 3	2	1
 * 0	0	0		CS fallback rejected by the UE
 * 0	0	1		CS fallback accepted by the UE
 */

typedef enum { CSFB_REJECTED_BY_UE, CSFB_ACCEPTED_BY_UE } Csfb_response;

#endif /* CSFB_RESPONSE_SEEN */
