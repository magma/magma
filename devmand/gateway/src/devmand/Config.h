// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <gflags/gflags.h>

namespace devmand {

DECLARE_string(listen_interface);
DECLARE_string(device_configuration_file);
DECLARE_uint64(poll_interval);
DECLARE_uint64(debug_print_interval);
DECLARE_bool(devices_readonly);

} // namespace devmand
