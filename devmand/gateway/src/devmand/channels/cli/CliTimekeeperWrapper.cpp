// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/CliTimekeeperWrapper.h>

using devmand::channels::cli::CancelableWTCallback;
using folly::Duration;
using folly::Future;
using folly::Unit;
using std::shared_ptr;

namespace devmand::channels::cli {

void CliTimekeeperWrapper::setCurrentSleepCallback(
    shared_ptr<CancelableWTCallback> _cb) {
  boost::mutex::scoped_lock scoped_lock(this->mutex);
  this->cb = _cb;
}

Future<Unit> CliTimekeeperWrapper::after(Duration dur) {
  shared_ptr<CancelableWTCallback> callback = timekeeper->cancelableSleep(dur);
  setCurrentSleepCallback(callback);
  return callback->getFuture();
}

void CliTimekeeperWrapper::cancelAll() {
  if (cb.use_count() >= 1) {
    cb->callbackCanceled();
  }
}

CliTimekeeperWrapper::~CliTimekeeperWrapper() {
  cancelAll();
}

CliTimekeeperWrapper::CliTimekeeperWrapper(
    const shared_ptr<CliThreadWheelTimekeeper>& _timekeeper)
    : timekeeper(_timekeeper) {}

const shared_ptr<CliThreadWheelTimekeeper>&
CliTimekeeperWrapper::getTimekeeper() const {
  return timekeeper;
}
} // namespace devmand::channels::cli