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

#include <list>
#include <string>

#include <folly/Synchronized.h>
#include <folly/dynamic.h>

namespace devmand {

// TODO add timestamps
class ErrorQueue final {
 public:
  ErrorQueue(unsigned int maxSize_ = 10) : maxSize(maxSize_) {}
  ~ErrorQueue() = default;
  ErrorQueue(const ErrorQueue&) = delete;
  ErrorQueue& operator=(const ErrorQueue&) = delete;
  ErrorQueue(ErrorQueue&&) = delete;
  ErrorQueue& operator=(ErrorQueue&&) = delete;

 public:
  void add(std::string&& error);
  folly::dynamic get();

 private:
  folly::Synchronized<std::list<std::string>> errors;
  unsigned int maxSize{0};
};

} // namespace devmand
