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

#include <stdint.h>
#include <string.h>
#include <iostream>
#include <string>

#include "3gpp_33.401.h"
#include "amf_app_messages_types.h"
#include "amf_client_proto_msg_to_itti_msg.h"
#include "common_types.h"
#include "security_types.h"

extern "C" {}

using magma::lte::M5GAuthenticationInformationAnswer;

namespace magma5g {

void convert_proto_msg_to_itti_m5g_auth_info_ans(
    M5GAuthenticationInformationAnswer msg,
    itti_amf_subs_auth_info_ans_t* itti_msg) {
  if (msg.eutran_vectors_size() > MAX_EPS_AUTH_VECTORS) {
    std::cout << "[ERROR] Number of eutran auth vectors received is:"
              << msg.eutran_vectors_size() << std::endl;
    return;
  }
  itti_msg->auth_info.nb_of_vectors = msg.eutran_vectors_size();
  uint8_t idx                       = 0;
  while (idx < itti_msg->auth_info.nb_of_vectors) {
    auto eutran_vector = msg.eutran_vectors(idx);
    eutran_vector_t* itti_eutran_vector =
        &(itti_msg->auth_info.eutran_vector[idx]);
    if (eutran_vector.rand().length() <= RAND_LENGTH_OCTETS) {
      memcpy(
          itti_eutran_vector->rand, eutran_vector.rand().c_str(),
          eutran_vector.rand().length());
    }
    uint8_t xres_len = 0;
    xres_len         = eutran_vector.xres().length();
    if ((xres_len > XRES_LENGTH_MIN) && (xres_len <= XRES_LENGTH_MAX)) {
      itti_eutran_vector->xres.size = eutran_vector.xres().length();
      memcpy(
          itti_eutran_vector->xres.data, eutran_vector.xres().c_str(),
          xres_len);
    } else {
      std::cout << "[ERROR] Invalid xres length " << xres_len << std::endl;
      return;
    }
    if (eutran_vector.autn().length() == AUTN_LENGTH_OCTETS) {
      memcpy(
          itti_eutran_vector->autn, eutran_vector.autn().c_str(),
          eutran_vector.autn().length());
    } else {
      std::cout << "[ERROR] Invalid AUTN length "
                << eutran_vector.autn().length() << std::endl;
      return;
    }
    if (eutran_vector.kasme().length() == KASME_LENGTH_OCTETS) {
      memcpy(
          itti_eutran_vector->kasme, eutran_vector.kasme().c_str(),
          eutran_vector.kasme().length());
    } else {
      std::cout << "[ERROR] Invalid KASME length "
                << eutran_vector.kasme().length() << std::endl;
      return;
    }
    if (eutran_vector.ck().length() == CK_LENGTH_OCTETS) {
      memcpy(
          itti_eutran_vector->ck, eutran_vector.ck().c_str(),
          eutran_vector.ck().length());
    } else {
      std::cout << "[ERROR] Invalid CK length " << eutran_vector.ck().length()
                << std::endl;
      return;
    }
    if (eutran_vector.ik().length() == IK_LENGTH_OCTETS) {
      memcpy(
          itti_eutran_vector->ik, eutran_vector.ik().c_str(),
          eutran_vector.ik().length());
    } else {
      std::cout << "[ERROR] Invalid IK length " << eutran_vector.ik().length()
                << std::endl;
      return;
    }
    ++idx;
  }
  return;
}

}  // namespace magma5g
