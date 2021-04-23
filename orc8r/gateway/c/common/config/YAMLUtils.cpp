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

#include "YAMLUtils.h"
#include <yaml-cpp/yaml.h>                     // IWYU pragma: keep
#include <boost/iterator/iterator_facade.hpp>  // for operator!=, iterator_f...
#include <string>                              // for string

namespace magma {

static const YAML::Node& cnode(const YAML::Node& node) {
  return node;
}

static void fill_in_new_fields(
    const YAML::Node& override_node, YAML::Node* merged_map) {
  for (auto new_node : override_node) {
    if (!new_node.first.IsScalar() ||
        !cnode(*merged_map)[new_node.first.Scalar()]) {
      (*merged_map)[new_node.first] = new_node.second;
    }
  }
}

static YAML::Node merge_maps(
    const YAML::Node& default_node, const YAML::Node& override_node) {
  auto merged_map = YAML::Node(YAML::NodeType::Map);
  for (auto base_child : default_node) {
    if (base_child.first.IsScalar()) {
      const std::string& key = base_child.first.Scalar();
      auto override_child    = YAML::Node(cnode(override_node)[key]);
      if (override_child) {
        merged_map[base_child.first] =
            YAMLUtils::merge_nodes(base_child.second, override_child);
      } else {
        merged_map[base_child.first] = base_child.second;
      }
    } else {
      merged_map[base_child.first] = base_child.second;
    }
  }
  // Add the mappings from override_node not already in 'merged_map'
  fill_in_new_fields(override_node, &merged_map);
  return merged_map;
}

/**
 * merge_nodes recursively merges YAML nodes of any type. The following logic
 * is implemented:
 * If override_node is not a map, merge result is override_node,
 * unless override_node is null
 * If default_node is not a map, merge result is override_node, which is a map
 * If default_node is a map, and override_node is an empty map, return
 * default_node
 * Otherwise, merge the two maps
 */
YAML::Node YAMLUtils::merge_nodes(
    const YAML::Node& default_node, const YAML::Node& override_node) {
  if (!override_node.IsMap()) {
    return override_node.IsNull() ? default_node : override_node;
  }
  if (!default_node.IsMap()) {
    return override_node;
  }
  if (!override_node.size()) {
    return default_node;
  }
  return merge_maps(default_node, override_node);
}

}  // namespace magma
