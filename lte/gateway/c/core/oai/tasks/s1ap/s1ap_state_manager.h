/*
 *
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

#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include "mme_config.h"
#include "s1ap_types.h"

#ifdef __cplusplus
}
#endif

#include "common_defs.h"
#include "state_manager.h"
#include "s1ap_state_converter.h"

namespace {
constexpr char S1AP_STATE_TABLE[] = "s1ap_state";
constexpr char S1AP_TASK_NAME[]   = "S1AP";
}  // namespace

namespace magma {
namespace lte {

/**
 * S1apStateManager is a thread safe singleton class that contains functions
 * to maintain S1AP task state, allocating and freeing related state structs.
 */
class S1apStateManager
    : public StateManager<
          s1ap_state_t, ue_description_t, magma::lte::oai::S1apState,
          magma::lte::oai::UeDescription, S1apStateConverter> {
 public:
  /**
   * Returns an instance of S1apStateManager, guaranteed to be thread safe and
   * initialized only once.
   * @return S1apStateManager instance
   */
  static S1apStateManager& getInstance();

  /**
   * Function to initialize member variables
   * @param max_ues number of max UEs in hashtable
   * @param max_enbs number of max eNBs in hashtable
   * @param persist_state should persist state in redis
   */
  void init(uint32_t max_ues, uint32_t max_enbs, bool persist_state);

  // Copy constructor and assignment operator are marked as deleted functions
  S1apStateManager(S1apStateManager const&) = delete;
  S1apStateManager& operator=(S1apStateManager const&) = delete;

  /**
   * Frees all memory allocated on s1ap_state cache struct
   */
  void free_state() override;

  /**
   * Reads S1AP context state for all UEs in db
   * @return operation response code
   */
  status_code_e read_ue_state_from_db() override;

  /**
   * Serializes s1ap_imsi_map to proto and saves it into data store
   */
  void write_s1ap_imsi_map_to_db();

  /**
   * Returns a pointer to s1ap_imsi_map
   */
  s1ap_imsi_map_t* get_s1ap_imsi_map();

 private:
  S1apStateManager();
  ~S1apStateManager() override;

  /**
   * Allocates new s1ap_state_t struct and its properties
   */
  void create_state() override;

  void create_s1ap_imsi_map();
  void clear_s1ap_imsi_map();

  uint32_t max_ues_;
  uint32_t max_enbs_;
  std::size_t s1ap_imsi_map_hash_;
  s1ap_imsi_map_t* s1ap_imsi_map_;
};
}  // namespace lte
}  // namespace magma
