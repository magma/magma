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
/****************************************************************************
  Source      ngap_state_manager.h
  Date        2020/07/28
  Author      Ashish Prajapati
  Subsystem   Access and Mobility Management Function
  Description Defines NG Application Protocol Messages

*****************************************************************************/

#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include "amf_config.h"
#include "ngap_types.h"

#ifdef __cplusplus
}
#endif

#include "common_defs.h"
#include "state_manager.h"
#include "ngap_state_converter.h"
using namespace magma::lte;
using namespace magma::lte::oai;

namespace magma5g {
constexpr char NGAP_STATE_TABLE[] = "ngap_state";
constexpr char NGAP_TASK_NAME[]   = "NGAP";
}  // namespace magma5g

namespace magma5g {

/**
 * NGapStateManager is a thread safe singleton class that contains functions
 * to maintain NGAP task state, allocating and freeing related state structs.
 */
class NgapStateManager
    : public StateManager<
          ngap_state_t, m5g_ue_description_t, oai::NgapState,
          magma::lte::oai::Ngap_UeDescription, NgapStateConverter> {
 public:
  /**
   * Returns an instance of NGapStateManager, guaranteed to be thread safe and
   * initialized only once.
   * @return NGapStateManager instance
   */
  static NgapStateManager& getInstance();

  /**
   * Function to initialize member variables
   * @param amf_config amf_config_t struct
   */
  void init(uint32_t max_ues, uint32_t max_enbs, bool use_stateless);

  // Copy constructor and assignment operator are marked as deleted functions
  NgapStateManager(NgapStateManager const&) = delete;
  NgapStateManager& operator=(NgapStateManager const&) = delete;

  /**
   * Frees all memory allocated on ngap_state cache struct
   */
  void free_state() override;

  /**
   * Reads NGAP context state for all UEs in db
   * @return operation response code
   */
  status_code_e read_ue_state_from_db() override;

  /**
   * Serializes ngap_imsi_map to proto and saves it into data store
   */
  void write_ngap_imsi_map_to_db();

  /**
   * Returns a pointer to ngap_imsi_map
   */
  ngap_imsi_map_t* get_ngap_imsi_map();

 private:
  NgapStateManager();
  ~NgapStateManager();

  /**
   * Allocates new ngap_state_t struct and its properties
   */
  void create_state() override;

  void create_ngap_imsi_map();
  void clear_ngap_imsi_map();

 public:
  void put_ngap_imsi_map();

 private:
  uint32_t max_ues_;
  uint32_t max_gnbs_;
  ngap_imsi_map_t* ngap_imsi_map_;
  std::size_t ngap_imsi_map_hash_;
};
}  // namespace magma5g
