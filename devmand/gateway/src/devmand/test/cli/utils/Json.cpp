// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <boost/algorithm/string/classification.hpp>
#include <boost/algorithm/string/join.hpp>
#include <boost/algorithm/string/replace.hpp>
#include <boost/algorithm/string/split.hpp>
#include <devmand/test/cli/utils/Json.h>
#include <vector>

namespace devmand {
namespace test {
namespace utils {
namespace json {

using namespace std;

string sortJson(const string& json) {
  vector<string> lines;
  // Split to lines
  boost::split(lines, json, boost::is_any_of("\n"), boost::token_compress_on);
  // Sort
  sort(lines.begin(), lines.end());
  auto joined = boost::algorithm::join(lines, "\n");
  // Remove comma
  boost::replace_all(joined, ",", "");

  return joined;
}

} // namespace json
} // namespace utils
} // namespace test
} // namespace devmand
