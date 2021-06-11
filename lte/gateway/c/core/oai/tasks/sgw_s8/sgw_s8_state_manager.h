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

#pragma once

#include <state_manager.h>
#include "sgw_s8_state.h"
#include "lte/protos/oai/sgw_state.pb.h"
#include "sgw_s8_state_converter.h"
#include "common_defs.h"

namespace {
constexpr int SGW_STATE_CONTEXT_HT_MAX_SIZE    = 512;
constexpr int MAX_PREDEFINED_PCC_RULES_HT_SIZE = 32;
constexpr char S11_BEARER_CONTEXT_INFO_HT_NAME[] =
    "s11_bearer_context_information_htbl";
constexpr char SGW_STATE_TABLE_NAME[] = "sgw_state";
constexpr char SGW_TASK_NAME[]        = "SGW";
}  // namespace

namespace magma {
namespace lte {

/**
 * SgwStateManager is a singleton (thread-safe, destruction guaranteed)
 * that contains functions to maintain SGW state, allocating and
 * freeing state structs, and writing / reading state to db.
 */
class SgwStateManager : public StateManager<
                            sgw_state_t, spgw_ue_context_t, oai::SgwState,
                            oai::SgwUeContext, SgwStateConverter> {
 public:
  /**
   * Returns an instance of SgwStateManager, guaranteed to be thread safe and
   * initialized only once.
   * @return SgwStateManager instance.
   */
  static SgwStateManager& getInstance();
  /**
   * Initialization function to initialize member variables.
   * @param persist_state should read and write state from db
   * @param config SGW config struct
   */
  void init(bool persist_state, const sgw_config_t* config);

  /**
   * Singleton class, copy constructor and assignment operator are marked
   * as deleted functions.
   */
  // Making them public for better debugging logging.
  SgwStateManager(SgwStateManager const&) = delete;
  SgwStateManager& operator=(SgwStateManager const&) = delete;

  /**
   * Frees all memory allocated on sgw_state_t.
   */
  void free_state() override;

  status_code_e read_ue_state_from_db() override;

 private:
  SgwStateManager();
  ~SgwStateManager();

  /**
   * Allocates a new sgw_state_t struct, and inits hashtables and state
   * structs to default values.
   */
  void create_state() override;

  const sgw_config_t* config_;
};

}  // namespace lte
}  // namespace magma
