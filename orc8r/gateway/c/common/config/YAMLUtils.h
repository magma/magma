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
 * YAMLUtils defines new utilities that the yaml-cpp library doesn't expose
 */
class YAMLUtils final {
public:

  /**
   * merge_nodes combines two YAML files together. override_node will
   * override any parameters it defines, and keep any existing parameters in
   * default_node that it doesn't define
   */
  static YAML::Node merge_nodes(
    const YAML::Node& default_node,
    const YAML::Node& override_node);
};

}
