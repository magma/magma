// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#pragma once

#include <string>

namespace devmand {

/* An abstraction of a class which handles metrics.
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
  virtual void setGauge(const std::string& key, double value) = 0;
};

} // namespace devmand
