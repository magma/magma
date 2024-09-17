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

map_uint32_sgw_eps_bearer_context_t* get_s8_state_teid_map() {
  return SgwStateManager::getInstance().get_s8_state_teid_map();
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

void sgw_s8_free_eps_bearer_context(
    sgw_eps_bearer_ctxt_t** sgw_eps_bearer_ctxt) {
  if (*sgw_eps_bearer_ctxt) {
    if ((*sgw_eps_bearer_ctxt)->pgw_cp_ip_port) {
      free_wrapper(
          reinterpret_cast<void**>(&(*sgw_eps_bearer_ctxt)->pgw_cp_ip_port));
    }
    free_cpp_wrapper(reinterpret_cast<void**>(sgw_eps_bearer_ctxt));
  }
}

void sgw_s8_free_pdn_connection(sgw_pdn_connection_t* pdn_connection_p) {
  if (pdn_connection_p) {
    if (pdn_connection_p->apn_in_use) {
      free_wrapper(reinterpret_cast<void**>(&pdn_connection_p->apn_in_use));
    }

    for (auto& ebix : pdn_connection_p->sgw_eps_bearers_array) {
      sgw_s8_free_eps_bearer_context(&ebix);
    }
  }
}
void sgw_free_s11_bearer_context_information(void** ptr) {
  if (!ptr) {
    return;
  }
  sgw_eps_bearer_context_information_t* sgw_eps_context =
      reinterpret_cast<sgw_eps_bearer_context_information_t*>(*ptr);
  if (sgw_eps_context) {
    sgw_s8_free_pdn_connection(&sgw_eps_context->pdn_connection);
    delete_pending_procedures(sgw_eps_context);
    free_cpp_wrapper(reinterpret_cast<void**>(ptr));
  }
  return;
}
