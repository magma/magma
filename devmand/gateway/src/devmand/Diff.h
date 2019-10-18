// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <set>
#include <algorithm>
#include <functional>

namespace devmand {

enum class DiffEvent {
  Add,
  Modify,
  Delete,
};

template <class Type>
using DiffEventHandler = std::function<void(DiffEvent, const Type&)>;

// TODO this can be made more generic with SFINAE

/* This implementation is based on std::set_difference, but I've gone ahead
 * re-implemented it because we're looking for modified elements (not just
 * added or removed elements) which requires slightly different logic, and
 * doing it with our own diff function let's us do it in a single pass.
 *
 * diffSorted() runs in O(n) time (linear) because it takes advantage of the
 * fact that both collections are sorted. It iterates a single time through
 * both.
 *
 * std::set is sorted by definition. Lists or vecs need to be sorted before
 * calling diffSorted(), making the whole thing O(n log n) time.
 */

template <class Iterator, class Type>
inline static void doHandle(
    DiffEvent event,
    Iterator& iter,
    const DiffEventHandler<Type>& handler) {
  handler(event, *iter);
  ++iter;
}

template <class Type>
void diffSorted(
    const std::set<Type>& oldC,
    const std::set<Type>& newC,
    const DiffEventHandler<Type>& handler) {
  // two iterators move through both collections together
  auto oldIt = oldC.begin();
  auto newIt = newC.begin();
  // keep going until you reach the end of both collections
  while (oldIt != oldC.end() || newIt != newC.end()) {
    if (oldIt == oldC.end()) {
      doHandle(DiffEvent::Add, newIt, handler);
    } else if (newIt == newC.end()) {
      doHandle(DiffEvent::Delete, oldIt, handler);
    } else {
      // not at the end of either collection
      if (*oldIt < *newIt) {
        doHandle(DiffEvent::Delete, oldIt, handler);
      } else if (*newIt < *oldIt) {
        doHandle(DiffEvent::Add, newIt, handler);
      } else {
        // element is unchanged, step forward
        ++oldIt;
        ++newIt;
      }
    }
  }
}

// separate definition for std::set to bypass sorting,
// std::set is sorted in it's implementation
template <class Type>
void diff(
    const std::set<Type>& oldSet,
    const std::set<Type>& newSet,
    const DiffEventHandler<Type>& handler) {
  diffSorted(oldSet, newSet, handler);
}

// TODO: generic definition for various containers, including std::list and
// std::vec. Make sure to sort the container before passing it to diffSorted.

} // namespace devmand
