#include "EventBaseWrapper.h"

namespace magma {
void EventBaseWrapper::loopForever() {
  evb_->loopForever();
};

void EventBaseWrapper::terminateLoopSoon() {
  evb_->terminateLoopSoon();
};

void EventBaseWrapper::runAfterDelay(
    folly::Function<void()> func, int32_t delayMs) {
  evb_->runAfterDelay(std::move(func), delayMs);
};

void EventBaseWrapper::runInEventBaseThread(folly::Cob&& cob) {
  evb_->runInEventBaseThread(std::move(cob));
}
};  // namespace magma