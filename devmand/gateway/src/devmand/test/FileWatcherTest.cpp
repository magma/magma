// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <experimental/filesystem>
#include <iostream>

#include <devmand/test/EventBaseTest.h>
#include <devmand/test/Notifier.h>
#include <devmand/utils/FileWatcher.h>

namespace devmand {
namespace test {

// TODO: make filewatcher tests compatible with the circleci
// VMs and re-enable them.
class FileWatcherTest : public EventBaseTest {
 public:
  FileWatcherTest() = default;

  ~FileWatcherTest() override = default;
  FileWatcherTest(const FileWatcherTest&) = delete;
  FileWatcherTest& operator=(const FileWatcherTest&) = delete;
  FileWatcherTest(FileWatcherTest&&) = delete;
  FileWatcherTest& operator=(FileWatcherTest&&) = delete;

 protected:
  std::function<void(FileWatchEvent)> eventHandler =
      [this](FileWatchEvent watchEvent) {
        std::cerr << "Event " << static_cast<int>(watchEvent.event) << " with #"
                  << expectedNumEvents << " expected remaining" << std::endl;
        EXPECT_LE(1, expectedNumEvents);
        events.push_back(watchEvent);
        if (--expectedNumEvents == 0) {
          eventNotifier.notify();
        }
      };

  std::function<void(FileWatchEvent)> failEventHandler = [](FileWatchEvent) {
    FAIL() << "Unexpected event";
  };

  void expectEvent(FileWatchEvent watchEvent) {
    EXPECT_EQ(watchEvent.event, events.front().event);
    EXPECT_EQ(watchEvent.filename, events.front().filename);
    events.pop_front();
  }

  std::experimental::filesystem::path getFileToWatch(
      const std::string& dir = "") const {
    std::string ret = "/tmp/";
    if (not dir.empty()) {
      ret += dir + "/";
    }
    auto* testInfo = ::testing::UnitTest::GetInstance()->current_test_info();
    ret += testInfo->name();
    return ret;
  }

 protected:
  Notifier eventNotifier;
  std::list<FileWatchEvent> events;
  unsigned int expectedNumEvents = 0;
};

TEST_F(FileWatcherTest, DISABLED_InitialEvent) {
  FileWatcher watcher(eventBase);
  auto filepath = getFileToWatch();
  EXPECT_TRUE(FileUtils::touch(filepath));

  expectedNumEvents = 3;
  EXPECT_TRUE(watcher.addWatch(filepath, eventHandler, true));
  eventNotifier.wait();
  EXPECT_EQ(3, events.size());
  expectEvent({FileEvent::Open, ""});
  expectEvent({FileEvent::Attrib, ""});
  expectEvent({FileEvent::CloseWrite, ""});
}

TEST_F(FileWatcherTest, DISABLED_InitialNoEvent) {
  FileWatcher watcher(eventBase);
  FileWatcher watcher2(eventBase);
  auto filepath = getFileToWatch();
  auto filepath2 = filepath.native() + "2";
  EXPECT_TRUE(FileUtils::touch(filepath));
  EXPECT_TRUE(FileUtils::touch(filepath2));

  expectedNumEvents = 3;
  EXPECT_TRUE(watcher.addWatch(filepath, failEventHandler, false));
  EXPECT_TRUE(watcher2.addWatch(filepath2, eventHandler, true));
  eventNotifier.wait();
  EXPECT_EQ(3, events.size());
  expectEvent({FileEvent::Open, ""});
  expectEvent({FileEvent::Attrib, ""});
  expectEvent({FileEvent::CloseWrite, ""});
}

TEST_F(FileWatcherTest, DISABLED_EventsAfterInitial) {
  FileWatcher watcher(eventBase);
  auto filepath = getFileToWatch();
  EXPECT_TRUE(FileUtils::touch(filepath));

  expectedNumEvents = 3;
  EXPECT_TRUE(watcher.addWatch(filepath, eventHandler, true));
  eventNotifier.wait();
  EXPECT_EQ(3, events.size());
  expectEvent({FileEvent::Open, ""});
  expectEvent({FileEvent::Attrib, ""});
  expectEvent({FileEvent::CloseWrite, ""});

  expectedNumEvents = 3;
  EXPECT_TRUE(FileUtils::touch(filepath));
  eventNotifier.wait();
  EXPECT_EQ(3, events.size());
  expectEvent({FileEvent::Open, ""});
  expectEvent({FileEvent::Attrib, ""});
  expectEvent({FileEvent::CloseWrite, ""});
}

TEST_F(FileWatcherTest, DISABLED_InitialEventOnDir) {
  FileWatcher watcher(eventBase);
  auto filepath = getFileToWatch("fwtest");
  auto dirname = filepath.parent_path();
  EXPECT_TRUE(FileUtils::mkdir(dirname));
  EXPECT_TRUE(FileUtils::touch(filepath));

  expectedNumEvents = 2;
  EXPECT_TRUE(watcher.addWatch(dirname, eventHandler, true));
  eventNotifier.wait();
  EXPECT_EQ(2, events.size());
  expectEvent({FileEvent::Attrib, ""});
  expectEvent({FileEvent::IsDir, ""});
}

TEST_F(FileWatcherTest, DISABLED_EventsAfterInitialOnDir) {
  FileWatcher watcher(eventBase);
  auto filepath = getFileToWatch("fwtest");
  auto dirname = filepath.parent_path();
  EXPECT_TRUE(FileUtils::mkdir(dirname));
  EXPECT_TRUE(FileUtils::touch(filepath));

  expectedNumEvents = 2;
  EXPECT_TRUE(watcher.addWatch(dirname, eventHandler, true));
  eventNotifier.wait();
  EXPECT_EQ(2, events.size());
  expectEvent({FileEvent::Attrib, ""});
  expectEvent({FileEvent::IsDir, ""});

  expectedNumEvents = 3;
  EXPECT_TRUE(FileUtils::touch(filepath));
  eventNotifier.wait();
  EXPECT_EQ(3, events.size());
  expectEvent({FileEvent::Open, filepath.filename()});
  expectEvent({FileEvent::Attrib, filepath.filename()});
  expectEvent({FileEvent::CloseWrite, filepath.filename()});
}

} // namespace test
} // namespace devmand
