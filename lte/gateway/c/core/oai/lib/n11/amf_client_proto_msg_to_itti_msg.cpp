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
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_38.413.h"
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/lib/n11/amf_client_proto_msg_to_itti_msg.hpp"

#include <stdint.h>
#include <string.h>
#include <iostream>
#include <string>

#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_33.401.h"

#include "lte/gateway/c/core/oai/include/amf_app_messages_types.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/common/security_types.h"

extern "C" {}

using magma::lte::M5GAuthenticationInformationAnswer;

namespace magma5g {

void convert_proto_msg_to_itti_m5g_auth_info_ans(
    M5GAuthenticationInformationAnswer msg,
    itti_amf_subs_auth_info_ans_t* itti_msg) {
  if (msg.m5gauth_vectors_size() > MAX_EPS_AUTH_VECTORS) {
    OAILOG_ERROR(LOG_AMF_APP, "Number of m5g auth vectors received is:%d\n",
                 msg.m5gauth_vectors_size());
    return;
  }
  itti_msg->auth_info.nb_of_vectors = msg.m5gauth_vectors_size();
  uint8_t idx = 0;
  while (idx < itti_msg->auth_info.nb_of_vectors) {
    auto m5gauth_vector = msg.m5gauth_vectors(idx);
    m5gauth_vector_t* itti_m5gauth_vector =
        &(itti_msg->auth_info.m5gauth_vector[idx]);
    if (m5gauth_vector.rand().length() <= RAND_LENGTH_OCTETS) {
      memcpy(itti_m5gauth_vector->rand, m5gauth_vector.rand().c_str(),
             m5gauth_vector.rand().length());
    }
    uint8_t xres_star_len = 0;
    xres_star_len = m5gauth_vector.xres_star().length();
    if ((xres_star_len > XRES_LENGTH_MIN) &&
        (xres_star_len <= XRES_LENGTH_MAX)) {
      itti_m5gauth_vector->xres_star.size = m5gauth_vector.xres_star().length();
      memcpy(itti_m5gauth_vector->xres_star.data,
             m5gauth_vector.xres_star().c_str(), xres_star_len);
    } else {
      OAILOG_ERROR(LOG_AMF_APP, "Invalid xres_star length %d\n", xres_star_len);
      return;
    }
    if (m5gauth_vector.autn().length() == AUTN_LENGTH_OCTETS) {
      memcpy(itti_m5gauth_vector->autn, m5gauth_vector.autn().c_str(),
             m5gauth_vector.autn().length());
    } else {
      OAILOG_ERROR(LOG_AMF_APP, "Invalid AUTN length %lu\n",
                   m5gauth_vector.autn().length());
      return;
    }
    if (m5gauth_vector.kseaf().length() == KSEAF_LENGTH_OCTETS) {
      memcpy(itti_m5gauth_vector->kseaf, m5gauth_vector.kseaf().c_str(),
             m5gauth_vector.kseaf().length());
    } else {
      OAILOG_ERROR(LOG_AMF_APP, "Invalid KSEAF length %lu\n",
                   m5gauth_vector.kseaf().length());
      return;
    }
    ++idx;
  }
  return;
}

void convert_proto_msg_to_itti_amf_decrypted_msin_info_ans(
    M5GSUCIRegistrationAnswer response,
    itti_amf_decrypted_msin_info_ans_t* amf_app_decrypted_msin_info_resp) {
  if (response.ue_msin_recv().length() <= 0) {
    OAILOG_ERROR(LOG_AMF_APP, "Decrypted MSIN response is invalid:%lu\n",
                 response.ue_msin_recv().length());
    return;
  }
  amf_app_decrypted_msin_info_resp->msin_length =
      response.ue_msin_recv().length();
  memcpy(amf_app_decrypted_msin_info_resp->msin,
         response.ue_msin_recv().c_str(),
         amf_app_decrypted_msin_info_resp->msin_length);
}

}  // namespace magma5g
