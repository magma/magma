/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include "StoredState.h"

namespace magma {
namespace lte {

class MeteringReporter {
 public:
  MeteringReporter();

  /**
   * Report all unreported traffic usage for a session.
   * All charging and monitoring keys are aggregated.
   */
  bool report_usage(
      const std::string& imsi, const std::string& session_id,
      SessionStateUpdateCriteria& update_criteria);

 private:
  /**
   * Report upload traffic usage for a session
   */
  void report_upload(
      const std::string& imsi, const std::string& session_id,
      double unreported_usage_bytes);

  /**
   * Report download traffic usage for a session
   */
  void report_download(
      const std::string& imsi, const std::string& session_id,
      double unreported_usage_bytes);

  /**
   * Report traffic usage for a session
   */
  void report_traffic(
      const std::string& imsi, const std::string& session_id,
      const std::string& traffic_direction, double unreported_usage_bytes);

  void increment_counter(
      const char* name, double increment, size_t n_labels, ...);
};

}  // namespace lte
}  // namespace magma
