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

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/channels/Channel.h>
#include <devmand/channels/cli/Cli.h>

namespace devmand {
namespace channels {
namespace cli {

class Channel : public channels::Channel, public devmand::channels::cli::Cli {
 public:
  Channel(string _id, const std::shared_ptr<devmand::channels::cli::Cli> _cli);
  Channel() = delete;
  virtual ~Channel();
  Channel(const Channel&) = delete;
  Channel& operator=(const Channel&) = delete;
  Channel(Channel&&) = delete;
  Channel& operator=(Channel&&) = delete;

  folly::SemiFuture<std::string> executeRead(const ReadCommand cmd) override;
  folly::SemiFuture<std::string> executeWrite(const WriteCommand cmd) override;

  folly::SemiFuture<folly::Unit> destroy() override;

 private:
  string id;
  std::shared_ptr<devmand::channels::cli::Cli> cli;
};

} // namespace cli
} // namespace channels
} // namespace devmand
