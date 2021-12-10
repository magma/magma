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
#pragma once

#include <yaml-cpp/node/node.h>  // for Node
#include <yaml-cpp/yaml.h>
#include <string>  // for string

namespace magma {

/**
 * ServiceConfigLoader is a helper class to parse proc files for process
 * information
 */
class ServiceConfigLoader final {
 public:
  /*
   * Load service configuration from file.
   *
   * @return YAML::Node a Node representation of the file.
   */
  YAML::Node load_service_config(const std::string& service_name);

 private:
  static constexpr const char* CONFIG_DIR = "/etc/magma/";
  static constexpr const char* OVERRIDE_DIR = "/var/opt/magma/configs/";
};

}  // namespace magma
