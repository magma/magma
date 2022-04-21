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

#include <glog/logging.h>
#include <lte/protos/mconfig/mconfigs.pb.h>
#include <orc8r/protos/common.pb.h>
#include <stdint.h>
#include <yaml-cpp/yaml.h>
#include <memory>
#include <ostream>
#include <string>

#include "lte/gateway/c/connection_tracker/src/EventTracker.hpp"
#include "lte/gateway/c/connection_tracker/src/PacketGenerator.hpp"
#include "orc8r/gateway/c/common/config/includes/MConfigLoader.hpp"
#include "orc8r/gateway/c/common/config/includes/ServiceConfigLoader.hpp"
#include "orc8r/gateway/c/common/logging/magma_logging.h"
#include "orc8r/gateway/c/common/logging/magma_logging_init.h"
#include "orc8r/gateway/c/common/sentry/includes/SentryWrapper.hpp"
#include "orc8r/gateway/c/common/service303/includes/MagmaService.hpp"

#define CONNECTION_SERVICE "connectiond"
#define CONNECTIOND_VERSION "1.0"

static magma::mconfig::ConnectionD get_default_mconfig() {
  magma::mconfig::ConnectionD mconfig;
  mconfig.set_log_level(magma::orc8r::LogLevel::INFO);
  return mconfig;
}

static magma::mconfig::ConnectionD load_mconfig() {
  magma::mconfig::ConnectionD mconfig;
  if (!magma::load_service_mconfig_from_file(CONNECTION_SERVICE, &mconfig)) {
    MLOG(MERROR) << "Unable to load mconfig for connectiond, using default";
    return get_default_mconfig();
  }
  return mconfig;
}

static uint32_t get_log_verbosity(const YAML::Node& config,
                                  magma::mconfig::ConnectionD mconfig) {
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

int main(void) {
  sentry_config_t sentry_config = construct_sentry_config_from_mconfig();
  initialize_sentry(SENTRY_TAG_CONNECTIOND, &sentry_config);
  magma::init_logging(CONNECTION_SERVICE);

  auto mconfig = load_mconfig();
  auto config =
      magma::ServiceConfigLoader{}.load_service_config(CONNECTION_SERVICE);
  magma::set_verbosity(get_log_verbosity(config, mconfig));
  MLOG(MINFO) << "Starting Connection Tracker";

  std::string interface_name = config["interface_name"].as<std::string>();
  std::string pkt_dst_mac = config["pkt_dst_mac"].as<std::string>();
  std::string pkt_src_mac = config["pkt_src_mac"].as<std::string>();
  int zone = config["zone"].as<int>();

  magma::service303::MagmaService server(CONNECTION_SERVICE,
                                         CONNECTIOND_VERSION);
  server.Start();

  auto pkt_generator = std::make_shared<magma::lte::PacketGenerator>(
      interface_name, pkt_dst_mac, pkt_src_mac);

  auto event_tracker =
      std::make_shared<magma::lte::EventTracker>(pkt_generator, zone);

  event_tracker->init_conntrack_event_loop();

  shutdown_sentry();
  return 0;
}
