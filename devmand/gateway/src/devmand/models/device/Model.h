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

#include <folly/dynamic.h>

#include <devmand/channels/ping/Engine.h>

namespace devmand {
namespace models {
namespace device {

class Model {
 public:
  Model() = delete;
  ~Model() = delete;
  Model(const Model&) = delete;
  Model& operator=(const Model&) = delete;
  Model(Model&&) = delete;
  Model& operator=(Model&&) = delete;

 public:
  static void init(folly::dynamic& state);

  static void addLatency(
      folly::dynamic& state,
      const std::string& type,
      const std::string& src,
      const std::string& dst,
      channels::ping::Rtt rtt);
};

} // namespace device
} // namespace models
} // namespace devmand
