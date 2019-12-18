// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <vector>

#include <folly/dynamic.h>
#include <folly/io/async/AsyncServerSocket.h>
#include <folly/json.h>

#include <devmand/channels/mikrotik/Channel.h>
#include <devmand/channels/mikrotik/LengthComputation.h>
#include <devmand/test/EventBaseTest.h>
#include <devmand/test/Notifier.h>
#include <devmand/test/TestUtils.h>

namespace devmand {
namespace test {

class MikrotikChannelTest : public folly::AsyncServerSocket::AcceptCallback,
                            public folly::AsyncTransportWrapper::ReadCallback,
                            public folly::AsyncWriter::WriteCallback,
                            public EventBaseTest {
 public:
  MikrotikChannelTest() = default;

  ~MikrotikChannelTest() override = default;
  MikrotikChannelTest(const MikrotikChannelTest&) = delete;
  MikrotikChannelTest& operator=(const MikrotikChannelTest&) = delete;
  MikrotikChannelTest(MikrotikChannelTest&&) = delete;
  MikrotikChannelTest& operator=(MikrotikChannelTest&&) = delete;

 protected:
  void listen() {
    lsocket = folly::AsyncServerSocket::newSocket(&eventBase);
    eventBase.runInEventBaseThread([this]() {
      lsocket->bind(1337);
      lsocket->addAcceptCallback(this, &eventBase);
      lsocket->listen(100);
      lsocket->startAccepting();
    });
  }

 public:
  void connectionAccepted(
      int fd,
      const folly::SocketAddress&) noexcept override {
    acceptNotifier.notify();
    asocket = folly::AsyncSocket::newSocket(
        &eventBase, folly::NetworkSocket::fromFd(fd));
    asocket->setReadCB(this);

    // Only accept once so tests are predicatable. Recall listen if needed.
    lsocket = nullptr;
  }

  void acceptError(const std::exception&) noexcept override {
    FAIL() << "accept error";
  }

 public:
  void getReadBuffer(void** bufReturn, size_t* lenReturn) override {
    *bufReturn = reinterpret_cast<void*>(buffer);
    *lenReturn = bufferLength;
  }

  void readDataAvailable(size_t len) noexcept override {
    if (not checkReads) {
      return;
    }

    LOG(INFO) << "Read of length " << len;

    std::string head{expectedReads.data(), len};
    EXPECT_EQ(head, std::string(buffer, len));
    expectedReads =
        std::string(&expectedReads.data()[len], expectedReads.length() - len);
  }

  void readEOF() noexcept override {}

  void readErr(const folly::AsyncSocketException& ex) noexcept override {
    (void)ex;
    FAIL() << __FUNCTION__ << ex.what();
  }

  // This is test code so just make async sync.
  void write(const std::string& buf) {
    assert(asocket != nullptr);
    asocket->write(this, buf.data(), buf.length());

    writeNotifier.wait();
  }

  void writeSuccess() noexcept override {
    writeNotifier.notify();
  }

  void writeErr(
      size_t bytesWritten,
      const folly::AsyncSocketException& ex) noexcept override {
    FAIL() << "writeErr " << ex.what() << " with " << bytesWritten
           << " bytes written";
  }

  std::string expectedReads;

 protected:
  Notifier acceptNotifier;
  Notifier writeNotifier;
  std::shared_ptr<folly::AsyncServerSocket> lsocket{nullptr};
  std::shared_ptr<folly::AsyncSocket> asocket{nullptr};
  constexpr static size_t maxBuffer{4096};
  size_t bufferLength{maxBuffer};
  char buffer[maxBuffer];
  bool checkReads{true};
};

TEST_F(MikrotikChannelTest, checkNotConnected) {
  checkReads = false;
  folly::SocketAddress address("127.0.0.1", 1337);
  auto channel = std::make_shared<channels::mikrotik::Channel>(
      eventBase, address, "foo", "bar");
  channel->connect();
  EXPECT_FALSE(channel->isConnected());
  stop();
}

TEST_F(MikrotikChannelTest, DISABLED_checkConnects) {
  checkReads = false;
  listen();
  folly::SocketAddress address("127.0.0.1", 1337);
  auto channel = std::make_shared<channels::mikrotik::Channel>(
      eventBase, address, "foo", "bar");
  channel->connect();
  acceptNotifier.wait();
  EXPECT_BECOMES_TRUE(channel->isConnected());
  stop();
}

template <typename... Args>
static inline void checkComputeLength(unsigned int wordLength, Args... chars) {
  std::string len = channels::mikrotik::computeLength(wordLength);
  auto lenIt = len.begin();
  for (const int& c : {chars...}) {
    if (lenIt == len.end()) {
      FAIL() << "cl -- over length on " << wordLength;
      break;
    }
    int ch = *reinterpret_cast<unsigned char*>(&*lenIt);
    EXPECT_EQ(c, ch) << " cl -- mismatch on " << wordLength;
    ++lenIt;
  }
  EXPECT_EQ(lenIt, len.end()) << " cl -- under length on " << wordLength;
}

TEST_F(MikrotikChannelTest, lengthComputation) {
  checkComputeLength(0x00, 0x00);
  checkComputeLength(0x01, 0x01);
  checkComputeLength(0x02, 0x02);
  checkComputeLength(0x7F, 0x7F);
  checkComputeLength(0x80, 0x80, 0x80);
  checkComputeLength(0x81, 0x80, 0x81);
  checkComputeLength(0x3FFF, 0xBF, 0xFF);
  checkComputeLength(0x4000, 0xC0, 0x40, 0x00);
  checkComputeLength(0x4001, 0xC0, 0x40, 0x01);
  checkComputeLength(0x1FFFFF, 0xDF, 0xFF, 0xFF);
  checkComputeLength(0x200000, 0xE0, 0x20, 0x00, 0x00);
  checkComputeLength(0x200001, 0xE0, 0x20, 0x00, 0x01);
  checkComputeLength(0xFFFFFFF, 0xEF, 0xFF, 0xFF, 0xFF);
  checkComputeLength(0x10000000, 0xF0, 0x10, 0x00, 0x00, 0x00);
  checkComputeLength(0x10000001, 0xF0, 0x10, 0x00, 0x00, 0x01);
  checkComputeLength(0xFFFFFFFF, 0xF0, 0xFF, 0xFF, 0xFF, 0xFF);
}

static inline void checkReadLength(
    size_t wordLength,
    std::vector<unsigned char> chars) {
  channels::mikrotik::ReadLength len = channels::mikrotik::readLength(
      reinterpret_cast<char*>(chars.data()), chars.size());
  EXPECT_EQ(wordLength, len.contentLength);
  EXPECT_EQ(chars.size(), len.lengthSize);
}

TEST_F(MikrotikChannelTest, lengthRead) {
  checkReadLength(0x00, {0x00});
  checkReadLength(0x01, {0x01});
  checkReadLength(0x02, {0x02});
  checkReadLength(0x7F, {0x7F});
  checkReadLength(0x80, {0x80, 0x80});
  checkReadLength(0x81, {0x80, 0x81});
  checkReadLength(0x3FFF, {0xBF, 0xFF});
  checkReadLength(0x4000, {0xC0, 0x40, 0x00});
  checkReadLength(0x4001, {0xC0, 0x40, 0x01});
  checkReadLength(0x1FFFFF, {0xDF, 0xFF, 0xFF});
  checkReadLength(0x200000, {0xE0, 0x20, 0x00, 0x00});
  checkReadLength(0x200001, {0xE0, 0x20, 0x00, 0x01});
  checkReadLength(0xFFFFFFF, {0xEF, 0xFF, 0xFF, 0xFF});
  checkReadLength(0x10000000, {0xF0, 0x10, 0x00, 0x00, 0x00});
  checkReadLength(0x10000001, {0xF0, 0x10, 0x00, 0x00, 0x01});
  checkReadLength(0xFFFFFFFF, {0xF0, 0xFF, 0xFF, 0xFF, 0xFF});
}

const std::string loginSentence = std::string(
    "\x6"
    "/login\x9"
    "=name=foo\xD"
    "=password=bar"
    "\0",
    32);

TEST_F(MikrotikChannelTest, checkLoginInit) {
  listen();
  folly::SocketAddress address("127.0.0.1", 1337);
  auto channel = std::make_shared<channels::mikrotik::Channel>(
      eventBase, address, "foo", "bar");
  expectedReads = loginSentence;
  channel->connect();
  acceptNotifier.wait();
  EXPECT_BECOMES_TRUE(channel->isConnected());
  EXPECT_NE(nullptr, asocket);
  stop();
}

TEST_F(MikrotikChannelTest, checkLoginBadResponse) {
  listen();
  folly::SocketAddress address("127.0.0.1", 1337);
  auto channel = std::make_shared<channels::mikrotik::Channel>(
      eventBase, address, "foo", "bar");
  expectedReads = loginSentence;
  channel->connect();
  acceptNotifier.wait();
  EXPECT_BECOMES_TRUE(channel->isConnected());
  EXPECT_NE(nullptr, asocket);
  write(
      "\x03"
      "foobar");
  // TODO fix this once it disconnects
  EXPECT_BECOMES_TRUE(channel->isConnected());
  stop();
}

TEST_F(MikrotikChannelTest, checkLoginInsufficientLength) {
  listen();
  folly::SocketAddress address("127.0.0.1", 1337);
  auto channel = std::make_shared<channels::mikrotik::Channel>(
      eventBase, address, "foo", "bar");
  expectedReads = loginSentence;
  channel->connect();
  acceptNotifier.wait();
  EXPECT_BECOMES_TRUE(channel->isConnected());
  EXPECT_NE(nullptr, asocket);
  write(
      "\xF0"
      "bar");
  EXPECT_TRUE(channel->isConnected());
  stop();
}

TEST_F(MikrotikChannelTest, checkLoginInsufficientData) {
  listen();
  folly::SocketAddress address("127.0.0.1", 1337);
  auto channel = std::make_shared<channels::mikrotik::Channel>(
      eventBase, address, "foo", "bar");
  expectedReads = loginSentence;
  channel->connect();
  acceptNotifier.wait();
  EXPECT_BECOMES_TRUE(channel->isConnected());
  EXPECT_NE(nullptr, asocket);
  write(
      "\xF0"
      "barfoofoofoofoofofooffofofoofofofofofoofof");
  EXPECT_TRUE(channel->isConnected());
  stop();
}

const std::string loginSentenceResponse = std::string("\x05!done\0", 7);

TEST_F(MikrotikChannelTest, checkLoginSuccess) {
  listen();
  folly::SocketAddress address("127.0.0.1", 1337);
  auto channel = std::make_shared<channels::mikrotik::Channel>(
      eventBase, address, "foo", "bar");
  expectedReads = loginSentence;
  channel->connect();
  acceptNotifier.wait();
  EXPECT_BECOMES_TRUE(channel->isConnected());
  EXPECT_NE(nullptr, asocket);
  write(loginSentenceResponse);
  EXPECT_BECOMES_TRUE(channel->isLoggedIn());
  EXPECT_TRUE(channel->isConnected());
  stop();
}

const std::string loginSentenceDeprecated = std::string(
    "\x6"
    "/login\x9"
    "=name=foo\x2C"
    "=response=0034FF0D6345057401A738C0B34BBD1744"
    "\0",
    63);

const std::string loginSentenceResponseDeprecated =
    std::string("\x05!done\x25=ret=e6ed6b07fcccc79a315677f0d1fce561\0", 45);

TEST_F(MikrotikChannelTest, checkLoginSuccessDeprecated) {
  listen();
  folly::SocketAddress address("127.0.0.1", 1337);
  auto channel = std::make_shared<channels::mikrotik::Channel>(
      eventBase, address, "foo", "bar");
  expectedReads = loginSentence;
  channel->connect();
  acceptNotifier.wait();
  EXPECT_BECOMES_TRUE(channel->isConnected());
  EXPECT_NE(nullptr, asocket);
  expectedReads += loginSentenceDeprecated;
  write(loginSentenceResponseDeprecated);
  EXPECT_BECOMES_TRUE(
      channel->getOperationalDatastore() ==
      channels::mikrotik::State::LoggingInDeprecated);
  write(loginSentenceResponse);
  EXPECT_BECOMES_TRUE(channel->isLoggedIn());
  EXPECT_TRUE(channel->isConnected());
  stop();
}

} // namespace test
} // namespace devmand
