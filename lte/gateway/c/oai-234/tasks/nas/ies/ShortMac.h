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

#ifndef SHORT_MAC_SEEN
#define SHORT_MAC_SEEN

#include <stdint.h>

#define SHORT_MAC_MINIMUM_LENGTH 3
#define SHORT_MAC_MAXIMUM_LENGTH 3

typedef uint16_t short_mac_t;

#define SHORT_MAC_FMT "%" PRIX16

int encode_short_mac(
    short_mac_t* shortmac, uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_short_mac(
    short_mac_t* shortmac, uint8_t iei, uint8_t* buffer, uint32_t len);

#endif /* SHORT MAC_SEEN */
