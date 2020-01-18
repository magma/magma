// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/engine/Engine.h>
#include <devmand/test/cli/utils/Ssh.h>
#include <folly/Conv.h>
#include <folly/futures/Future.h>
#include <regex>

namespace devmand {
namespace test {
namespace utils {
namespace ssh {

using namespace std;
using namespace devmand::channels::cli;

shared_ptr<CPUThreadPoolExecutor> testExecutor =
    make_shared<CPUThreadPoolExecutor>(2);

atomic_bool sshInitialized(false);

void initSsh() {
  bool f = false;
  if (sshInitialized.compare_exchange_strong(f, true)) {
    Engine::initSsh();
    MLOG(MDEBUG) << "Ssh for test initialized";
  }
}

static const auto sleep = regex(R"(sleep (\d+))");
static const auto echo = regex(R"(echo (.+))");
static const string newline = "\r\n";

static void handleCommand(ssh_channel channel, stringstream& inputBuf) {
  smatch match;
  string currentOutput = inputBuf.str();
  if (regex_search(currentOutput, match, sleep)) {
    int seconds = folly::to<int>(match[1].str());
    MLOG(MDEBUG) << "Executing sleep: " << seconds;
    this_thread::sleep_for(chrono::seconds(seconds));
    inputBuf.str(match.suffix().str());
  }

  match = smatch();
  currentOutput = inputBuf.str();
  if (regex_search(currentOutput, match, echo)) {
    string echoOut = match[1];
    MLOG(MDEBUG) << "Executing echo: " << echoOut;
    ssh_channel_write(channel, echoOut.c_str(), uint32_t(echoOut.length()));
    ssh_channel_write(channel, newline.c_str(), uint32_t(2));
    inputBuf.str(match.suffix().str());
  }
}

shared_ptr<server> startSshServer(
    shared_ptr<CPUThreadPoolExecutor> executor,
    string address,
    uint port,
    string rsaKey,
    string prompt) {
  MLOG(MDEBUG) << "Opening server: " << address << ":" << port;
  ssh_bind sshbind = ssh_bind_new();
  ssh_bind_set_blocking(sshbind, 0);
  ssh_session session = ssh_new();

  ssh_bind_options_set(sshbind, SSH_BIND_OPTIONS_RSAKEY, rsaKey.c_str());
  ssh_bind_options_set(sshbind, SSH_BIND_OPTIONS_BINDADDR, address.c_str());
  ssh_bind_options_set(sshbind, SSH_BIND_OPTIONS_BINDPORT, &port);

  std::stringstream fmt;
  fmt << address << ":" << port;
  auto retVal = std::make_shared<server>();
  retVal->id = fmt.str();
  retVal->sshbind = sshbind;
  retVal->session = session;

  // This flag makes this function return only after ssh server started
  // listening
  atomic_bool started;
  started.store(false);

  if (ssh_bind_listen(retVal->sshbind) < 0) {
    MLOG(MERROR) << "Cannot open server (socket listen error): " << retVal->id;
    retVal->close();
    throw runtime_error("Cannot open server");
  }

  Future<Unit> future =
      folly::via(executor.get(), [retVal, prompt, &started]() -> void {
        MLOG(MDEBUG) << "Waiting for session on: " << retVal->id;
        started.store(true);

        int r = ssh_bind_accept(retVal->sshbind, retVal->session);
        if (r == SSH_ERROR) {
          MLOG(MERROR) << "Cannot accept connection on server: " << retVal->id;
          return;
        }

        if (ssh_handle_key_exchange(retVal->session)) {
          MLOG(MERROR) << "Cannot kex on server: " << retVal->id << " due to "
                       << ssh_get_error(retVal->session);
          throw runtime_error("Cannot key exchange");
        }

        int auth = 0;
        ssh_message message;

        do {
          message = ssh_message_get(retVal->session);
          if (!message)
            break;
          switch (ssh_message_type(message)) {
            case SSH_REQUEST_AUTH:
              switch (ssh_message_subtype(message)) {
                case SSH_AUTH_METHOD_PASSWORD:
                  MLOG(MDEBUG)
                      << "User auth success: " << ssh_message_auth_user(message)
                      << ":" << ssh_message_auth_password(message);
                  auth = 1;
                  ssh_message_auth_reply_success(message, 0);
                  break;
                  // not authenticated, send default message
                case SSH_AUTH_METHOD_NONE:
                default:
                  ssh_message_auth_set_methods(
                      message, SSH_AUTH_METHOD_PASSWORD);
                  ssh_message_reply_default(message);
                  break;
              }
              break;
            default:
              ssh_message_reply_default(message);
          }
          ssh_message_free(message);
        } while (!auth);
        if (!auth) {
          MLOG(MERROR) << "User auth failed: " << ssh_message_auth_user(message)
                       << ":" << ssh_message_auth_password(message);
          throw runtime_error("Cannot authenticate");
        }

        ssh_channel chan = 0;
        do {
          message = ssh_message_get(retVal->session);
          if (message) {
            switch (ssh_message_type(message)) {
              case SSH_REQUEST_CHANNEL_OPEN:
                if (ssh_message_subtype(message) == SSH_CHANNEL_SESSION) {
                  chan = ssh_message_channel_request_open_reply_accept(message);
                  MLOG(MDEBUG) << "Channel opened";
                  break;
                }
              default:
                ssh_message_reply_default(message);
            }
            ssh_message_free(message);
          }
        } while (message && !chan);
        if (!chan) {
          MLOG(MERROR) << "Channel open failed: "
                       << ssh_get_error(retVal->session);
          throw runtime_error("Cannot open channel");
        }

        int pty = 0;

        do {
          message = ssh_message_get(retVal->session);
          if (message && ssh_message_type(message) == SSH_REQUEST_CHANNEL &&
              ssh_message_subtype(message) == SSH_CHANNEL_REQUEST_PTY) {
            MLOG(MDEBUG) << "PTY requested";
            pty = 1;
            ssh_message_channel_request_reply_success(message);
            ssh_message_free(message);
            break;
          }
          if (!pty) {
            ssh_message_reply_default(message);
          }
          ssh_message_free(message);
        } while (message && !pty);
        if (!pty) {
          MLOG(MERROR) << "PTY open failed: " << ssh_get_error(retVal->session);
          throw runtime_error("Cannot open pty");
        }

        int shell = 0;

        do {
          message = ssh_message_get(retVal->session);
          if (message && ssh_message_type(message) == SSH_REQUEST_CHANNEL &&
              ssh_message_subtype(message) == SSH_CHANNEL_REQUEST_SHELL) {
            MLOG(MDEBUG) << "Shell requested";
            shell = 1;
            ssh_message_channel_request_reply_success(message);
            ssh_message_free(message);
            break;
          }
          if (!shell) {
            ssh_message_reply_default(message);
          }
          ssh_message_free(message);
        } while (message && !shell);
        if (!shell) {
          MLOG(MERROR) << "Shell open failed: "
                       << ssh_get_error(retVal->session);
          throw runtime_error("Cannot open shell");
        }

        char buf[1];
        int i = 0;

        ssh_channel_write(chan, prompt.c_str(), uint32_t(prompt.length()));
        stringstream allInput;
        do {
          bool wasEnter = false;
          i = ssh_channel_read(chan, buf, 1, 0);
          if (i > 0) {
            string received = string(buf, uint(i));
            allInput << received;

            if (buf[0] == '\r' || buf[0] == '\n') {
              if (wasEnter && buf[0] == '\n') {
                // Handle \r\n as enter, throw away the \n
                wasEnter = false;
                continue;
              }

              wasEnter = true;
              MLOG(MDEBUG) << "Received: " << allInput.str();
              MLOG(MDEBUG) << "Received: newline";
              ssh_channel_write(chan, newline.c_str(), uint32_t(2));
              handleCommand(chan, allInput);
              ssh_channel_write(
                  chan, prompt.c_str(), uint32_t(prompt.length()));

              {
                lock_guard<std::mutex> lg(retVal->received_guard);
                retVal->received = allInput.str();
              }

            } else {
              wasEnter = false;
              ssh_channel_write(chan, buf, uint32_t(i));
            }
          }
        } while (i > 0);
        MLOG(MDEBUG) << "Session disconnected";
      });

  retVal->serverFuture = move(future);

  do {
    this_thread::sleep_for(chrono::milliseconds(500));
  } while (!started.load());

  return retVal;
}

} // namespace ssh
} // namespace utils
} // namespace test
} // namespace devmand
