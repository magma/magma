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

#include <string>
#include <memory>
#include <gmock/gmock.h>
#include <grpc++/grpc++.h>
#include <gtest/gtest.h>

#include "lte/protos/mobilityd.pb.h"
#include "lte/protos/mobilityd.grpc.pb.h"
#include "orc8r/gateway/c/common/async_grpc/includes/GRPCReceiver.h"

#include "lte/gateway/c/li_agent/src/MobilitydClient.h"
#include "lte/gateway/c/li_agent/src/ProxyConnector.h"
#include "lte/gateway/c/li_agent/src/Utilities.h"

using grpc::Status;
using ::testing::_;
using ::testing::Return;

namespace magma {
namespace lte {

magma::mconfig::LIAgentD create_liagentd_mconfig(const std::string& task_id,
                                                 const std::string& target_id) {
  auto mconfig = get_default_mconfig();
  magma::mconfig::NProbeTask np_task;
  np_task.set_task_id(task_id);
  np_task.set_target_id(target_id);

  auto task = mconfig.add_nprobe_tasks();
  task->CopyFrom(np_task);
  return mconfig;
}

class MockProxyConnector : public ProxyConnector {
 public:
  ~MockProxyConnector() {}

  MOCK_METHOD2(send_data, int(void* data, uint32_t size));
  MOCK_METHOD0(setup_proxy_socket, int());
  MOCK_METHOD0(cleanup, void());
};

class MockMobilitydClient : public MobilitydClient {
 public:
  ~MockMobilitydClient() {}

  MOCK_METHOD2(
      get_subscriber_id_from_ip,
      void(const struct in_addr& addr,
           std::function<void(Status, magma::lte::SubscriberID)> callback));
};

}  // namespace lte
}  // namespace magma
