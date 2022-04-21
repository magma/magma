/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include <iostream>
#include "lte/gateway/c/core/oai/test/amf/util_nas5g_pkt.hpp"

namespace magma5g {

//  API for testing decode UL NAS Transport request
bool decode_ul_nas_transport_msg(ULNASTransportMsg* ul_nas_pdu,
                                 const uint8_t* buffer, uint32_t len) {
  bool decode_success = true;
  uint8_t* decode_ul_nas_pdu_buffer =
      const_cast<uint8_t*>(reinterpret_cast<const uint8_t*>(buffer));

  if (ul_nas_pdu->DecodeULNASTransportMsg(ul_nas_pdu, decode_ul_nas_pdu_buffer,
                                          len) < 0) {
    decode_success = false;
  }

  return (decode_success);
}

//  API for testing decode UL NAS Deregister request
bool decode_ul_nas_deregister_request_msg(
    DeRegistrationRequestUEInitMsg* ul_nas_pdu, const uint8_t* buffer,
    uint32_t len) {
  bool decode_success = true;
  uint8_t* decode_ul_nas_dereg_pdu_buffer =
      const_cast<uint8_t*>(reinterpret_cast<const uint8_t*>(buffer));

  if (ul_nas_pdu->DecodeDeRegistrationRequestUEInitMsg(
          ul_nas_pdu, decode_ul_nas_dereg_pdu_buffer, len) < 0) {
    decode_success = false;
  }

  return (decode_success);
}

}  // namespace magma5g
