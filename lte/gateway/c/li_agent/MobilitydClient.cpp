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

#include <netinet/in.h>
#include <thread>

#include "MobilitydClient.h"
#include "ServiceRegistrySingleton.h"

#include "lte/protos/mobilityd.grpc.pb.h"
#include "lte/protos/mobilityd.pb.h"
#include "magma_logging.h"

using grpc::ClientContext;
using grpc::Status;
using magma::orc8r::Void;

namespace magma {
namespace lte {

int MobilitydClient::GetSubscriberIDFromIP(
    const struct in_addr& addr, std::string* imsi) {
  IPAddress ip_addr = IPAddress();
  ip_addr.set_version(IPAddress::IPV4);
  ip_addr.set_address(&addr, sizeof(struct in_addr));

  Void resp;
  SubscriberID match;
  ClientContext context;
  Status status = stub_->GetSubscriberIDFromIP(&context, ip_addr, &match);
  if (!status.ok()) {
    MLOG(MERROR) << "GetSubscriberIDFromIPv4 fails with code "
                 << status.error_code() << ", msg: " << status.error_message();
    return status.error_code();
  }
  imsi->assign(match.id());
  return 0;
}

MobilitydClient::MobilitydClient() {
  auto channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
      "mobilityd", ServiceRegistrySingleton::LOCAL);
  stub_ = MobilityService::NewStub(channel);
  std::thread resp_loop_thread([&]() { rpc_response_loop(); });
  resp_loop_thread.detach();
}

MobilitydClient& MobilitydClient::getInstance() {
  static MobilitydClient instance;
  return instance;
}

}  // namespace lte
}  // namespace magma
