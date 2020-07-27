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

#include <devmand/channels/cli/Cli.h>
#include <devmand/channels/cli/CliFlavour.h>
#include <devmand/channels/cli/TreeCache.h>
#include <folly/Executor.h>
#include <folly/futures/Future.h>

namespace devmand::channels::cli {

using namespace std;
using namespace folly;

class TreeCacheCli : public Cli {
 private:
  string id;
  shared_ptr<Cli> cli;
  shared_ptr<folly::Executor> executor;
  shared_ptr<CliFlavour> sharedCliFlavour;
  shared_ptr<TreeCache> cache;

 public:
  TreeCacheCli(
      string _id,
      shared_ptr<Cli> _cli,
      shared_ptr<folly::Executor> executor,
      shared_ptr<CliFlavour> sharedCliFlavour);

  // Visible for testing
  TreeCacheCli(
      string _id,
      shared_ptr<Cli> _cli,
      shared_ptr<folly::Executor> _executor,
      shared_ptr<CliFlavour> _sharedCliFlavour,
      shared_ptr<TreeCache> cache);

  SemiFuture<folly::Unit> destroy() override;

  ~TreeCacheCli() override;

  /*
   * Clear cache.
   */
  void clear();

  folly::SemiFuture<std::string> executeRead(const ReadCommand cmd) override;
  folly::SemiFuture<std::string> executeWrite(const WriteCommand cmd) override;
};
} // namespace devmand::channels::cli
