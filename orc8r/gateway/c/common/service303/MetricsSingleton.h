/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#pragma once

#include <stdarg.h>

#include <prometheus/registry.h>
#include <grpc++/grpc++.h>

#include "MetricsRegistry.h"

using grpc::Server;
using magma::service303::MetricsRegistry;
using prometheus::Counter;
using prometheus::Gauge;
using prometheus::Histogram;
using prometheus::Registry;
using prometheus::detail::CounterBuilder;
using prometheus::detail::GaugeBuilder;
using prometheus::detail::HistogramBuilder;

namespace magma {
namespace service303 {

// Forward decleration
class MetricsSingleton;

/*
 * MetricsSingleton is a singleton used to contain metrics registries and
 * interfaces to interact with unique prometheus timeseries each uniquely
 * defined by a family name, and a set of labels.
 */
class MetricsSingleton {
  friend class MagmaService;

 public:
  static MetricsSingleton& Instance();
  static void flush();  // destroy instance
  void RemoveCounter(const char* name, size_t label_count, va_list& args);
  void IncrementCounter(
      const char* name, double increment, size_t label_count, va_list& args);
  void IncrementGauge(
      const char* name, double increment, size_t label_count, va_list& args);
  void DecrementGauge(
      const char* name, double decrement, size_t label_count, va_list& args);
  void SetGauge(
      const char* name, double value, size_t label_count, va_list& args);
  void ObserveHistogram(
      const char* name, double observation, size_t label_count, va_list& args);

 private:
  MetricsSingleton();                         // Prevent construction
  MetricsSingleton(const MetricsSingleton&);  // Prevent construction by copying
  MetricsSingleton& operator=(const MetricsSingleton&);  // Prevent assignment
  void args_to_map(
      std::map<std::string, std::string>& labels, size_t label_count,
      va_list& args);  // Helper to convert variadic labels to map
  // Shared registry to store all our metrics
  std::shared_ptr<prometheus::Registry> registry_;
  // Dictionaries to store instances of our metrics and intialize new ones
  MetricsRegistry<Counter, CounterBuilder (&)()> counters_;
  MetricsRegistry<Gauge, GaugeBuilder (&)()> gauges_;
  MetricsRegistry<Histogram, HistogramBuilder (&)()> histograms_;
  static MetricsSingleton* instance_;
};

}  // namespace service303
}  // namespace magma
