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

#ifndef EPS_NETWORK_FEATURE_SUPPORT_SEEN
#define EPS_NETWORK_FEATURE_SUPPORT_SEEN

#include <stdint.h>

#define EPS_NETWORK_FEATURE_SUPPORT_MINIMUM_LENGTH 3
#define EPS_NETWORK_FEATURE_SUPPORT_MAXIMUM_LENGTH 3

typedef struct eps_network_feature_support_s {
  uint8_t b1;
  uint8_t b2;
} eps_network_feature_support_t;

#define EPS_NETWORK_FEATURE_SUPPORT_IMS_VOICE_OVER_PS_SESSION_IN_S1_NOT_SUPPORTED \
  0x00
#define EPS_NETWORK_FEATURE_SUPPORT_IMS_VOICE_OVER_PS_SESSION_IN_S1_SUPPORTED  \
  0x01
#define EPS_NETWORK_FEATURE_SUPPORT_IMS_VOICE_OVER_PS_SESSION_IN_S1_BITMASK 0x01

#define EPS_NETWORK_FEATURE_SUPPORT_EMERGENCY_BEARER_SERVICES_IN_S1_MODE_NOT_SUPPORTED \
  0x00
#define EPS_NETWORK_FEATURE_SUPPORT_EMERGENCY_BEARER_SERVICES_IN_S1_MODE_SUPPORTED \
  0x02
#define EPS_NETWORK_FEATURE_SUPPORT_EMERGENCY_BEARER_SERVICES_IN_S1_MODE_BITMASK \
  0x02

#define EPS_NETWORK_FEATURE_SUPPORT_LOCATION_SERVICES_VIA_EPC_NOT_SUPPORTED 0x00
#define EPS_NETWORK_FEATURE_SUPPORT_LOCATION_SERVICES_VIA_EPC_SUPPORTED 0x04
#define EPS_NETWORK_FEATURE_SUPPORT_LOCATION_SERVICES_VIA_EPC_BITMASK 0x04

#define EPS_NETWORK_FEATURE_SUPPORT_CS_LCS_NO_INFORMATION_ABOUT_SUPPORT_OF_LS_VIA_CS_DOMAIN_IS_AVAILABLE \
  0x00
#define EPS_NETWORK_FEATURE_SUPPORT_CS_LCS_LOCATION_SERVICES_VIA_CS_DOMAIN_SUPPORTED \
  0x08
#define EPS_NETWORK_FEATURE_SUPPORT_CS_LCS_LOCATION_SERVICES_VIA_CS_DOMAIN_NOT_SUPPORTED \
  0x10
#define EPS_NETWORK_FEATURE_SUPPORT_CS_LCS_BITMASK 0x18

#define EPS_NETWORK_FEATURE_SUPPORT_EXTENDED_SERVICE_REQUEST_NOT_SUPPORTED 0
#define EPS_NETWORK_FEATURE_SUPPORT_EXTENDED_SERVICE_REQUEST_SUPPORTED 0x20
#define EPS_NETWORK_FEATURE_SUPPORT_EXTENDED_SERVICE_REQUEST_BITMASK 0x20
#define EPS_NETWORK_FEATURE_SUPPORT_B2_UP_CIOT_SUPPORTED 0x01
#define EPS_NETWORK_FEATURE_SUPPORT_B2_S1_U_DATA_SUPPORTED 0x02
#define EPS_NETWORK_FEATURE_SUPPORT_B2_HC_CP_CIOT_SUPPORTED 0x04
#define EPS_NETWORK_FEATURE_SUPPORT_B2_EPCO_SUPPORTED 0x08
#define EPS_NETWORK_FEATURE_SUPPORT_B2_RESTRICTEC_SUPPORTED 0x10
#define EPS_NETWORK_FEATURE_SUPPORT_B2_RESTRICTDCNR_SUPPORTED 0x20
#define EPS_NETWORK_FEATURE_SUPPORT_B2_IWK_N26_SUPPORTED 0x40
#define EPS_NETWORK_FEATURE_SUPPORT_B2_15_BEARERS_SUPPORTED 0x80

int encode_eps_network_feature_support(
    eps_network_feature_support_t* epsnetworkfeaturesupport, uint8_t iei,
    uint8_t* buffer, uint32_t len);

int decode_eps_network_feature_support(
    eps_network_feature_support_t* epsnetworkfeaturesupport, uint8_t iei,
    uint8_t* buffer, uint32_t len);

#endif /* EPS_NETWORK_FEATURE_SUPPORT_SEEN */
