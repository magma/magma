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

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#include "TLVEncoder.h"
#include "TLVDecoder.h"
#include "oai::ProtocolConfigurationOptions.h"

int decode_ProtocolConfigurationOptions(
    oai::ProtocolConfigurationOptions* protocolconfigurationoptions,
    const uint8_t iei, const uint8_t* const buffer, const uint32_t len) {
  uint32_t decoded = 0;
  uint8_t ielen    = 0;
  int rv           = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);

  rv = decode_protocol_configuration_options(
      protocolconfigurationoptions, buffer + decoded, len - decoded);

  if (rv < 0) {
    return rv;
  }
  decoded += (uint32_t) rv;

#if NAS_DEBUG
  dump_ProtocolConfigurationOptions_xml(protocolconfigurationoptions, iei);
#endif
  return decoded;
}

int encode_ProtocolConfigurationOptions(
    oai::ProtocolConfigurationOptions* protocolconfigurationoptions,
    uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, PROTOCOL_CONFIGURATION_OPTIONS_MINIMUM_LENGTH, len);
#if NAS_DEBUG
  dump_ProtocolConfigurationOptions_xml(protocolconfigurationoptions, iei);
#endif

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;

  encoded += encode_protocol_configuration_options(
      protocolconfigurationoptions, buffer + encoded, len - encoded);

  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);
  return encoded;
}

void dump_ProtocolConfigurationOptions_xml(
    oai::ProtocolConfigurationOptions* protocolconfigurationoptions,
    uint8_t iei) {
  int i;

  OAILOG_DEBUG(LOG_NAS, "<Protocol Configuration Options>\n");

  if (iei > 0)
    /*
     * Don't display IEI if = 0
     */
    OAILOG_DEBUG(LOG_NAS, "    <IEI>0x%X</IEI>\n", iei);

  OAILOG_DEBUG(
      LOG_NAS, "    <Configuration protol>%u</Configuration protol>\n",
      protocolconfigurationoptions->configuration_protocol);
  i = 0;

  while (i < protocolconfigurationoptions->num_protocol_or_container_id) {
    OAILOG_DEBUG(
        LOG_NAS, "        <Protocol ID>%u</Protocol ID>\n",
        protocolconfigurationoptions->protocol_or_container_ids[i].id);
    OAILOG_DEBUG(
        LOG_NAS, "        <Length of protocol ID>%u</Length of protocol ID>\n",
        protocolconfigurationoptions->protocol_or_container_ids[i].length);
    bstring b = dump_bstring_xml(
        protocolconfigurationoptions->protocol_or_container_ids[i].contents);
    OAILOG_DEBUG(LOG_NAS, "        %s", bdata(b));
    bdestroy(b);
    i++;
  }

  OAILOG_DEBUG(LOG_NAS, "</Protocol Configuration Options>\n");
}
