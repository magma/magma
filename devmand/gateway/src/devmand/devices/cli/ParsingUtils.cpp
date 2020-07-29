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

#include <devmand/devices/cli/ParsingUtils.h>
#include <folly/Conv.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace std;

function<ydk::uint64(string)> toUI64 = [](auto s) { return stoull(s); };
function<ydk::uint16(string)> toUI16 = [](auto s) { return folly::to<int>(s); };

folly::Optional<string> extractValue(
    const string& output,
    const regex& pattern,
    const uint& groupToExtract) {
  std::stringstream ss(output);
  std::string line;

  while (std::getline(ss, line, '\n')) {
    boost::algorithm::trim(line);
    std::smatch match;
    if (std::regex_match(line, match, pattern) and
        match.size() > groupToExtract and match[groupToExtract].length() > 0) {
      return folly::Optional<string>(match[groupToExtract]);
    }
  }

  return folly::Optional<string>();
}

void parseValue(
    const string& output,
    const regex& pattern,
    const uint& groupToExtract,
    const std::function<void(string)>& setter) {
  const folly::Optional<string>& optValue =
      extractValue(output, pattern, groupToExtract);
  if (optValue) {
    setter(optValue.value());
  }
}

} // namespace cli
} // namespace devices
} // namespace devmand
