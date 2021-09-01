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

#ifndef SERVICE_TYPE_SEEN
#define SERVICE_TYPE_SEEN

#include <stdint.h>

#define SERVICE_TYPE_MINIMUM_LENGTH 1
#define SERVICE_TYPE_MAXIMUM_LENGTH 1

typedef uint8_t service_type_t;

int encode_service_type(
    service_type_t* servicetype, uint8_t iei, uint8_t* buffer, uint32_t len);

uint8_t encode_u8_service_type(service_type_t* servicetype);

int decode_service_type(
    service_type_t* servicetype, uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_u8_service_type(
    service_type_t* servicetype, uint8_t iei, uint8_t value, uint32_t len);

/*
  Service type value:reference 24301-af0:Table 9.9.3.27.1
* Bits
* 4	3	2	1
* 0	0	0	0	mobile originating CS fallback or 1xCS fallback
* 0	0	0	1	mobile terminating CS fallback or 1xCS fallback
* 0	0	1	0	mobile originating CS fallback emergency call or
1xCS fallback emergency call
* 0	0	1	1	unused; shall be interpreted as "mobile
originating CS fallback or 1xCS fallback", if received by the network
* 0	1	0	0	unused; shall be interpreted as "mobile
originating CS fallback or 1xCS fallback", if received by the network
* 1	0	0	0	packet services via S1
* 1	0	0	1	unused; shall be interpreted as "packet services
via S1", if received by the network
* 1	0	1	0	unused; shall be interpreted as "packet services
via S1", if received by the network
* 1	0	1	1	unused; shall be interpreted as "packet services
via S1", if received by the network
*/

typedef enum {
  MO_CS_FB = 0x0,
  MT_CS_FB,
  MO_CS_FB_EMRGNCY_CALL,
  MO_CS_FB1,
  MT_CS_FB1,
  PKT_SRV_VIA_S1,
  PKT_SRV_VIA_S1_1,
  PKT_SRV_VIA_S1_2,
  PKT_SRV_VIA_S1_3
} Service_Type;

#endif /* SERVICE_TYPE_SEEN */
