/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
