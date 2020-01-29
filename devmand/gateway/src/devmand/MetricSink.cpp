// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

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
