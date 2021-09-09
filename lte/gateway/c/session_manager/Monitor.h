/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#pragma once

#include "SessionCredit.h"
#include "StoredState.h"

namespace magma {
// Monitor is a structure to keep track of grants of bytes given used for
// monitoring. (With Federation, this corresponds to grants given by PCRF.)
// Since this grant is for monitoring only, the state of each
// monitor does NOT affect whether a session should be continued or not.
// At every point where the monitoring grant is exhausted, we will report the
// recorded usage.
struct Monitor {
  // Keep track of used/reported/allowed bytes
  SessionCredit credit;
  // Indicates whether the credit above is applied session-wide or per
  // monitoring key
  MonitoringLevel level;

  Monitor() {}

  explicit Monitor(const StoredMonitor& marshaled) {
    credit = SessionCredit(marshaled.credit);
    level = marshaled.level;
  }

  // Marshal into StoredMonitor structure used in SessionStore
  StoredMonitor marshal() {
    StoredMonitor marshaled{};
    marshaled.credit = credit.marshal();
    marshaled.level = level;
    return marshaled;
  }

  // Monitor will be deleted when the credit is exhausted and no more grants are
  // received or when we receive the explicit action
  bool should_delete_monitor() {
    return credit.current_grant_contains_zero() && credit.is_quota_exhausted(1);
  }

  bool should_send_update() {
    // update trigger by being the final report
    if (credit.is_report_last_credit()) {
      return true;
    }
    // updated trigger due to usage threshold (80%)
    if (!credit.current_grant_contains_zero() &&
        credit.is_quota_exhausted(SessionCredit::USAGE_REPORTING_THRESHOLD)) {
      // if grant contains zeros that means we are in final. Do no report
      return true;
    }
    return false;
  }
};

}  // namespace magma
