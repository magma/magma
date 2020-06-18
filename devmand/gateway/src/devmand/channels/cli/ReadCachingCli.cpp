// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/ReadCachingCli.h>
#include <folly/Optional.h>
#include <folly/Synchronized.h>
#include <folly/container/EvictingCacheMap.h>
#include <magma_logging.h>

using devmand::channels::cli::Cli;
using devmand::channels::cli::Command;
using folly::EvictingCacheMap;
using folly::Future;
using folly::Optional;
using folly::Synchronized;
using std::shared_ptr;
using std::string;
using CliCache = Synchronized<EvictingCacheMap<string, string>>;

namespace devmand::channels::cli {

folly::SemiFuture<std::string> ReadCachingCli::executeRead(
    const ReadCommand cmd) {
  if (!cmd.skipCache()) {
    Optional<string> cachedResult =
        cache->withWLock([cmd, this](auto& cache_) -> Optional<string> {
          if (cache_.exists(cmd.raw())) {
            MLOG(MDEBUG) << "[" << id << "] "
                         << "Found command: " << cmd << " in cache";
            return Optional<string>(cache_.get(cmd.raw()));
          } else {
            return Optional<string>(folly::none);
          }
        });

    if (cachedResult) {
      return folly::SemiFuture<string>(*cachedResult.get_pointer());
    }
  }

  return cli->executeRead(cmd)
      .via(executor.get())
      .thenValue([=](string output) {
        cache->wlock()->insert(cmd.raw(), output);
        return output;
      })
      .semi();
}

ReadCachingCli::ReadCachingCli(
    string _id,
    const std::shared_ptr<Cli>& _cli,
    const shared_ptr<CliCache>& _cache,
    const shared_ptr<folly::Executor> _executor)

    : id(_id), cli(_cli), cache(_cache), executor(_executor) {}

folly::SemiFuture<std::string> ReadCachingCli::executeWrite(
    const WriteCommand cmd) {
  return cli->executeWrite(cmd);
}

SemiFuture<folly::Unit> ReadCachingCli::destroy() {
  MLOG(MDEBUG) << "[" << id << "] "
               << "destroy: started";
  // call underlying destroy()
  SemiFuture<folly::Unit> innerDestroy = cli->destroy();
  MLOG(MDEBUG) << "[" << id << "] "
               << "destroy: done";
  return innerDestroy;
}

ReadCachingCli::~ReadCachingCli() {
  MLOG(MDEBUG) << "[" << id << "] "
               << "~RCcli: started";
  destroy().get();
  executor = nullptr;
  cache = nullptr;
  cli = nullptr;
  MLOG(MDEBUG) << "[" << id << "] "
               << "~RCcli: done";
}

shared_ptr<CliCache> ReadCachingCli::createCache() {
  return std::make_shared<CliCache>(EvictingCacheMap<string, string>(200, 10));
}

} // namespace devmand::channels::cli
