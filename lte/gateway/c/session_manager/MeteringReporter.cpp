/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "MagmaService.h"
#include "MeteringReporter.h"

namespace magma {
namespace lte {

const char* COUNTER_NAME       = "ue_traffic";
const char* LABEL_IMSI         = "IMSI";
const char* LABEL_SESSION_ID   = "session_id";
const char* LABEL_DIRECTION    = "direction";
const char* DIRECTION_UP       = "up";
const char* DIRECTION_DOWN     = "down";

MeteringReporter::MeteringReporter() {}

bool MeteringReporter::report_usage(
    const std::string& imsi, const std::string& session_id,
    SessionStateUpdateCriteria& update_criteria) {
  double total_tx = 0;
  double total_rx = 0;

  // Charging credit
  for (const auto& it : update_criteria.charging_credit_map) {
    auto credit_update = it.second;
    total_tx += (double) credit_update.bucket_deltas[USED_TX];
    total_rx += (double) credit_update.bucket_deltas[USED_RX];
  }

  // Monitoring credit
  for (const auto& it : update_criteria.monitor_credit_map) {
    auto credit_update  = it.second;
    total_tx += (double) credit_update.bucket_deltas[USED_TX];
    total_rx += (double) credit_update.bucket_deltas[USED_RX];
  }

  report_upload(imsi, session_id, total_tx);
  report_download(imsi, session_id, total_rx);
}

void MeteringReporter::report_upload(
    const std::string& imsi, const std::string& session_id,
    double unreported_usage_bytes) {
  report_traffic(
      imsi, session_id, DIRECTION_UP, unreported_usage_bytes);
}

void MeteringReporter::report_download(
    const std::string& imsi, const std::string& session_id,
    double unreported_usage_bytes) {
  report_traffic(
      imsi, session_id, DIRECTION_DOWN, unreported_usage_bytes);
}

void MeteringReporter::report_traffic(
    const std::string& imsi, const std::string& session_id,
    const std::string& traffic_direction, double unreported_usage_bytes) {
  increment_counter(
      COUNTER_NAME, unreported_usage_bytes, size_t(3), LABEL_IMSI, imsi.c_str(),
      LABEL_SESSION_ID, session_id.c_str(),
      LABEL_DIRECTION, traffic_direction.c_str());
}

void MeteringReporter::increment_counter(
    const char* name, double increment, size_t n_labels, ...) {
  va_list ap;
  va_start(ap, n_labels);
  MetricsSingleton::Instance().IncrementCounter(name, increment, n_labels, ap);
  va_end(ap);
}

}  // namespace lte
}  // namespace magma
