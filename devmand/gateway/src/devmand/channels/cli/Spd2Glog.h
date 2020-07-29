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

#include <spdlog/sinks/base_sink.h>
#include <mutex>

namespace devmand {
namespace channels {
namespace cli {

using namespace spdlog::level;

// Usage:
//  auto logger = spdlog::create<Spd2Glog>("spd_logger");
//  logger->info("log me");
class Spd2Glog : public spdlog::sinks::base_sink<std::mutex> {
 public:
  void toGlog(const spdlog::details::log_msg& msg);

 protected:
  void _sink_it(const spdlog::details::log_msg& msg) override;
  void flush() override;
};
} // namespace cli
} // namespace channels
} // namespace devmand
