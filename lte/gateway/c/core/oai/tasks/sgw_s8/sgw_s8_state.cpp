/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#include "lte/gateway/c/core/oai/include/sgw_s8_state.hpp"

#include <cstdlib>

#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/include/sgw_context_manager.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_procedures.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw_s8/sgw_s8_state_manager.hpp"
#include "lte/gateway/c/core/common/dynamic_memory_check.h"

extern "C" {
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
}

using magma::lte::SgwStateManager;

int sgw_state_init(bool persist_state, const sgw_config_t* config) {
  SgwStateManager::getInstance().init(persist_state, config);
  return RETURNok;
}

sgw_state_t* get_sgw_state(bool read_from_db) {
  return SgwStateManager::getInstance().get_state(read_from_db);
}

hash_table_ts_t* get_sgw_ue_state() {
  return SgwStateManager::getInstance().get_ue_state_ht();
}

int read_sgw_ue_state_db() {
  return SgwStateManager::getInstance().read_ue_state_from_db();
}

void sgw_state_exit() {
  SgwStateManager::getInstance().free_state();
  return;
}

void put_sgw_state() {
  SgwStateManager::getInstance().write_state_to_db();
  return;
}

void put_sgw_ue_state(sgw_state_t* sgw_state, imsi64_t imsi64) { return; }

void delete_sgw_ue_state(imsi64_t imsi64) { return; }

void sgw_free_s11_bearer_context_information(void** ptr) {
  if (!ptr) {
    return;
  }
  sgw_eps_bearer_context_information_t* sgw_eps_context =
      reinterpret_cast<sgw_eps_bearer_context_information_t*>(*ptr);
  if (sgw_eps_context) {
    sgw_free_pdn_connection(&sgw_eps_context->pdn_connection);
    delete_pending_procedures(sgw_eps_context);
    free_cpp_wrapper(reinterpret_cast<void**>(ptr));
  }
  return;
}
