// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

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
