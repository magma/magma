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

#include "lte/gateway/c/sctpd/src/sctpd.h"

#include <memory>
#include <grpcpp/grpcpp.h>
#include <signal.h>

#include "lte/gateway/c/sctpd/src/sctpd_downlink_impl.h"
#include "lte/gateway/c/sctpd/src/sctpd_event_handler.h"
#include "lte/gateway/c/sctpd/src/sctpd_uplink_client.h"
#include "lte/gateway/c/sctpd/src/util.h"
#include "orc8r/gateway/c/common/logging/magma_logging_init.h"

using grpc::Server;
using grpc::ServerBuilder;
using magma::sctpd::SctpdDownlinkImpl;
using magma::sctpd::SctpdEventHandler;
using magma::sctpd::SctpdUplinkClient;

int signalMask(void) {
  sigset_t set;
  sigemptyset(&set);
  sigaddset(&set, SIGSEGV);
  sigaddset(&set, SIGINT);
  sigaddset(&set, SIGTERM);

  if (sigprocmask(SIG_BLOCK, &set, NULL) < 0) {
    return -1;
  }
  return 0;
}

int signalHandler(int* end, std::unique_ptr<Server>& server,
                  SctpdDownlinkImpl& downLink) {
  int ret;
  siginfo_t info;
  sigset_t set;

  sigemptyset(&set);
  sigaddset(&set, SIGSEGV);
  sigaddset(&set, SIGINT);
  sigaddset(&set, SIGTERM);

  if (sigprocmask(SIG_BLOCK, &set, NULL) < 0) {
    perror("sigprocmask");
    return -1;
  }

  /*
   * Block till a signal is received.
   * NOTE: The signals defined by set are required to be blocked at the time
   * of the call to sigwait() otherwise sigwait() is not successful.
   */
  if ((ret = sigwaitinfo(&set, &info)) == -1) {
    perror("sigwait");
    return ret;
  }

  server->Shutdown();
  server->Wait();
  downLink.stop();
  *end = 1;
  return 0;
}

int main() {
  signalMask();

  magma::init_logging("sctpd");
  magma::set_verbosity(MDEBUG);

  auto channel =
      grpc::CreateChannel(UPSTREAM_SOCK, grpc::InsecureChannelCredentials());

  SctpdUplinkClient client(channel);
  SctpdEventHandler handler(client);
  SctpdDownlinkImpl service(handler);

  ServerBuilder builder;
  builder.AddListeningPort(DOWNSTREAM_SOCK, grpc::InsecureServerCredentials());
  builder.RegisterService(&service);

  std::unique_ptr<Server> sctpd_dl_server = builder.BuildAndStart();

  int end = 0;
  while (end == 0) {
    signalHandler(&end, sctpd_dl_server, service);
  }
  return 0;
}
