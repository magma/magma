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
#include "AuthenticationResponse.h"
#include "CommonDefs.h"

using namespace std;
namespace magma5g
{
  AuthenticationResponseMsg::AuthenticationResponseMsg()
  {
  };

  AuthenticationResponseMsg::~AuthenticationResponseMsg()
  {
  };

  // Decode AuthenticationResponse Messsage
  int AuthenticationResponseMsg::DecodeAuthenticationResponseMsg(AuthenticationResponseMsg *authenticationresponse, uint8_t* buffer, uint32_t len)
  {
    uint32_t decoded = 0;
    int decodedresult = 0;

    CHECK_PDU_POINTER_AND_LENGTH_DECODER (buffer, AUTHENTICATION_RESPONSE_MINIMUM_LENGTH, len);

    MLOG(MDEBUG) << " ---Decoding AuAuthentication Response Message---\n" << endl;
    if((decodedresult = authenticationresponse->extendedprotocoldiscriminator.DecodeExtendedProtocolDiscriminatorMsg(&authenticationresponse->extendedprotocoldiscriminator, 0, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;
    if((decodedresult = authenticationresponse->securityheadertype.DecodeSecurityHeaderTypeMsg (&authenticationresponse->securityheadertype, 0, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;
    if((decodedresult = authenticationresponse->sparehalfoctet.DecodeSpareHalfOctetMsg (&authenticationresponse->sparehalfoctet, 0, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;
    if((decodedresult = authenticationresponse->messagetype.DecodeMessageTypeMsg (&authenticationresponse->messagetype, 0, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;

    return decoded;
  };

  
  // Encode AuthenticationResponse Messsage
  int AuthenticationResponseMsg::EncodeAuthenticationResponseMsg(AuthenticationResponseMsg *authenticationresponse, uint8_t* buffer, uint32_t len)
  {
    uint32_t encoded = 0;
    //Not Implemented, Will be supported POST MVC
    return encoded;
  };
}
