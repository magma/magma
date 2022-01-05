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
  M5GMM_CAUSE                             = 0x58,
  REQUEST_TYPE                            = 0x80,
  PDU_SESSION_IDENTITY_2                  = 0x12,
  OLD_PDU_SESSION_IDENTITY_2              = 0x59,
  S_NSSAI                                 = 0x22,
  DNN                                     = 0x25,
  ADDITIONAL_INFORMATION                  = 0x24,
  MA_PDU_SESSION_INFORMATION              = 0xA0,
  PDU_ADDRESS                             = 0x29,
  RELEASE_ASSISTANCE_INDICATION           = 0xF0,
  QOS_FLOW_DESCRIPTIONS                   = 0x79,
  EXTENDED_PROTOCOL_CONFIGURATION_OPTIONS = 0x7B
};

enum class M5GMmCause : uint8_t {
  UNKNOWN_CAUSE,
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

enum class M5GSmCause : uint8_t {
  INVALID_CAUSE,
  OPERATOR_DETERMINED_BARRING                       = 0b00001000,
  INSUFFICIENT_RESOURCES                            = 0b00011010,
  MISSING_OR_UNKNOWN_DNN                            = 0b00011011,
  UNKNOWN_PDU_SESSION_TYPE                          = 0b00011100,
  USER_AUTHENTICATION_OR_AUTHORIZATION_FAILED       = 0b00011101,
  REQUEST_REJECTED_UNSPECIFIED                      = 0b00011111,
  SERVICE_OPTION_NOT_SUPPORTED                      = 0b00100000,
  REQUESTED_SERVICE_OPTION_NOT_SUBSCRIBED           = 0b00100001,
  PTI_ALREADY_IN_USE                                = 0b00100011,
  REGULAR_DEACTIVATION                              = 0b00100100,
  NETWORK_FAILURE                                   = 0b00100110,
  REACTIVATION_REQUESTED                            = 0b00100111,
  SEMANTIC_ERROR_IN_THE_TFT_OPERATION               = 0b00101001,
  SYNTACTICAL_ERROR_IN_THE_TFT_OPERATION            = 0b00101010,
  INVALID_PDU_SESSION_IDENTITY                      = 0b00101011,
  SEMANTIC_ERRORS_IN_PACKET_FILTER                  = 0b00101100,
  SYNTACTICAL_ERROR_IN_PACKET_FILTER                = 0b00101101,
  OUT_OF_LADN_SERVICE_AREA                          = 0b00101110,
  PTI_MISMATCH                                      = 0b00101111,
  PDU_SESSION_TYPE_IPV4_ONLY_ALLOWED                = 0b00110010,
  PDU_SESSION_TYPE_IPV6_ONLY_ALLOWED                = 0b00110011,
  PDU_SESSION_DOES_NOT_EXIST                        = 0b00110110,
  PDU_SESSION_TYPE_IPV4V6_ONLY_ALLOWED              = 0b00111001,
  PDU_SESSION_TYPE_UNSTRUCTURED_ONLY_ALLOWED        = 0b00111010,
  UNSUPPORTED_5QI_VALUE                             = 0b00111011,
  PDU_SESSION_TYPE_ETHERNET_ONLY_ALLOWED            = 0b00111101,
  INSUFFICIENT_RESOURCES_FOR_SPECIFIC_SLICE_AND_DNN = 0b01000011,
  NOT_SUPPORTED_SSC_MODE                            = 0b01000100,
  INSUFFICIENT_RESOURCES_FOR_SPECIFIC_SLICE         = 0b01000101,
  MISSING_OR_UNKNOWN_DNN_IN_A_SLICE                 = 0b01000110,
  INVALID_PTI_VALUE                                 = 0b01010001,
  MAXIMUM_DATA_RATE_PER_UE_FOR_USER_PLANE_INTEGRITY_PROTECTION_IS_TOO_LOW =
      0b01010010,
  SEMANTIC_ERROR_IN_THE_QOS_OPERATION                 = 0b01010011,
  SYNTACTICAL_ERROR_IN_THE_QOS_OPERATION              = 0b01010100,
  INVALID_MAPPED_EPS_BEARER_IDENTITY                  = 0b01010101,
  SEMANTICALLY_INCORRECT_MESSAGE                      = 0b01011111,
  INVALID_MANDATORY_INFORMATION                       = 0b01100000,
  MESSAGE_TYPE_NON_EXISTENT_OR_NOT_IMPLEMENTED        = 0b01100001,
  MESSAGE_TYPE_NOT_COMPATIBLE_WITH_THE_PROTOCOLSTATE  = 0b01100010,
  INFORMATION_ELEMENT_NON_EXISTENT_OR_NOT_IMPLEMENTED = 0b01100011,
  CONDITIONAL_IE_ERROR                                = 0b01100100,
  MESSAGE_NOT_COMPATIBLE_WITH_THE_PROTOCOL_STATE      = 0b01100101,
  PROTOCOL_ERROR_UNSPECIFIED                          = 0b01101111,
};

enum class M5GSessionAmbrUnit : uint8_t {
  VALUE_NOT_USED    = 000000000,
  MULTIPLES_1KBPS   = 0b00000001,
  MULTIPLES_4KBPS   = 0b00000010,
  MULTIPLES_16KBPS  = 0b00000011,
  MULTIPLES_64KBPS  = 0b00000100,
  MULTIPLES_256KBPS = 0b00000101,
  MULTIPLES_1MBPS   = 0b00000110,
  MULTIPLES_4MBPS   = 0b00000111,
  MULTIPLES_16MBPS  = 0b00001000,
  MULTIPLES_64MBPS  = 0b00001001,
  MULTIPLES_256MBPS = 0b00001010,
  MULTIPLES_1GBPS   = 0b00001011,
  MULTIPLES_4GBPS   = 0b00001100,
  MULTIPLES_16GBPS  = 0b00001101,
  MULTIPLES_64GBPS  = 0b00001110,
  MULTIPLES_256GBPS = 0b00001111,
  MULTIPLES_1TBPS   = 0b00010000,
  MULTIPLES_4TBPS   = 0b00010001,
  MULTIPLES_16TBPS  = 0b00010010,
  MULTIPLES_64TBPS  = 0b00010011,
  MULTIPLES_256TBPS = 0b00010100,
  MULTIPLES_1PBPS   = 0b00010101,
  MULTIPLES_4PBPS   = 0b00010110,
  MULTIPLES_16PBPS  = 0b00010111,
  MULTIPLES_64PBPS  = 0b00011000,
  MULTIPLES_256PBPS = 0b00011001,
};

enum class M5GRequestType : uint8_t {
  INITIAL_REQUEST                = 0b001,
  EXISTING_PDU_SESSION           = 0b010,
  INITIAL_EMERGENCY_REQUEST      = 0b011,
  EXISTING_EMERGENCY_PDU_SESSION = 0b100,
  MODIFICATION_REQUEST           = 0b101,
};

enum class M5GPduSessionType : uint8_t {
  IPV4         = 0b001,
  IPV6         = 0b010,
  IPV4V6       = 0b011,
  UNSTRUCTURED = 0b100,
  ETHERNET     = 0b101,
};

}  // namespace magma5g
