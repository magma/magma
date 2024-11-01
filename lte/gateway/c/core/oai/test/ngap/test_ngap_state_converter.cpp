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
#include "lte/gateway/c/core/oai/test/ngap/mock_utils.h"
#include "lte/gateway/c/core/oai/test/ngap/util_ngap_pkt.hpp"
#include <gtest/gtest.h>

extern "C" {
#include "lte/gateway/c/core/oai/common/log.h"
#include "Ngap_NGAP-PDU.h"
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_amf_handlers.h"
#include "lte/gateway/c/core/oai/include/amf_config.hpp"
}

#include "lte/gateway/c/core/oai/tasks/ngap/ngap_state_converter.hpp"
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_state_manager.hpp"
#include "lte/gateway/c/core/oai/include/map.h"

using ::testing::Test;

namespace magma5g {

class NgapStateConverterTest : public testing::Test {
 protected:
  void SetUp() {
    itti_init(TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info,
              NULL, NULL);

    amf_config_init(&amf_config);
    amf_config.use_stateless = true;
    ngap_state_init(amf_config.max_ues, amf_config.max_gnbs,
                    amf_config.use_stateless);

    state = get_ngap_state(false);
    imsi_map = get_ngap_imsi_map();
  }

  void TearDown() {
    ngap_state_exit();
    itti_free_desc_threads();
    amf_config_free(&amf_config);
    state = NULL;
    imsi_map = NULL;
    gNB_ref = NULL;
    ue_ref = NULL;
  }

  // This Function mocks NGAP task TearDown.
  void PseudoNgapTearDown() {
    ngap_state_exit();
    itti_free_desc_threads();
    amf_config_free(&amf_config);
    state = NULL;
    imsi_map = NULL;
    gNB_ref = NULL;
    ue_ref = NULL;
  }

  // This Function mocks NGAP task Setup.
  void PseudoNgapSetup() {
    itti_init(TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info,
              NULL, NULL);

    amf_config_init(&amf_config);
    amf_config.use_stateless = true;
    ngap_state_init(amf_config.max_ues, amf_config.max_gnbs,
                    amf_config.use_stateless);

    state = get_ngap_state(false);
    imsi_map = get_ngap_imsi_map();
  }

