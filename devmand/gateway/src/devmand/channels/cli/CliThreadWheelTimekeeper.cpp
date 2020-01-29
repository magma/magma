// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/CliThreadWheelTimekeeper.h>
#include <folly/futures/Future.h>

using devmand::channels::cli::CancelableWTCallback;
using folly::Duration;
using folly::Future;
using folly::Unit;
using std::shared_ptr;

shared_ptr<CancelableWTCallback>
devmand::channels::cli::CliThreadWheelTimekeeper::cancelableSleep(
    Duration dur) {
  auto cob = CancelableWTCallback::create(&eventBase_);
  eventBase_.runInEventBaseThread(
      [this, cob, dur] { wheelTimer_->scheduleTimeout(cob.get(), dur); });
  return cob;
}
