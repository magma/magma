// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <system_error>

#include <devmand/utils/FileWatcher.h>

namespace devmand {

FileWatcher::FileWatcher(folly::EventBase& _eventBase)
    : folly::EventHandler(&_eventBase),
      eventBase(_eventBase),
      inotifyFd(inotify_init1(IN_NONBLOCK | IN_CLOEXEC)) {
  if (inotifyFd < 0) {
    LOG(ERROR) << "failed creation of inotify";
    throw std::runtime_error("failed creation of inotify");
  } else {
    this->folly::EventHandler::changeHandlerFD(
        folly::NetworkSocket::fromFd(inotifyFd));

    registerHandler(folly::EventHandler::READ | folly::EventHandler::PERSIST);
  }
}

FileWatcher::~FileWatcher() {
  unregisterHandler();
  for (auto& watch : watches) {
    removeWatch(watch.second.fd);
  }
}

bool FileWatcher::addWatch(
    const std::string& filename,
    ActionCallback action,
    bool shouldReadInitial,
    uint32_t eventMask) {
  FileWatch watch;
  watch.filename = filename;
  watch.action = action;
  watch.fd = inotify_add_watch(inotifyFd, filename.c_str(), eventMask);
  if (watch.fd < 0) {
    LOG(ERROR) << "failed creation of watch " << filename << " error "
               << std::error_code(errno, std::system_category());
    return false;
  }

  if (not watches.emplace(watch.fd, watch).second) {
    LOG(ERROR) << "failed to add watch " << filename;
    return false;
  }

  if (shouldReadInitial) {
    if (not FileUtils::touch(filename)) {
      LOG(ERROR) << "failed initial read touch " << filename;
      return false;
    }
  }
  return true;
}

void FileWatcher::handlerReady(uint16_t events) noexcept {
  (void)events;

  static constexpr unsigned int bufferSize =
      (10 * (sizeof(struct inotify_event) + NAME_MAX + 1));
  char buf[bufferSize];
  ssize_t numRead = read(inotifyFd, buf, bufferSize);
  if (numRead <= 0) {
    LOG(ERROR) << "expected to find inotify event but none were present";
  }

  for (char* p = buf; p < buf + numRead;) {
    auto* event = reinterpret_cast<inotify_event*>(p);

    auto watch = watches.find(event->wd);
    if (watch == watches.end()) {
      LOG(ERROR) << "watch not found " << event->wd;
    } else {
      uint32_t mask{1};
      for (int i = 0; i < 32; ++i) {
        if (event->mask & mask) {
          if (mask == IN_IGNORED) { // Ignore watch removes
            continue;
          }
          FileWatchEvent watchEvent;
          watchEvent.event = static_cast<FileEvent>(mask);
          watchEvent.filename = event->len == 0 ? "" : std::string(event->name);
          // TODO this breaks the set down into many events. Maybe it is best
          // to send a bitset abstraction...
          watch->second.action(std::move(watchEvent));
        }
        mask = mask << 1;
      }
    }

    p += sizeof(struct inotify_event) + event->len;
  }
}

void FileWatcher::removeWatch(int watchFd) {
  if (inotify_rm_watch(inotifyFd, watchFd) < 0) {
    LOG(ERROR) << "removed watch " << watchFd;
  }
}

} // namespace devmand
