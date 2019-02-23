/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <cpp_redis/cpp_redis>

namespace magma {

enum ObjectMapResult {
  SUCCESS = 0,
  CLIENT_ERROR = 1,
  KEY_NOT_FOUND = 2,
  INCORRECT_VALUE_TYPE = 3,
  DESERIALIZE_FAIL = 4,
  SERIALIZE_FAIL = 5,
};

/**
 * ObjectMap is an abstract class to represent any class that can store key,
 * object pairs
 */
template <typename ObjectType>
class ObjectMap {
  virtual ObjectMapResult set(
    const std::string& key,
    const ObjectType& object) = 0;

  virtual ObjectMapResult get(
    const std::string& key,
    ObjectType& object_out) = 0;

  virtual ObjectMapResult getall(std::vector<ObjectType>& values_out) = 0;
};

}
