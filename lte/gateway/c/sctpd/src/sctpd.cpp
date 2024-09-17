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

#include "lte/gateway/c/sctpd/src/sctpd.hpp"

#include <systemd/sd-daemon.h>
#include <bits/types/siginfo_t.h>
#include <glog/logging.h>
#include <grpcpp/grpcpp.h>
#include <grpcpp/security/credentials.h>
#include <grpcpp/security/server_credentials.h>
#include <lte/protos/mconfig/mconfigs.pb.h>
#include <orc8r/protos/common.pb.h>
#include <signal.h>
#include <stdint.h>
#include <stdio.h>
#include <yaml-cpp/yaml.h>
#include <memory>
#include <ostream>
#include <string>

#include "lte/gateway/c/sctpd/src/sctpd_downlink_impl.hpp"
#include "lte/gateway/c/sctpd/src/sctpd_event_handler.hpp"
#include "lte/gateway/c/sctpd/src/sctpd_uplink_client.hpp"
#include "orc8r/gateway/c/common/config/MConfigLoader.hpp"
#include "orc8r/gateway/c/common/config/ServiceConfigLoader.hpp"
#include "orc8r/gateway/c/common/logging/magma_logging.hpp"
#include "orc8r/gateway/c/common/logging/magma_logging_init.hpp"
#include "orc8r/gateway/c/common/sentry/SentryWrapper.hpp"

#define SCTPD_SERVICE "sctpd"
#define SHARED_MCONFIG "shared_mconfig"

using grpc::Server;
using grpc::ServerBuilder;
using magma::sctpd::SctpdDownlinkImpl;
using magma::sctpd::SctpdEventHandler;
using magma::sctpd::SctpdUplinkClient;

int signalMask(void) {
  sigset_t set;
  sigemptyset(&set);
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
  sigaddset(&set, SIGINT);
  sigaddset(&set, SIGTERM);

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

static magma::mconfig::SctpD load_sctpd_mconfig() {
  magma::mconfig::SctpD mconfig;
  if (!magma::load_service_mconfig_from_file(SCTPD_SERVICE, &mconfig)) {
    mconfig.set_log_level(magma::orc8r::LogLevel::INFO);
  }
  return mconfig;
}

static uint32_t get_log_verbosity(const YAML::Node& config,
                                  magma::mconfig::SctpD mconfig) {
  if (!config["log_level"].IsDefined()) {
    return magma::get_log_verbosity_from_mconfig(mconfig.log_level());
  }
  std::string log_level = config["log_level"].as<std::string>();
  if (log_level == "DEBUG") {
    return MDEBUG;
  } else if (log_level == "INFO") {
    return MINFO;
  } else if (log_level == "WARNING") {
    return MWARNING;
  } else if (log_level == "ERROR") {
    return MERROR;
  } else if (log_level == "FATAL") {
    return MFATAL;
  } else {
    MLOG(MINFO) << "Invalid log level in config: "
                << config["log_level"].as<std::string>();
    return MINFO;
  }
}

int main() {
  signalMask();

  auto sctpd_mconfig = load_sctpd_mconfig();
  auto config = magma::ServiceConfigLoader{}.load_service_config(SCTPD_SERVICE);
  magma::init_logging(SCTPD_SERVICE);
  magma::set_verbosity(get_log_verbosity(config, sctpd_mconfig));

  sentry_config_t sentry_config = construct_sentry_config_from_mconfig();
  initialize_sentry(SENTRY_TAG_SCTPD, &sentry_config);

  auto channel =
      grpc::CreateChannel(UPSTREAM_SOCK, grpc::InsecureChannelCredentials());

  SctpdUplinkClient client(channel);
  SctpdEventHandler handler(client);
  SctpdDownlinkImpl service(handler);

  ServerBuilder builder;
  builder.AddListeningPort(DOWNSTREAM_SOCK, grpc::InsecureServerCredentials());
  builder.RegisterService(&service);

  std::unique_ptr<Server> sctpd_dl_server = builder.BuildAndStart();

  sd_notify(0, "READY=1");

  int end = 0;
  while (end == 0) {
    signalHandler(&end, sctpd_dl_server, service);
  }
  shutdown_sentry();
  return 0;
}
