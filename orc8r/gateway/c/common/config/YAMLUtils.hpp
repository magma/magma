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

namespace magma {

/**
 * YAMLUtils defines new utilities that the yaml-cpp library doesn't expose
 */
class YAMLUtils final {
 public:
  /**
   * merge_nodes combines two YAML files together. override_node will
   * override any parameters it defines, and keep any existing parameters in
   * default_node that it doesn't define
   */
  static YAML::Node merge_nodes(const YAML::Node& default_node,
                                const YAML::Node& override_node);
};

}  // namespace magma
