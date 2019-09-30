// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/models/interface/Model.h>

namespace devmand {
namespace models {
namespace interface {

void Model::init(folly::dynamic& state) {
  auto& interfaces = state["openconfig-interfaces:interfaces"] =
      folly::dynamic::object;
  interfaces["interface"] = folly::dynamic::array;
}

void Model::updateInterface(
    folly::dynamic& state,
    int index,
    const YangPath& path,
    const folly::dynamic& value) {
  auto& interfaces = state["openconfig-interfaces:interfaces"];
  for (auto& interface : interfaces["interface"]) {
    if (interface["state"]["ifindex"] == index) {
      YangUtils::set(interface, path, value);
      return;
    }
  }

  folly::dynamic interface = folly::dynamic::object;
  auto& istate = interface["state"] = folly::dynamic::object;
  istate["counters"] = folly::dynamic::object;
  auto& config = interface["config"] = folly::dynamic::object;
  interface["name"] = config["name"] = folly::to<std::string>(index);
  istate["ifindex"] = index;
  YangUtils::set(interface, path, value);
  interfaces["interface"].push_back(interface);
}

} // namespace interface
} // namespace models
} // namespace devmand
