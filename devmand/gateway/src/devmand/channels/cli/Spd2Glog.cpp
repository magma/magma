// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <devmand/channels/cli/Spd2Glog.h>
#include <magma_logging.h>
#include <spdlog/details/log_msg.h>

void devmand::channels::cli::Spd2Glog::_sink_it(
    const spdlog::details::log_msg& msg) {
  toGlog(msg);
}

void devmand::channels::cli::Spd2Glog::toGlog(
    const spdlog::details::log_msg& msg) {
  if (msg.level == trace || msg.level == debug) {
    MLOG(MDEBUG) << msg.formatted.str();
  } else if (msg.level == info || msg.level == warn) {
    MLOG(MINFO) << msg.formatted.str();
  } else {
    MLOG(MERROR) << msg.formatted.str();
  }
}

void devmand::channels::cli::Spd2Glog::flush() {}
