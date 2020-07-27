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

#include <experimental/filesystem>
#include <string>

#include <folly/dynamic.h>

namespace devmand {

// So this could potentially be much better by using ydk or some similar
// construct. For the time being just use a using.
using YangPath = std::experimental::filesystem::path;

class YangUtils final {
 public:
  YangUtils() = delete;
  ~YangUtils() = delete;
  YangUtils(const YangUtils&) = delete;
  YangUtils& operator=(const YangUtils&) = delete;
  YangUtils(YangUtils&&) = delete;
  YangUtils& operator=(YangUtils&&) = delete;

 public:
  // TODO convert lookup to return a ptr so we dont need to copy
  // Looks up the value in the yang given as specified by the path.
  static folly::dynamic lookup(
      const folly::dynamic& yang,
      const YangPath& path);

  // Sets the value in the yang given at the specified path. This does not
  // create any of the parent objects in the yang.
  static void
  set(folly::dynamic& yang, const YangPath& path, const folly::dynamic& value);
};

} // namespace devmand
