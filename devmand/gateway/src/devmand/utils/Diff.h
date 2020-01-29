// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <experimental/type_traits>

#include <algorithm>
#include <functional>
#include <list>
#include <set>
#include <vector>

namespace devmand {

enum class DiffEvent {
  Add,
  Modify,
  Delete,
};

template <class Type>
using DiffEventHandler = std::function<void(DiffEvent, const Type&)>;

/* This implementation is based on std::set_difference, but I've gone ahead
 * re-implemented it because 1) we need to detect modified configs, not just
 * added or removed, and 2) we might want to customize the logic, and 3) it
 * works in a single pass, and 4) it works on sets and lists.
 *
 * diffSorted() runs in O(n) time (linear) because it takes advantage of the
 * fact that both collections are sorted. It iterates a single time through
 * both.
 *
 * std::set is sorted by definition. Lists or vecs need to be sorted before
 * calling diffSorted(), making the whole thing O(n log n) time. There's an
 * overloaded "diff" function which sorts lists and skips that for sets.
 *
 * Right now diff() is tested and works on std::set and std::list. std::vec
 * does not work - it's an unsorted type without a sort function.
 */

template <class Iterator, class Type>
inline static void doHandle(
    DiffEvent event,
    Iterator& iter,
    const DiffEventHandler<Type>& handler) {
  handler(event, *iter);
  ++iter;
}

template <class Container, class Type>
void diffSorted(
    const Container& oldC,
    const Container& newC,
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
        if (*oldIt != *newIt) {
          doHandle(DiffEvent::Modify, newIt, handler);
        }
        // element is unchanged, step forward
        ++oldIt;
        ++newIt;
      }
    }
  }
}

// declared type that is a sort member function
template <typename T>
using sort_t = decltype(std::declval<T&>().sort());

// constexpr that determines if a class has a sort member function
template <typename T>
constexpr bool has_sort = std::experimental::is_detected_v<sort_t, T>;

/* Containers are copy by value because in the case of a list or vec the
 * algorithm needs to sort the collection, which requires modifying it.
 * Generally this shouldn't be too expensive because if the Container has
 * large types they'll just be references anyways.
 *
 * This template applies if Container has a sort() member function.
 */
template <class Container, class Type>
std::enable_if_t<has_sort<Container>, void> diff(
    Container oldContainer,
    Container newContainer,
    const DiffEventHandler<Type>& handler) {
  oldContainer.sort();
  newContainer.sort();
  diffSorted(oldContainer, newContainer, handler);
}

/* Overloaded instantiation of diff() for containers that don't have a sort()
 * function. It checks to make sure that the container is in fact sorted.
 * Unsorted containers without a sort() function throw a runtime error.
 */
template <class Container, class Type>
std::enable_if_t<!has_sort<Container>, void> diff(
    const Container& oldContainer,
    const Container& newContainer,
    const DiffEventHandler<Type>& handler) {
  auto oldIsSorted = std::is_sorted(oldContainer.begin(), oldContainer.end());
  auto newIsSorted = std::is_sorted(newContainer.begin(), newContainer.end());
  if (not oldIsSorted or not newIsSorted) {
    throw "Diff.h diff() function was given an unsorted collection without"
        " a built in sort() function.";
  }
  diffSorted(oldContainer, newContainer, handler);
}

} // namespace devmand
