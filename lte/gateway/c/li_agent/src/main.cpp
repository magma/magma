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
#include <thread>

#include "orc8r/gateway/c/common/service303/includes/MagmaService.h"
#include "orc8r/gateway/c/common/config/includes/MConfigLoader.h"
#include "orc8r/gateway/c/common/service_registry/includes/ServiceRegistrySingleton.h"

#include "lte/gateway/c/li_agent/src/InterfaceMonitor.h"
#include "lte/gateway/c/li_agent/src/PDUGenerator.h"
#include "lte/gateway/c/li_agent/src/ProxyConnector.h"
#include "lte/gateway/c/li_agent/src/Utilities.h"
#include "orc8r/gateway/c/common/logging/magma_logging_init.h"

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

  auto mconfig = magma::lte::load_mconfig();
  auto config  = magma::ServiceConfigLoader{}.load_service_config(LIAGENTD);
  magma::set_verbosity(get_log_verbosity(config, mconfig));

  // Ignoring SIGPIPE due to ssl write throwing it occasionally
  sigset_t blockedSignal;
  sigemptyset(&blockedSignal);
  sigaddset(&blockedSignal, SIGPIPE);
  pthread_sigmask(SIG_BLOCK, &blockedSignal, NULL);

  bool enable = config["enable"].as<bool>();
  if (!enable) {
    MLOG(MINFO) << "LI Agent service disabled";
    return 0;
  }

  MLOG(MINFO) << "Starting LI Agent service " << config;

  std::string interface_name = config["interface_name"].as<std::string>();
  std::string pkt_dst_mac    = config["pkt_dst_mac"].as<std::string>();
  std::string pkt_src_mac    = config["pkt_src_mac"].as<std::string>();
  std::string proxy_addr     = config["proxy_addr"].as<std::string>();
  std::string cert_file      = config["cert_file"].as<std::string>();
  std::string key_file       = config["key_file"].as<std::string>();
  int proxy_port             = config["proxy_port"].as<int>();
  int sync_interval          = config["sync_interval"].as<int>();
  int inactivity_time        = config["inactivity_time"].as<int>();

  auto mobilityd_client = std::make_unique<magma::lte::AsyncMobilitydClient>();
  std::thread mobilitydd_response_handling_thread([&]() {
    MLOG(MINFO) << "Started MobilityD response thread";
    mobilityd_client->rpc_response_loop();
  });

  magma::service303::MagmaService server(LIAGENTD, LIAGENTD_VERSION);
  server.Start();

  auto proxy_connector = std::make_unique<magma::lte::ProxyConnectorImpl>(
      proxy_addr, proxy_port, cert_file, key_file);
  if (proxy_connector->setup_proxy_socket() < 0) {
    MLOG(MERROR) << "Coudn't setup proxy socket, terminating";
    return -1;
  }

  auto pkt_generator = std::make_unique<magma::lte::PDUGenerator>(
      pkt_dst_mac, pkt_src_mac, sync_interval, inactivity_time,
      std::move(proxy_connector), std::move(mobilityd_client), mconfig);

  auto interface_watcher = std::make_unique<magma::lte::InterfaceMonitor>(
      interface_name, std::move(pkt_generator));
  if (interface_watcher->init_interface_monitor() < 0) {
    MLOG(MERROR) << "Coudn't setup interface sniffing, terminating";
    return -1;
  }

  if (interface_watcher->start_capture() < 0) {
    MLOG(MERROR) << "Coudn't start interface sniffing, terminating";
    return -1;
  }

  return 0;
}
