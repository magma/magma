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

#ifndef FILE_ATTACH_REQUEST_SEEN
#define FILE_ATTACH_REQUEST_SEEN

#include <stdint.h>

#include "SecurityHeaderType.h"
#include "MessageType.h"
#include "EpsAttachType.h"
#include "NasKeySetIdentifier.h"
#include "EpsMobileIdentity.h"
#include "UeNetworkCapability.h"
#include "EsmMessageContainer.h"
#include "TrackingAreaIdentity.h"
#include "AdditionalUpdateType.h"
#include "GutiType.h"
#include "UeAdditionalSecurityCapability.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"
#include "3gpp_24.301.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define ATTACH_REQUEST_MINIMUM_LENGTH                                          \
  (EPS_ATTACH_TYPE_MINIMUM_LENGTH + NAS_KEY_SET_IDENTIFIER_MINIMUM_LENGTH +    \
   EPS_MOBILE_IDENTITY_MINIMUM_LENGTH + UE_NETWORK_CAPABILITY_MINIMUM_LENGTH + \
   ESM_MESSAGE_CONTAINER_MINIMUM_LENGTH)

/* Maximum length macro. Formed by maximum length of each field */
#define ATTACH_REQUEST_MAXIMUM_LENGTH                                          \
  (EPS_ATTACH_TYPE_MAXIMUM_LENGTH + NAS_KEY_SET_IDENTIFIER_MAXIMUM_LENGTH +    \
   EPS_MOBILE_IDENTITY_MAXIMUM_LENGTH + UE_NETWORK_CAPABILITY_MAXIMUM_LENGTH + \
   ESM_MESSAGE_CONTAINER_MAXIMUM_LENGTH + P_TMSI_SIGNATURE_MAXIMUM_LENGTH +    \
   EPS_MOBILE_IDENTITY_MAXIMUM_LENGTH +                                        \
   TRACKING_AREA_IDENTITY_MAXIMUM_LENGTH + DRX_PARAMETER_MAXIMUM_LENGTH +      \
   MS_NETWORK_CAPABILITY_MAXIMUM_LENGTH +                                      \
   MS_NETWORK_FEATURE_SUPPORT_MAXIMUM_LENGTH +                                 \
   LOCATION_AREA_IDENTIFICATION_MAXIMUM_LENGTH + TMSI_STATUS_MAXIMUM_LENGTH +  \
   MOBILE_STATION_CLASSMARK_2_MAXIMUM_LENGTH +                                 \
   MOBILE_STATION_CLASSMARK_3_MAXIMUM_LENGTH +                                 \
   SUPPORTED_CODEC_LIST_MAXIMUM_LENGTH +                                       \
   ADDITIONAL_UPDATE_TYPE_MAXIMUM_LENGTH + GUTI_TYPE_MAXIMUM_LENGTH +          \
   UE_ADDITIONAL_SECURITY_CAPABILITY_MAXIMUM_LENGTH)

/* If an optional value is present and should be encoded, the corresponding
 * Bit mask should be set to 1.
 */
#define ATTACH_REQUEST_OLD_PTMSI_SIGNATURE_PRESENT (1 << 0)
#define ATTACH_REQUEST_ADDITIONAL_GUTI_PRESENT (1 << 1)
#define ATTACH_REQUEST_LAST_VISITED_REGISTERED_TAI_PRESENT (1 << 2)
#define ATTACH_REQUEST_DRX_PARAMETER_PRESENT (1 << 3)
#define ATTACH_REQUEST_MS_NETWORK_CAPABILITY_PRESENT (1 << 4)
#define ATTACH_REQUEST_OLD_LOCATION_AREA_IDENTIFICATION_PRESENT (1 << 5)
#define ATTACH_REQUEST_TMSI_STATUS_PRESENT (1 << 6)
#define ATTACH_REQUEST_MOBILE_STATION_CLASSMARK_2_PRESENT (1 << 7)
#define ATTACH_REQUEST_MOBILE_STATION_CLASSMARK_3_PRESENT (1 << 8)
#define ATTACH_REQUEST_SUPPORTED_CODECS_PRESENT (1 << 9)
#define ATTACH_REQUEST_ADDITIONAL_UPDATE_TYPE_PRESENT (1 << 10)
#define ATTACH_REQUEST_OLD_GUTI_TYPE_PRESENT (1 << 11)
#define ATTACH_REQUEST_VOICE_DOMAIN_PREFERENCE_AND_UE_USAGE_SETTING_PRESENT    \
  (1 << 12)
#define ATTACH_REQUEST_MS_NETWORK_FEATURE_SUPPORT_PRESENT (1 << 13)
#define ATTACH_REQUEST_NETWORK_RESOURCE_IDENTIFIER_CONTAINER_PRESENT (1 << 14)
#define ATTACH_REQUEST_UE_ADDITIONAL_SECURITY_CAPABILITY_PRESENT (1 << 15)

