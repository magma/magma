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

#include <gmock/gmock.h>
#include <grpc++/grpc++.h>
#include <gtest/gtest.h>

#include <folly/io/async/EventBase.h>

#include "DirectorydClient.h"
#include "ProxyConnector.h"

using grpc::Status;
using ::testing::_;
using ::testing::Return;

namespace magma {

class MockProxyConnector : public ProxyConnector {
 public:
  ~MockProxyConnector() {}

  MOCK_METHOD2(send_data, int(void* data, uint32_t size));
  MOCK_METHOD0(setup_proxy_socket, int());
  MOCK_METHOD0(cleanup, void());
};

class MockDirectorydClient : public DirectorydClient {
 public:
  ~MockDirectorydClient() {}

  MOCK_METHOD2(
      get_directoryd_xid_field,
      void(
          const std::string& ip,
          std::function<void(Status status, DirectoryField)> callback));
};

}  // namespace magma
