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
#include <gtest/gtest.h>
#include "yaml-cpp/yaml.h"

#include "YAMLUtils.h"

using ::testing::Test;
using YAML::Node;

namespace magma {

TEST(test_simple_field_overrides, test_yaml_utils) {
  Node default_node;
  default_node["foo"]     = "bar";
  default_node["default"] = "default";

  Node override_node;
  override_node["foo"]      = "barbar";
  override_node["override"] = "override";

  auto merge = YAMLUtils::merge_nodes(default_node, override_node);
  EXPECT_EQ("barbar", merge["foo"].as<std::string>());
  EXPECT_EQ("default", merge["default"].as<std::string>());
  EXPECT_EQ("override", merge["override"].as<std::string>());
}

TEST(test_override_with_nothing, test_yaml_utils) {
  Node default_node;
  default_node["foo"] = "bar";
  Node override_node;
  auto merge = YAMLUtils::merge_nodes(default_node, override_node);
  EXPECT_EQ("bar", merge["foo"].as<std::string>());
}

TEST(test_nested_overrides, test_yaml_utils) {
  Node default_node;
  default_node["foo"] = "bar";
  Node nested_node;
  nested_node["a"]     = "b";
  nested_node["c"]     = "d";
  default_node["nest"] = nested_node;

  Node override_node;
  Node nest_override;
  Node new_nest;
  nest_override["a"]    = "e";
  new_nest["new_nest"]  = "other";
  override_node["nest"] = nest_override;
  override_node["foo"]  = new_nest;

  auto merge = YAMLUtils::merge_nodes(default_node, override_node);
  EXPECT_TRUE(merge["foo"].IsMap());
  EXPECT_EQ("other", merge["foo"]["new_nest"].as<std::string>());
  EXPECT_EQ("e", merge["nest"]["a"].as<std::string>());
  EXPECT_EQ("d", merge["nest"]["c"].as<std::string>());
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
}  // namespace magma
