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
#pragma once

#include "M5GAuthenticationServiceClient.h"

#include <iomanip>
#include <memory>
#include <string>

namespace magma5g {

/**
 * This class is single place holder for all client related services.
 * For instance : subscriberdb, sessiond, mobilityd
 */
class AmfClientServicer {
 public:
  AmfClientServicer() {}

  AmfClientServicer(
      std::shared_ptr<M5GAuthenticationServiceClient> m5g_auth_client)
      : m5g_auth_client_(m5g_auth_client) {}

  void update_auth_servicer(
      std::shared_ptr<M5GAuthenticationServiceClient> m5g_auth_client);

 private:
  std::shared_ptr<M5GAuthenticationServiceClient> m5g_auth_client_;

 public:
  bool get_subscriber_authentication_info(
      const std::string& imsi, uint8_t imsi_length, const char* snni,
      amf_ue_ngap_id_t ue_id);

  bool get_subscriber_authentication_info_resync(
      const std::string& imsi, uint8_t imsi_length, const char* snni,
      const void* resync_info, uint8_t resync_info_len, amf_ue_ngap_id_t ue_id);
};

}  // namespace magma5g