  ngap_state_t* state = NULL;
  ngap_imsi_map_t* imsi_map = NULL;
  gnb_description_t* gNB_ref = NULL;
  m5g_ue_description_t* ue_ref = NULL;
  const unsigned int AMF_UE_NGAP_ID = 0x05;
  const unsigned int gNB_UE_NGAP_ID = 0x09;
  bool is_task_state_same = false;
  bool is_ue_state_same = false;
  imsi64_t imsi64;
};

TEST_F(NgapStateConverterTest, NgapStateConversionSuccess) {
  sctp_assoc_id_t assoc_id = 1;
  ngap_state_t* init_state = create_ngap_state(2, 2);
  ngap_state_t* final_state = create_ngap_state(2, 2);

  gnb_description_t* gnb_association = ngap_new_gnb(init_state);
  gnb_association->sctp_assoc_id = assoc_id;
  gnb_association->gnb_id = 0xFFFFFFFF;
  gnb_association->ng_state = NGAP_READY;
  gnb_association->instreams = 1;
  gnb_association->outstreams = 2;
  gnb_association->next_sctp_stream = 3;
  gnb_association->nb_ue_associated = 2;
  gnb_association->default_paging_drx = 100;
  memset(gnb_association->gnb_name, 'A', sizeof(gnb_association->gnb_name));
  uint64_t data = 0;
  void* associd1 = NULL;
  void* associd2 = NULL;

  // filling ue_id_coll
  hashtable_uint64_ts_insert(&gnb_association->ue_id_coll, (const hash_key_t)1,
                             17);
  hashtable_uint64_ts_insert(&gnb_association->ue_id_coll, (const hash_key_t)2,
                             25);

  // filling supported_tai_items
  m5g_supported_ta_list_t* gnb_ta_list = &gnb_association->supported_ta_list;
  gnb_ta_list->list_count = 1;
  gnb_ta_list->supported_tai_items[0].tac = 1;
  gnb_ta_list->supported_tai_items[0].bplmnlist_count = 1;
  gnb_ta_list->supported_tai_items[0].bplmn_list[0].plmn_id.mcc_digit1 = 1;
  gnb_ta_list->supported_tai_items[0].bplmn_list[0].plmn_id.mcc_digit2 = 1;
  gnb_ta_list->supported_tai_items[0].bplmn_list[0].plmn_id.mcc_digit3 = 1;
  gnb_ta_list->supported_tai_items[0].bplmn_list[0].plmn_id.mnc_digit1 = 1;
  gnb_ta_list->supported_tai_items[0].bplmn_list[0].plmn_id.mnc_digit2 = 1;
  gnb_ta_list->supported_tai_items[0].bplmn_list[0].plmn_id.mnc_digit3 = 1;

  // Inserting 1 enb association
  hashtable_ts_insert(&init_state->gnbs,
                      (const hash_key_t)gnb_association->sctp_assoc_id,
                      reinterpret_cast<void*>(gnb_association));
  init_state->num_gnbs = 1;

  hashtable_ts_insert(&init_state->amfid2associd, (const hash_key_t)1,
                      reinterpret_cast<void**>(&assoc_id));
  oai::NgapState state_proto;
  NgapStateConverter::state_to_proto(init_state, &state_proto);
  NgapStateConverter::proto_to_state(state_proto, final_state);

  EXPECT_EQ(init_state->num_gnbs, final_state->num_gnbs);
  gnb_description_t* gnbd = nullptr;
  gnb_description_t* gnbd_final = nullptr;
  EXPECT_EQ(hashtable_ts_get(&init_state->gnbs, (const hash_key_t)assoc_id,
                             reinterpret_cast<void**>(&gnbd)),
            HASH_TABLE_OK);
  EXPECT_EQ(hashtable_ts_get(&final_state->gnbs, (const hash_key_t)assoc_id,
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
      hashtable_uint64_ts_get(&gnbd->ue_id_coll, (const hash_key_t)1, &data),
      HASH_TABLE_OK);
  EXPECT_EQ(data, 17);
  data = 0;
  EXPECT_EQ(hashtable_uint64_ts_get(&gnbd_final->ue_id_coll,
                                    (const hash_key_t)2, &data),
            HASH_TABLE_OK);
  EXPECT_EQ(data, 25);
  EXPECT_EQ(hashtable_ts_get(&init_state->amfid2associd, (const hash_key_t)1,
                             reinterpret_cast<void**>(&associd1)),
            HASH_TABLE_OK);
  EXPECT_EQ(hashtable_ts_get(&final_state->amfid2associd, (const hash_key_t)1,
                             reinterpret_cast<void**>(&associd2)),
            HASH_TABLE_OK);
  EXPECT_EQ(gnbd->supported_ta_list.list_count,
            gnbd_final->supported_ta_list.list_count);
  EXPECT_EQ(gnbd->supported_ta_list.supported_tai_items[0].tac,
            gnbd_final->supported_ta_list.supported_tai_items[0].tac);

  free_ngap_state(init_state);
  free_ngap_state(final_state);
}

TEST(NgapStateConversionUeContext, NgapStateConversionUeContext) {
  m5g_ue_description_t* ue = reinterpret_cast<m5g_ue_description_t*>(
      calloc(1, sizeof(m5g_ue_description_t)));
  m5g_ue_description_t* final_ue = reinterpret_cast<m5g_ue_description_t*>(
      calloc(1, sizeof(m5g_ue_description_t)));

  // filling with test values
  ue->ng_ue_state = NGAP_UE_CONNECTED;
  ue->gnb_ue_ngap_id = 7;
  ue->amf_ue_ngap_id = 4;
  ue->sctp_assoc_id = 3;
  ue->comp_ngap_id =
      ngap_get_comp_ngap_id(ue->sctp_assoc_id, ue->gnb_ue_ngap_id);
  ue->sctp_stream_recv = 5;
  ue->sctp_stream_send = 5;
  ue->ngap_ue_context_rel_timer.id = 1;
  ue->ngap_ue_context_rel_timer.msec = 1000;

  oai::Ngap_UeDescription ue_proto;
  NgapStateConverter::ue_to_proto(ue, &ue_proto);
  NgapStateConverter::proto_to_ue(ue_proto, final_ue);

  EXPECT_EQ(ue->ng_ue_state, final_ue->ng_ue_state);
  EXPECT_EQ(ue->gnb_ue_ngap_id, final_ue->gnb_ue_ngap_id);
  EXPECT_EQ(ue->amf_ue_ngap_id, final_ue->amf_ue_ngap_id);
  EXPECT_EQ(ue->sctp_assoc_id, final_ue->sctp_assoc_id);
  EXPECT_EQ(ue->comp_ngap_id, final_ue->comp_ngap_id);
  EXPECT_EQ(ue->sctp_stream_recv, final_ue->sctp_stream_recv);
  EXPECT_EQ(ue->sctp_stream_send, final_ue->sctp_stream_send);

  EXPECT_EQ(ue->ngap_ue_context_rel_timer.id,
            final_ue->ngap_ue_context_rel_timer.id);
  EXPECT_EQ(ue->ngap_ue_context_rel_timer.msec,
            final_ue->ngap_ue_context_rel_timer.msec);

  free_wrapper(reinterpret_cast<void**>(&ue));
  free_wrapper(reinterpret_cast<void**>(&final_ue));
}

TEST(NgapStateConversionNgapImsimap, NgapStateConversionNgapImsimap) {
  ngap_imsi_map_t* ngap_imsi_map =
      reinterpret_cast<ngap_imsi_map_t*>(calloc(1, sizeof(ngap_imsi_map_t)));
  ngap_imsi_map_t* final_ngap_imsi_map =
      reinterpret_cast<ngap_imsi_map_t*>(calloc(1, sizeof(ngap_imsi_map_t)));
  amf_ue_ngap_id_t ue_id = 1;
  imsi64_t imsi64 = 1010000000001;
  imsi64_t final_imsi64;
  uint32_t max_ues_ = 3;

  ngap_imsi_map->amf_ue_id_imsi_htbl =
      hashtable_uint64_ts_create(max_ues_, nullptr, nullptr);
  final_ngap_imsi_map->amf_ue_id_imsi_htbl =
      hashtable_uint64_ts_create(max_ues_, nullptr, nullptr);

  EXPECT_EQ(hashtable_uint64_ts_get(ngap_imsi_map->amf_ue_id_imsi_htbl,
                                    (const hash_key_t)ue_id, &imsi64),
            HASH_TABLE_KEY_NOT_EXISTS);

  EXPECT_EQ(hashtable_uint64_ts_insert(ngap_imsi_map->amf_ue_id_imsi_htbl,
                                       (const hash_key_t)ue_id, imsi64),
            HASH_TABLE_OK);

  EXPECT_EQ(hashtable_uint64_ts_insert(ngap_imsi_map->amf_ue_id_imsi_htbl,
                                       (const hash_key_t)ue_id, imsi64),
            HASH_TABLE_SAME_KEY_VALUE_EXISTS);

  oai::NgapImsiMap ngap_imsi_proto;
  NgapStateConverter::ngap_imsi_map_to_proto(ngap_imsi_map, &ngap_imsi_proto);
  NgapStateConverter::proto_to_ngap_imsi_map(ngap_imsi_proto,
                                             final_ngap_imsi_map);

  EXPECT_EQ(hashtable_uint64_ts_get(final_ngap_imsi_map->amf_ue_id_imsi_htbl,
                                    (const hash_key_t)2, &imsi64),
            HASH_TABLE_KEY_NOT_EXISTS);

  hashtable_uint64_ts_get(final_ngap_imsi_map->amf_ue_id_imsi_htbl,
                          (const hash_key_t)ue_id, &final_imsi64);
  EXPECT_EQ(imsi64, final_imsi64);

  EXPECT_EQ(hashtable_uint64_ts_remove(final_ngap_imsi_map->amf_ue_id_imsi_htbl,
                                       (const hash_key_t)2),
            HASH_TABLE_KEY_NOT_EXISTS);

  EXPECT_EQ(hashtable_uint64_ts_remove(final_ngap_imsi_map->amf_ue_id_imsi_htbl,
                                       (const hash_key_t)ue_id),
            HASH_TABLE_OK);

  EXPECT_EQ(hashtable_uint64_ts_get(final_ngap_imsi_map->amf_ue_id_imsi_htbl,
                                    (const hash_key_t)ue_id, &imsi64),
            HASH_TABLE_KEY_NOT_EXISTS);

  hashtable_uint64_ts_destroy(ngap_imsi_map->amf_ue_id_imsi_htbl);
  hashtable_uint64_ts_destroy(final_ngap_imsi_map->amf_ue_id_imsi_htbl);
  free_wrapper(reinterpret_cast<void**>(&ngap_imsi_map));
  free_wrapper(reinterpret_cast<void**>(&final_ngap_imsi_map));
}

unsigned char initial_ue_message_hexbuf[] = {
    0x00, 0x0f, 0x40, 0x48, 0x00, 0x00, 0x05, 0x00, 0x55, 0x00, 0x02,
    0x00, 0x01, 0x00, 0x26, 0x00, 0x1a, 0x19, 0x7e, 0x00, 0x41, 0x79,
    0x00, 0x0d, 0x01, 0x22, 0x62, 0x54, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x01, 0x2e, 0x04, 0xf0, 0xf0, 0xf0, 0xf0, 0x00,
    0x79, 0x00, 0x13, 0x50, 0x22, 0x42, 0x65, 0x00, 0x00, 0x00, 0x01,
    0x00, 0x22, 0x42, 0x65, 0x00, 0x00, 0x01, 0xe4, 0xf7, 0x04, 0x44,
    0x00, 0x5a, 0x40, 0x01, 0x18, 0x00, 0x70, 0x40, 0x01, 0x00};

unsigned char dl_nas_auth_req_msg[] = {
    0x7e, 0x00, 0x56, 0x00, 0x02, 0x00, 0x00, 0x21, 0xb4, 0x74, 0x3d,
    0x51, 0x76, 0xb8, 0xe5, 0x45, 0xe1, 0xdc, 0x03, 0x68, 0x25, 0x9a,
    0x67, 0x6c, 0x20, 0x10, 0xa4, 0x9b, 0x6b, 0x3d, 0x65, 0x6d, 0x80,
    0x00, 0x41, 0xc5, 0x72, 0x9e, 0xd9, 0xe1, 0xf0, 0xd6};

// 1.Stateless triggered after SCTP ASSOCIATION
TEST_F(NgapStateConverterTest, TestAfterSctpAssociation) {
  bstring ran_cp_ipaddr;
  sctp_new_peer_t peerInfo;
  ran_cp_ipaddr = bfromcstr("\xc0\xa8\x3c\x8d");
  peerInfo = {
      .instreams = 1,
      .outstreams = 2,
      .assoc_id = 3,
      .ran_cp_ipaddr = ran_cp_ipaddr,
  };
  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // This flag is triggered to check stateless after SCTP ASSOCIATION.
  is_ue_state_same = true;
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  // Initial UE message
  Ngap_NGAP_PDU_t decoded_pdu = {};
  uint16_t length_initial_ue_message_hexbuf =
      sizeof(initial_ue_message_hexbuf) / sizeof(unsigned char);

  bstring ngap_initial_ue_msg =
      blk2bstr(initial_ue_message_hexbuf, length_initial_ue_message_hexbuf);

  // Check if the pdu can be decoded
  ASSERT_EQ(ngap_amf_decode_pdu(&decoded_pdu, ngap_initial_ue_msg), RETURNok);

  // Check if initial UE message is handled successfully
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id,
                                    peerInfo.instreams, &decoded_pdu),
            RETURNok);

