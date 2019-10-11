// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/cartography/Cartographer.h>

namespace devmand {
namespace cartography {

Cartographer::Cartographer(
    const AddHandler& addHandler,
    const DeleteHandler& deleteHandler)
    : add(addHandler), del(deleteHandler) {
  assert(add != nullptr);
  assert(del != nullptr);
}

void Cartographer::addDeviceDiscoveryMethod(
    const std::shared_ptr<Method>& method) {
  assert(method != nullptr);
  auto result = methods.emplace(method);
  if (result.second) {
    method->setHandlers(add, del);
    method->enable();
  } else {
    LOG(ERROR) << "Failed to add device discovery method";
    throw std::runtime_error("Failed to add device discovery method.");
  }
}

} // namespace cartography
} // namespace devmand
