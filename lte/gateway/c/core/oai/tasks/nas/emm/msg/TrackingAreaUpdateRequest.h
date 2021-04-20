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

#ifndef FILE_TRACKING_AREA_UPDATE_REQUEST_SEEN
#define FILE_TRACKING_AREA_UPDATE_REQUEST_SEEN
#include <stdint.h>

#include "SecurityHeaderType.h"
#include "MessageType.h"
#include "EpsUpdateType.h"
#include "NasKeySetIdentifier.h"
#include "EpsMobileIdentity.h"
#include "Nonce.h"
#include "UeNetworkCapability.h"
#include "TrackingAreaIdentity.h"
#include "UeRadioCapabilityInformationUpdateNeeded.h"
#include "EpsBearerContextStatus.h"
#include "AdditionalUpdateType.h"
#include "GutiType.h"
#include "UeAdditionalSecurityCapability.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"
#include "3gpp_24.301.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define TRACKING_AREA_UPDATE_REQUEST_MINIMUM_LENGTH                            \
  (EPS_UPDATE_TYPE_MINIMUM_LENGTH + NAS_KEY_SET_IDENTIFIER_MINIMUM_LENGTH +    \
   EPS_MOBILE_IDENTITY_MINIMUM_LENGTH)

/* Maximum length macro. Formed by maximum length of each field */
#define TRACKING_AREA_UPDATE_REQUEST_MAXIMUM_LENGTH                            \
  (EPS_UPDATE_TYPE_MAXIMUM_LENGTH + NAS_KEY_SET_IDENTIFIER_MAXIMUM_LENGTH +    \
   EPS_MOBILE_IDENTITY_MAXIMUM_LENGTH +                                        \
   NAS_KEY_SET_IDENTIFIER_MAXIMUM_LENGTH +                                     \
   CIPHERING_KEY_SEQUENCE_NUMBER_MAXIMUM_LENGTH +                              \
   P_TMSI_SIGNATURE_MAXIMUM_LENGTH + EPS_MOBILE_IDENTITY_MAXIMUM_LENGTH +      \
   NONCE_MAXIMUM_LENGTH + UE_NETWORK_CAPABILITY_MAXIMUM_LENGTH +               \
   TRACKING_AREA_IDENTITY_MAXIMUM_LENGTH + DRX_PARAMETER_MAXIMUM_LENGTH +      \
   UE_RADIO_CAPABILITY_INFORMATION_UPDATE_NEEDED_MAXIMUM_LENGTH +              \
   EPS_BEARER_CONTEXT_STATUS_MAXIMUM_LENGTH +                                  \
   MS_NETWORK_CAPABILITY_MAXIMUM_LENGTH +                                      \
   LOCATION_AREA_IDENTIFICATION_MAXIMUM_LENGTH + TMSI_STATUS_MAXIMUM_LENGTH +  \
   MOBILE_STATION_CLASSMARK_2_MAXIMUM_LENGTH +                                 \
   MOBILE_STATION_CLASSMARK_3_MAXIMUM_LENGTH +                                 \
   SUPPORTED_CODEC_LIST_MAXIMUM_LENGTH +                                       \
   ADDITIONAL_UPDATE_TYPE_MAXIMUM_LENGTH +                                     \
   UE_ADDITIONAL_SECURITY_CAPABILITY_MAXIMUM_LENGTH)

/* If an optional value is present and should be encoded, the corresponding
 * Bit mask should be set to 1.
 */
#define TRACKING_AREA_UPDATE_REQUEST_NONCURRENT_NATIVE_NAS_KEY_SET_IDENTIFIER_PRESENT \
  (1 << 0)
#define TRACKING_AREA_UPDATE_REQUEST_GPRS_CIPHERING_KEY_SEQUENCE_NUMBER_PRESENT \
  (1 << 1)
#define TRACKING_AREA_UPDATE_REQUEST_OLD_PTMSI_SIGNATURE_PRESENT (1 << 2)
#define TRACKING_AREA_UPDATE_REQUEST_ADDITIONAL_GUTI_PRESENT (1 << 3)
#define TRACKING_AREA_UPDATE_REQUEST_NONCEUE_PRESENT (1 << 4)
#define TRACKING_AREA_UPDATE_REQUEST_UE_NETWORK_CAPABILITY_PRESENT (1 << 5)
#define TRACKING_AREA_UPDATE_REQUEST_LAST_VISITED_REGISTERED_TAI_PRESENT       \
  (1 << 6)
