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
#include "RegistrationRequest.h"
#include "CommonDefs.h"
namespace magma5g
{

  RegistrationRequestMsg::RegistrationRequestMsg()
  {
  };

  RegistrationRequestMsg::~RegistrationRequestMsg()
  {
  };

  // Decode RegistrationRequest Message
  int RegistrationRequestMsg::DecodeRegistrationRequestMsg(RegistrationRequestMsg *registrationrequest, uint8_t* buffer, uint32_t len)
  {
    uint32_t decoded = 0;
    int decodedresult = 0;

    CHECK_PDU_POINTER_AND_LENGTH_DECODER (buffer, REGISTRATION_REQUEST_MINIMUM_LENGTH, len);

    MLOG(MDEBUG) << "DecodeRegistrationRequestMsg : \n";
    if((decodedresult = registrationrequest->extendedprotocoldiscriminator.DecodeExtendedProtocolDiscriminatorMsg(&registrationrequest->extendedprotocoldiscriminator, 0, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    if((decodedresult = registrationrequest->securityheadertype.DecodeSecurityHeaderTypeMsg (&registrationrequest->securityheadertype, 0, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    if((decodedresult = registrationrequest->messagetype.DecodeMessageTypeMsg (&registrationrequest->messagetype, 0, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    if((decodedresult = registrationrequest->m5gsregistrationtype.DecodeM5GSRegistrationTypeMsg (&registrationrequest->m5gsregistrationtype, 0, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    if((decodedresult = registrationrequest->naskeysetidentifier.DecodeNASKeySetIdentifierMsg (&registrationrequest->naskeysetidentifier, 0, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    if((decodedresult = registrationrequest->m5gsmobileidentity.DecodeM5GSMobileIdentityMsg (&registrationrequest->m5gsmobileidentity, 0, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    // TBD
#if 0
    if((decodedresult = registrationrequest->m5gmmcapability.DecodeM5GMMCapabilityMsg (&registrationrequest->m5gmmcapability, M5GMMCAPABILITY, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    if((decodedresult = registrationrequest->uesecuritycapability.DecodeUESecurityCapabilityMsg (&registrationrequest->uesecuritycapability, UESECURITYCAPABILITY, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;
    //MLOG(DEBUG) << "__func__: "After DecodeUESecurityCapabilityMsg buffer " << hex << int(*buffer) << "decoded " << decodedresult << endl;

    if((decodedresult = registrationrequest->pdusessionstatus.DecodePDUSessionStatusMsg (&registrationrequest->pdusessionstatus, PDUSESSIONSTATUS, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    if((decodedresult = registrationrequest->nasmessagecontainer.DecodeNASMessageContainerMsg (&registrationrequest->nasmessagecontainer, NASMESSAGECONTAINER, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    if((decodedresult = registrationrequest->nssai.DecodeNSSAIMsg (&registrationrequest->nssai, REQUESTEDNSSAI, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    if((decodedresult = registrationrequest->m5gstrackingareaidentity.DecodeM5GSTrackingAreaIdentityMsg (&registrationrequest->m5gstrackingareaidentity, TAI, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    /*    if((decodedresult = registrationrequest->s1uenetworkcapability.DecodeS1UENetworkCapabilityMsg (&registrationrequest->s1uenetworkcapability, 0x17, buffer+decoded, len-decoded))<0)
          return decodedresult;
          else
          decoded += decodedresult;
     */
    if((decodedresult = registrationrequest->uplinkdatastatus.DecodeUplinkDataStatusMsg (&registrationrequest->uplinkdatastatus, UPLINKDATASTATUS, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    if((decodedresult = registrationrequest->pdusessionstatus.DecodePDUSessionStatusMsg (&registrationrequest->pdusessionstatus, PDUSESSIONSTATUS, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    if((decodedresult = registrationrequest->micoindication.DecodeMICOIndicationMsg (&registrationrequest->micoindication, MICOINDICATION, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    if((decodedresult = registrationrequest->uestatus.DecodeUEStatusMsg (&registrationrequest->uestatus, UESTATUS, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    if((decodedresult = registrationrequest->allowedpdusessionstatus.DecodeAllowedPDUSessionStatusMsg (&registrationrequest->allowedpdusessionstatus, ALLOWEDPDUSESSIONSTATUS, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    if((decodedresult = registrationrequest->uesusagesetting.DecodeUESUsageSettingMsg (&registrationrequest->uesusagesetting, UEUSSAGESETTING, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    if((decodedresult = registrationrequest->m5gsdrxparameters.DecodeM5GSDRXParametersMsg (&registrationrequest->m5gsdrxparameters, M5GSDRXPARAMETERS, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    if((decodedresult = registrationrequest->epsnasmessagecontainer.DecodeEPSNASMessageContainerMsg (&registrationrequest->epsnasmessagecontainer, EPSNASMESSAGECONTAINER, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    if((decodedresult = registrationrequest->ladnindication.DecodeLADNIndicationMsg (&registrationrequest->ladnindication, LADNINFORMATION, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;
    /*
       if((decodedresult = registrationrequest->payloadcontainer.DecodePayloadContainerMsg (&registrationrequest->payloadcontainer, PAYLOADCONTAINER, buffer+decoded, len-decoded))<0)
       return decodedresult;
       else
       decoded += decodedresult;
     */
    if((decodedresult = registrationrequest->networkslicingindication.DecodeNetworkSlicingIndicationMsg (&registrationrequest->networkslicingindication, NETWORKSLICINGINDICATION, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    if((decodedresult = registrationrequest->m5gsupdatetype.DecodeM5GSUpdateTypeMsg (&registrationrequest->m5gsupdatetype, M5GSUPDATETYPE, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    if((decodedresult = registrationrequest->nasmessagecontainer.DecodeNASMessageContainerMsg (&registrationrequest->nasmessagecontainer, NASMESSAGECONTAINER, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;
#endif

    return decoded;
  }

  // TBD
  // Encode Registration Request Message
  int RegistrationRequestMsg::EncodeRegistrationRequestMsg( RegistrationRequestMsg *registrationrequest, uint8_t* buffer, uint32_t len)
  {
    uint32_t encoded = 0;
    
    MLOG(MDEBUG) << "EncodeRegistrationRequestMsg:";
#if 0
    int encodedresult = 0;

    // Check if we got a NULL pointer and if buffer length is >= minimum length expected for the message.
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER (buffer, REGISTRATION_REQUEST_MINIMUM_LENGTH, len);

    if((encodedresult = registrationrequest->EncodeExtendedProtocolDiscriminatorMsg (registrationrequest->extendedprotocoldiscriminator, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = registrationrequest->EncodeSecurityHeaderTypeMsg (registrationrequest->securityheadertype, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = registrationrequest->EncodeMessageTypeMsg (registrationrequest->messagetype, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = registrationrequest->Encode5GSRegistrationTypeMsg (registrationrequest->_5gsregistrationtype, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = registrationrequest->EncodeNasKeySetIdentifierMsg (registrationrequest->naskeysetidentifier, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = registrationrequest->Encode5GSMobileIdentityMsg (registrationrequest->_5gsmobileidentity, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = registrationrequest->Encode5GMMCapabilityMsg (registrationrequest->_5gmmcapability, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = Encodeue_security_capability (registrationrequest->uesecuritycapability, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = Encodenssai (registrationrequest->nssai, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = Encode_5gs_tracking_area_identity (registrationrequest->_5gstrackingareaidentity, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = Encodes1_ue_network_capability (registrationrequest->s1uenetworkcapability, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = Encodeuplink_data_status (registrationrequest->uplinkdatastatus, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = Encodepdu_session_status (registrationrequest->pdusessionstatus, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = Encodemico_indication (registrationrequest->micoindication, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = Encodeue_status (registrationrequest->uestatus, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = Encodeallowed_pdu_session_status (registrationrequest->allowedpdusessionstatus, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = Encodeues_usage_setting (registrationrequest->uesusagesetting, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = Encode_5gsdrx_parameters (registrationrequest->_5gsdrxparameters, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = Encodeepsnas_message_container (registrationrequest->epsnasmessagecontainer, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = Encodeladn_indication (registrationrequest->ladnindication, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = Encodepayload_container (registrationrequest->payloadcontainer, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = Encodenetwork_slicing_indication (registrationrequest->networkslicingindication, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = Encode_5gs_update_type (registrationrequest->_5gsupdatetype, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = Encodenas_message_container (registrationrequest->nasmessagecontainer, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

#endif

    return encoded;
  }
}
