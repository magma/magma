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

#include <stdint.h>
#include <functional>
#include <memory>

namespace magma5g {

class M5GMobilityServiceClient {
 public:
  virtual ~M5GMobilityServiceClient() {}
  virtual int allocate_ipv4_address(
      const char* subscriber_id, const char* apn, uint32_t pdu_session_id,
      uint8_t pti, uint32_t pdu_session_type, uint32_t gnb_gtp_teid,
      uint8_t* gnb_gtp_teid_ip_addr, uint8_t gnb_gtp_teid_ip_addr_len) = 0;

  virtual int release_ipv4_address(
      const char* subscriber_id, const char* apn,
      const struct in_addr* addr) = 0;
};

class AsyncM5GMobilityServiceClient : public M5GMobilityServiceClient {
 public:
  int allocate_ipv4_address(
      const char* subscriber_id, const char* apn, uint32_t pdu_session_id,
      uint8_t pti, uint32_t pdu_session_type, uint32_t gnb_gtp_teid,
      uint8_t* gnb_gtp_teid_ip_addr, uint8_t gnb_gtp_teid_ip_addr_len);

  int release_ipv4_address(
      const char* subscriber_id, const char* apn, const struct in_addr* addr);

  static AsyncM5GMobilityServiceClient& getInstance();

  AsyncM5GMobilityServiceClient(AsyncM5GMobilityServiceClient const&) = delete;
  void operator=(AsyncM5GMobilityServiceClient const&) = delete;

 private:
  AsyncM5GMobilityServiceClient();
};
}  // namespace magma5g
