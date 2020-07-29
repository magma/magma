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

#include <devmand/channels/Channel.h>
#include <devmand/channels/cli/Cli.h>
#include <libssh/libssh.h>

namespace devmand {
namespace channels {
namespace cli {
namespace sshsession {

using std::runtime_error;
using std::string;

class SshSession {
 private:
  struct SshSessionState {
    string ip;
    int port;
    string username;
    string password;
    std::atomic<ssh_channel> channel;
    std::atomic<ssh_session> session;
  } sessionState;
  string id;
  int verbosity;

  bool checkSuccess(int return_code, int OK_RETURN_CODE);
  template <typename E>
  void terminate();

 public:
  SshSession(string _id);
  ~SshSession();

  socket_t getSshFd();
  void openShell(
      const string& ip,
      int port,
      const string& username,
      const string& password,
      const long timeout);
  void close();
  bool isOpen();
  void write(const string& command);
  string read();
};

} // namespace sshsession
} // namespace cli
} // namespace channels
} // namespace devmand
