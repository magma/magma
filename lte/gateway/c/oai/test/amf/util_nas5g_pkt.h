/*
 * Copyright 2020 The Magma Authors.
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#pragma once

#include "M5GRegistrationRequest.h"
#include "M5GAuthenticationFailure.h"
#include "M5gNasMessage.h"
#include "M5GULNASTransport.h"

namespace magma5g {

class NAS5GPktSnapShot {
 public:
  static uint8_t reg_req_buffer[38];
  static uint8_t reg_resync_buffer[20];
  static uint8_t guti_based_registration[91];
  static uint8_t pdu_session_est_req_type1[131];
  static uint8_t pdu_session_est_req_type2[47];

  uint32_t get_reg_req_buffer_len() {
    return sizeof(reg_req_buffer) / sizeof(unsigned char);
  }

  uint32_t get_reg_resync_buffer_len() {
    return sizeof(reg_resync_buffer) / sizeof(unsigned char);
  }

  uint32_t get_guti_based_registration_len() {
    return sizeof(guti_based_registration) / sizeof(unsigned char);
  }

  uint32_t get_pdu_session_est_type1_len() {
    return sizeof(pdu_session_est_req_type1) / sizeof(unsigned char);
  }

  uint32_t get_pdu_session_est_type2_len() {
    return sizeof(pdu_session_est_req_type2) / sizeof(unsigned char);
  }

  NAS5GPktSnapShot() {}
};

//  API for testing decode registration request
bool decode_registration_request_msg(
    RegistrationRequestMsg* reg_request, const uint8_t* buffer, uint32_t len);

bool decode_auth_failure_decode_msg(
    AuthenticationFailureMsg* auth_failure, const uint8_t* buffer,
    uint32_t len);

bool decode_pdu_session_est_req_msg(
    ULNASTransportMsg* ul_pdu_sess_est_req, const uint8_t* buffer,
    uint32_t len);
}  // namespace magma5g
