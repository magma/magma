// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <assert.h>

namespace devmand {
namespace utils {

template <class T>
class LifetimeTracker {
 public:
  LifetimeTracker() {
    ++allocations;
  }

  virtual ~LifetimeTracker() {
    ++deallocations;
  }

  LifetimeTracker(const LifetimeTracker&) = default;
  LifetimeTracker& operator=(const LifetimeTracker&) = default;
  LifetimeTracker(LifetimeTracker&&) = delete;
  LifetimeTracker& operator=(LifetimeTracker&&) = default;

  static unsigned int getAllocations() {
    return allocations;
  }

  static unsigned int getDeallocations() {
    return deallocations;
  }

  static unsigned int getLivingCount() {
    assert(allocations > deallocations);

    return allocations - deallocations;
  }

 private:
  static unsigned int allocations;
  static unsigned int deallocations;
};

template <class T>
unsigned int LifetimeTracker<T>::allocations;
template <class T>
unsigned int LifetimeTracker<T>::deallocations;

} // namespace utils
} // namespace devmand
