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

#include <lte/protos/mconfig/mconfigs.pb.h>
#include <orc8r/protos/mconfig/mconfigs.pb.h>

#include <memory>
#include <grpcpp/grpcpp.h>
#include <signal.h>

#include "lte/gateway/c/sctpd/src/sctpd_downlink_impl.h"
#include "lte/gateway/c/sctpd/src/sctpd_event_handler.h"
#include "lte/gateway/c/sctpd/src/sctpd_uplink_client.h"
#include "lte/gateway/c/sctpd/src/util.h"
#include "orc8r/gateway/c/common/config/includes/MConfigLoader.h"
#include "orc8r/gateway/c/common/logging/magma_logging_init.h"
#include "orc8r/gateway/c/common/sentry/includes/SentryWrapper.h"
#include "orc8r/gateway/c/common/service_registry/includes/ServiceRegistrySingleton.h"

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

static magma::mconfig::SctpD load_sctpd_mconfig() {
  magma::mconfig::SctpD mconfig;
  if (!magma::load_service_mconfig_from_file(SCTPD_SERVICE, &mconfig)) {
    mconfig.set_log_level(magma::orc8r::LogLevel::INFO);
  }
  return mconfig;
}

static magma::mconfig::SharedMconfig load_shared_mconfig() {
  magma::mconfig::SharedMconfig mconfig;
  magma::load_service_mconfig_from_file(SHARED_MCONFIG, &mconfig);
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

  auto sentry_mconfig = load_shared_mconfig().sentry_config();
  sentry_config_t sentry_config;
  sentry_config.sample_rate = sentry_mconfig.sample_rate();
  strncpy(sentry_config.url_native, sentry_mconfig.dsn_native().c_str(),
          MAX_URL_LENGTH - 1);
  sentry_config.url_native[MAX_URL_LENGTH - 1] = '\0';
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

  int end = 0;
  while (end == 0) {
    signalHandler(&end, sctpd_dl_server, service);
  }
  return 0;
}
