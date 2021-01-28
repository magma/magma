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

#ifndef EPS_QUALITY_OF_SERVICE_SEEN
#define EPS_QUALITY_OF_SERVICE_SEEN

#include <stdint.h>

#include "common_types.h"

#define EPS_QUALITY_OF_SERVICE_MINIMUM_LENGTH 2
#define EPS_QUALITY_OF_SERVICE_MAXIMUM_LENGTH 10

typedef struct {
  uint8_t maxBitRateForUL;
  uint8_t maxBitRateForDL;
  uint8_t guarBitRateForUL;
  uint8_t guarBitRateForDL;
} EpsQoSBitRates;

typedef struct {
  uint8_t bitRatesPresent : 1;
  uint8_t bitRatesExtPresent : 1;
  uint8_t qci;
  EpsQoSBitRates bitRates;
  EpsQoSBitRates bitRatesExt;
} EpsQualityOfService;

int encode_eps_quality_of_service(
    EpsQualityOfService* epsqualityofservice, uint8_t iei, uint8_t* buffer,
    uint32_t len);

int decode_eps_quality_of_service(
    EpsQualityOfService* epsqualityofservice, uint8_t iei, uint8_t* buffer,
    uint32_t len);

int eps_qos_bit_rate_value(uint8_t br);
int eps_qos_bit_rate_ext_value(uint8_t br);
int qos_params_to_eps_qos(
    const qci_t qci, const bitrate_t mbr_dl, const bitrate_t mbr_ul,
    const bitrate_t gbr_dl, const bitrate_t gbr_ul,
    EpsQualityOfService* const eps_qos, bool is_default_bearer);

#endif /* EPS_QUALITY_OF_SERVICE_SEEN */
