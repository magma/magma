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

#include <boost/algorithm/string/replace.hpp>
#include <folly/futures/Future.h>
#include <atomic>
#include <iostream>

using std::string;
using std::vector;

namespace devmand {
namespace channels {
namespace cli {
/*
 * Command struct encapsulating a string to be executed on a device.
 */
static std::atomic<int> commandCounter;

class Command {
 protected:
  explicit Command(std::string _command, int idx, bool skipCache);
  string command;
  int idx;
  bool skipCache_;

 public:
  Command() = delete;
  bool isMultiCommand();
  vector<Command> splitMultiCommand();

  string raw() const {
    return command;
  }

  bool skipCache() const {
    return skipCache_;
  }

  static string escape(string cmd) {
    boost::replace_all(cmd, "\n", "\\n");
    boost::replace_all(cmd, "\r", "\\r");
    boost::replace_all(cmd, "\t", "\\t");
    boost::replace_all(cmd, "\"", "\\\"");
    return cmd;
  }

  friend std::ostream& operator<<(std::ostream& _stream, Command const& c) {
    _stream << "[" << c.idx << "] " << Command::escape(c.raw());
    return _stream;
  }
};

class WriteCommand : public Command {
 public:
  static WriteCommand create(const std::string& cmd, bool skipCache = false);
  static WriteCommand create(const Command& cmd);

  WriteCommand(const WriteCommand& wc);
  WriteCommand& operator=(const WriteCommand& other);

 private:
  WriteCommand(const string& command, int idx, bool skipCache);
};

class ReadCommand : public Command {
 public:
  static ReadCommand create(const std::string& cmd, bool skipCache = false);
  static ReadCommand create(const Command& cmd);

  ReadCommand& operator=(const ReadCommand& other);
  ReadCommand(const ReadCommand& rc);

 private:
  ReadCommand(const string& command, int idx, bool skipCache);
};

} // namespace cli
} // namespace channels
} // namespace devmand
