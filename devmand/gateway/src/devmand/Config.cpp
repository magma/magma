// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#include <devmand/Config.h>

namespace devmand {

DEFINE_string(listen_interface, "eth0", "The interface to listen on.");
DEFINE_string(
    device_configuration_file,
    "/etc/devmand/devices.yml",
    "Accepts .yml or .mconfig files. Inotify watches the file, and applies necessary changes.");
DEFINE_uint64(poll_interval, 120, "The polling interval in seconds.");

} // namespace devmand
