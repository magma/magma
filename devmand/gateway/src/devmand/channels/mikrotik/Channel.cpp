// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <boost/algorithm/hex.hpp>
#include <openssl/md5.h>

#include <chrono>
#include <stdexcept>

#include <folly/GLog.h>

#include <devmand/channels/mikrotik/Channel.h>
#include <devmand/channels/mikrotik/LengthComputation.h>
#include <devmand/utils/EventBaseUtils.h>
#include <devmand/utils/StringUtils.h>

namespace devmand {
namespace channels {
namespace mikrotik {

const constexpr std::chrono::milliseconds connectTimeout(1);
const constexpr std::chrono::seconds retryTime{10};

Channel::Channel(
    folly::EventBase& _eventBase,
    folly::SocketAddress addr,
    const std::string& _username,
    const std::string& _password)
    : eventBase(_eventBase),
      address(addr),
      socket(nullptr),
      state(State::Disconnected),
      reconnect(true),
      username(_username),
      password(_password) {}

Channel::~Channel() {
  reconnect = false;
}

std::string sentenceTerminator{};

void Channel::terminateSentence() {
  writeWordAndLength(sentenceTerminator);
}

void Channel::login() {
  state = State::LoggingIn;
  writeWordAndLength("/login");
  writeWordAndLength(folly::sformat("=name={}", username));
  writeWordAndLength(folly::sformat("=password={}", password));
  terminateSentence();
}

void Channel::loginDeprecated(const std::string& code) {
  LOG(ERROR) << "Logging in with deprecated version";
  state = State::LoggingInDeprecated;

  // Convert the code from a hex string to a binary string
  std::string binaryCode;
  boost::algorithm::unhex(code, std::back_inserter(binaryCode));

  // Concat it with a null and the password.
  std::string concat{"\x00", 1};
  concat.append(password);
  concat.append(binaryCode.c_str(), binaryCode.length());

  // MD5 Sum the concatenation
  unsigned char buffer[MD5_DIGEST_LENGTH];
  MD5(reinterpret_cast<const unsigned char*>(concat.data()),
      concat.size(),
      buffer);

  // Convert the md5 sum into a hex string
  std::string md5;
  boost::algorithm::hex(
      std::string(reinterpret_cast<char*>(buffer), MD5_DIGEST_LENGTH),
      std::back_inserter(md5));

  writeWordAndLength("/login");
  writeWordAndLength(folly::sformat("=name={}", username));
  writeWordAndLength(folly::sformat("=response=00{}", md5));
  writeWordAndLength(sentenceTerminator);
}

void Channel::writeWordAndLength(const Word& word) {
  if (socket != nullptr) {
    write(computeLength(static_cast<uint32_t>(word.length())));
    write(word);
  } else {
    // TODO handle this
  }
}

void Channel::write(const Word& word) {
  debugPrintHex("Outgoing Mikrotik data: ", word);
  assert(socket != nullptr);
  WriteTask task{*this, word};
  auto result = writeTasks.emplace(task.getId(), std::move(task));
  if (result.second) {
    result.first->second.writeTo(socket);
  } else {
    LOG(ERROR) << "failed to write to device";
    // TODO handle error
  }
}

folly::Future<Reply> Channel::write(const Sentence&& sentence) {
  outstandingRequests.emplace_back(folly::Promise<Reply>());
  auto& promise = outstandingRequests.back();
  for (auto& word : sentence) {
    writeWordAndLength(word);
  }
  terminateSentence();
  return promise.getFuture();
}

void Channel::complete(WriteTask::Id id) {
  if (writeTasks.erase(id) != 1) {
    LOG(ERROR) << "failed to completed write task " << id;
  }
}

void Channel::disconnect() {
  if (socket != nullptr) {
    socket = nullptr;
  }
  currentIn.clear();
  state = State::Disconnected;
  clearOutstandingRequests();
}

bool Channel::isConnected() const {
  return state != State::Disconnected;
}

bool Channel::isLoggedIn() const {
  return state == State::FullyConnected;
}

void Channel::connect() {
  if (socket != nullptr) {
    LOG(INFO) << "Socket already connected so not connecting.";
  }

  socket = folly::AsyncSocket::newSocket(&eventBase);
  if (socket != nullptr) {
    socket->connect(this, address, connectTimeout.count());
  } else {
    LOG(ERROR) << "socket new failed";
  }
}

void Channel::connectSuccess() noexcept {
  socket->setReadCB(this);
  login();
}

void Channel::connectErr(const folly::AsyncSocketException& ex) noexcept {
  LOG(ERROR) << "connection error " << ex.what();
  if (socket != nullptr) {
    disconnect();
    tryReconnect();
  }
}

void Channel::tryReconnect() {
  if (reconnect) {
    std::weak_ptr<Channel> weak(this->shared_from_this());
    EventBaseUtils::scheduleIn(
        eventBase,
        [weak]() {
          if (auto shared = weak.lock()) {
            shared->connect();
          }
        },
        retryTime);
  }
}

void Channel::clearOutstandingRequests() {
  for (auto& request : outstandingRequests) {
    request.setValue(Reply{});
  }
  outstandingRequests.clear();
}

void Channel::getReadBuffer(void** bufReturn, size_t* lenReturn) {
  currentBuffer.reserve(maxBuffer);
  *bufReturn = reinterpret_cast<void*>(currentBuffer.data());
  *lenReturn = maxBuffer;
}

void Channel::readDataAvailable(size_t len) noexcept {
  data += std::string(currentBuffer.data(), len);

  handleData();
}

void Channel::debugPrintHex(const std::string& prefix, const std::string& buf) {
  LOG(INFO) << prefix << " " << StringUtils::asHexString(buf, " ");
}

bool Channel::readWord(Word& wordOut) {
  ReadLength length = readLength(data.data(), data.length());
  if (length.lengthSize == 0) {
    LOG(INFO) << "insufficient data to read length";
    return false;
  } else if ((length.contentLength + length.lengthSize) > data.length()) {
    LOG(INFO) << "insufficient data to read entire msg";
    return false;
  } else {
    wordOut =
        std::string{data.data() + length.lengthSize, length.contentLength};
    data = data.substr(
        length.lengthSize + length.contentLength, std::string::npos);
    return true;
  }
}

void Channel::handleData() {
  debugPrintHex("Incoming Mikrotik data: ", data);
  Word word;
  while (readWord(word)) {
    if (word.empty()) {
      LOG(INFO) << "End of sentence";
      if (not handle(currentIn)) {
        LOG(ERROR) << "login failed";
        if (socket != nullptr) {
          disconnect();
          tryReconnect();
        }
      } else {
        currentIn.clear();
      }
    } else {
      LOG(INFO) << "Read word [" << word << "]";
      currentIn.emplace_back(word);
    }
  }
}

bool Channel::handle(const Sentence& sentence) {
  if (sentence.empty()) {
    return false;
  }

  switch (state) {
    case State::LoggingIn:
    case State::LoggingInDeprecated:
      if (sentence.front() == "!done") {
        if (sentence.size() == 1) {
          state = State::FullyConnected;
          writeSentence(currentOut);
          currentOut.clear();
          return true;
        } else if (sentence.size() == 2 and sentence[1].size() > 5) {
          std::string code = sentence[1];
          code = code.substr(5, code.size());
          loginDeprecated(code);
          return true;
        }
      } else {
        // TODO print error msg (=message=error message)
      }
      return false;
    case State::FullyConnected:
      LOG(ERROR) << "handle msg!";
      return true;
    case State::Disconnected:
    default:
      return false;
  }
}

void Channel::readEOF() noexcept {
  LOG(ERROR) << "read eof";
  if (socket != nullptr) {
    disconnect();
    tryReconnect();
  }
}

void Channel::readErr(const folly::AsyncSocketException& ex) noexcept {
  LOG(ERROR) << "read error " << ex.what();
  if (socket != nullptr) {
    disconnect();
    tryReconnect();
  }
}

State Channel::getOperationalDatastore() const {
  return state;
}

void Channel::enableSyslog() {
  // TODO don't just assume success
  writeWordAndLength("/system/logging/action/add");
  writeWordAndLength("=name=syslog");
  writeWordAndLength("=target=remote");
  writeWordAndLength(folly::sformat("=remote={}", "192.168.90.100"));
  terminateSentence();
}

void Channel::writeSentence(const Sentence& sentence) {
  if (state == State::FullyConnected) {
    for (auto& word : sentence) {
      writeWordAndLength(word);
    }
    terminateSentence();
  } else {
    currentOut.insert(currentOut.end(), sentence.begin(), sentence.end());
  }
}

} // namespace mikrotik
} // namespace channels
} // namespace devmand
