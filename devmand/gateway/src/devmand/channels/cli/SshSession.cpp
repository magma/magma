// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/SshSession.h>
#include <libssh/libssh.h>
#include <magma_logging.h>

namespace devmand {
namespace channels {
namespace cli {
namespace sshsession {

using std::string;

void SshSession::close() {
  auto chnl = sessionState.channel.load();
  auto sess = sessionState.session.load();

  if (chnl != nullptr) {
    MLOG(MINFO) << "[" << id << "] "
                << "Disconnecting from host: " << sessionState.ip
                << " port: " << sessionState.port;

    if (ssh_channel_is_eof(chnl) != 0) {
      ssh_channel_close(chnl);
    }
    if (ssh_channel_is_open(chnl) != 0) {
      ssh_channel_free(chnl);
    }
  }

  if (sess != nullptr) {
    if (ssh_is_connected(sess) == 1) {
      ssh_disconnect(sess);
    }
    ssh_free(sess);
  }

  sessionState.channel.store(nullptr);
  sessionState.session.store(nullptr);
}

bool SshSession::isOpen() {
  return sessionState.session.load() != nullptr &&
      ssh_is_connected(sessionState.session.load());
}

void SshSession::openShell(
    const string& ip,
    int port,
    const string& username,
    const string& password,
    const long timeout) {
  MLOG(MINFO) << "[" << id << "] "
              << "Connecting to host: " << ip << " port: " << port;
  sessionState.ip = ip;
  sessionState.port = port;
  sessionState.username = username;
  sessionState.username = password;
  sessionState.session.store(ssh_new());
  ssh_options_set(sessionState.session, SSH_OPTIONS_USER, username.c_str());
  ssh_options_set(sessionState.session, SSH_OPTIONS_HOST, ip.c_str());
  ssh_options_set(sessionState.session, SSH_OPTIONS_LOG_VERBOSITY, &verbosity);
  ssh_options_set(sessionState.session, SSH_OPTIONS_PORT, &port);
  // Connection timeout in seconds
  ssh_options_set(sessionState.session, SSH_OPTIONS_TIMEOUT, &timeout);

  checkSuccess(ssh_connect(sessionState.session), SSH_OK);

  int rc = ssh_userauth_password(
      sessionState.session, username.c_str(), password.c_str());
  checkSuccess(rc, SSH_AUTH_SUCCESS);

  sessionState.channel.store(ssh_channel_new(sessionState.session));
  if (sessionState.channel == nullptr) {
    terminate<CliException>();
  }

  rc = ssh_channel_open_session(sessionState.channel);
  checkSuccess(rc, SSH_OK);

  rc = ssh_channel_request_pty(sessionState.channel);
  checkSuccess(rc, SSH_OK);

  rc = ssh_channel_change_pty_size(sessionState.channel, 0, 0);
  checkSuccess(rc, SSH_OK);

  rc = ssh_channel_request_shell(sessionState.channel);
  checkSuccess(rc, SSH_OK);
}

bool SshSession::checkSuccess(int return_code, int OK_RETURN_CODE) {
  if (return_code == OK_RETURN_CODE) {
    return true;
  }
  terminate<CliException>(); // TODO is this an appropriate reaction to every
                             // problem??
  return false;
}

template <typename E>
void SshSession::terminate() {
  const char* error_message = sessionState.session != nullptr
      ? ssh_get_error(sessionState.session)
      : "unknown";
  MLOG(MERROR) << "[" << id << "] "
               << "Error in SSH connection to host: " << sessionState.ip
               << " port: " << sessionState.port
               << " with error: " << error_message;
  string error = "Error with SSH: ";
  throw E(error + error_message);
}

string SshSession::read() {
  char buffer[2048];
  string result;

  auto chnl = sessionState.channel.load();

  while (ssh_channel_is_open(chnl) && !ssh_channel_is_eof(chnl)) {
    int bytes_read =
        ssh_channel_read_nonblocking(chnl, buffer, sizeof(buffer), 0);

    if (bytes_read < 0) {
      MLOG(MERROR) << "[" << id << "] "
                   << "Error reading data from SSH connection, read bytes: "
                   << bytes_read;
      terminate<CommandExecutionException>();
    } else if (bytes_read == 0) {
      return result;
    } else {
      result.append(buffer, (unsigned int)bytes_read);
    }
  }

  return "";
}

void SshSession::write(const string& command) {
  // If we are executing empty string, we can complete right away
  // Why ? because keepalive is empty string
  if (command.empty()) {
    return;
  }
  const char* data = command.c_str();
  int bytes = ssh_channel_write(
      sessionState.channel.load(),
      data,
      (unsigned int)command.length() * sizeof(data[0]));

  if (bytes == SSH_ERROR) {
    MLOG(MERROR) << "[" << id << "] "
                 << "Error while executing command " << command;
    terminate<CommandExecutionException>();
  }
}

SshSession::~SshSession() {
  close();
}

SshSession::SshSession(string _id) : id(_id), verbosity(SSH_LOG_NOLOG) {
  sessionState.channel.store(nullptr);
  sessionState.session.store(nullptr);
}

socket_t SshSession::getSshFd() {
  return ssh_get_fd(sessionState.session.load());
}

} // namespace sshsession
} // namespace cli
} // namespace channels
} // namespace devmand
