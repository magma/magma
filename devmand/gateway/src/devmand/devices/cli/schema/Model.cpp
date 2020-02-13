// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/cli/schema/Model.h>

namespace devmand {
namespace devices {
namespace cli {

const Model Model::OPENCONFIG_0_1_6 = Model("/usr/share/openconfig@0.1.6");
const Model Model::IETF_0_1_5 = Model("/usr/share/ietf@0.1.5");

} // namespace cli
} // namespace devices
} // namespace devmand
