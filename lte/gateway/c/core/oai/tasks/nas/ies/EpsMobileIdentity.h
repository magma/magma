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

#ifndef EPS_MOBILE_IDENTITY_SEEN
#define EPS_MOBILE_IDENTITY_SEEN

#include <stdint.h>

#define EPS_MOBILE_IDENTITY_MINIMUM_LENGTH 3
#define EPS_MOBILE_IDENTITY_MAXIMUM_LENGTH 13

typedef struct guti_eps_mobile_identity_s {
  uint8_t spare : 4;
#define EPS_MOBILE_IDENTITY_EVEN 0
#define EPS_MOBILE_IDENTITY_ODD 1
  uint8_t oddeven : 1;
  uint8_t typeofidentity : 3;
  uint8_t mcc_digit2 : 4;
  uint8_t mcc_digit1 : 4;
  uint8_t mnc_digit3 : 4;
  uint8_t mcc_digit3 : 4;
  uint8_t mnc_digit2 : 4;
  uint8_t mnc_digit1 : 4;
  uint16_t mme_group_id;
  uint8_t mme_code;
  uint32_t m_tmsi;
} guti_eps_mobile_identity_t;

typedef struct imsi_eps_mobile_identity_s {
  uint8_t identity_digit1 : 4;
  uint8_t oddeven : 1;
  uint8_t typeofidentity : 3;
  uint8_t identity_digit2 : 4;
  uint8_t identity_digit3 : 4;
  uint8_t identity_digit4 : 4;
  uint8_t identity_digit5 : 4;
  uint8_t identity_digit6 : 4;
  uint8_t identity_digit7 : 4;
  uint8_t identity_digit8 : 4;
  uint8_t identity_digit9 : 4;
  uint8_t identity_digit10 : 4;
  uint8_t identity_digit11 : 4;
  uint8_t identity_digit12 : 4;
  uint8_t identity_digit13 : 4;
  uint8_t identity_digit14 : 4;
  uint8_t identity_digit15 : 4;
  // because of union put this extra attribute at the end
  uint8_t num_digits;
} imsi_eps_mobile_identity_t;

typedef imsi_eps_mobile_identity_t imei_eps_mobile_identity_t;

typedef union eps_mobile_identity_s {
#define EPS_MOBILE_IDENTITY_IMSI 0b001
#define EPS_MOBILE_IDENTITY_GUTI 0b110
#define EPS_MOBILE_IDENTITY_IMEI 0b011
  imsi_eps_mobile_identity_t imsi;
  guti_eps_mobile_identity_t guti;
  imei_eps_mobile_identity_t imei;
} eps_mobile_identity_t;

#define EPS_MOBILE_IDENTITY_XML_STR "eps_mobile_identity"
#define TYPE_OF_IDENTITY_ATTR_XML_STR "type_of_identity"

int encode_eps_mobile_identity(
    eps_mobile_identity_t* epsmobileidentity, uint8_t iei, uint8_t* buffer,
    uint32_t len);

int decode_eps_mobile_identity(
    eps_mobile_identity_t* epsmobileidentity, uint8_t iei, uint8_t* buffer,
    uint32_t len);

#endif /* EPS_MOBILE_IDENTITY_SEEN */
