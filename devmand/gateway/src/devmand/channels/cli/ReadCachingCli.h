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
#include <folly/Synchronized.h>
#include <folly/container/EvictingCacheMap.h>

namespace devmand::channels::cli {

using folly::EvictingCacheMap;
using folly::Future;
using folly::SemiFuture;
using folly::Synchronized;
using std::shared_ptr;
using std::string;
using CliCache = Synchronized<EvictingCacheMap<string, string>>;

class ReadCachingCli : public Cli {
 private:
  string id;
  shared_ptr<Cli> cli{};
  shared_ptr<CliCache> cache;
  shared_ptr<folly::Executor> executor;

 public:
  ReadCachingCli(
      string _id,
      const shared_ptr<Cli>& _cli,
      const shared_ptr<CliCache>& _cache,
      const shared_ptr<folly::Executor> executor);

  SemiFuture<folly::Unit> destroy() override;

  ~ReadCachingCli() override;

  static shared_ptr<CliCache> createCache();
  folly::SemiFuture<std::string> executeRead(const ReadCommand cmd) override;
  folly::SemiFuture<std::string> executeWrite(const WriteCommand cmd) override;
};
} // namespace devmand::channels::cli
