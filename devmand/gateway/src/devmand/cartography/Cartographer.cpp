// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

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
