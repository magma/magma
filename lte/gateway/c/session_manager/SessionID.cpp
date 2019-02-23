/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include <sstream>
#include <string>
#include <time.h>

#include "SessionID.h"

SessionIDGenerator::SessionIDGenerator() {
  // init random seed
  srand(time(NULL));
}

std::string SessionIDGenerator::gen_session_id(const std::string& imsi) {
  // imsi- + random 6 digit number
  return imsi + "-" + std::to_string(rand() % 1000000);
}

bool SessionIDGenerator::get_imsi_from_session_id(
    const std::string& session_id,
    std::string& imsi_out) {
  std::istringstream ss(session_id);
  if (std::getline(ss, imsi_out, '-')) {
    return true;
  }
  return false;
}
