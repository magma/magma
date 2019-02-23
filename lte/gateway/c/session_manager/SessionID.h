/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#pragma once

class SessionIDGenerator {
public:
  SessionIDGenerator();

  /**
   * Generates a random session id from the IMSI in the form
   * "<IMSI>-<RANDOM NUM>"
   */
  std::string gen_session_id(const std::string& imsi);

  /**
   * Parses an IMSI value from a session_id
   */
  bool get_imsi_from_session_id(
      const std::string& session_id,
      std::string& imsi_out);
};
