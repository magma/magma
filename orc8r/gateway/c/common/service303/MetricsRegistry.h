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

#include <unordered_map>

#include <prometheus/registry.h>
#include <prometheus/family.h>
#include <prometheus/counter.h>
#include <prometheus/gauge.h>
#include <prometheus/histogram.h>

#include <prometheus/metrics.pb.h>
#include <orc8r/protos/metricsd.pb.h>

using prometheus::Family;
using prometheus::Registry;

namespace magma {
namespace service303 {
using namespace orc8r;

/**
 * MetricsRegistry is a dictionary for metrics instances. It ensures we
 * constuct a single instance of a metric family per name and a single
 * instance for each label set in that family.
 */
template<typename T, typename MetricFamilyFactory>
class MetricsRegistry {
 public:
  MetricsRegistry(
      const std::shared_ptr<prometheus::Registry>& registry,
      const MetricFamilyFactory& factory);

  /**
   * Get or create a metric instance matching this name and label set
   *
   * @param name: the metric name
   * @param labels: list of tuples denoting label key value pairs
   * @param args...: other arguments the Metric constructor may need
   * @return prometheus T instance
   */
  template<typename... Args>
  T& Get(
      const std::string& name, const std::map<std::string, std::string>& labels,
      Args&&... args);

  /**
   * Remove a metric instance specified by name/labels
   * @param name
   * @param labels
   */
  void Remove(
      const std::string& name,
      const std::map<std::string, std::string>& labels);

  const std::size_t SizeFamilies() { return families_.size(); }

  const std::size_t SizeMetrics() { return metrics_.size(); }

 private:
  static std::size_t hash_name_and_labels(
      const std::string& name,
      const std::map<std::string, std::string>& labels);
  // Convert labels to enums if applicable
  static void parse_labels(
      const std::map<std::string, std::string>& labels,
      std::map<std::string, std::string>& parsed_labels);
  std::unordered_map<std::size_t, Family<T>*> families_;
  std::unordered_map<std::size_t, T*> metrics_;
  const std::shared_ptr<prometheus::Registry>& registry_;
  const MetricFamilyFactory& factory_;
};

template<typename T, typename MetricFamilyFactory>
MetricsRegistry<T, MetricFamilyFactory>::MetricsRegistry(
    const std::shared_ptr<prometheus::Registry>& registry,
    const MetricFamilyFactory& factory)
    : registry_(registry), factory_(factory) {}

template<typename T, typename MetricFamilyFactory>
template<typename... Args>
T& MetricsRegistry<T, MetricFamilyFactory>::Get(
    const std::string& name, const std::map<std::string, std::string>& labels,
    Args&&... args) {
  // Create the family if we haven't seen it before
  Family<T>* family;
  size_t name_hash = std::hash<std::string>{}(name);
  auto family_it   = families_.find(name_hash);
  if (family_it != families_.end()) {
    family = family_it->second;
  } else {
    // If the name is a defined MetricName, use the enum value instead
    MetricName name_value;
    const std::string& proto_name =
        MetricName_Parse(name, &name_value) ? std::to_string(name_value) : name;
    // Factory constructs the metric on the heap and adds it to registry_
    family = &factory_().Name(proto_name).Register(*registry_);
    families_.insert({{name_hash, family}});
  }

  // Create the metric if we haven't seen it before
  T* metric;
  size_t metric_hash = hash_name_and_labels(name, labels);
  auto metric_it     = metrics_.find(metric_hash);
  if (metric_it != metrics_.end()) {
    metric = metric_it->second;
  } else {
    std::map<std::string, std::string> converted_labels;
    parse_labels(labels, converted_labels);
    metric = &family->Add(converted_labels, std::forward<Args>(args)...);
    metrics_.insert({{metric_hash, metric}});
  }
  return *metric;
}

template<typename T, typename MetricFamilyFactory>
void MetricsRegistry<T, MetricFamilyFactory>::Remove(
    const std::string& name, const std::map<std::string, std::string>& labels) {
  Family<T>* family;
  size_t name_hash = std::hash<std::string>{}(name);
  auto family_it   = families_.find(name_hash);
  if (family_it == families_.end()) {
    return;
  }
  family = family_it->second;

  T* metric;
  size_t metric_hash = hash_name_and_labels(name, labels);
  auto metric_it     = metrics_.find(metric_hash);
  if (metric_it == metrics_.end()) {
    return;
  }
  metric = metric_it->second;
  family->Remove(metric);
  metrics_.erase(metric_hash);
}

template<typename T, typename MetricFamilyFactory>
std::size_t MetricsRegistry<T, MetricFamilyFactory>::hash_name_and_labels(
    const std::string& name, const std::map<std::string, std::string>& labels) {
  auto combined = std::accumulate(
      labels.begin(), labels.end(), std::string{name},
      [](const std::string& acc,
         const std::pair<std::string, std::string>& label_pair) {
        return acc + label_pair.first + label_pair.second;
      });
  return std::hash<std::string>{}(combined);
}

template<typename T, typename MetricFamilyFactory>
void MetricsRegistry<T, MetricFamilyFactory>::parse_labels(
    const std::map<std::string, std::string>& labels,
    std::map<std::string, std::string>& parsed_labels) {
  for (const auto& label_pair : labels) {
    // convert label name
    MetricLabelName label_name_enum;
    const std::string& label_name =
        MetricLabelName_Parse(label_pair.first, &label_name_enum) ?
            std::to_string(label_name_enum) :
            label_pair.first;
    // insert into new map
    parsed_labels.insert({{label_name, label_pair.second}});
  }
}

}  // namespace service303
}  // namespace magma
