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

#pragma once

#include <algorithm>
#include <iomanip>
#include <sstream>
#include <string>

#if __cplusplus <= 201703L
// This function will be provided by c++20 but here until then. Wrapped to
// provide compiler errors when we switch.
static inline bool startsWith(
    const std::string& str,
    const std::string& prefix) {
  return (
      (prefix.size() <= str.size()) and
      std::equal(prefix.begin(), prefix.end(), str.begin()));
}

static inline bool endsWith(const std::string& str, const std::string& suffix) {
  return str.size() >= suffix.size() and
      0 == str.compare(str.size() - suffix.size(), suffix.size(), suffix);
}
#endif

namespace devmand {

class StringUtils final {
 public:
  StringUtils() = delete;
  ~StringUtils() = delete;
  StringUtils(const StringUtils&) = delete;
  StringUtils& operator=(const StringUtils&) = delete;
  StringUtils(StringUtils&&) = delete;
  StringUtils& operator=(StringUtils&&) = delete;

 public:
  static inline bool iequals(const std::string& rhs, const std::string& lhs) {
    return std::equal(
        rhs.begin(),
        rhs.end(),
        lhs.begin(),
        lhs.end(),
        [](char crhs, char clhs) { return tolower(crhs) == tolower(clhs); });
  }

  static inline std::string unquote(const std::string& in) {
    std::stringstream ss;
    std::string out;
    ss << in;
    ss >> std::quoted(out);
    return out;
  }

  static inline std::string asHexString(
      const std::string& buf,
      const std::string& delim) {
    std::stringstream msg;
    msg << std::hex;
    for (char c : buf) {
      msg << static_cast<int>(c & 0xFF) << delim;
    }
    msg << std::dec;

    std::string ret{msg.str()};
    if (not buf.empty()) {
      for (unsigned int i = 0; i < delim.size(); ++i) {
        ret.pop_back();
      }
    }
    return ret;
  }
};

} // namespace devmand
