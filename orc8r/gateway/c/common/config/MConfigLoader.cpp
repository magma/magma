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

#include "orc8r/gateway/c/common/config/includes/MConfigLoader.h"

#include <bits/exception.h>
#include <google/protobuf/stubs/status.h>    // for Status
#include <google/protobuf/util/json_util.h>  // for JsonStringToMessage
#include <nlohmann/json.hpp>
#include <cstdlib>  // for getenv
#include <fstream>  // IWYU pragma: keep

#include "orc8r/gateway/c/common/logging/magma_logging.h"  // for MLOG

namespace google {
namespace protobuf {
class Message;
}  // namespace protobuf
}  // namespace google

namespace {
static constexpr const char* DYNAMIC_MCONFIG_PATH =
    "/var/opt/magma/configs/gateway.mconfig";
static constexpr const char* CONFIG_DIR = "/etc/magma";
static constexpr const char* MCONFIG_FILE_NAME = "gateway.mconfig";

bool check_file_exists(const std::string filename) {
  std::ifstream f(filename.c_str());
  return f.is_open();
}

void open_mconfig_file(std::ifstream* file) {
  // Load from /var/opt/magma if config exists, else read from /etc/magma
  if (check_file_exists(DYNAMIC_MCONFIG_PATH)) {
    file->open(DYNAMIC_MCONFIG_PATH);
    return;
  }
  const char* cfg_dir = std::getenv("MAGMA_CONFIG_LOCATION");
  if (cfg_dir == nullptr) {
    cfg_dir = CONFIG_DIR;
  }
  auto file_path = std::string(cfg_dir) + "/" + std::string(MCONFIG_FILE_NAME);
  file->open(file_path.c_str());
  return;
}

}  // namespace

namespace magma {

using json = nlohmann::json;

bool load_service_mconfig_from_file(const std::string& service_name,
                                    google::protobuf::Message* message) {
  // TODO(smoeller): Should use deffered file.close() here, e.g. absl::Cleanup
  std::ifstream file;
  open_mconfig_file(&file);
  if (!file.is_open()) {
    MLOG(MERROR) << "Couldn't load mconfig file";
    return false;
  }
  bool success = load_service_mconfig(service_name, &file, message);
  file.close();
  return success;
}

bool load_service_mconfig(const std::string& service_name,
                          std::istream* config_stream,
                          google::protobuf::Message* message) {
  json mconfig_json;
  try {
    *config_stream >> mconfig_json;
  } catch (const std::exception& e) {
    MLOG(MERROR) << "Parsing failure of config stream " << e.what();
    return false;
  }

  // config is located at mconfig_json["configs_by_key"][service_name]
  auto configs_it = mconfig_json.find("configs_by_key");
  if (configs_it == mconfig_json.end()) {
    configs_it = mconfig_json.find("configsByKey");
    if (configs_it == mconfig_json.end()) {
      MLOG(MERROR) << "Could not find configs_by_key in mconfig";
      return false;
    }
  }

  // Check if service exists
  auto service_it = configs_it->find(service_name);
  if (service_it == configs_it->end()) {
    MLOG(MERROR) << "Couldn't find " << service_name << " config";
    return false;
  }
  service_it->erase("@type");  // @type param makes parsing fail

  // Parse to message and return
  auto status =
      google::protobuf::util::JsonStringToMessage(service_it->dump(), message);
  if (!status.ok()) {
    MLOG(MERROR) << "Couldn't parse " << service_name
                 << " config, error: " << status.ToString();
  }
  return status.ok();
}

}  // namespace magma
