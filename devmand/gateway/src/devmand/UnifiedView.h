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

#include <map>
#include <string>

#include <folly/dynamic.h>

#include <devmand/devices/Id.h>

namespace devmand {

// TODO convert this to ydk?
using YangModelBundle = folly::dynamic;
using UnifiedView = std::map<devices::Id, YangModelBundle>;
using SharedUnifiedView = folly::Synchronized<UnifiedView>;

} // namespace devmand
