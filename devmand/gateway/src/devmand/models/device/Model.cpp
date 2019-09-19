// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/models/device/Model.h>

namespace devmand {
namespace models {
namespace device {

void Model::init(folly::dynamic& state) {
  auto& system = state["fbc-symphony-device:system"] = folly::dynamic::object;

  // ietf-geo-location ########################################################
  // Inits all of these to defaults. Some are left out as they are to be filled
  // in by device but these are left in the file to document them.
  auto& geol = system["geo-location"] = folly::dynamic::object;
  auto& rf = geol["reference-frame"] = folly::dynamic::object;
  rf["astronomical-body"] = "earth";
  auto& gs = rf["geodetic-system"] = folly::dynamic::object;
  gs["geodetic-datum"] = "wgs-84";
  // gs["coord-accuracy"] = 0;
  // gs["height-accuracy"] = 0;
  geol["latitude"] = 0;
  geol["longitude"] = 0;
  geol["height"] = 0;
  // auto& vel = geol["velocity"] = folly::dynamic::object;
  // vel["v-north"] = 0;
  // vel["v-east"] = 0;
  // vel["v-up"] = 0;
  // geol["timestamp"] = 0;
}

} // namespace device
} // namespace models
} // namespace devmand
