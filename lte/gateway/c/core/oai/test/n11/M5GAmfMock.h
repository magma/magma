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

#include <gmock/gmock.h>
#include <grpc++/grpc++.h>
#include <gtest/gtest.h>

#include <memory>
#include <string>
#include <vector>

using grpc::Status;
using ::testing::_;
using ::testing::Return;

namespace magma5g {

class MockM5GAuthenticationServiceClient
    : public M5GAuthenticationServiceClient {
 public:
  MockM5GAuthenticationServiceClient() {
    ON_CALL(*this, GetSubscriberAuthInfoRPC(_, _)).WillByDefault(Return());
  }
  MOCK_METHOD4(
      get_subs_auth_info, bool(
                              const std::string& imsi, uint8_t imsi_length,
                              const char* snni, amf_ue_ngap_id_t ue_id));

  MOCK_METHOD6(
      get_subs_auth_info_resync,
      bool(
          const std::string& imsi, uint8_t imsi_length, const char* snni,
          const void* resync_info, uint8_t resync_info_len,
          amf_ue_ngap_id_t ue_id));

 private:
  MOCK_METHOD2(
      GetSubscriberAuthInfoRPC,
      void(
          M5GAuthenticationInformationRequest& request,
          const std::function<void(
              grpc::Status, M5GAuthenticationInformationAnswer)>& callback));
};
}  // namespace magma5g
