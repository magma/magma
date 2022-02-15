/**
 * Copyright 2021 The Magma Authors.
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
#include <gtest/gtest.h>

extern "C" {
#include "lte/gateway/c/core/oai/common/log.h"
#include "S1ap_S1AP-PDU.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_handlers.h"
}

#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_state_converter.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_state_manager.h"

using ::testing::Test;

namespace magma {
namespace lte {

class S1APStateConverterTest : public ::testing::Test {
  virtual void SetUp() {}

  virtual void TearDown() {}
};

TEST_F(S1APStateConverterTest, S1apStateConversionSuccess) {
  sctp_assoc_id_t assoc_id = 1;
  s1ap_state_t* init_state = create_s1ap_state(2, 2);
  s1ap_state_t* final_state = create_s1ap_state(2, 2);

  enb_description_t* enb_association = s1ap_new_enb();
  enb_association->sctp_assoc_id = assoc_id;
  enb_association->enb_id = 0xFFFFFFFF;
  enb_association->s1_state = S1AP_READY;

  // filling ue_id_coll
  hashtable_uint64_ts_insert(&enb_association->ue_id_coll, (const hash_key_t)1,
                             17);
  hashtable_uint64_ts_insert(&enb_association->ue_id_coll, (const hash_key_t)2,
                             25);

  // filling supported_tai_items
  supported_ta_list_t* enb_ta_list = &enb_association->supported_ta_list;
  enb_ta_list->list_count = 1;
  enb_ta_list->supported_tai_items[0].tac = 1;
  enb_ta_list->supported_tai_items[0].bplmnlist_count = 1;
  enb_ta_list->supported_tai_items[0].bplmns[0].mcc_digit1 = 1;
  enb_ta_list->supported_tai_items[0].bplmns[0].mcc_digit2 = 1;
  enb_ta_list->supported_tai_items[0].bplmns[0].mcc_digit3 = 1;
  enb_ta_list->supported_tai_items[0].bplmns[0].mnc_digit1 = 1;
  enb_ta_list->supported_tai_items[0].bplmns[0].mnc_digit2 = 1;
  enb_ta_list->supported_tai_items[0].bplmns[0].mnc_digit3 = 1;

  // Inserting 1 enb association
  hashtable_ts_insert(&init_state->enbs,
                      (const hash_key_t)enb_association->sctp_assoc_id,
                      (void*)enb_association);
  init_state->num_enbs = 1;

  hashtable_ts_insert(&init_state->mmeid2associd, (const hash_key_t)1,
                      (void**)&assoc_id);

  oai::S1apState state_proto;
  S1apStateConverter::state_to_proto(init_state, &state_proto);
  S1apStateConverter::proto_to_state(state_proto, final_state);

  EXPECT_EQ(init_state->num_enbs, final_state->num_enbs);
  enb_description_t* enbd = nullptr;
  enb_description_t* enbd_final = nullptr;
  EXPECT_EQ(hashtable_ts_get(&init_state->enbs, (const hash_key_t)assoc_id,
                             reinterpret_cast<void**>(&enbd)),
            HASH_TABLE_OK);
  EXPECT_EQ(hashtable_ts_get(&final_state->enbs, (const hash_key_t)assoc_id,
                             reinterpret_cast<void**>(&enbd_final)),
            HASH_TABLE_OK);

  EXPECT_EQ(enbd->sctp_assoc_id, enbd_final->sctp_assoc_id);
  EXPECT_EQ(enbd->enb_id, enbd_final->enb_id);
  EXPECT_EQ(enbd->s1_state, enbd_final->s1_state);

  EXPECT_EQ(enbd->supported_ta_list.list_count,
            enbd_final->supported_ta_list.list_count);
  EXPECT_EQ(enbd->supported_ta_list.supported_tai_items[0].tac,
            enbd_final->supported_ta_list.supported_tai_items[0].tac);

  free_s1ap_state(init_state);
  free_s1ap_state(final_state);
}

TEST_F(S1APStateConverterTest, S1apStateConversionExpectedEnbCount) {
  sctp_assoc_id_t assoc_id = 1;
  s1ap_state_t* init_state = create_s1ap_state(2, 2);
  s1ap_state_t* final_state = create_s1ap_state(2, 2);

  enb_description_t* enb_association = s1ap_new_enb();
  enb_association->sctp_assoc_id = assoc_id;
  enb_association->enb_id = 0xFFFFFFFF;
  enb_association->s1_state = S1AP_READY;
  // Inserting 1 enb association
  hashtable_ts_insert(&init_state->enbs,
                      (const hash_key_t)enb_association->sctp_assoc_id,
                      (void*)enb_association);
  // state_to_proto should update num_enbs to match expected eNB count on the
  // hashtable
  init_state->num_enbs = 5;

  oai::S1apState state_proto;
  S1apStateConverter::state_to_proto(init_state, &state_proto);
  EXPECT_EQ(init_state->num_enbs, 1);

  S1apStateConverter::proto_to_state(state_proto, final_state);
  EXPECT_EQ(final_state->num_enbs, 1);

  free_s1ap_state(init_state);
  free_s1ap_state(final_state);
}

TEST_F(S1APStateConverterTest, S1apStateConversionUeContext) {
  ue_description_t* ue = (ue_description_t*)calloc(1, sizeof(ue_description_t));
  ue_description_t* final_ue =
      (ue_description_t*)calloc(1, sizeof(ue_description_t));

  // filling with test values
  ue->mme_ue_s1ap_id = 1;
  ue->enb_ue_s1ap_id = 1;
  ue->sctp_assoc_id = 1;
  ue->comp_s1ap_id = S1AP_GENERATE_COMP_S1AP_ID(1, 1);
  ue->s1ap_handover_state.mme_ue_s1ap_id = 1;
  ue->s1ap_handover_state.source_enb_id = 1;
  ue->s1ap_ue_context_rel_timer.id = 1;
  ue->s1ap_ue_context_rel_timer.msec = 1000;

  oai::UeDescription ue_proto;
  S1apStateConverter::ue_to_proto(ue, &ue_proto);
  S1apStateConverter::proto_to_ue(ue_proto, final_ue);

  EXPECT_EQ(ue->comp_s1ap_id, final_ue->comp_s1ap_id);
  EXPECT_EQ(ue->mme_ue_s1ap_id, final_ue->mme_ue_s1ap_id);
  EXPECT_EQ(ue->s1ap_ue_context_rel_timer.id,
            final_ue->s1ap_ue_context_rel_timer.id);

  free_wrapper((void**)&ue);
  free_wrapper((void**)&final_ue);
}

}  // namespace lte
}  // namespace magma
