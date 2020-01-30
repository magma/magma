// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <devmand/channels/cli/CancelableWTCallback.h>
#include <folly/futures/ThreadWheelTimekeeper.h>

namespace devmand::channels::cli {

using devmand::channels::cli::CancelableWTCallback;
using folly::Duration;
using folly::Future;
using folly::SemiFuture;
using folly::ThreadWheelTimekeeper;
using folly::Unit;
using std::shared_ptr;

class CliThreadWheelTimekeeper : public ThreadWheelTimekeeper {
 public:
  shared_ptr<CancelableWTCallback> cancelableSleep(Duration);
};

} // namespace devmand::channels::cli
