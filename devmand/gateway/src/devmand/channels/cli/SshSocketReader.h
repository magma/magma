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

#pragma once

#define LOG_WITH_GLOG

#include <devmand/channels/cli/SshSessionAsync.h>
#include <magma_logging.h>

using devmand::channels::cli::sshsession::SshSessionAsync;

namespace devmand {
namespace channels {
namespace cli {

class SshSocketReader {
 private:
  struct event_base* base;
  std::thread notificationThread;

 public:
  static SshSocketReader& getInstance(); // singleton
  SshSocketReader();
  SshSocketReader(SshSocketReader const&) = delete; // singleton
  void operator=(SshSocketReader const&) = delete; // singleton
  virtual ~SshSocketReader();
  struct event*
  addSshReader(event_callback_fn callbackFn, socket_t fd, void* ptr);
};

} // namespace cli
} // namespace channels
} // namespace devmand
