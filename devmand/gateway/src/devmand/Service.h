// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

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
