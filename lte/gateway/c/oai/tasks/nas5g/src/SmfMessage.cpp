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
#include "SmfMessage.h"
#include "M5gNasMessage.h"
#include "CommonDefs.h"

using namespace std;
namespace magma5g {
SmfMsg::SmfMsg(){};
SmfMsg::~SmfMsg(){};

// Decode SMF Message Header
int SmfMsg::SmfMsgDecodeHeaderMsg(
    SmfMsgHeader* hdr, uint8_t* buffer, uint32_t len) {
  int size = 0;

  MLOG(MDEBUG) << "SmfMsgDecodeHeaderMsg:" << endl;
  if (len > 0 || buffer != NULL) {
    DECODE_U8(buffer + size, hdr->extendedprotocoldiscriminator, size);
    DECODE_U8(buffer + size, hdr->pdusessionid, size);
    DECODE_U8(buffer + size, hdr->proceduretractionid, size);
    DECODE_U8(buffer + size, hdr->messagetype, size);
    MLOG(MDEBUG) << "epd = 0x" << hex << int(hdr->extendedprotocoldiscriminator)
                 << "pdu session id = 0x" << hex << int(hdr->pdusessionid)
                 << " proceduretractionid = 0x" << hex
                 << int(hdr->proceduretractionid) << " messagetype = 0x" << hex
                 << int(hdr->messagetype);
  } else {
    MLOG(MERROR) << "Error : Buffer is Empty" << endl;
    return (RETURN_ERROR);
  }

  if (hdr->extendedprotocoldiscriminator != M5G_SESSION_MANAGEMENT_MESSAGES) {
    MLOG(MERROR) << "Error : TLV not supported" << endl;
    return (TLV_PROTOCOL_NOT_SUPPORTED);
  }
  return (size);
}

// Encode SMF Message Header
int SmfMsg::SmfMsgEncodeHeaderMsg(
    SmfMsgHeader* hdr, uint8_t* buffer, uint32_t len) {
  int size = 0;

  MLOG(MDEBUG) << "SmfMsgEncodeHeaderMsg:";
  if (len > 0 || buffer != NULL) {
    ENCODE_U8(buffer + size, hdr->extendedprotocoldiscriminator, size);
    ENCODE_U8(buffer + size, hdr->pdusessionid, size);
    ENCODE_U8(buffer + size, hdr->proceduretractionid, size);
    ENCODE_U8(buffer + size, hdr->messagetype, size);
    MLOG(MDEBUG) << "epd = 0x" << hex << int(hdr->extendedprotocoldiscriminator)
                 << "pdu session id = 0x" << hex << int(hdr->pdusessionid)
                 << " proceduretractionid = 0x" << hex
                 << int(hdr->proceduretractionid) << " messagetype = 0x" << hex
                 << int(hdr->messagetype);
  } else {
    MLOG(MERROR) << "Error : Buffer is Empty ";
    return (RETURN_ERROR);
  }
  if (hdr->extendedprotocoldiscriminator != M5G_SESSION_MANAGEMENT_MESSAGES) {
    MLOG(MERROR) << "Error : TLV not supported";
    return (TLV_PROTOCOL_NOT_SUPPORTED);
  }
  return (size);
}

// Decode SMF Message
int SmfMsg::SmfMsgDecodeMsg(SmfMsg* msg, uint8_t* buffer, uint32_t len) {
  int decoderesult = 0;
  int headerresult = 0;

  MLOG(MDEBUG) << "SmfMsgDecodeMsg:" << endl;
  if (len <= 0 || buffer == NULL) {
    MLOG(MERROR) << "Error : Buffer is Empty" << endl;
    return (RETURN_ERROR);
  }

  headerresult = msg->SmfMsgDecodeHeaderMsg(&msg->header, buffer, len);
  if (headerresult <= 0) {
    MLOG(MERROR) << "   Error : Header Decoding Failed" << std::dec
                 << RETURN_ERROR;
    return (RETURN_ERROR);
  }

  buffer       = buffer + headerresult;
  decoderesult = decoderesult + headerresult;

  MLOG(MDEBUG) << "msg type = 0x" << hex << int(msg->header.messagetype);
  switch (int(msg->header.messagetype)) {
    case PDU_SESSION_ESTABLISHMENT_REQUEST:
      MLOG(MDEBUG) << "PDU Session Establishment request msg" << endl;
      decoderesult = msg->pdusessionestablishmentrequest
                         .DecodePDUSessionEstablishmentRequestMsg(
                             &msg->pdusessionestablishmentrequest, buffer, len);
      break;
    case PDU_SESSION_RELEASE_REQUEST:
      MLOG(MDEBUG) << "PDU Session Release request msg" << endl;
      decoderesult = msg->pdusessionreleaserequest
                         .DecodePDUSessionReleaseRequestMsg(
                             &msg->pdusessionreleaserequest, buffer, len);
      break;
    default:
      decoderesult = TLV_WRONG_MESSAGE_TYPE;
  }
  return (decoderesult);
}

// Encode SMF Message
int SmfMsg::SmfMsgEncodeMsg(SmfMsg* msg, uint8_t* buffer, uint32_t len) {
  int encoderesult = 0;
  int headerresult = 0;

  MLOG(MDEBUG) << " SmfMsgEncodeMsg : " << endl;
  if (len <= 0 || buffer == NULL) {
    MLOG(MERROR) << "Error : Buffer is Empty";
    return (RETURN_ERROR);
  }

  headerresult = msg->SmfMsgEncodeHeaderMsg(&msg->header, buffer, len);
  if (headerresult <= 0) {
    MLOG(MERROR) << "   Error : Header Encoding Failed" << std::dec
                 << RETURN_ERROR;
    return (RETURN_ERROR);
  }

  buffer       = buffer + headerresult;
  encoderesult = encoderesult + headerresult;

  MLOG(MDEBUG) << "msg type = 0x" << hex << int(msg->header.messagetype);
  switch ((unsigned char) msg->header.messagetype) {
    case PDU_SESSION_ESTABLISHMENT_ACCEPT:
      MLOG(MDEBUG) << "PDU Session Establishment accept msg" << endl;
      encoderesult = msg->pdusessionestablishmentaccept
                         .EncodePDUSessionEstablishmentAcceptMsg(
                             &msg->pdusessionestablishmentaccept, buffer, len);
      break;
    default:
      encoderesult = TLV_WRONG_MESSAGE_TYPE;
  }
  return (encoderesult);
}
}  // namespace magma5g
