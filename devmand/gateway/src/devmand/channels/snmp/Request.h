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

#include <folly/futures/Future.h>

#include <devmand/channels/snmp/Snmp.h>

namespace devmand {
namespace channels {
namespace snmp {

class Channel;

class Request final {
 public:
  Request(Channel* channel_, Oid oid_);
  Request() = delete;
  ~Request() = default;
  Request(const Request&) = delete;
  Request& operator=(const Request&) = delete;
  Request(Request&&) = default;
  Request& operator=(Request&&) = delete;

 public:
  Channel* channel{nullptr};
  Oid oid;
  folly::Promise<Response> responsePromise{};
};

} // namespace snmp
} // namespace channels
} // namespace devmand
