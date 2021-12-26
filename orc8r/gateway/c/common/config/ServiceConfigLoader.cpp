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

#include "orc8r/gateway/c/common/config/includes/ServiceConfigLoader.h"

#include <iostream>  // for operator<<, basic_ostream
#include <string>    // for allocator, operator+, char_traits

#include "orc8r/gateway/c/common/config/YAMLUtils.h"       // for YAMLUtils
#include "orc8r/gateway/c/common/logging/magma_logging.h"  // for MLOG

namespace magma {

YAML::Node ServiceConfigLoader::load_service_config(
    const std::string& service_name) {
  auto file_path = std::string(CONFIG_DIR) + service_name + ".yml";
  YAML::Node base_config = YAML::LoadFile(file_path);

  // Try to override original file, if an override exists
  try {
    auto override_file = std::string(OVERRIDE_DIR) + service_name + ".yml";
    return YAMLUtils::merge_nodes(base_config, YAML::LoadFile(override_file));
  } catch (YAML::BadFile&) {
    MLOG(MDEBUG) << "Override file not found for service " << service_name;
  }
  return base_config;
}

}  // namespace magma
