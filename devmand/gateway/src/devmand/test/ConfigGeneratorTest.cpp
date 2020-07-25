/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#include <gtest/gtest.h>

#include <experimental/filesystem>
#include <list>
#include <set>

#include <devmand/utils/ConfigGenerator.h>
#include <devmand/utils/FileUtils.h>

namespace devmand {
namespace utils {
namespace test {

static std::experimental::filesystem::path getFilename() {
  std::string ret = "/tmp/";
  auto* testInfo = ::testing::UnitTest::GetInstance()->current_test_info();
  ret += testInfo->name();
  return ret;
}

TEST(ConfigGeneratorTest, singleAdd) {
  auto filename = getFilename();
  ConfigGenerator cg{filename, "s{}e"};
  cg.add("{}", "a");

  EXPECT_EQ("sae", FileUtils::readContents(filename));
}

TEST(ConfigGeneratorTest, doubleAdd) {
  auto filename = getFilename();
  ConfigGenerator cg{filename, "s{}e"};
  cg.add("{}", "a");
  cg.add("{}", "b");

  EXPECT_EQ("sabe", FileUtils::readContents(filename));
}

TEST(ConfigGeneratorTest, newlinesAdd) {
  auto filename = getFilename();
  ConfigGenerator cg{filename, "s\n{}\ne\n"};
  cg.add("{}", "a\nb");
  cg.add("{}", "c\nd");

  EXPECT_EQ("s\na\nbc\nd\ne\n", FileUtils::readContents(filename));
}

TEST(ConfigGeneratorTest, doubleAddThenRm) {
  auto filename = getFilename();
  ConfigGenerator cg{filename, "s{}e"};
  cg.add("{}", "a");
  cg.add("{}", "b");

  EXPECT_EQ("sabe", FileUtils::readContents(filename));

  cg.remove("{}", "b");

  EXPECT_EQ("sae", FileUtils::readContents(filename));
}

TEST(ConfigGeneratorTest, doubleAddOrder) {
  auto filename = getFilename();
  ConfigGenerator cg{filename, "s{}e"};
  cg.add("{}", "b");
  cg.add("{}", "a");
  cg.add("{}", "c");

  EXPECT_EQ("sabce", FileUtils::readContents(filename));
}

} // namespace test
} // namespace utils
} // namespace devmand
