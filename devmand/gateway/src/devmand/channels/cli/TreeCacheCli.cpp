// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/TreeCacheCli.h>

namespace devmand::channels::cli {

using namespace std;
using namespace folly;

TreeCacheCli::TreeCacheCli(
    string _id,
    shared_ptr<Cli> _cli,
    shared_ptr<folly::Executor> _executor,
    shared_ptr<CliFlavour> _sharedCliFlavour)
    : id(_id),
      cli(_cli),
      executor(_executor),
      sharedCliFlavour(_sharedCliFlavour) {
  cache = make_unique<TreeCache>(_sharedCliFlavour);
}

TreeCacheCli::TreeCacheCli(
    string _id,
    shared_ptr<Cli> _cli,
    shared_ptr<folly::Executor> _executor,
    shared_ptr<CliFlavour> _sharedCliFlavour,
    shared_ptr<TreeCache> _cache)
    : id(_id),
      cli(_cli),
      executor(_executor),
      sharedCliFlavour(_sharedCliFlavour),
      cache(_cache) {}

SemiFuture<folly::Unit> TreeCacheCli::destroy() {
  MLOG(MDEBUG) << "[" << id << "] "
               << "destroy: started";
  // call underlying destroy()
  SemiFuture<folly::Unit> innerDestroy = cli->destroy();
  MLOG(MDEBUG) << "[" << id << "] "
               << "destroy: done";
  return innerDestroy;
}

TreeCacheCli::~TreeCacheCli() {
  MLOG(MDEBUG) << "[" << id << "] "
               << "~TCcli: started";
  destroy().get();
  executor = nullptr;
  cli = nullptr;
  MLOG(MDEBUG) << "[" << id << "] "
               << "~TCcli: done";
}

void TreeCacheCli::clear() {
  cache->clear();
}

folly::SemiFuture<std::string> TreeCacheCli::executeRead(
    const ReadCommand cmd) {
  if (!cmd.skipCache()) {
    Optional<pair<string, vector<string>>> maybeParsedCommand =
        cache->parseCommand(cmd.raw());
    if (maybeParsedCommand) {
      pair<string, vector<string>> parsedCommand = maybeParsedCommand.value();
      // if it is just base command, just run it on inner cli where
      // it will be hopefully cached
      if (parsedCommand.second.size() > 0) {
        Optional<string> cachedResult = cache->getSection(parsedCommand);
        if (cachedResult) {
          MLOG(MDEBUG) << "[" << id << "] "
                       << "Found in cache:" << cmd;
          return folly::SemiFuture<string>(cachedResult.value());
        } else if (cache->isEmpty()) {
          MLOG(MDEBUG) << "[" << id << "] "
                       << "Cache is empty, running base of:" << cmd;
          // the cache is empty or there is missing functionality in section
          // parsing
          ReadCommand baseCmd = ReadCommand::create(parsedCommand.first, false);
          // run the base command to populate cache
          return cli->executeRead(baseCmd)
              .via(executor.get())
              .thenValue([_cache = this->cache,
                          _id = this->id,
                          parsedCommand,
                          cmd,
                          _cli = this->cli](string output) {
                _cache->update(output);
                MLOG(MDEBUG) << "Cache keys: " << _cache->toString();
                Optional<string> cachedResult2 =
                    _cache->getSection(parsedCommand);
                if (cachedResult2) {
                  MLOG(MDEBUG) << "[" << _id << "] "
                               << "Found in cache:" << cmd;
                  return folly::SemiFuture<string>(cachedResult2.value());
                }
                MLOG(MWARNING)
                    << "[" << _id << "] "
                    << "Cache was updated, unexpected cache miss for:" << cmd;

                // worst case scenario: we must run the exact command
                return _cli->executeRead(cmd);
              })
              .semi();
        } else {
          MLOG(MWARNING) << "[" << id << "] "
                         << "Unexpected cache miss for:" << cmd;
          // there might be missing functionality in section parsing
          // or the command is wrong. In any case just run the command
        }
      } else { // no subcommands, hit underlying cache
        MLOG(MDEBUG) << "[" << id << "] "
                     << "Pass through:" << cmd;
      }
    }
  }
  // Just run the command: Either this is not supported command,
  // or there is missing functionality in section parsing
  // or this is just supported base command and should be handled by underlying
  // cache.
  return cli->executeRead(cmd);
}

folly::SemiFuture<std::string> TreeCacheCli::executeWrite(
    const WriteCommand cmd) {
  return cli->executeWrite(cmd);
}
} // namespace devmand::channels::cli
