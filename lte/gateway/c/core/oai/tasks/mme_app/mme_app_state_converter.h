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

#pragma once
extern "C" {
#include "mme_app_desc.h"
#include "mme_app_ue_context.h"
}

#include <sstream>
#include "lte/protos/oai/mme_nas_state.pb.h"
#include "state_converter.h"

/******************************************************************************
 * This is a helper class to encapsulate all functions for converting in-memory
 * state of MME and NAS task to/from proto for persisting UE state in data
 * store. The class does not have any member variables. The class does not
 * allocate any memory, but calls NAS state converter, which dynamically
 * allocates memory for EMM procedures. All the allocated memory is cleared by
 * the caller class MmeNasStateManager
 ******************************************************************************/
namespace magma {
namespace lte {

class MmeNasStateConverter : public StateConverter {
 public:
  // Constructor
  MmeNasStateConverter();

  // Destructor
  ~MmeNasStateConverter();

  // Serialize mme_app_desc_t to oai::MmeNasState proto
  static void state_to_proto(
      const mme_app_desc_t* mme_nas_state_p, oai::MmeNasState* state_proto);

  // Deserialize mme_app_desc_t from oai::MmeNasState proto
  static void proto_to_state(
      const oai::MmeNasState& state_proto, mme_app_desc_t* mme_nas_state_p);

  static void ue_to_proto(
      const ue_mm_context_t* ue_ctxt, oai::UeContext* ue_ctxt_proto);

  static void proto_to_ue(
      const oai::UeContext& ue_ctxt_proto, ue_mm_context_t* ue_ctxt);

  static void mme_app_convert_string_to_guti(
      guti_t* guti_p, const std::string& guti_str);

  static char* mme_app_convert_guti_to_string(guti_t* guti_p);

 private:
  /***********************************************************
   *                 Hashtable <-> Proto
   * Functions to serialize/deserialize in-memory hashtables
   * for MME task. Only MME task inserts/removes elements in
   * the hashtables, so these calls are thread-safe.
   * We only need to lock the UE context structure as it can
   * also be accessed by the NAS task. If hashtable is empty
   * the proto field is also empty
   ***********************************************************/

  static void hashtable_ts_to_proto(
      hash_table_ts_t* state_htbl,
      google::protobuf::Map<unsigned long, oai::UeContext>* proto_map);

  static void proto_to_hashtable_ts(
      const google::protobuf::Map<unsigned long, oai::UeContext>& proto_map,
      hash_table_ts_t* state_htbl);

  static void guti_table_to_proto(
      const obj_hash_table_uint64_t* guti_htbl,
      google::protobuf::Map<std::string, unsigned long>* proto_map);

  static void proto_to_guti_table(
      const google::protobuf::Map<std::string, unsigned long>& proto_map,
      obj_hash_table_uint64_t* guti_htbl);

  /**********************************************************
   *                 UE Context <-> Proto                    *
   *  Functions to serialize/desearialize UE context         *
   *  The caller needs to acquire a lock on UE context       *
   **********************************************************/

  static void mme_app_timer_to_proto(
      const mme_app_timer_t& state_mme_timer, oai::Timer* timer_proto);

  static void proto_to_mme_app_timer(
      const oai::Timer& timer_proto, mme_app_timer_t* state_mme_app_timer);

  static void sgs_context_to_proto(
      sgs_context_t* state_sgs_context, oai::SgsContext* sgs_context_proto);

  static void proto_to_sgs_context(
      const oai::SgsContext& sgs_context_proto,
      sgs_context_t* state_sgs_context);

  static void fteid_to_proto(
      const fteid_t& state_fteid, oai::Fteid* fteid_proto);

  static void proto_to_fteid(
      const oai::Fteid& fteid_proto, fteid_t* state_fteid);

  static void bearer_context_to_proto(
      const bearer_context_t& state_bearer_context,
      oai::BearerContext* bearer_context_proto);

  static void proto_to_bearer_context(
      const oai::BearerContext& bearer_context_proto,
      bearer_context_t* state_bearer_context);

  static void bearer_context_list_to_proto(
      const ue_mm_context_t& state_ue_context,
      oai::UeContext* ue_context_proto);

  static void proto_to_bearer_context_list(
      const oai::UeContext& ue_context_proto,
      ue_mm_context_t* state_ue_context);

  static void esm_pdn_to_proto(
      const esm_pdn_t& state_esm_pdn, oai::EsmPdn* esm_pdn_proto);

  static void proto_to_esm_pdn(
      const oai::EsmPdn& esm_pdn_proto, esm_pdn_t* state_esm_pdn);

  static void pdn_context_to_proto(
      const pdn_context_t& state_pdn_context,
      oai::PdnContext* pdn_context_proto);

  static void proto_to_pdn_context(
      const oai::PdnContext& pdn_context_proto,
      pdn_context_t* state_pdn_context);

  static void pdn_context_list_to_proto(
      const ue_mm_context_t& state_ue_context,
      oai::UeContext* ue_context_proto);

  static void proto_to_pdn_context_list(
      const oai::UeContext& ue_context_proto,
      ue_mm_context_t* state_ue_context);

  static void ue_context_to_proto(
      const ue_mm_context_t* ue_ctxt, oai::UeContext* ue_ctxt_proto);

  static void proto_to_ue_mm_context(
      const oai::UeContext& ue_context_proto,
      ue_mm_context_t* state_ue_mm_context);

  static void regional_subscription_to_proto(
      const ue_mm_context_t& state_ue_context,
      oai::UeContext* ue_context_proto);

  static void proto_to_regional_subscription(
      const oai::UeContext& ue_context_proto,
      ue_mm_context_t* state_ue_context);
};
}  // namespace lte
}  // namespace magma
