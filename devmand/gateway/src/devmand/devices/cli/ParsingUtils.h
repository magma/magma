// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <boost/algorithm/string.hpp>
#include <folly/Optional.h>
#include <ydk/types.hpp>
#include <regex>
#include <vector>

namespace devmand {
namespace devices {
namespace cli {

using namespace std;

template <typename T>
vector<T> parseKeys(
    const string& output,
    const regex& pattern,
    const uint& groupToExtract,
    const int& skipLines = 0,
    const function<T(string)>& postProcess = [](auto str) { return str; });

template <typename T>
vector<T> parseLineKeys(
    const string& output,
    const regex& pattern,
    const function<T(string)>& postProcess = [](auto str) { return str; });

folly::Optional<string> extractValue(
    const string& output,
    const regex& pattern,
    const uint& groupToExtract);

void parseValue(
    const string& output,
    const regex& pattern,
    const uint& groupToExtract,
    const std::function<void(string)>& setter);

template <typename T>
void parseLeaf(
    const string& output,
    const regex& pattern,
    ydk::YLeaf& leaf,
    const uint& groupToExtract = 1,
    const function<T(string)>& postProcess = [](auto str) { return str; });

extern function<ydk::uint64(string)> toUI64;
extern function<ydk::uint16(string)> toUI16;

// Templated functions implemented in header

template <typename T>
void parseLeaf(
    const string& output,
    const regex& pattern,
    ydk::YLeaf& leaf,
    const uint& groupToExtract,
    const function<T(string)>& postProcess) {
  parseValue(output, pattern, groupToExtract, [&postProcess, &leaf](string v) {
    leaf = postProcess(v);
  });
}

template <typename T>
vector<T> parseKeys(
    const string& output,
    const regex& pattern,
    const uint& groupToExtract,
    const int& skipLines,
    const function<T(string)>& postProcess) {
  std::stringstream ss(output);
  std::string line;
  vector<T> retval;
  int counter = 0;

  while (std::getline(ss, line, '\n')) {
    counter++;
    if (counter <= skipLines) {
      continue;
    }

    boost::algorithm::trim(line);
    smatch match;
    if (regex_search(line, match, pattern) and match.size() > groupToExtract) {
      T processed = postProcess(match[groupToExtract]);
      retval.push_back(processed);
    }
  }

  return retval;
}

template <typename T>
vector<T> parseLineKeys(
    const string& output,
    const regex& pattern,
    const function<T(string)>& postProcess) {
  vector<T> retval;
  smatch match;
  string currentOutput = output;
  while (regex_search(currentOutput, match, pattern)) {
    T processed = postProcess(match[0]);
    retval.push_back(processed);
    currentOutput = match.suffix().str();
  }

  return retval;
}

} // namespace cli
} // namespace devices
} // namespace devmand
