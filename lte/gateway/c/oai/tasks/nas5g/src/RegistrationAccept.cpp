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
#include <sstream>
#include "RegistrationAccept.h"
#include "CommonDefs.h"

namespace magma5g
{
  RegistrationAcceptMsg::RegistrationAcceptMsg()
  {
  };

  RegistrationAcceptMsg::~RegistrationAcceptMsg()
  {
  };

  int RegistrationAcceptMsg::DecodeRegistrationAcceptMsg( RegistrationAcceptMsg *registrationaccept, uint8_t* buffer, uint32_t len)
  {
    uint32_t decoded = 0;

    int decoded_result = 0;

    CHECK_PDU_POINTER_AND_LENGTH_DECODER (buffer, REGISTRATION_ACCEPT_MINIMUM_LENGTH, len);

    if((decoded_result = registrationaccept -> extendedprotocoldiscriminator.DecodeExtendedProtocolDiscriminatorMsg(&registrationaccept->extendedprotocoldiscriminator, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = registrationaccept -> securityheadertype.DecodeSecurityHeaderTypeMsg(&registrationaccept->securityheadertype, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = registrationaccept -> messagetype.DecodeMessageTypeMsg(&registrationaccept->messagetype, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = registrationaccept -> m5gsregistrationresult.DecodeM5GSRegistrationResultMsg(&registrationaccept->m5gsregistrationresult, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;
#if 0
    if((decoded_result = Decode_5gs_mobile_identity (&registrationaccept->_5gsmobileidentity, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decodeplmn_list (&registrationaccept->plmnlist, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decode_5gs_tracking_area_identity_list (&registrationaccept->_5gstrackingareaidentitylist, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decodenssai (&registrationaccept->nssai, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decoderejected_nssai (&registrationaccept->rejectednssai, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decode_5gs_network_feature_support (&registrationaccept->_5gsnetworkfeaturesupport, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decodepdu_session_status (&registrationaccept->pdusessionstatus, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decodepdu_session_reactivation_result (&registrationaccept->pdusessionreactivationresult, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decodepdu_session_reactivation_result_error_cause (&registrationaccept->pdusessionreactivationresulterrorcause, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decodeladn_information (&registrationaccept->ladninformation, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decodemico_indication (&registrationaccept->micoindication, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decodenetwork_slicing_indication (&registrationaccept->networkslicingindication, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decodeservice_area_list (&registrationaccept->servicearealist, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decodegprs_timer3 (&registrationaccept->gprstimer3, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decodegprs_timer2 (&registrationaccept->gprstimer2, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decodeemergency_number_list (&registrationaccept->emergencynumberlist, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decodeextended_emergency_number_list (&registrationaccept->extendedemergencynumberlist, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decodesor_transparent_container (&registrationaccept->sortransparentcontainer, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decodeeap_message (&registrationaccept->eapmessage, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decodenssai_inclusion_mode (&registrationaccept->nssaiinclusionmode, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decodeoperator_defined_access_category_definitions (&registrationaccept->operatordefinedaccesscategorydefinitions, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;

    if((decoded_result = Decode_5gsdrx_parameters (&registrationaccept->_5gsdrxparameters, 0, buffer+decoded,len-decoded)) <0)
      return decoded_result;
    else
      decoded+=decoded_result;
#endif

    return decoded;
  }

  int RegistrationAcceptMsg::EncodeRegistrationAcceptMsg( RegistrationAcceptMsg *registrationaccept, uint8_t* buffer, uint32_t len)
  {
    uint32_t encoded = 0;
    int encodedresult = 0;

    CHECK_PDU_POINTER_AND_LENGTH_ENCODER (buffer, REGISTRATION_ACCEPT_MINIMUM_LENGTH, len);

    if((encodedresult = registrationaccept->extendedprotocoldiscriminator.EncodeExtendedProtocolDiscriminatorMsg (&registrationaccept->extendedprotocoldiscriminator, 0, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = registrationaccept->securityheadertype.EncodeSecurityHeaderTypeMsg (&registrationaccept->securityheadertype, 0, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = registrationaccept->messagetype.EncodeMessageTypeMsg (&registrationaccept->messagetype, 0, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = registrationaccept->m5gsregistrationresult.EncodeM5GSRegistrationResultMsg (&registrationaccept->m5gsregistrationresult, 0, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;
#if 0

    if((encodedresult = registrationaccept->m5gsmobileidentity.EncodeM5GSMobileIdentityMsg (&registrationaccept->m5gsmobileidentity, M5GSMOBILEIDENTITY, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult; 

    if((encodedresult = registrationaccept->plmnlist.EncodePLMNListMsg (&registrationaccept->plmnlist, PLMNLIST, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = registrationaccept->m5gstrackingareaidentitylist.EncodeM5GSTrackingAreaIdentityListMsg (&registrationaccept->m5gstrackingareaidentitylist, TAILIST, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = registrationaccept->nssai.EncodeNSSAIMsg (&registrationaccept->nssai, ALLOWEDNSSAI, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;

    /*if((encodedresult = registrationaccept->rejectednssai.EncodeRejectedNSSAIMsg (&registrationaccept->rejectednssai, REJECTEDNSSAI, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;
*/
    if((encodedresult = registrationaccept->m5gsnetworkfeaturesupport.EncodeM5GSNetworkFeatureSupportMsg (&registrationaccept->m5gsnetworkfeaturesupport, M5GSNETWORKFEATURESUPPORT, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = registrationaccept->gprstimer3.EncodeGPRSTimer3Msg (&registrationaccept->gprstimer3, GPRSTIMER3, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;
     return 0;

    if((encodedresult = registrationaccept->pdusessionstatus.EncodePDUSessionStatusMsg (&registrationaccept->pdusessionstatus, PDUSESSIONSTATUS, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;
/*
    if((encodedresult = registrationaccept->pdusessionreactivationresult.EncodePDUSessionReactivationResultMsg (&registrationaccept->pdusessionreactivationresult, PDUSESSIONREACTIVATIONRESULT, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = registrationaccept->pdusessionreactivationresulterrorcause.EncodePDUSessionReactivationResultErrorCauseMsg (&registrationaccept->pdusessionreactivationresulterrorcause, PDUSESSIONREACTIVATIONRESULTERROR, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;
*/
    if((encodedresult = registrationaccept->ladninformation.EncodeLADNInformationMsg (&registrationaccept->ladninformation, LADNINFORMATION, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = registrationaccept->micoindication.EncodeMICOIndicationMsg (&registrationaccept->micoindication, MICOINDICATION, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = registrationaccept->networkslicingindication.EncodeNetworkSlicingIndicationMsg (&registrationaccept->networkslicingindication, NETWORKSLICINGINDICATION, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;
/*
    if((encodedresult = registrationaccept->servicearealist.EncodeServiceAreaListMsg (&registrationaccept->servicearealist, SERVICEAREALIST, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;
*/
    if((encodedresult = registrationaccept->gprstimer3.EncodeGPRSTimer3Msg (&registrationaccept->gprstimer3, GPRSTIMER3, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = registrationaccept->gprstimer2.EncodeGPRSTimer2Msg (&registrationaccept->gprstimer2, GPRSTIMER2, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = registrationaccept->emergencynumberlist.EncodeEmergencyNumberListMsg (&registrationaccept->emergencynumberlist, EMERGENCYNUMBERLIST, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;

/*
    if((encodedresult = registrationaccept->extendedemergencynumberlist.EncodeExtendedEmergencyNumberListMsg (&registrationaccept->extendedemergencynumberlist, EXTENDEDNUMBEREMERGENCYLIST, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;
    if((encodedresult = registrationaccept->sortransparentcontainer.EncodeSORTransparentContainerMsg (&registrationaccept->sortransparentcontainer, SORTRANSPARANTCONTAINER, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;
*/
    if((encodedresult = registrationaccept->eapmessage.EncodeEAPMessageMsg (&registrationaccept->eapmessage, EAPMESSAGE, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = registrationaccept->nssaiinclusionmode.EncodeNSSAIInclusionModeMsg (&registrationaccept->nssaiinclusionmode, NSSAIINCLUSIONMODE, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;
/*
    if((encodedresult = registrationaccept->operatordefinedaccesscategorydefinitions.EncodeOperatorDefinedAccessCategoryDefinitionsMsg (&registrationaccept->operatordefinedaccesscategorydefinitions, OPERATORDEFINEDACCESSCATEGORYDEF, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;
*/

    if((encodedresult = registrationaccept->m5gsdrxparameters.EncodeM5GSDRXParametersMsg (&registrationaccept->m5gsdrxparameters, M5GSDRXPARAMETERS, buffer+encoded, len-encoded)) <0)
      return encodedresult;
    else
      encoded += encodedresult;
#endif
    return encoded;
  }
}
