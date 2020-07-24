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

#ifndef PROTOCOL_CONFIGURATION_OPTIONS_H_
#define PROTOCOL_CONFIGURATION_OPTIONS_H_
#include <stdint.h>
#include "3gpp_24.008.h"

#define PROTOCOL_CONFIGURATION_OPTIONS_MINIMUM_LENGTH PCO_MIN_LENGTH
#define PROTOCOL_CONFIGURATION_OPTIONS_MAXIMUM_LENGTH PCO_MAX_LENGTH

typedef protocol_configuration_options_t oai::ProtocolConfigurationOptions;

int encode_ProtocolConfigurationOptions(
    oai::ProtocolConfigurationOptions* protocolconfigurationoptions,
    uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_ProtocolConfigurationOptions(
    oai::ProtocolConfigurationOptions* protocolconfigurationoptions,
    const uint8_t iei, const uint8_t* const buffer, const uint32_t len);

void dump_ProtocolConfigurationOptions_xml(
    oai::ProtocolConfigurationOptions* protocolconfigurationoptions,
    uint8_t iei);

#endif /* PROTOCOL CONFIGURATION OPTIONS_H_ */
