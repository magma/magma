// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <boost/thread/mutex.hpp>
#include <devmand/channels/cli/CancelableWTCallback.h>
#include <devmand/channels/cli/CliThreadWheelTimekeeper.h>
#include <folly/Unit.h>

namespace devmand::channels::cli {

using devmand::channels::cli::CancelableWTCallback;
using folly::Duration;
using folly::Future;
using folly::SemiFuture;
using folly::Timekeeper;
using folly::Unit;
using folly::unit;
using std::shared_ptr;

class CliTimekeeperWrapper : public Timekeeper {
 private:
  shared_ptr<CliThreadWheelTimekeeper> timekeeper;
  boost::mutex mutex;
  shared_ptr<CancelableWTCallback> cb;

 public:
  CliTimekeeperWrapper(const shared_ptr<CliThreadWheelTimekeeper>& timekeeper);

  ~CliTimekeeperWrapper();

 public:
  const shared_ptr<CliThreadWheelTimekeeper>& getTimekeeper() const;

  void cancelAll();

  void setCurrentSleepCallback(shared_ptr<CancelableWTCallback> _cb);

  Future<Unit> after(Duration dur) override;
};

} // namespace devmand::channels::cli