  Ngap_InitialUEMessage_t* container;

  container =
      &(decoded_pdu.choice.initiatingMessage.value.choice.InitialUEMessage);
  Ngap_InitialUEMessage_IEs_t* ie = NULL;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_InitialUEMessage_IEs_t, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID);

  // Check if Ran_UE_NGAP_ID is present in initial message
  ASSERT_TRUE(ie != NULL);

  // Check if gNB exists
  gNB_ref = ngap_state_get_gnb(state, peerInfo.assoc_id);
  ASSERT_TRUE(gNB_ref != NULL);

  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  uint64_t comp_ngap_id;
  gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID);

  // Check for UE associated with gNB
  ue_ref = ngap_state_get_ue_gnbid(gNB_ref->sctp_assoc_id, gnb_ue_ngap_id);
  ASSERT_TRUE(ue_ref != NULL);

  // Check if UE is pointing to invalid ID if it is initial ue message
  EXPECT_EQ(ue_ref->amf_ue_ngap_id, INVALID_AMF_UE_NGAP_ID);

  itti_amf_app_ngap_amf_ue_id_notification_t notification_p;
  memset(&notification_p, 0,
         sizeof(itti_amf_app_ngap_amf_ue_id_notification_t));
  notification_p.gnb_ue_ngap_id = gnb_ue_ngap_id;
  notification_p.amf_ue_ngap_id = 1;
  notification_p.sctp_assoc_id = gNB_ref->sctp_assoc_id;
  comp_ngap_id = ngap_get_comp_ngap_id(notification_p.sctp_assoc_id,
                                       notification_p.gnb_ue_ngap_id);

  ngap_handle_amf_ue_id_notification(state, &notification_p);

  // Check for UE associated with gNB
  ue_ref = ngap_state_get_ue_gnbid(gNB_ref->sctp_assoc_id, gnb_ue_ngap_id);
  ASSERT_TRUE(ue_ref != NULL);

  // Check if UE is pointing to amf_ue_ngap_id
  EXPECT_EQ(ue_ref->amf_ue_ngap_id, 1);

  MessageDef* message_p = NULL;
  bstring buffer;
  unsigned int len = sizeof(dl_nas_auth_req_msg) / sizeof(unsigned char);

  message_p = itti_alloc_new_message(TASK_AMF_APP, NGAP_NAS_DL_DATA_REQ);

  NGAP_NAS_DL_DATA_REQ(message_p).gnb_ue_ngap_id = gnb_ue_ngap_id;
  NGAP_NAS_DL_DATA_REQ(message_p).amf_ue_ngap_id = ue_ref->amf_ue_ngap_id;
  message_p->ittiMsgHeader.imsi = 0x311480000000001;
  buffer = bfromcstralloc(len, "\0");
  memcpy(buffer->data, dl_nas_auth_req_msg, len);
  buffer->slen = len;
  NGAP_NAS_DL_DATA_REQ(message_p).nas_msg = bstrcpy(buffer);

  EXPECT_EQ(ngap_generate_downlink_nas_transport(
                state, gnb_ue_ngap_id, ue_ref->amf_ue_ngap_id,
                &NGAP_NAS_DL_DATA_REQ(message_p).nas_msg,
                message_p->ittiMsgHeader.imsi),
            RETURNok);

  NGAPClientServicer::getInstance().map_ngap_state_proto_str.clear();
  NGAPClientServicer::getInstance().map_ngap_uestate_proto_str.clear();
  NGAPClientServicer::getInstance().map_imsi_table_proto_str.clear();

  EXPECT_EQ(hashtable_uint64_ts_get(imsi_map->amf_ue_id_imsi_htbl,
                                    (const hash_key_t)ue_ref->amf_ue_ngap_id,
                                    &imsi64),
            HASH_TABLE_OK);
  EXPECT_EQ(
      hashtable_ts_get(&state->gnbs, (const hash_key_t)gNB_ref->sctp_assoc_id,
                       reinterpret_cast<void**>(gNB_ref)),
      HASH_TABLE_OK);
  EXPECT_EQ(hashtable_ts_get(&state->amfid2associd,
                             (const hash_key_t)ue_ref->amf_ue_ngap_id,
                             reinterpret_cast<void**>(&gNB_ref->sctp_assoc_id)),
            HASH_TABLE_OK);
  EXPECT_EQ(hashtable_uint64_ts_get(&gNB_ref->ue_id_coll,
                                    (const hash_key_t)ue_ref->amf_ue_ngap_id,
                                    &comp_ngap_id),
            HASH_TABLE_OK);

  if (!is_task_state_same) {
    put_ngap_state();
  }
  if (!is_ue_state_same) {
    put_ngap_imsi_map();
    put_ngap_ue_state(imsi64);
  }

  EXPECT_FALSE(
      NGAPClientServicer::getInstance().map_ngap_state_proto_str.isEmpty());
  EXPECT_TRUE(
      NGAPClientServicer::getInstance().map_ngap_uestate_proto_str.isEmpty());
  EXPECT_TRUE(
      NGAPClientServicer::getInstance().map_imsi_table_proto_str.isEmpty());

  EXPECT_EQ(NGAPClientServicer::getInstance().map_ngap_state_proto_str.size(),
            1);
  EXPECT_EQ(NGAPClientServicer::getInstance().map_ngap_uestate_proto_str.size(),
            0);
  EXPECT_EQ(NGAPClientServicer::getInstance().map_imsi_table_proto_str.size(),
            0);

  // Calling PseudoNgapTearDown() and PseudoNgapSetup() simulates a service
  // restart.

  NgapStateConverterTest::PseudoNgapTearDown();

  EXPECT_EQ(state, nullptr);
  EXPECT_EQ(imsi_map, nullptr);

  NgapStateConverterTest::PseudoNgapSetup();

  // Check if gNB exists
  gNB_ref = ngap_state_get_gnb(state, peerInfo.assoc_id);
  ASSERT_TRUE(gNB_ref != NULL);

  EXPECT_NE(state, nullptr);
  EXPECT_NE(imsi_map, nullptr);

  // Only state->gnbs hashtable will sync in this test case
  EXPECT_EQ(hashtable_ts_is_key_exists(
                &state->gnbs, (const hash_key_t)gNB_ref->sctp_assoc_id),
            HASH_TABLE_OK);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);
  itti_free_msg_content(message_p);
  free(message_p);
  bdestroy_wrapper(&ran_cp_ipaddr);
  bdestroy(buffer);
  bdestroy(ngap_initial_ue_msg);
}

