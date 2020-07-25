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

#include "StoredState.h"
#include "SessionCredit.h"

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

  // Marshal into StoredMonitor structure used in SessionStore
  StoredMonitor marshal() {
    StoredMonitor marshaled{};
    marshaled.credit = credit.marshal();
    marshaled.level = level;
    return marshaled;
  }

  // Unmarshal from StoredMonitor structure used in SessionStore
  static std::unique_ptr<Monitor> unmarshal(const StoredMonitor &marshaled) {
    Monitor monitor;
    monitor.credit = SessionCredit::unmarshal(marshaled.credit);
    monitor.level = marshaled.level;
    return std::make_unique<Monitor>(monitor);
  }
};

}  // namespace magma
