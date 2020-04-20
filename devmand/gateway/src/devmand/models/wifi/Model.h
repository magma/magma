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

#include <devmand/utils/YangUtils.h>

namespace devmand {
namespace models {
namespace wifi {

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
  static void updateRadio(
      folly::dynamic& state,
      int index,
      const YangPath& path,
      const folly::dynamic& value);

  static void updateSsid(
      folly::dynamic& state,
      int index,
      const YangPath& path,
      const folly::dynamic& value);

  static void updateSsidBssid(
      folly::dynamic& state,
      int indexSsid,
      int indexBssid,
      const YangPath& path,
      const folly::dynamic& value);
};

} // namespace wifi
} // namespace models
} // namespace devmand
