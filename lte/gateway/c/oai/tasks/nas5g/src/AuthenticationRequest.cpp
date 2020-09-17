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
#include "AuthenticationRequest.h"
#include "CommonDefs.h"

using namespace std;
namespace magma5g
{
  AuthenticationRequestMsg::AuthenticationRequestMsg()
  {
  };

  AuthenticationRequestMsg::~AuthenticationRequestMsg()
  {
  };

  // Decode AuthenticationRequest Messsage
  int AuthenticationRequestMsg::DecodeAuthenticationRequestMsg(AuthenticationRequestMsg *authenticationrequest, uint8_t* buffer, uint32_t len)
  {
    uint32_t decoded = 0;
    /*** Not Implemented, will be supported POST MVC ***/
    return decoded;
  };

  // Encode AuthenticationRequest Messsage
  int AuthenticationRequestMsg::EncodeAuthenticationRequestMsg(AuthenticationRequestMsg *authenticationrequest, uint8_t* buffer, uint32_t len)
  {
    uint32_t encoded = 0;
    int encodedresult = 0;

    // Check if we got a NULL pointer and if buffer length is >= minimum length expected for the message.
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER (buffer, AUTHENTICATION_REQUEST_MINIMUM_LENGTH, len);

    if((encodedresult = authenticationrequest->extendedprotocoldiscriminator.EncodeExtendedProtocolDiscriminatorMsg (&authenticationrequest->extendedprotocoldiscriminator, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;
    if((encodedresult = authenticationrequest->securityheadertype.EncodeSecurityHeaderTypeMsg (&authenticationrequest->securityheadertype, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;
    if((encodedresult = authenticationrequest->sparehalfoctet.EncodeSpareHalfOctetMsg (&authenticationrequest->sparehalfoctet, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;
    if((encodedresult = authenticationrequest->messagetype.EncodeMessageTypeMsg (&authenticationrequest->messagetype, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;
    if((encodedresult = authenticationrequest->naskeysetidentifier.EncodeNASKeySetIdentifierMsg (&authenticationrequest->naskeysetidentifier, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;
    if((encodedresult = authenticationrequest->abba.EncodeABBAMsg (&authenticationrequest->abba, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;
    if((encodedresult = authenticationrequest->authrand.EncodeAuthenticationParameterRANDMsg (&authenticationrequest->authrand, AUTH_PARAM_RAND, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;

    if((encodedresult = authenticationrequest->authautn.EncodeAuthenticationParameterAUTNMsg (&authenticationrequest->authautn, AUTH_PARAM_AUTN, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;
    #ifdef HANDLE_POST_MVC
    if((encodedresult = authenticationrequest->eap.EncodeEAPMessageMsg (&authenticationrequest->eap, EAPMESSAGE, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;
    #endif

    return encoded;
  };
}
