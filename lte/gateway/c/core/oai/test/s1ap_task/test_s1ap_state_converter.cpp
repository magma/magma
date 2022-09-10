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
#include "lte/gateway/c/core/oai/include/amf_config.hpp"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
}

#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_handlers.hpp"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_state_converter.hpp"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_state_manager.hpp"
#include "lte/gateway/c/core/oai/test/mock_tasks/mock_tasks.hpp"

using ::testing::Test;

namespace magma {
namespace lte {

class S1APStateConverterTest : public ::testing::Test {
  void SetUp() {
    itti_init(TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info,
              NULL, NULL);
    mme_config_init(&mme_config);
    s1ap_state_init(amf_config.max_ues, amf_config.max_gnbs,
                    amf_config.use_stateless);
  }

  void TearDown() {
    s1ap_state_exit();
    free_mme_config(&mme_config);
    itti_free_desc_threads();
  }
};

TEST_F(S1APStateConverterTest, S1apStateConversionSuccess) {
  sctp_assoc_id_t assoc_id = 1;
  S1apState* init_state = create_s1ap_state();
  S1apState* final_state = create_s1ap_state();

  EnbDescription* enb_association = s1ap_new_enb();
  enb_association->set_sctp_assoc_id(assoc_id);
  enb_association->set_enb_id(0xFFFFFFFF);
  enb_association->set_s1_enb_state(magma::lte::oai::S1AP_READY);

  // filling ue_id_coll
  magma::proto_map_uint32_uint64_t ue_id_coll;
  ue_id_coll.map = enb_association->mutable_ue_id_map();
  ue_id_coll.insert(1, 17);
  ue_id_coll.insert(2, 25);

  // filling supported_tai_items

  magma::lte::oai::SupportedTaList* enb_ta_list =
      enb_association->mutable_supported_ta_list();
  enb_ta_list->set_list_count(1);

  for (int tai_idx = 0; tai_idx < enb_ta_list->list_count(); tai_idx++) {
    magma::lte::oai::SupportedTaiItems* supported_tai_items =
        enb_ta_list->add_supported_tai_items();
    supported_tai_items->set_tac(1);

    supported_tai_items->set_bplmnlist_count(1);
#define PLMN_BYTES 7
    char plmn_array[PLMN_BYTES] = {1, 2, 3, 4, 5, 6};
    for (int plmn_idx = 0; plmn_idx < supported_tai_items->bplmnlist_count();
         plmn_idx++) {
      supported_tai_items->add_bplmns(plmn_array);
    }
  }
  // Inserting 1 enb association
  map_uint32_enb_description_t enb_map.map = init_state->mutable_enbs();
  enb_map.insert(enb_association->sctp_assoc_id(), enb_association);
  init_state->set_num_enbs(1);

  proto_map_uint32_uint32_t mmeid2associd_map =
      init_state->mutable_mmeid2associd();
  mmeid2associd_map.insert(1, assoc_id);

  oai::S1apState state_proto;
  S1apStateConverter::state_to_proto(init_state, &state_proto);
  S1apStateConverter::proto_to_state(state_proto, final_state);

  EXPECT_EQ(init_state->num_enbs(), final_state->num_enbs());
  EnbDescription* enbd = nullptr;
  EnbDescription* enbd_final = nullptr;
  EXPECT_EQ(init_state->mutable_enbs().get(assoc_id, &enbd),
            magma::PROTO_MAP_OK);
  EXPECT_EQ(final_state->mutable_enbs().get(assoc_id, &enbd_final),
            magma::PROTO_MAP_OK);

  EXPECT_EQ(enbd->sctp_assoc_id(), enbd_final->sctp_assoc_id());
  EXPECT_EQ(enbd->enb_id(), enbd_final->enb_id());
  EXPECT_EQ(enbd->s1_enb_state(), enbd_final->s1_enb_state());

  EXPECT_EQ(enbd->supported_ta_list().list_count(),
            enbd_final->supported_ta_list().list_count());
  EXPECT_EQ(enbd->supported_ta_list().supported_tai_items(0).tac(),
            enbd_final->supported_ta_list().supported_tai_items(0).tac());

  free_s1ap_state(init_state);
  free_s1ap_state(final_state);
}

TEST_F(S1APStateConverterTest, S1apStateConversionExpectedEnbCount) {
  sctp_assoc_id_t assoc_id = 1;
  S1apState* init_state = create_s1ap_state();
  S1apState* final_state = create_s1ap_state();

  EnbDescription* enb_association = s1ap_new_enb();
  enb_association->set_sctp_assoc_id(assoc_id);
  enb_association->set_enb_id(0xFFFFFFFF);
  enb_association->set_s1_enb_state(magma::lte::oai::S1AP_READY);

  map_uint32_enb_description_t enb_map.map = init_state->mutable_enbs();
  // Inserting 1 enb association
  enb_map.insert(enb_association->sctp_assoc_id(), enb_association);
  // state_to_proto should update num_enbs to match expected eNB count on the
  // map
  init_state->set_num_enbs(5);

  oai::S1apState state_proto;
  S1apStateConverter::state_to_proto(init_state, &state_proto);
  EXPECT_EQ(init_state->num_enbs(), 1);

  S1apStateConverter::proto_to_state(state_proto, final_state);
  EXPECT_EQ(final_state->num_enbs(), 1);

  free_s1ap_state(init_state);
  free_s1ap_state(final_state);
}

TEST_F(S1APStateConverterTest, S1apStateConversionUeContext) {
  oai::UeDescription* ue = new oai::UeDescription();
  oai::UeDescription* final_ue = new oai::UeDescription();

  // filling with test values
  ue->set_mme_ue_s1ap_id(1);
  ue->set_enb_ue_s1ap_id(1);
  ue->set_sctp_assoc_id(1);
  ue->set_comp_s1ap_id(S1AP_GENERATE_COMP_S1AP_ID(1, 1));
  ue->mutable_s1ap_handover_state()->set_mme_ue_s1ap_id(1);
  ue->mutable_s1ap_handover_state()->set_source_enb_id(1);
  ue->mutable_s1ap_ue_context_rel_timer()->set_id(1);
  ue->mutable_s1ap_ue_context_rel_timer()->set_msec(1000);

  oai::UeDescription ue_proto;
  S1apStateConverter::ue_to_proto(ue, &ue_proto);
  S1apStateConverter::proto_to_ue(ue_proto, final_ue);

  EXPECT_EQ(ue->comp_s1ap_id(), final_ue->comp_s1ap_id());
  EXPECT_EQ(ue->mme_ue_s1ap_id(), final_ue->mme_ue_s1ap_id());
  EXPECT_EQ(ue->s1ap_ue_context_rel_timer().id(),
            final_ue->s1ap_ue_context_rel_timer().id());

  delete ue;
  delete final_ue;
}

}  // namespace lte
}  // namespace magma
