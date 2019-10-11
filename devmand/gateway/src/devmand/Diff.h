// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <set>

namespace devmand {

enum class DiffEvent {
  Add,
  Modify,
  Delete,
};

template <class Type>
using DiffEventHandler = std::function<void(DiffEvent, const Type&)>;

// TODO First this can be made more generic with SFINAE second the imlementation
// below is quite poor algorithmically but it works for now.
template <class Type>
void diff(
    const std::set<Type>& oldCollection,
    const std::set<Type>& newCollection,
    const DiffEventHandler<Type>& handler) {
  for (auto& oldElem : oldCollection) {
    const auto& newElem = std::find(
        std::cbegin(newCollection), std::cend(newCollection), oldElem);
    if (newElem != newCollection.end()) {
      if (oldElem != *newElem) {
        handler(DiffEvent::Modify, *newElem);
      } else {
        // NOTE handler(DiffEvent::Same, *newElem);
      }
    } else {
      handler(DiffEvent::Delete, oldElem);
    }
  }
  for (auto& newElem : newCollection) {
    const auto& oldElem = std::find(
        std::cbegin(oldCollection), std::cend(oldCollection), newElem);
    if (oldElem == oldCollection.end()) {
      handler(DiffEvent::Add, newElem);
    }
  }
}

} // namespace devmand
