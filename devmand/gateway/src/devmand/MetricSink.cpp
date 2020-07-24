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

#include <devmand/MetricSink.h>

namespace devmand {

void MetricSink::setGauge(const std::string& key, int value) {
  setGauge(key, static_cast<double>(value), "", "");
}

void MetricSink::setGauge(const std::string& key, size_t value) {
  setGauge(key, static_cast<double>(value), "", "");
}

void MetricSink::setGauge(const std::string& key, unsigned int value) {
  setGauge(key, static_cast<double>(value), "", "");
}

void MetricSink::setGauge(
    const std::string& key,
    long long unsigned int value) {
  setGauge(key, static_cast<double>(value), "", "");
}

void MetricSink::setGauge(const std::string& key, double value) {
  setGauge(key, value, "", "");
}

} // namespace devmand
