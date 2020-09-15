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
#include <iostream>
#include <sstream>
#include "DeRegistrationRequestUEInit.h"
#include "CommonDefs.h"

using namespace std;
namespace magma5g
{
  DeRegistrationRequestUEInitMsg::DeRegistrationRequestUEInitMsg()
  {
  };

  DeRegistrationRequestUEInitMsg::~DeRegistrationRequestUEInitMsg()
  {
  };

  //Decode De Registration Request(UE) Message and its IEs
  int DeRegistrationRequestUEInitMsg::DecodeDeRegistrationRequestUEInitMsg(DeRegistrationRequestUEInitMsg *deregistrationrequest, uint8_t* buffer, uint32_t len)
  {
    uint32_t decoded = 0;
    int decodedresult = 0;

    CHECK_PDU_POINTER_AND_LENGTH_DECODER (buffer, DEREGISTRATION_REQUEST_UEINIT_MINIMUM_LENGTH, len);

    MLOG(MDEBUG) << "\n\n---Decoding De-Registration Request Message---\n" << endl;
    if((decodedresult = deregistrationrequest->extendedprotocoldiscriminator.DecodeExtendedProtocolDiscriminatorMsg(&deregistrationrequest->extendedprotocoldiscriminator, 0, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;
    if((decodedresult = deregistrationrequest->securityheadertype.DecodeSecurityHeaderTypeMsg (&deregistrationrequest->securityheadertype, 0, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;
    if((decodedresult = deregistrationrequest->sparehalfoctet.DecodeSpareHalfOctetMsg (&deregistrationrequest->sparehalfoctet, 0, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;
    if((decodedresult = deregistrationrequest->messagetype.DecodeMessageTypeMsg (&deregistrationrequest->messagetype, 0, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;
    if((decodedresult = deregistrationrequest->m5gsderegistrationtype.DecodeM5GSDeRegistrationTypeMsg (&deregistrationrequest->m5gsderegistrationtype, 0, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;
    if((decodedresult = deregistrationrequest->naskeysetidentifier.DecodeNASKeySetIdentifierMsg (&deregistrationrequest->naskeysetidentifier, 0, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;
    if((decodedresult = deregistrationrequest->m5gsmobileidentity.DecodeM5GSMobileIdentityMsg (&deregistrationrequest->m5gsmobileidentity, 0, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;
    return decoded;
  };

  //Encode De Registration Request(UE) Message and its IEs
  int DeRegistrationRequestUEInitMsg::EncodeDeRegistrationRequestUEInitMsg( DeRegistrationRequestUEInitMsg *deregistrationrequest, uint8_t* buffer, uint32_t len)
  {
    uint32_t encoded = 0;
    int encodedresult = 0;

    // Check if we got a NULL pointer and if buffer length is >= minimum length expected for the message.
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER (buffer, DEREGISTRATION_REQUEST_UEINIT_MINIMUM_LENGTH, len);

    if((encodedresult = deregistrationrequest->extendedprotocoldiscriminator.EncodeExtendedProtocolDiscriminatorMsg (&deregistrationrequest->extendedprotocoldiscriminator, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;
    if((encodedresult = deregistrationrequest->securityheadertype.EncodeSecurityHeaderTypeMsg (&deregistrationrequest->securityheadertype, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;
    if((encodedresult = deregistrationrequest->sparehalfoctet.EncodeSpareHalfOctetMsg (&deregistrationrequest->sparehalfoctet, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;
    if((encodedresult = deregistrationrequest->messagetype.EncodeMessageTypeMsg (&deregistrationrequest->messagetype, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;
    if((encodedresult = deregistrationrequest->m5gsderegistrationtype.EncodeM5GSDeRegistrationTypeMsg (&deregistrationrequest->m5gsderegistrationtype, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;
    if((encodedresult = deregistrationrequest->naskeysetidentifier.EncodeNASKeySetIdentifierMsg (&deregistrationrequest->naskeysetidentifier, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;
    if((encodedresult = deregistrationrequest->m5gsmobileidentity.EncodeM5GSMobileIdentityMsg (&deregistrationrequest->m5gsmobileidentity, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;
    return encoded;
  };
}
