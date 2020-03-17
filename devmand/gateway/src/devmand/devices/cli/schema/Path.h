// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <folly/dynamic.h>
#include <ostream>

namespace devmand {
namespace devices {
namespace cli {

using namespace folly;
using namespace std;

class Path {
 public:
  typedef dynamic Keys;

 private:
  string path;

 public:
  // Path separator for YANG is a single char: '/'
  //  using string type for convenience
  static const string PATH_SEPARATOR;
  static const Path ROOT;

  explicit Path(const string& _path);
  Path(const char* _path) : Path(string(_path)){};

  const Path unkeyed() const;
  vector<string> getSegments() const;
  string getLastSegment() const;
  const Path prefixAllSegments() const;
  const Path unprefixAllSegments() const;
  bool isChildOfUnprefixed(const Path& parent) const;
  bool isLastSegmentKeyed() const;
  Optional<string> getFirstModuleName() const;
  u_long getDepth() const;
  bool isChildOf(const Path& parent) const;
  const Path getParent() const;
  const Path getChild(string childSegment) const;

  const Path addKeys(Keys keys) const;
  const Path addKeysToSegment(string segment, Keys keys) const;
  Keys getKeys() const;
  Keys getKeysFromSegment(string segment) const;
  unsigned int segmentDistance(Path path) const;

  static const Path joinSegments(vector<string> segments);
  static const string serializeKeys(Keys keys);
  static const Keys parseKeys(string keys);

  string str() const;
  friend ostream& operator<<(ostream& os, const Path& path);

  bool operator==(const Path& rhs) const;
  bool operator!=(const Path& rhs) const;
  friend bool operator<(const Path& lhs, const Path& rhs);
  friend bool operator>(const Path& lhs, const Path& rhs);
  friend bool operator<=(const Path& lhs, const Path& rhs);
  friend bool operator>=(const Path& lhs, const Path& rhs);
  friend Path operator+(const Path& lhs, const string& rhs);
};

class InvalidPathException : public runtime_error {
 public:
  InvalidPathException(string path, string reason)
      : runtime_error("Invalid path: " + path + " due to: " + reason){};
  InvalidPathException(string reason) : runtime_error(reason){};
};

} // namespace cli
} // namespace devices
} // namespace devmand
