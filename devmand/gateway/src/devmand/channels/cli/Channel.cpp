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

#include <devmand/channels/cli/Channel.h>

namespace devmand {
namespace channels {
namespace cli {

using folly::SemiFuture;
using folly::Unit;
using folly::unit;

Channel::Channel(
    string _id,
    const std::shared_ptr<devmand::channels::cli::Cli> _cli)
    : id(_id), cli(_cli) {}

SemiFuture<Unit> Channel::destroy() {
  // idempotency
  if (cli == nullptr) {
    return folly::makeSemiFuture(unit);
  }
  // call underlying destroy()
  SemiFuture<Unit> innerDestroy = cli->destroy();
  return innerDestroy;
}

Channel::~Channel() {
  MLOG(MDEBUG) << "[" << id << "] "
               << "~Channel: started";
  destroy().get();
  cli = nullptr;
  MLOG(MDEBUG) << "[" << id << "] "
               << "~Channel: finished";
}

folly::SemiFuture<std::string> Channel::executeRead(const ReadCommand cmd) {
  MLOG(MDEBUG) << "[" << id << "] "
               << "Executing command and reading: "
               << "\"" << cmd << "\"";

  return cli->executeRead(cmd);
}

folly::SemiFuture<std::string> Channel::executeWrite(const WriteCommand cmd) {
  MLOG(MDEBUG) << "[" << id << "]"
               << "Executing command"
               << "\"" << cmd << "\"";

  return cli->executeWrite(cmd);
}

} // namespace cli
} // namespace channels
} // namespace devmand