typedef enum attach_request_iei_tag {
  ATTACH_REQUEST_OLD_PTMSI_SIGNATURE_IEI         = GMM_PTMSI_SIGNATURE_IEI,
  ATTACH_REQUEST_ADDITIONAL_GUTI_IEI             = 0x50, /* 0x50 = 80  */
  ATTACH_REQUEST_LAST_VISITED_REGISTERED_TAI_IEI = 0x52, /* 0x52 = 82  */
  ATTACH_REQUEST_DRX_PARAMETER_IEI               = GMM_DRX_PARAMETER_IEI,
  ATTACH_REQUEST_MS_NETWORK_FEATURE_SUPPORT_IEI =
      C_MS_NETWORK_FEATURE_SUPPORT_IEI,
  ATTACH_REQUEST_MS_NETWORK_CAPABILITY_IEI = GMM_MS_NETWORK_CAPABILITY_IEI,
  ATTACH_REQUEST_OLD_LOCATION_AREA_IDENTIFICATION_IEI =
      C_LOCATION_AREA_IDENTIFICATION_IEI,
  ATTACH_REQUEST_TMSI_STATUS_IEI = GMM_TMSI_STATUS_IEI,
  ATTACH_REQUEST_MOBILE_STATION_CLASSMARK_2_IEI =
      C_MOBILE_STATION_CLASSMARK_2_IEI,
  ATTACH_REQUEST_MOBILE_STATION_CLASSMARK_3_IEI =
      C_MOBILE_STATION_CLASSMARK_3_IEI,
  ATTACH_REQUEST_SUPPORTED_CODECS_IEI       = CC_SUPPORTED_CODEC_LIST_IE,
  ATTACH_REQUEST_ADDITIONAL_UPDATE_TYPE_IEI = 0xF0, /* 0xF0 = 240 */
  ATTACH_REQUEST_OLD_GUTI_TYPE_IEI          = 0xE0, /* 0xE0 = 224 */
  ATTACH_REQUEST_VOICE_DOMAIN_PREFERENCE_AND_UE_USAGE_SETTING_IEI =
      GMM_VOICE_DOMAIN_PREFERENCE_AND_UE_USAGE_SETTING_IEI,
  ATTACH_REQUEST_NETWORK_RESOURCE_IDENTIFIER_CONTAINER_IEI = 0x10,
  ATTACH_REQUEST_DEVICE_PROPERTIES_IEI                     = 0xD0,
  ATTACH_REQUEST_DEVICE_PROPERTIES_LOW_PRIO_IEI            = 0xD1,
  ATTACH_REQUEST_UE_ADDITIONAL_SECURITY_CAPABILITY_IEI     = 0x6F
} attach_request_iei;

/*
 * Message name: Attach request
 * Description: This message is sent by the UE to the network in order to
 * perform an attach procedure. See table 8.2.4.1. Significance: dual Direction:
 * UE to network
 */

typedef struct attach_request_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  security_header_type_t securityheadertype : 4;
  message_type_t messagetype;
  eps_attach_type_t epsattachtype;
  NasKeySetIdentifier naskeysetidentifier;
  eps_mobile_identity_t oldgutiorimsi;
  ue_network_capability_t uenetworkcapability;
  EsmMessageContainer esmmessagecontainer;
  /* Optional fields */
  uint32_t presencemask;
  p_tmsi_signature_t oldptmsisignature;
  eps_mobile_identity_t additionalguti;
  tai_t lastvisitedregisteredtai;
  drx_parameter_t drxparameter;
  ms_network_capability_t msnetworkcapability;
  location_area_identification_t oldlocationareaidentification;
  tmsi_status_t tmsistatus;
  mobile_station_classmark2_t mobilestationclassmark2;
  mobile_station_classmark3_t mobilestationclassmark3;
  supported_codec_list_t supportedcodecs;
  additional_update_type_t additionalupdatetype;
  guti_type_t oldgutitype;
  voice_domain_preference_and_ue_usage_setting_t
      voicedomainpreferenceandueusagesetting;
  ms_network_feature_support_t msnetworkfeaturesupport;
  network_resource_identifier_container_t networkresourceidentifiercontainer;
  ue_additional_security_capability_t ueadditionalsecuritycapability;
} attach_request_msg;

int decode_attach_request(
    attach_request_msg* attachrequest, uint8_t* buffer, uint32_t len);

int encode_attach_request(
    attach_request_msg* attachrequest, uint8_t* buffer, uint32_t len);

#endif /* ! defined(FILE_ATTACH_REQUEST_SEEN) */
