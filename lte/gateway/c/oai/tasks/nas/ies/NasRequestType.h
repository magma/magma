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

#ifndef NAS_REQUEST_TYPE_SEEN
#define NAS_REQUEST_TYPE_SEEN

#include <stdint.h>

#define REQUEST_TYPE_MINIMUM_LENGTH 1
#define REQUEST_TYPE_MAXIMUM_LENGTH 1

#define REQUEST_TYPE_INITIAL_REQUEST 0b001
#define REQUEST_TYPE_HANDOVER 0b010
#define REQUEST_TYPE_UNUSED 0b011
#define REQUEST_TYPE_EMERGENCY 0b100

typedef uint8_t request_type_t;

int encode_request_type(
    request_type_t* requesttype, uint8_t iei, uint8_t* buffer, uint32_t len);

uint8_t encode_u8_request_type(request_type_t* requesttype);

int decode_request_type(
    request_type_t* requesttype, uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_u8_request_type(
    request_type_t* requesttype, uint8_t iei, uint8_t value, uint32_t len);

#endif /* NAS_REQUEST_TYPE_SEEN */
