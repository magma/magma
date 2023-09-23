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

#include "lte/gateway/c/core/oai/include/mme_config.hpp"
#include "lte/gateway/c/core/oai/include/s1ap_types.hpp"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/include/state_utility.hpp"

namespace {
constexpr char S1AP_STATE_TABLE[] = "s1ap_state";
constexpr char S1AP_TASK_NAME[] = "S1AP";
}  // namespace

namespace magma {
namespace lte {

/**
 * create_s1ap_state allocates a new S1apState struct and initializes
 * its properties.
 */
oai::S1apState* create_s1ap_state(void);
/**
 * free_s1ap_state deallocates a S1apState struct and its properties.
 */
void free_s1ap_state(oai::S1apState* state_cache_p);

/**
 * S1apStateManager is a thread safe singleton class that contains functions
 * to maintain S1AP task state, allocating and freeing related state structs.
 */
class S1apStateManager : public StateUtility {
 public:
  /**
   * Returns an instance of S1apStateManager, guaranteed to be thread safe and
   * initialized only once.
   * @return S1apStateManager instance
   */
  static S1apStateManager& getInstance();

  /**
   * Function to initialize member variables
   * @param persist_state should persist state in redis
   */
  void init(bool persist_state);

  // Copy constructor and assignment operator are marked as deleted functions
  S1apStateManager(S1apStateManager const&) = delete;
  S1apStateManager& operator=(S1apStateManager const&) = delete;

  /**
   * Frees all memory allocated on s1ap_state cache struct
   */
  void free_state();

  /**
   * Reads S1AP context state for all UEs in db
   * @return operation response code
   */
  status_code_e read_ue_state_from_db();

  /**
   * Serializes s1ap_imsi_map to proto and saves it into data store
   */
  void write_s1ap_imsi_map_to_db();

  /**
   * Returns a pointer to s1ap_imsi_map
   */
  oai::S1apImsiMap* get_s1ap_imsi_map();
  map_uint64_ue_description_t* get_s1ap_ue_state();
  oai::S1apState* get_state(bool read_from_db);
  void write_s1ap_state_to_db();
  void s1ap_write_ue_state_to_db(const oai::UeDescription* ue_context,
                                 const std::string& imsi_str);
  status_code_e read_state_from_db();

 private:
  S1apStateManager();
  ~S1apStateManager();

  /**
   * Allocates new S1apState struct and its properties
   */
  void create_state();

  void create_s1ap_imsi_map();
  void clear_s1ap_imsi_map();

  std::size_t s1ap_imsi_map_hash_;
  oai::S1apImsiMap* s1ap_imsi_map_;
  map_uint64_ue_description_t state_ue_map;
  oai::S1apState* state_cache_p;
  // State version counters for task and ue context
  uint64_t task_state_version;
  std::unordered_map<std::string, uint64_t> ue_state_version;
  // Last written hash values for task and ue context
  std::size_t task_state_hash;
  std::unordered_map<std::string, std::size_t> ue_state_hash;
};

}  // namespace lte
}  // namespace magma
