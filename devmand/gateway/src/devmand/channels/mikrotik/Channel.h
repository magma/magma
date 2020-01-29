// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <list>
#include <map>
#include <memory>
#include <string>
#include <vector>

#include <folly/futures/Future.h>
#include <folly/io/async/AsyncSocket.h>
#include <folly/io/async/EventBase.h>
#include <folly/io/async/EventHandler.h>

#include <devmand/channels/Channel.h>
#include <devmand/channels/mikrotik/WriteTask.h>

namespace devmand {
namespace channels {
namespace mikrotik {

// TODO I need to make this use a real buffer class but eh this works for now
// Make this a bootcamp task to find which buffer class should be used and to
// convert it.
using Buffer = std::vector<char>;

using Word = std::string;
using Sentence = std::vector<Word>;

enum class State {
  Disconnected,
  LoggingIn,
  LoggingInDeprecated,
  FullyConnected
};

using Reply = std::list<Sentence>;
using OutstandingRequests = std::list<folly::Promise<Reply>>;

// This is implemented for post-v6.43
class Channel : public channels::Channel,
                public folly::AsyncSocket::ConnectCallback,
                public folly::AsyncTransportWrapper::ReadCallback,
                public std::enable_shared_from_this<Channel> {
  // TOOD make all channels shared from this.
 public:
  Channel(
      folly::EventBase& _eventBase,
      folly::SocketAddress addr,
      const std::string& _username,
      const std::string& _password);
  Channel() = delete;
  virtual ~Channel();
  Channel(const Channel&) = delete;
  Channel& operator=(const Channel&) = delete;
  Channel(Channel&&) = delete;
  Channel& operator=(Channel&&) = delete;

 public:
  void connect();
  void disconnect();
  void tryReconnect();
  bool isConnected() const;
  bool isLoggedIn() const;

  State getOperationalDatastore() const;

  void complete(WriteTask::Id id);

  void login();
  void loginDeprecated(const std::string& code);

  folly::Future<Reply> write(const Sentence&& sentence);

  void writeSentence(const Sentence& sentence);

 private:
  void terminateSentence();
  void writeWordAndLength(const Word& word);
  void write(const Word& word);

  void connectSuccess() noexcept override;
  void connectErr(const folly::AsyncSocketException& ex) noexcept override;

  void getReadBuffer(void** bufReturn, size_t* lenReturn) override;
  void readDataAvailable(size_t len) noexcept override;
  void readEOF() noexcept override;
  void readErr(const folly::AsyncSocketException& ex) noexcept override;

  void handleData();
  bool handle(const Sentence& sentence);
  bool readWord(Word& wordOut);

  void clearOutstandingRequests();

  void debugPrintHex(const std::string& prefix, const std::string& buf);

  void enableSyslog();

 private:
  folly::EventBase& eventBase;
  folly::SocketAddress address;
  std::shared_ptr<folly::AsyncSocket> socket;
  std::map<WriteTask::Id, WriteTask> writeTasks; // TODO make this bounded
  State state{State::Disconnected};
  bool reconnect{true};
  std::string username;
  std::string password;
  Sentence currentIn;
  Sentence currentOut; // TODO do we need to sync this?
  OutstandingRequests outstandingRequests;

  constexpr static size_t maxBuffer{4096};
  Buffer currentBuffer;
  std::string data;
};

} // namespace mikrotik
} // namespace channels
} // namespace devmand
