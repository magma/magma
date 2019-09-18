// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/mikrotik/Channel.h>
#include <devmand/channels/mikrotik/WriteTask.h>

namespace devmand {
namespace channels {
namespace mikrotik {

WriteTask::Id WriteTask::guid{0};

WriteTask::Id WriteTask::getId() const {
  return id;
}

WriteTask::Id WriteTask::getNextId() {
  return guid++;
}

WriteTask::WriteTask(Channel& channel_, const std::string& buffer_)
    : channel(channel_), buffer(buffer_), id(getNextId()) {}

void WriteTask::writeSuccess() noexcept {
  LOG(INFO) << "write success";
  channel.complete(getId());
}

void WriteTask::writeTo(std::shared_ptr<folly::AsyncSocket>& socket) {
  assert(socket != nullptr);
  socket->write(this, buffer.c_str(), buffer.length());
}

void WriteTask::writeErr(
    size_t bytesWritten,
    const folly::AsyncSocketException& ex) noexcept {
  LOG(ERROR) << "write error @" << bytesWritten << " " << ex.what();
  channel.disconnect();
  channel.tryReconnect();
  channel.complete(getId());
}

} // namespace mikrotik
} // namespace channels
} // namespace devmand
