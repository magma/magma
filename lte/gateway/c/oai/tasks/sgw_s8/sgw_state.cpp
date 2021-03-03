/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include "sgw_s8_state.h"

#include <cstdlib>
#include <conversions.h>

extern "C" {
#include "assertions.h"
#include "bstrlib.h"
#include "dynamic_memory_check.h"
#include "sgw_context_manager.h"
}

#include "sgw_state_manager.h"

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
}

void put_sgw_state() {
  SgwStateManager::getInstance().write_state_to_db();
}

void put_sgw_ue_state(sgw_state_t* sgw_state, imsi64_t imsi64) {}

void delete_sgw_ue_state(imsi64_t imsi64) {}
