/*
   Copyright 2020 The Magma Authors.
   This source code is licensed under the BSD-style license found in the
   LICENSE file in the root directory of this source tree.
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

#pragma once
#include <sstream>
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GExtendedProtocolDiscriminator.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSpareHalfOctet.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSecurityHeaderType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GMessageType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPayloadContainerType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPayloadContainer.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GRequestType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GDNN.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GNSSAI.hpp"

namespace magma5g {
// ULNASTransport Message Class
class ULNASTransportMsg {
 public:
  ExtendedProtocolDiscriminatorMsg extended_protocol_discriminator;
  SpareHalfOctetMsg spare_half_octet;
  SecurityHeaderTypeMsg sec_header_type;
  MessageTypeMsg message_type;
  PayloadContainerTypeMsg payload_container_type;
  PayloadContainerMsg payload_container;
  // optinal parameters
  RequestType request_type;
#define UL_NAS_TRANSPORT_MINIMUM_LENGTH 7
  DNNMsg dnn;
  NSSAIMsg nssai;

  ULNASTransportMsg();
  ~ULNASTransportMsg();
  int DecodeULNASTransportMsg(ULNASTransportMsg* ul_nas_transport,
                              uint8_t* buffer, uint32_t len);
  int EncodeULNASTransportMsg(ULNASTransportMsg* ul_nas_transport,
                              uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g

/******************************************************************************
   SPEC TS-24501  Table 8.2.10.1.1: UL NAS TRANSPORT message content
-------------------------------------------------------------------------------
|IEI|   Information Element  |    Type/Reference      |Presence|Format|Length |
|---|------------------------|------------------------|--------|------|-------|
|   |Extended protocol descr-|Extended Protocol descr-|    M   |  V   |  1    |
|   |-iminator               |-iminator 9.2           |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|   |Security header type    |Security header type 9.3|    M   |  V   |  1/2  |
|---|------------------------|------------------------|--------|------|-------|
|   |Spare half octet        |Spare half octet 9.5    |    M   |  V   |  1/2  |
|---|------------------------|------------------------|--------|------|-------|
|   |UL NAS TRANSPORT message|Message type 9.7        |    M   |  V   |  1    |
|   |        identity        |                        |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|   |Payload container type  |Payload container type  |    M   |  V   |  1/2  |
|   |                        |9.11.3.40               |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|   |Spare half octet        |Spare half octet 9.5    |    M   |  V   |  1/2  |
|---|------------------------|------------------------|--------|------|-------|
|   |Payload container       |Payload container       |    M   | LV-E |3-65537|
|   |                        |9.11.3.39               |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
******************************************************************************/
