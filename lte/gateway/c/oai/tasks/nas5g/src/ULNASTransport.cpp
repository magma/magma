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
#include "ULNASTransport.h"
#include "CommonDefs.h"

namespace magma5g {
ULNASTransportMsg::ULNASTransportMsg(){};
ULNASTransportMsg::~ULNASTransportMsg(){};

// Decode ULNASTransport Message and its IEs
int ULNASTransportMsg::DecodeULNASTransportMsg(
    ULNASTransportMsg* ulnastransport, uint8_t* buffer, uint32_t len) {
  uint32_t decoded  = 0;
  int decodedresult = 0;

  // Checking Pointer
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, DL_NAS_TRANSPORT_MINIMUM_LENGTH, len);

  MLOG(MDEBUG) << "DecodeULNASTransportMsg : \n";
  if ((decodedresult = ulnastransport->extendedprotocoldiscriminator
                           .DecodeExtendedProtocolDiscriminatorMsg(
                               &ulnastransport->extendedprotocoldiscriminator,
                               0, buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = ulnastransport->sparehalfoctet.DecodeSpareHalfOctetMsg(
           &ulnastransport->sparehalfoctet, 0, buffer + decoded,
           len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           ulnastransport->securityheadertype.DecodeSecurityHeaderTypeMsg(
               &ulnastransport->securityheadertype, 0, buffer + decoded,
               len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = ulnastransport->messagetype.DecodeMessageTypeMsg(
           &ulnastransport->messagetype, 0, buffer + decoded, len - decoded)) <
      0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           ulnastransport->payloadcontainertype.DecodePayloadContainerTypeMsg(
               &ulnastransport->payloadcontainertype, 0, buffer + decoded,
               len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           ulnastransport->payloadcontainer.DecodePayloadContainerMsg(
               &ulnastransport->payloadcontainer, 0, buffer + decoded,
               len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;

  return decoded;
}

// Encode DL NAS Transport Message and its IEs
int ULNASTransportMsg::EncodeULNASTransportMsg(
    ULNASTransportMsg* ulnastransport, uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  MLOG(MDEBUG) << "EncodeULNASTransportMsg:";
  int encodedresult = 0;

  // Check if we got a NDLL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, DL_NAS_TRANSPORT_MINIMUM_LENGTH, len);

  if ((encodedresult = ulnastransport->extendedprotocoldiscriminator
                           .EncodeExtendedProtocolDiscriminatorMsg(
                               &ulnastransport->extendedprotocoldiscriminator,
                               0, buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = ulnastransport->sparehalfoctet.EncodeSpareHalfOctetMsg(
           &ulnastransport->sparehalfoctet, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           ulnastransport->securityheadertype.EncodeSecurityHeaderTypeMsg(
               &ulnastransport->securityheadertype, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = ulnastransport->messagetype.EncodeMessageTypeMsg(
           &ulnastransport->messagetype, 0, buffer + encoded, len - encoded)) <
      0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = ulnastransport->sparehalfoctet.EncodeSpareHalfOctetMsg(
           &ulnastransport->sparehalfoctet, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           ulnastransport->payloadcontainertype.EncodePayloadContainerTypeMsg(
               &ulnastransport->payloadcontainertype, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           ulnastransport->payloadcontainer.EncodePayloadContainerMsg(
               &ulnastransport->payloadcontainer, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;

  return encoded;
}
}  // namespace magma5g
