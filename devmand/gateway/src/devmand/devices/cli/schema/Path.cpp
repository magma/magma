// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <boost/algorithm/string/classification.hpp>
#include <boost/algorithm/string/join.hpp>
#include <boost/algorithm/string/split.hpp>
#include <devmand/devices/cli/schema/Path.h>
#include <regex>

namespace devmand {
namespace devices {
namespace cli {

using namespace std;
using namespace folly;

const string Path::PATH_SEPARATOR = "/";
const Path Path::ROOT = Path(PATH_SEPARATOR);

Path::Path(const string& _path) : path(_path) {
  if (_path.empty()) {
    throw InvalidPathException(path, "Empty path");
  }
  // equivalent to _path.startsWith(PATH_SEPARATOR)
  if (_path.rfind(PATH_SEPARATOR, 0) != 0) {
    throw InvalidPathException(path, "Not an absolute path");
  }

  // TODO
  // path should not contain leading/trailing whitespace
  // pre-cache segments and unkeyed version
}

// TODO description
vector<string> Path::getSegments() const {
  if (ROOT == *this) {
    return vector<string>();
  }

  // Keys could contain path separators, so split segments of unkeyed version
  //  and use them to extract segments with keys
  vector<string> unkeyedSegments;
  boost::split(
      unkeyedSegments,
      unkeyed().path,
      boost::is_any_of(Path::PATH_SEPARATOR),
      boost::token_compress_on);

  string pathCopy = path;
  vector<string> segments;

  // skip the first empty segment since the path is always absolute
  for (unsigned int i = 1; i < unkeyedSegments.size(); ++i) {
    if (i == unkeyedSegments.size() - 1) {
      segments.push_back(pathCopy.substr(1));
      break;
    }

    pathCopy = pathCopy.substr(unkeyedSegments[i].length() + 1);

    unsigned long nextSegmentStart =
        pathCopy.find(PATH_SEPARATOR + unkeyedSegments[i + 1]);
    segments.push_back(
        unkeyedSegments[i] + pathCopy.substr(0, nextSegmentStart));

    pathCopy = pathCopy.substr(nextSegmentStart);
  }
  return segments;
}

static const auto PREFIXED_SEGMENT = regex("([^:]+):([^:]+)");

const Path Path::prefixAllSegments() const {
  if (ROOT == *this) {
    return ROOT;
  }

  vector<string> unkeyedSegments = unkeyed().getSegments();

  string pathCopy = path;
  stringstream newPath;
  newPath << PATH_SEPARATOR;

  string lastPrefix = "";
  for (unsigned int i = 0; i < unkeyedSegments.size(); ++i) {
    smatch match;
    if (regex_match(unkeyedSegments[i], match, PREFIXED_SEGMENT)) {
      lastPrefix = match[1];
      newPath << unkeyedSegments[i];
    } else {
      newPath << lastPrefix << ":" << unkeyedSegments[i];
    }

    if (i == unkeyedSegments.size() - 1) {
      newPath << pathCopy.substr(unkeyedSegments[i].length() + 1);
      break;
    }

    pathCopy = pathCopy.substr(unkeyedSegments[i].length() + 1);

    unsigned long nextSegmentStart =
        pathCopy.find(PATH_SEPARATOR + unkeyedSegments[i + 1]);
    newPath << pathCopy.substr(0, nextSegmentStart) << PATH_SEPARATOR;

    pathCopy = pathCopy.substr(nextSegmentStart);
  }

  return Path(newPath.str());
}

static const auto KEYS_IN_PATH = regex("\\[([^\\]]+)\\]");

const Path Path::unkeyed() const {
  return Path(regex_replace(path, KEYS_IN_PATH, ""));
}

string Path::getLastSegment() const {
  if (ROOT == *this) {
    throw InvalidPathException("Invalid operation on root path");
  }
  return getSegments().back();
}

u_long Path::getDepth() const {
  return unkeyed().getSegments().size();
}

bool Path::isChildOf(const Path& parent) const {
  auto thisSegments = unkeyed().getSegments();
  auto parentSegments = parent.unkeyed().getSegments();
  return includes(
      thisSegments.begin(),
      thisSegments.end(),
      parentSegments.begin(),
      parentSegments.end());
}

const Path Path::getParent() const {
  if (ROOT == *this) {
    throw InvalidPathException("Invalid operation on root path");
  }
  vector<string> segments = getSegments();

  if (segments.size() == 1) {
    return ROOT;
  }

  return joinSegments(vector<string>(segments.begin(), segments.end() - 1));
}

const Path Path::getChild(string childSegment) const {
  if (childSegment.empty()) {
    throw InvalidPathException(path, "Unable to add empty child segment");
  }
  // equivalent to childSegment.startsWith(PATH_SEPARATOR)
  if (childSegment.rfind(PATH_SEPARATOR, 0) == 0) {
    throw InvalidPathException(
        path,
        "Unable to add child segment starting with path separator: " +
            childSegment);
  }

  if (ROOT == *this) {
    return Path(path + childSegment);
  }

  return Path(path + PATH_SEPARATOR + childSegment);
}

const Path Path::addKeys(Keys keys) const {
  if (ROOT == *this) {
    throw InvalidPathException("Invalid operation on root path");
  }

  if (keys.empty()) {
    throw InvalidPathException("Unable to add empty keys to path: " + path);
  }

  if (!getKeys().empty()) {
    throw InvalidPathException("Unable to add keys to path with keys: " + path);
  }

  return Path(path + serializeKeys(keys));
}

Path::Keys Path::getKeys() const {
  if (ROOT == *this) {
    throw InvalidPathException("Invalid operation on root path");
  }

  smatch match;
  string lastSegment = getLastSegment();
  regex_search(lastSegment, match, KEYS_IN_PATH);
  string keysAsString = match[0];
  return parseKeys(keysAsString);
}

Path::Keys Path::getKeysFromSegment(string segment) const {
  if (ROOT == *this) {
    throw InvalidPathException("Invalid operation on root path");
  }

  auto unkeyedPath = unkeyed().getSegments();
  auto it = find(unkeyedPath.begin(), unkeyedPath.end(), segment);
  if (it == unkeyedPath.end()) {
    throw InvalidPathException(path, "Cannot find segment: " + segment);
  }
  long index = distance(unkeyedPath.begin(), it);
  if (index < 0) {
    throw InvalidPathException(path, "Cannot find segment: " + segment);
  }

  smatch match;
  string lastSegment = getSegments()[u_int(index)];
  regex_search(lastSegment, match, KEYS_IN_PATH);
  string keysAsString = match[0];
  return parseKeys(keysAsString);
}

// TODO make Keys a proper class and move serialize/parse there

static const auto KEY_IN_PATH = regex("([^=]*)=\'([^\']*)\'");
static const string KEY_SEPARATOR = ",";

const Path::Keys Path::parseKeys(string keys) {
  smatch removeArrayBrackets;
  regex_match(keys, removeArrayBrackets, KEYS_IN_PATH);
  keys = removeArrayBrackets[1];

  vector<string> individualKeys;
  boost::split(
      individualKeys,
      keys,
      boost::is_any_of(KEY_SEPARATOR),
      boost::token_compress_on);

  dynamic parsedKeys = dynamic::object();
  for (const auto& oneKeyAsString : individualKeys) {
    smatch match;
    if (regex_match(oneKeyAsString, match, KEY_IN_PATH)) {
      string k = match[1];
      string v = match[2];
      parsedKeys[k] = v;
    }
  }

  return parsedKeys;
}

const string Path::serializeKeys(Keys keys) {
  stringstream keysAsString;
  keysAsString << "[";
  for (const auto& keyName : keys.keys()) {
    keysAsString << keyName << "='" << keys[keyName] << "'";
  }
  keysAsString << "]";
  return keysAsString.str();
}

const Path Path::joinSegments(vector<string> segments) {
  return Path(
      Path::PATH_SEPARATOR +
      boost::algorithm::join(segments, Path::PATH_SEPARATOR));
}

ostream& operator<<(ostream& os, const Path& path) {
  os << path.path;
  return os;
}

bool Path::operator==(const Path& rhs) const {
  return path == rhs.path;
}

bool Path::operator!=(const Path& rhs) const {
  return !(rhs == *this);
}

bool operator<(const Path& lhs, const Path& rhs) {
  return lhs.path < rhs.path;
}

bool operator>(const Path& lhs, const Path& rhs) {
  return rhs < lhs;
}

bool operator<=(const Path& lhs, const Path& rhs) {
  return !(rhs < lhs);
}

bool operator>=(const Path& lhs, const Path& rhs) {
  return !(lhs < rhs);
}

Path operator+(const Path& lhs, const string& rhs) {
  return Path(lhs.str() + rhs);
}

string Path::str() const {
  return path;
}

} // namespace cli
} // namespace devices
} // namespace devmand
