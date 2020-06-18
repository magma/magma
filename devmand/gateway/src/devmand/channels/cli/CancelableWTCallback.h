// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <folly/futures/Future.h>
#include <folly/io/async/EventBase.h>

using folly::EventBase;
using folly::exception_wrapper;
using folly::Future;
using folly::FutureNoTimekeeper;
using folly::Promise;
using folly::Unit;

namespace devmand::channels::cli {

struct CancelableWTCallback
    : public std::enable_shared_from_this<CancelableWTCallback>,
      public folly::HHWheelTimer::Callback {
  struct PrivateConstructorTag {};

 public:
  CancelableWTCallback(PrivateConstructorTag, EventBase* base) : base_(base) {}

  // Only allow creation by this factory, to ensure heap allocation.
  static std::shared_ptr<CancelableWTCallback> create(EventBase* base) {
    // optimization opportunity: memory pool
    auto cob =
        std::make_shared<CancelableWTCallback>(PrivateConstructorTag{}, base);
    // Capture shared_ptr of cob in lambda so that Core inside Promise will
    // hold a ref count to it. The ref count will be released when Core goes
    // away which happens when both Promise and Future go away
    cob->promise_.setInterruptHandler(
        [cob](exception_wrapper ew) { cob->interruptHandler(std::move(ew)); });
    return cob;
  }

  Future<Unit> getFuture() {
    return promise_.getFuture();
  }

  FOLLY_NODISCARD Promise<Unit> stealPromise() {
    // Don't need promise anymore. Break the circular reference as promise_
    // is holding a ref count to us via Core. Core won't go away until both
    // Promise and Future go away.
    return std::move(promise_);
  }

  void callbackCanceled() noexcept override {
    base_ = nullptr;
    // Don't need Promise anymore, break the circular reference
    auto promise = stealPromise();
    if (!promise.isFulfilled()) {
      promise.setException(FutureNoTimekeeper{});
    }
  }

 protected:
  folly::Synchronized<EventBase*> base_;
  Promise<Unit> promise_;

  void timeoutExpired() noexcept override {
    base_ = nullptr;
    // Don't need Promise anymore, break the circular reference
    auto promise = stealPromise();
    if (!promise.isFulfilled()) {
      promise.setValue();
    }
  }

  void interruptHandler(exception_wrapper ew) {
    auto rBase = base_.rlock();
    if (!*rBase) {
      return;
    }
    // Capture shared_ptr of self in lambda, if we don't do this, object
    // may go away before the lambda is executed from event base thread.
    // This is not racing with timeoutExpired anymore because this is called
    // through Future, which means Core is still alive and keeping a ref count
    // on us, so what timeouExpired is doing won't make the object go away
    (*rBase)->runInEventBaseThread(
        [me = shared_from_this(), ew = std::move(ew)]() mutable {
          me->cancelTimeout();
          // Don't need Promise anymore, break the circular reference
          auto promise = me->stealPromise();
          if (!promise.isFulfilled()) {
            promise.setException(std::move(ew));
          }
        });
  }
};
} // namespace devmand::channels::cli