#define TRACKING_AREA_UPDATE_REQUEST_DRX_PARAMETER_PRESENT (1 << 7)
#define TRACKING_AREA_UPDATE_REQUEST_UE_RADIO_CAPABILITY_INFORMATION_UPDATE_NEEDED_PRESENT \
  (1 << 8)
#define TRACKING_AREA_UPDATE_REQUEST_EPS_BEARER_CONTEXT_STATUS_PRESENT (1 << 9)
#define TRACKING_AREA_UPDATE_REQUEST_MS_NETWORK_CAPABILITY_PRESENT (1 << 10)
#define TRACKING_AREA_UPDATE_REQUEST_OLD_LOCATION_AREA_IDENTIFICATION_PRESENT  \
  (1 << 11)
#define TRACKING_AREA_UPDATE_REQUEST_TMSI_STATUS_PRESENT (1 << 12)
#define TRACKING_AREA_UPDATE_REQUEST_MOBILE_STATION_CLASSMARK_2_PRESENT        \
  (1 << 13)
#define TRACKING_AREA_UPDATE_REQUEST_MOBILE_STATION_CLASSMARK_3_PRESENT        \
  (1 << 14)
#define TRACKING_AREA_UPDATE_REQUEST_SUPPORTED_CODECS_PRESENT (1 << 15)
#define TRACKING_AREA_UPDATE_REQUEST_ADDITIONAL_UPDATE_TYPE_PRESENT (1 << 16)
#define TRACKING_AREA_UPDATE_REQUEST_OLD_GUTI_TYPE_PRESENT (1 << 17)
#define TRACKING_AREA_UPDATE_REQUEST_VOICE_DOMAIN_PREFERENCE_PRESENT (1 << 18)
#define TRACKING_AREA_UPDATE_REQUEST_MS_NETWORK_FEATURE_SUPPORT_PRESENT        \
  (1 << 19)
#define TRACKING_AREA_UPDATE_REQUEST_UE_ADDITIONAL_SECURITY_CAPABILITY_PRESENT \
  (1 << 20)

