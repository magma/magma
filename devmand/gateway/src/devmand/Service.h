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

#include <devmand/MetricSink.h>

namespace devmand {

class Application;

/* An abstraction of services the devmand application provides separated for
 * linkage. A service in devmand takes the unified view of devices and exports
 * them to a north bound layer.
 */
class Service : public MetricSink {
 public:
  Service(Application& application);
  Service() = delete;
  virtual ~Service() = default;
  Service(const Service&) = delete;
  Service& operator=(const Service&) = delete;
  Service(Service&&) = delete;
  Service& operator=(Service&&) = delete;

 public:
  virtual void start() = 0;
  virtual void wait() = 0;
  virtual void stop() = 0;

 protected:
  Application& app;
};

} // namespace devmand
