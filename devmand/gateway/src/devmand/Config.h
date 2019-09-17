// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#pragma once

#include <gflags/gflags.h>

namespace devmand {

DECLARE_string(listen_interface);
DECLARE_string(device_configuration_file);
DECLARE_uint64(poll_interval);

} // namespace devmand
