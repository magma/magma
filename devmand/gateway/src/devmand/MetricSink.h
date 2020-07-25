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

#include <string>

namespace devmand {

/* An abstraction of a class which handles metrics. This represents a place for
 * metrics updates to go such as a time series database or a log.
 */
class MetricSink {
 public:
  MetricSink() = default;
  virtual ~MetricSink() = default;
  MetricSink(const MetricSink&) = delete;
  MetricSink& operator=(const MetricSink&) = delete;
  MetricSink(MetricSink&&) = delete;
  MetricSink& operator=(MetricSink&&) = delete;

 public:
  virtual void setGauge(
      const std::string& key,
      double value,
      const std::string& labelName,
      const std::string& labelValue) = 0;

  // Overloads
  void setGauge(const std::string& key, int value);
  void setGauge(const std::string& key, size_t value);
  void setGauge(const std::string& key, unsigned int value);
  void setGauge(const std::string& key, long long unsigned int value);
  void setGauge(const std::string& key, double value);
};

} // namespace devmand
