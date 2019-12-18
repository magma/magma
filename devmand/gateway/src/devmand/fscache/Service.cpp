// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

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
