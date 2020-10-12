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

#pragma once
#include <sstream>
#include "PDUSessionEstablishmentRequest.h"
#include "PDUSessionEstablishmentAccept.h"
#include "PDUSessionReleaseRequest.h"

using namespace std;
namespace magma5g {
// Smf NAS Msg Class
class SmfMsgHeader {
 public:
  uint8_t extendedprotocoldiscriminator;
  uint8_t pdusessionid;
  uint8_t proceduretractionid;
  uint8_t messagetype;
};

// Smf NAS Msg Class
class SmfMsg {
 public:
  SmfMsgHeader header;
  PDUSessionEstablishmentRequestMsg pdusessionestablishmentrequest;
  PDUSessionEstablishmentAcceptMsg pdusessionestablishmentaccept;
  PDUSessionReleaseRequestMsg pdusessionreleaserequest;

  SmfMsg();
  ~SmfMsg();
  int SmfMsgDecodeHeaderMsg(SmfMsgHeader* hdr, uint8_t* buffer, uint32_t len);
  int SmfMsgEncodeHeaderMsg(SmfMsgHeader* hdr, uint8_t* buffer, uint32_t len);
  int SmfMsgDecodeMsg(SmfMsg* msg, uint8_t* buffer, uint32_t len);
  int SmfMsgEncodeMsg(SmfMsg* msg, uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g
