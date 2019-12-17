// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <folly/GLog.h>

#include <devmand/utils/YangUtils.h>

namespace devmand {

folly::dynamic YangUtils::lookup(
    const folly::dynamic& yang,
    const YangPath& path) {
  // TODO handle arrays
  const folly::dynamic* cur{&yang};
  for (auto& elem : path) {
    if (not cur->isObject()) {
      LOG(ERROR) << "yang path lookup failed on elem '" << elem
                 << "' as its parent was not an object ("
                 << (cur->isNull() ? "null" : cur->asString()) << ")";
      return nullptr;
    }
    cur = cur->get_ptr(std::string(elem));
    if (cur == nullptr) {
      LOG(ERROR) << "yang path lookup failed on elem '" << elem
                 << "' as it did not exist";
      break;
    }
  }
  return cur == nullptr ? nullptr : *cur;
}

void YangUtils::set(
    folly::dynamic& yang,
    const YangPath& path,
    const folly::dynamic& value) {
  // TODO handle arrays
  folly::dynamic* cur{&yang};
  folly::dynamic* last{nullptr};
  for (auto& elem : path) {
    if (cur == nullptr) {
      LOG(ERROR) << "yang path set failed on elem previous to '" << elem
                 << "' as it did not exist";
      return;
    }

    if (not cur->isObject()) {
      LOG(ERROR) << "yang path set failed on elem '" << elem
                 << "' as its parent was not an object ("
                 << (cur->isNull() ? "null" : cur->asString()) << ")";
      return;
    }
    last = cur;
    cur = cur->get_ptr(std::string(elem));
  }

  if (last != nullptr) {
    (*last)[path.filename().native()] = value;
  } else {
    LOG(ERROR) << "empty yang path provided";
  }
}

} // namespace devmand
