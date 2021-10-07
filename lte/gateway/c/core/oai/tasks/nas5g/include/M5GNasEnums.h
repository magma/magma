/**
 * Copyright 2021 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#pragma once

namespace magma5g {

enum class M5GIei : uint8_t {
  M5GMM_CAUSE                   = 0x58,
  REQUEST_TYPE                  = 0x80,
  PDU_SESSION_IDENTITY_2        = 0x12,
  OLD_PDU_SESSION_IDENTITY_2    = 0x59,
  S_NSSA                        = 0x22,
  DNN                           = 0x25,
  ADDITIONAL_INFORMATION        = 0x24,
  MA_PDU_SESSION_INFORMATION    = 0xA0,
  RELEASE_ASSISTANCE_INDICATION = 0xF0
};

enum class M5GMmCause : uint8_t {
  ILLEGAL_UE                                          = 0b00000011,
  PEI_NOT_ACCEPTED                                    = 0b00000101,
  ILLEGAL_ME                                          = 0b00000110,
  FIVEG_SERVICES_NOT_ALLOWED                          = 0b00000111,
  UE_IDENTITY_CANNOT_BE_DERIVED_FROM_NETWORK          = 0b00001001,
  IMPLICITLY_DEREGISTERED                             = 0b00001010,
  PLMN_NOT_ALLOWED                                    = 0b00001011,
  TA_NOT_ALLOWED                                      = 0b00001100,
  ROAMING_NOT_ALLOWED_IN_TA                           = 0b00001101,
  NO_SUITIBLE_CELLS_IN_TA                             = 0b00001111,
  MAC_FAILURE                                         = 0b00010100,
  SYNCH_FAILURE                                       = 0b00010101,
  CONGESTION                                          = 0b00010110,
  UE_SECURITY_CAP_MISMATCH                            = 0b00010111,
  SEC_MODE_REJECTED_UNSPECIFIED                       = 0b00011000,
  NON_5G_AUTHENTICATION_UNACCEPTABLE                  = 0b00011010,
  N1_MODE_NOT_ALLOWED                                 = 0b00011011,
  RESTRICTED_SERVICE_AREA                             = 0b00011100,
  LADN_NOT_AVAILABLE                                  = 0b00101011,
  MAX_PDU_SESSIONS_REACHED                            = 0b01000001,
  INSUFFICIENT_RESOURCES_FOR_SLICE_AND_DNN            = 0b01000011,
  INSUFFICIENT_RESOURCES_FOR_SLICE                    = 0b01000101,
  NGKSI_ALREADY_IN_USE                                = 0b01000111,
  NON_3GPP_ACCESS_TO_CN_NOT_ALLOWED                   = 0b01001000,
  SERVING_NETWORK_NOT_AUTHORIZED                      = 0b01001001,
  PAYLOAD_NOT_FORWARDED                               = 0b01011010,
  DNN_NOT_SUPPORTED_OR_NOT_SUBSCRIBED                 = 0b01011011,
  INSUFFICIENT_USER_PLANE_RESOURCES                   = 0b01011100,
  SEMANTICALLY_INCORRECT_MESSAGE                      = 0b01011111,
  INVALID_MANDATORY_INFORMATION                       = 0b01100000,
  MESSAGE_TYPE_NON_EXISTENT_OR_NOT_IMPLEMENTED        = 0b01100001,
  MESSAGE_TYPE_NOT_COMPATIBLE_WITH_PROTOCOL_STATE     = 0b01100010,
  INFORMATION_ELEMENT_NON_EXISTENT_OR_NOT_IMPLEMENTED = 0b01100011,
  CONDITIONAL_IE_ERROR                                = 0b01100100,
  MESSAGE_NOT_COMPATIBLE_WITH_PROTOCOL_STATE          = 0b01100101,
  UNSPECIFIED_PROTOCOL_ERROR                          = 0b01101111,
};

enum class M5GRequestType : uint8_t {
  INITIAL_REQUEST                = 0b001,
  EXISTING_PDU_SESSION           = 0b010,
  INITIAL_EMERGENCY_REQUEST      = 0b011,
  EXISTING_EMERGENCY_PDU_SESSION = 0b100,
  MODIFICATION_REQUEST           = 0b101,
};

}  // namespace magma5g
