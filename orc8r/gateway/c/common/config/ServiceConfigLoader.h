/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <string>
#include "yaml-cpp/yaml.h"

namespace magma {

/**
 * ServiceConfigLoader is a helper class to parse proc files for process
 * information
 */
class ServiceConfigLoader final {

  public:
    /*
     * Load service configuration from file.
     *
     * @return YAML::Node a Node representation of the file.
     */
    YAML::Node load_service_config(const std::string& service_name);


  private:
    static constexpr const char* CONFIG_DIR = "/etc/magma/";
    static constexpr const char* OVERRIDE_DIR = "/var/opt/magma/configs/";

};

}
