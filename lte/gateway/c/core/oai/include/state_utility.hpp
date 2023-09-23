/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#pragma once

#include <cstdlib>

#include <unordered_map>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/redis_utils/redis_client.hpp"

constexpr char IMSI_STR_PREFIX[] = "IMSI";

namespace magma {
namespace lte {

class StateUtility {
 public:
  std::string get_imsi_str(imsi64_t imsi64) {
    AssertFatal(
        is_initialized,
        "StateUtility init() function should be called to initialize state");

    char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
    IMSI64_TO_STRING(imsi64, (char*)imsi_str, IMSI_BCD_DIGITS_MAX);
    return imsi_str;
  }

  void clear_ue_state_db(const std::string& imsi_str) {
    AssertFatal(
        is_initialized,
        "StateUtility init() function should be called to initialize state");

    if (persist_state_enabled) {
      std::vector<std::string> keys = {IMSI_STR_PREFIX + imsi_str + ":" +
                                       task_name};
      if (redis_client->clear_keys(keys) != RETURNok) {
        OAILOG_ERROR(log_task, "Failed to remove UE state from db");
        return;
      }
      OAILOG_DEBUG(log_task, "Removing UE state for IMSI %s", imsi_str.c_str());
    }
  }

  bool is_persist_state_enabled() const { return persist_state_enabled; }

 protected:
  StateUtility()
      : redis_client(nullptr),
        is_initialized(false),
        state_dirty(false),
        persist_state_enabled(false),
        log_task(LOG_UTIL) {}
  virtual ~StateUtility() = default;

  imsi64_t get_imsi_from_key(const std::string& key) const {
    imsi64_t imsi64;
    std::string imsi_str_prefix = key.substr(0, key.find(':'));
    std::string imsi_str = imsi_str_prefix.substr(4, imsi_str_prefix.length());
    IMSI_STRING_TO_IMSI64(imsi_str.c_str(), &imsi64);
    return imsi64;
  }

  std::unique_ptr<RedisClient> redis_client;
  // Flag for check asserting if the state has been initialized.
  bool is_initialized;
  // Flag for check asserting that write should be done after read.
  bool state_dirty;
  // Flag for enabling writing and reading to db.
  bool persist_state_enabled;

  std::string table_key;
  std::string task_name;
  log_proto_t log_task;
};

}  // namespace lte
}  // namespace magma
