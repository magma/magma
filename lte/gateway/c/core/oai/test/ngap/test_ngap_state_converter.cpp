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
#include "Ngap_NGAP-PDU.h"
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_amf_handlers.h"
}

#include "lte/gateway/c/core/oai/tasks/ngap/ngap_state_converter.h"
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_state_manager.h"

using ::testing::Test;

namespace magma5g {

class NGAPStateConverterTest : public ::testing::Test {
  virtual void SetUp() {}

  virtual void TearDown() {}
};

TEST_F(NGAPStateConverterTest, NgapStateConversionSuccess) {
  sctp_assoc_id_t assoc_id  = 1;
  ngap_state_t* init_state  = create_ngap_state(2, 2);
  ngap_state_t* final_state = create_ngap_state(2, 2);

  gnb_description_t* gnb_association  = ngap_new_gnb(init_state);
  gnb_association->sctp_assoc_id      = assoc_id;
  gnb_association->gnb_id             = 0xFFFFFFFF;
  gnb_association->ng_state           = NGAP_READY;
  gnb_association->instreams          = 1;
  gnb_association->outstreams         = 2;
  gnb_association->next_sctp_stream   = 3;
  gnb_association->nb_ue_associated   = 2;
  gnb_association->default_paging_drx = 100;
  memset(gnb_association->gnb_name, 'A', sizeof(gnb_association->gnb_name));
  uint64_t data  = 0;
  void* associd1 = NULL;
  void* associd2 = NULL;

  // filling ue_id_coll
  hashtable_uint64_ts_insert(
      &gnb_association->ue_id_coll, (const hash_key_t) 1, 17);
  hashtable_uint64_ts_insert(
      &gnb_association->ue_id_coll, (const hash_key_t) 2, 25);

  // filling supported_tai_items
  m5g_supported_ta_list_t* gnb_ta_list    = &gnb_association->supported_ta_list;
  gnb_ta_list->list_count                 = 1;
  gnb_ta_list->supported_tai_items[0].tac = 1;
  gnb_ta_list->supported_tai_items[0].bplmnlist_count                  = 1;
  gnb_ta_list->supported_tai_items[0].bplmn_list[0].plmn_id.mcc_digit1 = 1;
  gnb_ta_list->supported_tai_items[0].bplmn_list[0].plmn_id.mcc_digit2 = 1;
  gnb_ta_list->supported_tai_items[0].bplmn_list[0].plmn_id.mcc_digit3 = 1;
  gnb_ta_list->supported_tai_items[0].bplmn_list[0].plmn_id.mnc_digit1 = 1;
  gnb_ta_list->supported_tai_items[0].bplmn_list[0].plmn_id.mnc_digit2 = 1;
  gnb_ta_list->supported_tai_items[0].bplmn_list[0].plmn_id.mnc_digit3 = 1;

  // Inserting 1 enb association
  hashtable_ts_insert(
      &init_state->gnbs, (const hash_key_t) gnb_association->sctp_assoc_id,
      reinterpret_cast<void*>(gnb_association));
  init_state->num_gnbs = 1;

  hashtable_ts_insert(
      &init_state->amfid2associd, (const hash_key_t) 1,
      reinterpret_cast<void**>(&assoc_id));

  oai::NgapState state_proto;
  NgapStateConverter::state_to_proto(init_state, &state_proto);
  NgapStateConverter::proto_to_state(state_proto, final_state);

  EXPECT_EQ(init_state->num_gnbs, final_state->num_gnbs);
  gnb_description_t* gnbd       = nullptr;
  gnb_description_t* gnbd_final = nullptr;
  EXPECT_EQ(
      hashtable_ts_get(
          &init_state->gnbs, (const hash_key_t) assoc_id,
          reinterpret_cast<void**>(&gnbd)),
      HASH_TABLE_OK);
  EXPECT_EQ(
      hashtable_ts_get(
          &final_state->gnbs, (const hash_key_t) assoc_id,
          reinterpret_cast<void**>(&gnbd_final)),
      HASH_TABLE_OK);

  EXPECT_EQ(gnbd->sctp_assoc_id, gnbd_final->sctp_assoc_id);
  EXPECT_EQ(gnbd->gnb_id, gnbd_final->gnb_id);
  EXPECT_EQ(gnbd->ng_state, gnbd_final->ng_state);
  EXPECT_EQ(gnbd->instreams, gnbd_final->instreams);
  EXPECT_EQ(gnbd->outstreams, gnbd_final->outstreams);
  EXPECT_EQ(gnbd->next_sctp_stream, gnbd_final->next_sctp_stream);
  EXPECT_EQ(gnbd->nb_ue_associated, gnbd_final->nb_ue_associated);
  EXPECT_EQ(gnbd->default_paging_drx, gnbd_final->default_paging_drx);
  EXPECT_EQ(
      memcmp(gnbd->gnb_name, gnbd_final->gnb_name, sizeof(gnbd->gnb_name)), 0);

  EXPECT_EQ(
      hashtable_uint64_ts_get(&gnbd->ue_id_coll, (const hash_key_t) 1, &data),
      HASH_TABLE_OK);
  EXPECT_EQ(data, 17);
  data = 0;
  EXPECT_EQ(
      hashtable_uint64_ts_get(
          &gnbd_final->ue_id_coll, (const hash_key_t) 2, &data),
      HASH_TABLE_OK);
  EXPECT_EQ(data, 25);
  EXPECT_EQ(
      hashtable_ts_get(
          &init_state->amfid2associd, (const hash_key_t) 1,
          reinterpret_cast<void**>(&associd1)),
      HASH_TABLE_OK);
  EXPECT_EQ(
      hashtable_ts_get(
          &final_state->amfid2associd, (const hash_key_t) 1,
          reinterpret_cast<void**>(&associd2)),
      HASH_TABLE_OK);
  EXPECT_EQ(
      gnbd->supported_ta_list.list_count,
      gnbd_final->supported_ta_list.list_count);
  EXPECT_EQ(
      gnbd->supported_ta_list.supported_tai_items[0].tac,
      gnbd_final->supported_ta_list.supported_tai_items[0].tac);

  free_ngap_state(init_state);
  free_ngap_state(final_state);
}
}  // namespace magma5g
