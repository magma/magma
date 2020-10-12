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
#include "DLNASTransport.h"
#include "CommonDefs.h"

namespace magma5g {
DLNASTransportMsg::DLNASTransportMsg(){};
DLNASTransportMsg::~DLNASTransportMsg(){};

// Decode DLNASTransport Message and its IEs
int DLNASTransportMsg::DecodeDLNASTransportMsg(
    DLNASTransportMsg* dlnastransport, uint8_t* buffer, uint32_t len) {
  uint32_t decoded  = 0;
  int decodedresult = 0;

  // Checking Pointer
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, DL_NAS_TRANSPORT_MINIMUM_LENGTH, len);

  MLOG(MDEBUG) << "DecodeDLNASTransportMsg : \n";
  if ((decodedresult = dlnastransport->extendedprotocoldiscriminator
                           .DecodeExtendedProtocolDiscriminatorMsg(
                               &dlnastransport->extendedprotocoldiscriminator,
                               0, buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = dlnastransport->sparehalfoctet.DecodeSpareHalfOctetMsg(
           &dlnastransport->sparehalfoctet, 0, buffer + decoded,
           len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           dlnastransport->securityheadertype.DecodeSecurityHeaderTypeMsg(
               &dlnastransport->securityheadertype, 0, buffer + decoded,
               len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = dlnastransport->messagetype.DecodeMessageTypeMsg(
           &dlnastransport->messagetype, 0, buffer + decoded, len - decoded)) <
      0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = dlnastransport->sparehalfoctet.DecodeSpareHalfOctetMsg(
           &dlnastransport->sparehalfoctet, 0, buffer + decoded,
           len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           dlnastransport->payloadcontainertype.DecodePayloadContainerTypeMsg(
               &dlnastransport->payloadcontainertype, 0, buffer + decoded,
               len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           dlnastransport->payloadcontainer.DecodePayloadContainerMsg(
               &dlnastransport->payloadcontainer, 0, buffer + decoded,
               len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;

  return decoded;
}

// Encode DL NAS Transport Message and its IEs
int DLNASTransportMsg::EncodeDLNASTransportMsg(
    DLNASTransportMsg* dlnastransport, uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  MLOG(MDEBUG) << "EncodeDLNASTransportMsg:";
  int encodedresult = 0;

  // Check if we got a NDLL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, DL_NAS_TRANSPORT_MINIMUM_LENGTH, len);

  if ((encodedresult = dlnastransport->extendedprotocoldiscriminator
                           .EncodeExtendedProtocolDiscriminatorMsg(
                               &dlnastransport->extendedprotocoldiscriminator,
                               0, buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = dlnastransport->sparehalfoctet.EncodeSpareHalfOctetMsg(
           &dlnastransport->sparehalfoctet, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           dlnastransport->securityheadertype.EncodeSecurityHeaderTypeMsg(
               &dlnastransport->securityheadertype, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = dlnastransport->messagetype.EncodeMessageTypeMsg(
           &dlnastransport->messagetype, 0, buffer + encoded, len - encoded)) <
      0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = dlnastransport->sparehalfoctet.EncodeSpareHalfOctetMsg(
           &dlnastransport->sparehalfoctet, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           dlnastransport->payloadcontainertype.EncodePayloadContainerTypeMsg(
               &dlnastransport->payloadcontainertype, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           dlnastransport->payloadcontainer.EncodePayloadContainerMsg(
               &dlnastransport->payloadcontainer, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;

  return encoded;
}
}  // namespace magma5g
