// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <functional>
#include <memory>

#include <devmand/cartography/DeviceConfig.h>
#include <devmand/devices/Device.h>

namespace devmand {
namespace cartography {

using AddHandler = std::function<void(const DeviceConfig&)>;
using DeleteHandler = std::function<void(const DeviceConfig&)>;
// TODO no support for modify yet

/*
 * An abstract class which represents a single way to map devices on a network.
 */
class Method : public std::enable_shared_from_this<Method> {
 public:
  Method() = default;
  virtual ~Method() = default;
  Method(const Method&) = delete;
  Method& operator=(const Method&) = delete;
  Method(Method&&) = delete;
  Method& operator=(Method&&) = delete;

 public:
  virtual void enable() = 0;
  void setHandlers(
      const AddHandler& addHandler,
      const DeleteHandler& deleteHandler);

 protected:
  AddHandler add;
  DeleteHandler del;
};

} // namespace cartography
} // namespace devmand
