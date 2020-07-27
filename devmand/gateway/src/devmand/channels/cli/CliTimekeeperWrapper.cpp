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

#include <devmand/channels/cli/CliTimekeeperWrapper.h>

using devmand::channels::cli::CancelableWTCallback;
using folly::Duration;
using folly::Future;
using folly::Unit;
using std::shared_ptr;

namespace devmand::channels::cli {

void CliTimekeeperWrapper::setCurrentSleepCallback(
    shared_ptr<CancelableWTCallback> _cb) {
  boost::mutex::scoped_lock scoped_lock(this->mutex);
  this->cb = _cb;
}

Future<Unit> CliTimekeeperWrapper::after(Duration dur) {
  shared_ptr<CancelableWTCallback> callback = timekeeper->cancelableSleep(dur);
  setCurrentSleepCallback(callback);
  return callback->getFuture();
}

void CliTimekeeperWrapper::cancelAll() {
  if (cb.use_count() >= 1) {
    cb->callbackCanceled();
  }
}

CliTimekeeperWrapper::~CliTimekeeperWrapper() {
  cancelAll();
}

CliTimekeeperWrapper::CliTimekeeperWrapper(
    const shared_ptr<CliThreadWheelTimekeeper>& _timekeeper)
    : timekeeper(_timekeeper) {}

const shared_ptr<CliThreadWheelTimekeeper>&
CliTimekeeperWrapper::getTimekeeper() const {
  return timekeeper;
}
} // namespace devmand::channels::cli