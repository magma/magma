/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <stdarg.h>

#include <prometheus/registry.h>
#include <grpc++/grpc++.h>

#include "MetricsRegistry.h"

using magma::service303::MetricsRegistry;
using prometheus::Registry;
using prometheus::Counter;
using prometheus::detail::CounterBuilder;
using prometheus::Gauge;
using prometheus::detail::GaugeBuilder;
using prometheus::Histogram;
using prometheus::detail::HistogramBuilder;
using grpc::Server;

namespace magma { namespace service303 {

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
    static void flush(); // destroy instance
    void IncrementCounter(const char* name,
      double increment,
      size_t label_count,
      va_list& args);
    void IncrementGauge(const char* name,
      double increment,
      size_t label_count,
      va_list& args);
    void DecrementGauge(const char* name,
      double decrement,
      size_t label_count,
      va_list& args);
    void SetGauge(const char* name,
      double value,
      size_t label_count,
      va_list& args);
    void ObserveHistogram(const char* name,
      double observation,
      size_t label_count,
      va_list& args);
  private:
    MetricsSingleton(); // Prevent construction
    MetricsSingleton(const MetricsSingleton&); // Prevent construction by copying
    MetricsSingleton& operator=(const MetricsSingleton&); // Prevent assignment
    void args_to_map(std::map<std::string, std::string>& labels, size_t label_count, va_list& args); // Helper to convert variadic labels to map
    // Shared registry to store all our metrics
    std::shared_ptr<prometheus::Registry> registry_;
    // Dictionaries to store instances of our metrics and intialize new ones
    MetricsRegistry<Counter, CounterBuilder (&)()> counters_;
    MetricsRegistry<Gauge, GaugeBuilder (&)()> gauges_;
    MetricsRegistry<Histogram, HistogramBuilder (&)()> histograms_;
    static MetricsSingleton* instance_;
};

}}
