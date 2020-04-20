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

#include <devmand/channels/cli/CancelableWTCallback.h>
#include <folly/futures/ThreadWheelTimekeeper.h>

namespace devmand::channels::cli {

using devmand::channels::cli::CancelableWTCallback;
using folly::Duration;
using folly::Future;
using folly::SemiFuture;
using folly::ThreadWheelTimekeeper;
using folly::Unit;
using std::shared_ptr;

class CliThreadWheelTimekeeper : public ThreadWheelTimekeeper {
 public:
  shared_ptr<CancelableWTCallback> cancelableSleep(Duration);
};

} // namespace devmand::channels::cli
