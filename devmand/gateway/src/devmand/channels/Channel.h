// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

namespace devmand {
namespace channels {

/*
 * This is the common interface all channels must implement. A channel is a way
 * to communicate with a device. The channel's engine maintains any state needed
 * outside of an individual connection.
 */
class Channel {
 public:
  Channel() = default;
  virtual ~Channel() = default;
  Channel(const Channel&) = delete;
  Channel& operator=(const Channel&) = delete;
  Channel(Channel&&) = delete;
  Channel& operator=(Channel&&) = delete;
};

} // namespace channels
} // namespace devmand
