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

#include <chrono>
#include <ctime>
#include <iomanip>
#include <ostream>

namespace devmand {

namespace utils {

using Clock = std::chrono::steady_clock;
using SteadyTimePoint = std::chrono::steady_clock::time_point;
using SystemTimePoint = std::chrono::system_clock::time_point;
using TimePoint = SteadyTimePoint;

namespace Time {

extern TimePoint now();

extern TimePoint from(struct timeval tv);

extern std::string toString(const TimePoint& time);

} // namespace Time

} // namespace utils

} // namespace devmand

extern std::ostream& operator<<(
    std::ostream& stream,
    const devmand::utils::TimePoint& time);
