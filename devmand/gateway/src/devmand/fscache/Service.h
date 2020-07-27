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

#include <devmand/Service.h>

namespace devmand {
namespace fscache {

class Service : public ::devmand::Service {
 public:
  Service(Application& application);
  Service() = delete;
  virtual ~Service() = default;
  Service(const Service&) = delete;
  Service& operator=(const Service&) = delete;
  Service(Service&&) = delete;
  Service& operator=(Service&&) = delete;

 public:
  void start() override;
  void wait() override;
  void stop() override;
  void setGauge(
      const std::string& key,
      double value,
      const std::string& labelName,
      const std::string& labelValue) override;

 private:
};

} // namespace fscache
} // namespace devmand
