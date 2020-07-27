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

#include <devmand/utils/Time.h>

#include <sstream>

namespace devmand {
namespace utils {
namespace Time {

TimePoint now() {
  return Clock::now();
}

TimePoint from(struct timeval tv) {
  return TimePoint{std::chrono::seconds{tv.tv_sec} +
                   std::chrono::microseconds{tv.tv_usec}};
}

std::string toString(const TimePoint& time) {
  std::stringstream stream;
  stream << time;
  return stream.str();
}

} // namespace Time
} // namespace utils
} // namespace devmand

std::ostream& operator<<(
    std::ostream& stream,
    const devmand::utils::TimePoint& time) {
  std::time_t epoch =
      std::chrono::duration_cast<std::chrono::seconds>(time.time_since_epoch())
          .count();
  return stream << std::put_time(std::localtime(&epoch), "%c");
}
