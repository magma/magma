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

#include <devmand/fscache/Service.h>

namespace devmand {
namespace fscache {

Service::Service(Application& application)
    : ::devmand::Service(application) { // TODO
}

void Service::setGauge(
    const std::string&,
    double,
    const std::string&,
    const std::string&) {}

void Service::start() {}

void Service::wait() {}

void Service::stop() {}

} // namespace fscache
} // namespace devmand