// 2. Stateless feature Unit Test for NGAP
TEST_F(NgapStateConverterTest, TestNgapServiceRestart) {
  bstring ran_cp_ipaddr;
  sctp_new_peer_t peerInfo;
  ran_cp_ipaddr = bfromcstr("\xc0\xa8\x3c\x8d");
  peerInfo = {
      .instreams = 1,
      .outstreams = 2,
      .assoc_id = 3,
      .ran_cp_ipaddr = ran_cp_ipaddr,
  };
  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  // Initial UE message
  Ngap_NGAP_PDU_t decoded_pdu = {};
  uint16_t length_initial_ue_message_hexbuf =
      sizeof(initial_ue_message_hexbuf) / sizeof(unsigned char);

  bstring ngap_initial_ue_msg =
      blk2bstr(initial_ue_message_hexbuf, length_initial_ue_message_hexbuf);

  // Check if the pdu can be decoded
  ASSERT_EQ(ngap_amf_decode_pdu(&decoded_pdu, ngap_initial_ue_msg), RETURNok);

  // check if initial UE message is handled successfully
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id,
                                    peerInfo.instreams, &decoded_pdu),
            RETURNok);

  Ngap_InitialUEMessage_t* container;

  container =
      &(decoded_pdu.choice.initiatingMessage.value.choice.InitialUEMessage);
  Ngap_InitialUEMessage_IEs_t* ie = NULL;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_InitialUEMessage_IEs_t, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID);

  // Check if Ran_UE_NGAP_ID is present in initial message
  ASSERT_TRUE(ie != NULL);

  // Check if gNB exists
  gNB_ref = ngap_state_get_gnb(state, peerInfo.assoc_id);
  ASSERT_TRUE(gNB_ref != NULL);

  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  uint64_t comp_ngap_id;
  gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID);

  // Check for UE associated with gNB
  ue_ref = ngap_state_get_ue_gnbid(gNB_ref->sctp_assoc_id, gnb_ue_ngap_id);
  ASSERT_TRUE(ue_ref != NULL);

  // Check if UE is pointing to invalid ID if it is initial ue message
  EXPECT_EQ(ue_ref->amf_ue_ngap_id, INVALID_AMF_UE_NGAP_ID);

  itti_amf_app_ngap_amf_ue_id_notification_t notification_p;
  memset(&notification_p, 0,
         sizeof(itti_amf_app_ngap_amf_ue_id_notification_t));
  notification_p.gnb_ue_ngap_id = gnb_ue_ngap_id;
  notification_p.amf_ue_ngap_id = 1;
  notification_p.sctp_assoc_id = gNB_ref->sctp_assoc_id;
  comp_ngap_id = ngap_get_comp_ngap_id(notification_p.sctp_assoc_id,
                                       notification_p.gnb_ue_ngap_id);

  ngap_handle_amf_ue_id_notification(state, &notification_p);

  // Check for UE associated with gNB
  ue_ref = ngap_state_get_ue_gnbid(gNB_ref->sctp_assoc_id, gnb_ue_ngap_id);
  ASSERT_TRUE(ue_ref != NULL);

  // Check if UE is pointing to amf_ue_ngap_id
  EXPECT_EQ(ue_ref->amf_ue_ngap_id, 1);

  MessageDef* message_p = NULL;
  bstring buffer;
  unsigned int len = sizeof(dl_nas_auth_req_msg) / sizeof(unsigned char);

  message_p = itti_alloc_new_message(TASK_AMF_APP, NGAP_NAS_DL_DATA_REQ);

  NGAP_NAS_DL_DATA_REQ(message_p).gnb_ue_ngap_id = gnb_ue_ngap_id;
  NGAP_NAS_DL_DATA_REQ(message_p).amf_ue_ngap_id = ue_ref->amf_ue_ngap_id;
  message_p->ittiMsgHeader.imsi = 0x311480000000001;
  buffer = bfromcstralloc(len, "\0");
  memcpy(buffer->data, dl_nas_auth_req_msg, len);
  buffer->slen = len;
  NGAP_NAS_DL_DATA_REQ(message_p).nas_msg = bstrcpy(buffer);

  EXPECT_EQ(ngap_generate_downlink_nas_transport(
                state, gnb_ue_ngap_id, ue_ref->amf_ue_ngap_id,
                &NGAP_NAS_DL_DATA_REQ(message_p).nas_msg,
                message_p->ittiMsgHeader.imsi),
            RETURNok);

  NGAPClientServicer::getInstance().map_ngap_state_proto_str.clear();
  NGAPClientServicer::getInstance().map_ngap_uestate_proto_str.clear();
  NGAPClientServicer::getInstance().map_imsi_table_proto_str.clear();

  EXPECT_EQ(hashtable_uint64_ts_get(imsi_map->amf_ue_id_imsi_htbl,
                                    (const hash_key_t)ue_ref->amf_ue_ngap_id,
                                    &imsi64),
            HASH_TABLE_OK);
  EXPECT_EQ(
      hashtable_ts_get(&state->gnbs, (const hash_key_t)gNB_ref->sctp_assoc_id,
                       reinterpret_cast<void**>(gNB_ref)),
      HASH_TABLE_OK);
  EXPECT_EQ(hashtable_ts_get(&state->amfid2associd,
                             (const hash_key_t)ue_ref->amf_ue_ngap_id,
                             reinterpret_cast<void**>(&gNB_ref->sctp_assoc_id)),
            HASH_TABLE_OK);
  EXPECT_EQ(hashtable_uint64_ts_get(&gNB_ref->ue_id_coll,
                                    (const hash_key_t)ue_ref->amf_ue_ngap_id,
                                    &comp_ngap_id),
            HASH_TABLE_OK);

  if (!is_task_state_same) {
    put_ngap_state();
  }
  if (!is_ue_state_same) {
    put_ngap_imsi_map();
    put_ngap_ue_state(imsi64);
  }

  EXPECT_FALSE(
      NGAPClientServicer::getInstance().map_ngap_state_proto_str.isEmpty());
  EXPECT_FALSE(
      NGAPClientServicer::getInstance().map_ngap_uestate_proto_str.isEmpty());
  EXPECT_FALSE(
      NGAPClientServicer::getInstance().map_imsi_table_proto_str.isEmpty());

  EXPECT_EQ(NGAPClientServicer::getInstance().map_ngap_state_proto_str.size(),
            1);
  EXPECT_EQ(NGAPClientServicer::getInstance().map_ngap_uestate_proto_str.size(),
            1);
  EXPECT_EQ(NGAPClientServicer::getInstance().map_imsi_table_proto_str.size(),
            1);

  // Calling PseudoNgapTearDown() and PseudoNgapSetup() simulates a service
  // restart.

  NgapStateConverterTest::PseudoNgapTearDown();

  EXPECT_EQ(state, nullptr);
  EXPECT_EQ(imsi_map, nullptr);

  NgapStateConverterTest::PseudoNgapSetup();

  // Check if gNB exists
  gNB_ref = ngap_state_get_gnb(state, peerInfo.assoc_id);
  ASSERT_TRUE(gNB_ref != NULL);

  // Check for UE associated with gNB
  ue_ref = ngap_state_get_ue_gnbid(gNB_ref->sctp_assoc_id, gnb_ue_ngap_id);
  ASSERT_TRUE(ue_ref != NULL);

  // Check if UE is pointing to amf_ue_ngap_id
  EXPECT_EQ(ue_ref->amf_ue_ngap_id, 1);

  EXPECT_NE(state, nullptr);
  EXPECT_NE(imsi_map, nullptr);
  EXPECT_EQ(hashtable_ts_is_key_exists(
                &state->gnbs, (const hash_key_t)gNB_ref->sctp_assoc_id),
            HASH_TABLE_OK);
  EXPECT_EQ(
      hashtable_ts_is_key_exists(&state->amfid2associd,
                                 (const hash_key_t)ue_ref->amf_ue_ngap_id),
      HASH_TABLE_OK);
  EXPECT_EQ(hashtable_uint64_ts_is_key_exists(
                &gNB_ref->ue_id_coll, (const hash_key_t)ue_ref->amf_ue_ngap_id),
            HASH_TABLE_OK);
  EXPECT_EQ(hashtable_uint64_ts_is_key_exists(
                imsi_map->amf_ue_id_imsi_htbl,
                (const hash_key_t)ue_ref->amf_ue_ngap_id),
            HASH_TABLE_OK);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);
  itti_free_msg_content(message_p);
  free(message_p);
  bdestroy_wrapper(&ran_cp_ipaddr);
  bdestroy(buffer);
  bdestroy(ngap_initial_ue_msg);
}
}  // namespace magma5g
