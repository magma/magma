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
#include "M5GExtendedProtocolDiscriminator.h"
#include "M5GSpareHalfOctet.h"
#include "M5GSecurityHeaderType.h"
#include "M5GMessageType.h"
#include "M5GNASSecurityAlgorithms.h"
#include "M5GNASKeySetIdentifier.h"
#include "M5GUESecurityCapability.h"
#include "M5GIMEISVRequest.h"

using namespace std;
namespace magma5g {
// SecurityModeCommand Message Class
class SecurityModeCommandMsg {
 public:
  ExtendedProtocolDiscriminatorMsg extended_protocol_discriminator;
  SpareHalfOctetMsg spare_half_octet;
  SecurityHeaderTypeMsg sec_header_type;
  MessageTypeMsg message_type;
  NASSecurityAlgorithmsMsg nas_sec_algorithms;
  NASKeySetIdentifierMsg nas_key_set_identifier;
  UESecurityCapabilityMsg ue_sec_capability;
  ImeisvRequestMsg imeisv_request;
#define SECURITY_MODE_COMMAND_MINIMUM_LENGTH 8

  SecurityModeCommandMsg();
  ~SecurityModeCommandMsg();
  int DecodeSecurityModeCommandMsg(
      SecurityModeCommandMsg* sec_mode_command, uint8_t* buffer, uint32_t len);
  int EncodeSecurityModeCommandMsg(
      SecurityModeCommandMsg* sec_mode_command, uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g

/******************************************************************************
   SPEC TS-24501 Table 8.2.25.1.1 SECURITY MODE COMMAND message content
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
|   |Security mode command   |Message type 9.7        |    M   |  V   |  1    |
|   |message identity        |                        |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|   |Selected NAS security   |NAS security algorithms |    M   |  V   |  1    |
|   |algorithms              |9.11.3.34               |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|   |ngKSI                   |NAS key set identifier  |    M   |  V   |  1/2  |
|   |                        |9.11.3.32               |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|   |Spare half octet        |Spare half octet 9.5    |    M   |  V   |  1/2  |
|---|------------------------|------------------------|--------|------|-------|
|   |Replayed UE security    |UE security capability  |    M   |  LV  |  3-9  |
|   |capabilities            |9.11.3.54               |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|E- |IMEISV request          |IMEISV request 9.11.3.28|    O   |  TV  |  1    |
|---|------------------------|------------------------|--------|------|-------|
******************************************************************************/