typedef enum tracking_area_update_request_iei_tag {
  TRACKING_AREA_UPDATE_REQUEST_NONCURRENT_NATIVE_NAS_KEY_SET_IDENTIFIER_IEI =
      0xB0, /* 0xB0 = 176 */
  TRACKING_AREA_UPDATE_REQUEST_GPRS_CIPHERING_KEY_SEQUENCE_NUMBER_IEI =
      C_CIPHERING_KEY_SEQUENCE_NUMBER_IEI,
  TRACKING_AREA_UPDATE_REQUEST_OLD_PTMSI_SIGNATURE_IEI =
      GMM_PTMSI_SIGNATURE_IEI,
  TRACKING_AREA_UPDATE_REQUEST_ADDITIONAL_GUTI_IEI       = 0x50, /* 0x50 = 80 */
  TRACKING_AREA_UPDATE_REQUEST_NONCEUE_IEI               = 0x55, /* 0x55 = 85 */
  TRACKING_AREA_UPDATE_REQUEST_UE_NETWORK_CAPABILITY_IEI = 0x58, /* 0x58 = 88 */
  TRACKING_AREA_UPDATE_REQUEST_LAST_VISITED_REGISTERED_TAI_IEI =
      0x52, /* 0x52 = 82 */
  TRACKING_AREA_UPDATE_REQUEST_DRX_PARAMETER_IEI = GMM_DRX_PARAMETER_IEI,
  TRACKING_AREA_UPDATE_REQUEST_UE_RADIO_CAPABILITY_INFORMATION_UPDATE_NEEDED_IEI =
      0xA0, /* 0xA0 = 160 */
  TRACKING_AREA_UPDATE_REQUEST_EPS_BEARER_CONTEXT_STATUS_IEI =
      0x57, /* 0x57 = 87 */
  TRACKING_AREA_UPDATE_REQUEST_MS_NETWORK_CAPABILITY_IEI =
      GMM_MS_NETWORK_CAPABILITY_IEI,
  TRACKING_AREA_UPDATE_REQUEST_OLD_LOCATION_AREA_IDENTIFICATION_IEI =
      C_LOCATION_AREA_IDENTIFICATION_IEI,
  TRACKING_AREA_UPDATE_REQUEST_TMSI_STATUS_IEI = GMM_TMSI_STATUS_IEI,
  TRACKING_AREA_UPDATE_REQUEST_MOBILE_STATION_CLASSMARK_2_IEI =
      C_MOBILE_STATION_CLASSMARK_2_IEI,
  TRACKING_AREA_UPDATE_REQUEST_MOBILE_STATION_CLASSMARK_3_IEI =
      C_MOBILE_STATION_CLASSMARK_3_IEI,
  TRACKING_AREA_UPDATE_REQUEST_SUPPORTED_CODECS_IEI = 0x40, /* 0x40 = 64 */
  TRACKING_AREA_UPDATE_REQUEST_ADDITIONAL_UPDATE_TYPE_IEI =
      0xF0,                                              /* 0xF0 = 240 */
  TRACKING_AREA_UPDATE_REQUEST_OLD_GUTI_TYPE_IEI = 0xE0, /* 0xE0 = 224 */
  TRACKING_AREA_UPDATE_REQUEST_VOICE_DOMAIN_PREFERENCE_IEI =
      0x5D, /* 0x5D = 93 */
  TRACKING_AREA_UPDATE_REQUEST_MS_NETWORK_FEATURE_SUPPORT_IEI =
      C_MS_NETWORK_FEATURE_SUPPORT_IEI,
  TRACKING_AREA_UPDATE_REQUEST_UE_ADDITIONAL_SECURITY_CAPABILITY_IEI = 0x6F
} tracking_area_update_request_iei;

/*
 * Message name: Tracking area update request
 * Description: The purposes of sending the tracking area update request by the
 * UE to the network are described in subclause 5.5.3.1. See table 8.2.29.1.
 * Significance: dual
 * Direction: UE to network
 */

typedef struct tracking_area_update_request_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  security_header_type_t securityheadertype : 4;
  message_type_t messagetype;
  EpsUpdateType epsupdatetype;
  NasKeySetIdentifier naskeysetidentifier;
  eps_mobile_identity_t oldguti;
  /* Optional fields */
  uint32_t presencemask;
  NasKeySetIdentifier noncurrentnativenaskeysetidentifier;
  ciphering_key_sequence_number_t gprscipheringkeysequencenumber;
  p_tmsi_signature_t oldptmsisignature;
  eps_mobile_identity_t additionalguti;
  nonce_t nonceue;
  ue_network_capability_t uenetworkcapability;
  tai_t lastvisitedregisteredtai;
  drx_parameter_t drxparameter;
  ue_radio_capability_information_update_needed_t
      ueradiocapabilityinformationupdateneeded;
  eps_bearer_context_status_t epsbearercontextstatus;
  ms_network_capability_t msnetworkcapability;
  location_area_identification_t oldlocationareaidentification;
  tmsi_status_t tmsistatus;
  mobile_station_classmark2_t mobilestationclassmark2;
  mobile_station_classmark3_t mobilestationclassmark3;
  supported_codec_list_t supportedcodecs;
  additional_update_type_t additionalupdatetype;
  voice_domain_preference_and_ue_usage_setting_t
      voicedomainpreferenceandueusagesetting;
  guti_type_t oldgutitype;
  ms_network_feature_support_t msnetworkfeaturesupport;
  ue_additional_security_capability_t ueadditionalsecuritycapability;
} tracking_area_update_request_msg;

int decode_tracking_area_update_request(
    tracking_area_update_request_msg* trackingareaupdaterequest,
    uint8_t* buffer, uint32_t len);

int encode_tracking_area_update_request(
    tracking_area_update_request_msg* trackingareaupdaterequest,
    uint8_t* buffer, uint32_t len);

#endif /* ! defined(FILE_TRACKING_AREA_UPDATE_REQUEST_SEEN) */
