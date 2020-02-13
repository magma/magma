// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <string>

namespace devmand {
namespace devices {
namespace cli {

using namespace std;

class Model {
 public:
  Model() = delete;
  ~Model() = default;

 protected:
  explicit Model(string _dir) : dir(_dir) {}

 private:
  const string dir;

 public:
  const string& getDir() const {
    return dir;
  }
  bool operator<(const Model& x) const {
    return dir < x.dir;
  }

  static const Model OPENCONFIG_0_1_6;
  static const Model IETF_0_1_5;
};

} // namespace cli
} // namespace devices
} // namespace devmand
