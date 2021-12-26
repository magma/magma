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
#include <google/protobuf/util/time_util.h>
#include <time.h>
#include <algorithm>
#include <chrono>
#include <iomanip>
#include <sstream>

#include "Utilities.h"

namespace google {
namespace protobuf {
class Timestamp;
}  // namespace protobuf
}  // namespace google

namespace magma {

std::string bytes_to_hex(const std::string& s) {
  std::ostringstream ret;
  unsigned int c;
  for (std::string::size_type i = 0; i < s.length(); ++i) {
    c = (unsigned int)(unsigned char)s[i];
    ret << " " << std::hex << std::setfill('0') << std::setw(2)
        << (std::nouppercase) << c;
  }
  return ret.str();
}

uint64_t get_time_in_sec_since_epoch() {
  auto now = std::chrono::system_clock::now();
  return std::chrono::duration_cast<std::chrono::seconds>(
             now.time_since_epoch())
      .count();
}

std::chrono::milliseconds time_difference_from_now(
    const google::protobuf::Timestamp& timestamp) {
  const auto rule_time_sec =
      google::protobuf::util::TimeUtil::TimestampToSeconds(timestamp);
  const auto now = time(NULL);
  const auto delta = std::max(rule_time_sec - now, 0L);
  std::chrono::seconds sec(delta);
  return std::chrono::duration_cast<std::chrono::milliseconds>(sec);
}

std::chrono::milliseconds time_difference_from_now(
    const std::time_t timestamp) {
  const auto now = time(nullptr);
  const auto delta = std::max(timestamp - now, 0L);
  std::chrono::seconds sec(delta);
  return std::chrono::duration_cast<std::chrono::milliseconds>(sec);
}

}  // namespace magma
