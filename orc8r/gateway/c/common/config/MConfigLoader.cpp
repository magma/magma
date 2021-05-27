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

#include "includes/MConfigLoader.h"
#include <google/protobuf/stubs/status.h>    // for Status
#include <google/protobuf/util/json_util.h>  // for JsonStringToMessage
#include <cstdlib>                           // for getenv
#include <fstream>                           // for operator<<, char_traits
#include <json.hpp>                          // for basic_json<>::iterator
#include "magma_logging.h"                   // for MLOG
namespace google {
namespace protobuf {
class Message;
}
}  // namespace google

using json = nlohmann::json;

namespace magma {

static bool check_file_exists(const std::string filename) {
  std::ifstream f(filename.c_str());
  return f.is_open();
}

bool MConfigLoader::load_service_mconfig(
    const std::string& service_name, google::protobuf::Message* message) {
  std::ifstream file;
  get_mconfig_file(&file);
  if (!file.is_open()) {
    MLOG(MERROR) << "Couldn't load mconfig file";
    return false;
  }

  json mconfig_json;
  file >> mconfig_json;
  file.close();

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
    MLOG(MERROR) << "Couldn't parse " << service_name << " config";
  }
  return status.ok();
}

void MConfigLoader::get_mconfig_file(std::ifstream* file) {
  // Load from /var/opt/magma if config exists, else read from /etc/magma
  if (check_file_exists(MConfigLoader::DYNAMIC_MCONFIG_PATH)) {
    file->open(MConfigLoader::DYNAMIC_MCONFIG_PATH);
    return;
  }
  const char* cfg_dir = std::getenv("MAGMA_CONFIG_LOCATION");
  if (cfg_dir == nullptr) {
    cfg_dir = MConfigLoader::CONFIG_DIR;
  }
  auto file_path = std::string(cfg_dir) + "/" +
                   std::string(MConfigLoader::MCONFIG_FILE_NAME);
  file->open(file_path.c_str());
  return;
}

}  // namespace magma
