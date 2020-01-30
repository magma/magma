// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once
#include <devmand/channels/cli/CliFlavour.h>

namespace devmand::channels::cli {

using namespace std;
using namespace folly;

/*
 * Thread unsafe structured cache that can serve subsets of running
 * configuration. When running configuration is returned from the device,
 * calling update function will split it into sections. Consequent calls to get
 * method
 */
class TreeCache {
 private:
  shared_ptr<CliFlavour> cliFlavour;
  // section cache
  map<vector<string>, string> treeCache;

 public:
  // TODO: cache invalidation

  TreeCache(shared_ptr<CliFlavour> sharedCliFlavour);

  /*
   * Clear cache.
   */
  void clear();

  /*
   * Parse command. If supported, first part will be base command,
   * second part are parameters.
   */
  Optional<pair<string, vector<string>>> parseCommand(string cmd);

  /*
   * Try to update cache with actual command output. Return true
   * iif cmd does not contain subsections and thus cache
   * was populated.
   */
  // FIXME: currently only one 'show running-config' command is supported.
  void update(string output);

  /*
   * Try to get result of supported command subsection from cache.
   * If show running command with subcommands is passed, return result from
   * cache.
   * It is an error to call getSection with command pair with no second part.
   * Argument cmd must be produced by parseCommand, so only supported
   * (show running-config) command can be passed.
   */
  Optional<string> getSection(pair<string, vector<string>> cmd);

  /*
   * Return true iif cache is empty.
   */
  bool isEmpty();

  string toString() {
    string result = "(" + to_string(treeCache.size()) + ")[";
    for (const auto& entry : treeCache) {
      for (const auto& subkey : entry.first) {
        result += "{" + subkey + "}";
      }
      result += ",";
    }
    result += "]";
    return result;
  }

  // visible for testing:
  size_t size();
  // parsing start
  map<vector<string>, string> readConfigurationToMap(string showRunningOutput);

  map<vector<string>, string> readConfigurationToMap(
      string showRunningOutput,
      unsigned int indentationLevel);

  /*
   * Iterate over sections. If previous match is detected, content will be
   * set to start after found section starting with newline.
   */
  static bool hasNextSection(regex& regexMatch, smatch& sm, string& content);

  /*
   * Assuming command output starts and ends with \n and has lines separated by
   * \n, this regex pattern will match first section of it.
   */
  string createSectionPattern(unsigned int indentationLevel);

  static string createSectionPattern(
      Optional<char> maybeIndentChar,
      string configSubsectionEnd,
      unsigned int indentationLevel);
  // parsing end
  // reading from cache start

  /*
   * If supported show running command supplied, return pair of base, remainder.
   */
  Optional<pair<string, string>> splitSupportedCommand(string cmd);
  /*
   * Parse command output, return only section that was specified by subset
   * argument. This can be one or more subcommands separated by indentation
   * characters.
   */
  static Optional<string> findMatchingSubset(
      vector<string> subcommands,
      map<vector<string>, string> treeCache);
};

} // namespace devmand::channels::cli
