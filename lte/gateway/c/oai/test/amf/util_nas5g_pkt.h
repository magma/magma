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
#include "M5gNasMessage.h"

namespace magma5g {

class NAS5GPktSnapShot {
   public:
       static uint8_t reg_req_buffer[38];

       uint32_t get_reg_req_buffer_len() {
           return sizeof(reg_req_buffer) / sizeof (unsigned char);
       }

   NAS5GPktSnapShot() {};
};

// Base test function
bool decode_registration_request_msg(RegistrationRequestMsg* reg_request,
                                     const uint8_t* buffer, uint32_t len);

} //magma5g
