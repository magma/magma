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
#include "DeRegistrationAcceptUEInit.h"
#include "CommonDefs.h"

using namespace std;
namespace magma5g
{

  DeRegistrationAcceptUEInitMsg::DeRegistrationAcceptUEInitMsg()
  {
  };

  DeRegistrationAcceptUEInitMsg::~DeRegistrationAcceptUEInitMsg()
  {
  };

  //Decoding De Registration Accept Message and its IEs
  int DeRegistrationAcceptUEInitMsg::DecodeDeRegistrationAcceptUEInitMsg(DeRegistrationAcceptUEInitMsg *deregistrationaccept, uint8_t* buffer, uint32_t len)
  {
    uint32_t decoded = 0;
    int decodedresult = 0;

    CHECK_PDU_POINTER_AND_LENGTH_DECODER (buffer, DEREGISTRATION_ACCEPT_UEINIT_MINIMUM_LENGTH, len);

    if((decodedresult = deregistrationaccept->extendedprotocoldiscriminator.DecodeExtendedProtocolDiscriminatorMsg(&deregistrationaccept->extendedprotocoldiscriminator, 0, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;
    if((decodedresult = deregistrationaccept->securityheadertype.DecodeSecurityHeaderTypeMsg (&deregistrationaccept->securityheadertype, 0, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;
    if((decodedresult = deregistrationaccept->sparehalfoctet.DecodeSpareHalfOctetMsg (&deregistrationaccept->sparehalfoctet, 0, buffer+decoded, len-decoded))<0)
	    return decodedresult;
    else
	    decoded += decodedresult;
    if((decodedresult = deregistrationaccept->messagetype.DecodeMessageTypeMsg (&deregistrationaccept->messagetype, 0, buffer+decoded, len-decoded))<0)
      return decodedresult;
    else
      decoded += decodedresult;
    return decoded;
  };

  //Encoding De Registration Accept Message and its IEs
  int DeRegistrationAcceptUEInitMsg::EncodeDeRegistrationAcceptUEInitMsg( DeRegistrationAcceptUEInitMsg *deregistrationaccept, uint8_t* buffer, uint32_t len)
  {
    uint32_t encoded = 0;
    int encodedresult = 0;

    // Check if we got a NULL pointer and if buffer length is >= minimum length expected for the message.
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER (buffer, DEREGISTRATION_ACCEPT_UEINIT_MINIMUM_LENGTH, len);

    if((encodedresult = deregistrationaccept->extendedprotocoldiscriminator.EncodeExtendedProtocolDiscriminatorMsg (&deregistrationaccept->extendedprotocoldiscriminator, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;
    if((encodedresult = deregistrationaccept->securityheadertype.EncodeSecurityHeaderTypeMsg (&deregistrationaccept->securityheadertype, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;
    if((encodedresult = deregistrationaccept->sparehalfoctet.EncodeSpareHalfOctetMsg (&deregistrationaccept->sparehalfoctet, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;
    if((encodedresult = deregistrationaccept->messagetype.EncodeMessageTypeMsg (&deregistrationaccept->messagetype, 0, buffer+encoded, len-encoded))<0)
      return encodedresult;
    else
      encoded += encodedresult;
    return encoded;
  };
}
