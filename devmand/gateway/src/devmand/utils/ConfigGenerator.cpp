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

#include <algorithm>
#include <experimental/filesystem>

#include <folly/Format.h>
#include <folly/GLog.h>

#include <devmand/utils/ConfigGenerator.h>
#include <devmand/utils/FileUtils.h>

#include <boost/algorithm/string/join.hpp>

namespace devmand {

namespace utils {

ConfigGenerator::ConfigGenerator(
    const std::string& configFile_,
    const std::string& fileTemplate_)
    : configFile(configFile_), fileTemplate(fileTemplate_) {
  FileUtils::mkdir(
      std::experimental::filesystem::path(configFile).parent_path());
}

void ConfigGenerator::rewrite() {
  // TODO make this async
  FileUtils::write(
      configFile,
      folly::sformat(fileTemplate, boost::algorithm::join(entries, "")));
}

} // namespace utils

} // namespace devmand
