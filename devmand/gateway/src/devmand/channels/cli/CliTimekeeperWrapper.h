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

#include <boost/thread/mutex.hpp>
#include <devmand/channels/cli/CancelableWTCallback.h>
#include <devmand/channels/cli/CliThreadWheelTimekeeper.h>
#include <folly/Unit.h>

namespace devmand::channels::cli {

using devmand::channels::cli::CancelableWTCallback;
using folly::Duration;
using folly::Future;
using folly::SemiFuture;
using folly::Timekeeper;
using folly::Unit;
using folly::unit;
using std::shared_ptr;

class CliTimekeeperWrapper : public Timekeeper {
 private:
  shared_ptr<CliThreadWheelTimekeeper> timekeeper;
  boost::mutex mutex;
  shared_ptr<CancelableWTCallback> cb;

 public:
  CliTimekeeperWrapper(const shared_ptr<CliThreadWheelTimekeeper>& timekeeper);

  ~CliTimekeeperWrapper();

 public:
  const shared_ptr<CliThreadWheelTimekeeper>& getTimekeeper() const;

  void cancelAll();

  void setCurrentSleepCallback(shared_ptr<CancelableWTCallback> _cb);

  Future<Unit> after(Duration dur) override;
};

} // namespace devmand::channels::cli
