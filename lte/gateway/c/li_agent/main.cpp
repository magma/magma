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

#include <stdio.h>
#include <stdlib.h>
#include <lte/protos/mconfig/mconfigs.pb.h>
#include <thread>

#include "MagmaService.h"
#include "MConfigLoader.h"
#include "ServiceRegistrySingleton.h"

#include "InterfaceMonitor.h"
#include "PDUGenerator.h"
#include "ProxyConnector.h"
#include "magma_logging_init.h"

#define LIAGENTD "liagentd"
#define LIAGENTD_VERSION "1.0"

static magma::mconfig::LIAgentD get_default_mconfig() {
  magma::mconfig::LIAgentD mconfig;
  mconfig.set_log_level(magma::orc8r::LogLevel::INFO);
  return mconfig;
}

static magma::mconfig::LIAgentD load_mconfig() {
  magma::mconfig::LIAgentD mconfig;
  magma::MConfigLoader loader;
  if (!loader.load_service_mconfig(LIAGENTD, &mconfig)) {
    MLOG(MERROR) << "Unable to load mconfig for liagentd, using default";
    return get_default_mconfig();
  }
  return mconfig;
}

static uint32_t get_log_verbosity(
    const YAML::Node& config, magma::mconfig::LIAgentD mconfig) {
  if (!config["log_level"].IsDefined()) {
    if (mconfig.log_level() < 0 || mconfig.log_level() > 4) {
      return MINFO;
    }
    return mconfig.log_level();
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

int main(void) {
  magma::init_logging(LIAGENTD);

  auto mconfig = load_mconfig();
  auto config  = magma::ServiceConfigLoader{}.load_service_config(LIAGENTD);
  magma::set_verbosity(get_log_verbosity(config, mconfig));
  MLOG(MINFO) << "Starting LI Agent service";

  // Ignoring SIGPIPE due to ssl write throwing it occasionally
  sigset_t blockedSignal;
  sigemptyset(&blockedSignal);
  sigaddset(&blockedSignal, SIGPIPE);
  pthread_sigmask(SIG_BLOCK, &blockedSignal, NULL);

  auto directoryd_client = std::make_unique<magma::AsyncDirectorydClient>();
  std::thread directoryd_response_handling_thread([&]() {
    MLOG(MINFO) << "Started DirectoryD response thread";
    directoryd_client->rpc_response_loop();
  });

  std::string interface_name = config["interface_name"].as<std::string>();
  std::string pkt_dst_mac    = config["pkt_dst_mac"].as<std::string>();
  std::string pkt_src_mac    = config["pkt_src_mac"].as<std::string>();
  std::string proxy_addr     = config["proxy_addr"].as<std::string>();
  int proxy_port             = config["proxy_port"].as<int>();
  std::string cert_file      = config["cert_file"].as<std::string>();
  std::string key_file       = config["key_file"].as<std::string>();

  magma::service303::MagmaService server(LIAGENTD, LIAGENTD_VERSION);
  server.Start();

  auto proxy_connector = std::make_unique<magma::ProxyConnectorImpl>(
      proxy_addr, proxy_port, cert_file, key_file);
  if (proxy_connector->setup_proxy_socket() < 0) {
    MLOG(MERROR) << "Coudn't setup proxy socket, terminating";
    return -1;
  }
  auto pkt_generator = std::make_unique<magma::PDUGenerator>(
      std::move(proxy_connector), std::move(directoryd_client), pkt_dst_mac,
      pkt_src_mac);

  auto interface_watcher = std::make_unique<magma::InterfaceMonitor>(
      interface_name, std::move(pkt_generator));
  if (interface_watcher->init_iface_pcap_monitor() < 0) {
    MLOG(MERROR) << "Coudn't setup interface sniffing, terminating";
    return -1;
  }

  proxy_connector->cleanup();
  directoryd_response_handling_thread.join();

  return 0;
}
