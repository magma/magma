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

#include "include/amf_client_servicer.h"
#include <memory>

namespace magma5g {

AmfClientServicer amf_client_servicer_g;

// Fetch the amf client servicer reference
AmfClientServicer& get_amf_client_server_ref() {
  return amf_client_servicer_g;
}

// Initialize the client servicer layer
void amf_client_servicer_init() {
  auto authentication_client =
      std::make_shared<AsyncM5GAuthenticationServiceClient>();

  amf_client_servicer_g = AmfClientServicer(authentication_client);
}

// For authentication to subscriberdb
bool AmfClientServicer::get_subscriber_authentication_info(
    const std::string& imsi, uint8_t imsi_length, const char* snni,
    amf_ue_ngap_id_t ue_id) {
  bool result = false;

  result = m5g_auth_client_->get_subs_auth_info(imsi, imsi_length, snni, ue_id);

  return (result);
}

// For authentication to subscriberdb after resync
bool AmfClientServicer::get_subscriber_authentication_info_resync(
    const std::string& imsi, uint8_t imsi_length, const char* snni,
    const void* resync_info, uint8_t resync_info_len, amf_ue_ngap_id_t ue_id) {
  bool result = false;

  result = m5g_auth_client_->get_subs_auth_info_resync(
      imsi, imsi_length, snni, resync_info, resync_info_len, ue_id);

  return (result);
}

}  // namespace magma5g
