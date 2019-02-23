/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <string>
#include <iostream>


#include "ServiceConfigLoader.h"
#include "YAMLUtils.h"
#include "magma_logging.h"

namespace magma {

YAML::Node ServiceConfigLoader::load_service_config(
    const std::string& service_name){
  auto file_path = std::string(CONFIG_DIR) + service_name + ".yml";
  YAML::Node base_config = YAML::LoadFile(file_path);

  // Try to override original file, if an override exists
  try {
    auto override_file = std::string(OVERRIDE_DIR) + service_name + ".yml";
    return YAMLUtils::merge_nodes(base_config, YAML::LoadFile(override_file));
  } catch (YAML::BadFile e) {
    MLOG(MDEBUG) << "Override file not found for service " << service_name;
  }
  return base_config;
}

}
