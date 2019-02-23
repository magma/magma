/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include "MetricsSingleton.h"

using magma::service303::MetricsRegistry;
using prometheus::Registry;
using prometheus::BuildCounter;
using prometheus::BuildGauge;
using prometheus::BuildHistogram;
using magma::service303::MetricsSingleton;

MetricsSingleton* MetricsSingleton::instance_ = NULL;

MetricsSingleton& MetricsSingleton::Instance() {
  if (instance_ == NULL) {
    instance_ = new MetricsSingleton();
  }
  return *instance_;
}

void MetricsSingleton::flush() {
  delete instance_;
  instance_ = new MetricsSingleton();
}

MetricsSingleton::MetricsSingleton() :
  registry_(std::make_shared<Registry>()),
  counters_(registry_, BuildCounter),
  gauges_(registry_, BuildGauge),
  histograms_(registry_, BuildHistogram) {

}

void MetricsSingleton::args_to_map(std::map<std::string, std::string>& labels, size_t label_count, va_list& args) {
  for (size_t i = 0; i < label_count; i++) {
    labels.insert({{va_arg(args, const char *), va_arg(args, const char *)}});
  }
}

void MetricsSingleton::IncrementCounter(const char* name,
  double increment,
  size_t label_count,
  va_list& args) {
  std::map<std::string, std::string> labels;
  args_to_map(labels, label_count, args);
  counters_.Get(name, labels).Increment(increment);
}

void MetricsSingleton::IncrementGauge(const char* name,
  double increment,
  size_t label_count,
  va_list& args) {
  std::map<std::string, std::string> labels;
  args_to_map(labels, label_count, args);
  gauges_.Get(name, labels).Increment(increment);
}

void MetricsSingleton::DecrementGauge(const char* name,
  double decrement,
  size_t label_count,
  va_list& args) {
  std::map<std::string, std::string> labels;
  args_to_map(labels, label_count, args);
  gauges_.Get(name, labels).Decrement(decrement);
}

void MetricsSingleton::SetGauge(const char* name,
  double value,
  size_t label_count,
  va_list& args) {
  std::map<std::string, std::string> labels;
  args_to_map(labels, label_count, args);
  gauges_.Get(name, labels).Set(value);
}

void MetricsSingleton::ObserveHistogram(const char* name,
  double observation,
  size_t label_count,
  va_list &args) {
  std::map<std::string, std::string> labels;
  args_to_map(labels, label_count, args);

  size_t boundary_count = va_arg(args, size_t);
  std::vector<double> boundaries;
  for (size_t i = 0; i < boundary_count; i++) {
    boundaries.push_back(va_arg(args, double));
  }
  histograms_.Get(name, labels, Histogram::BucketBoundaries(boundaries)).Observe(observation);
}
