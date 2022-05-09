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
  Source      ngap_state_converter.cpp
  Date        2020/07/28
  Author      Ashish Prajapati
  Subsystem   Access and Mobility Management Function
  Description Defines NG Application Protocol Messages

*****************************************************************************/

#include "lte/gateway/c/core/oai/tasks/ngap/ngap_state_converter.hpp"
using namespace std;
using namespace magma::lte;
using namespace magma::lte::oai;

using magma::lte::oai::GnbDescription;
using magma::lte::oai::Ngap_UeDescription;
using magma::lte::oai::NgapImsiMap;
using magma::lte::oai::NgapState;

namespace magma5g {

NgapStateConverter::~NgapStateConverter() = default;
NgapStateConverter::NgapStateConverter() = default;

void NgapStateConverter::state_to_proto(ngap_state_t* state, NgapState* proto) {
  OAILOG_FUNC_IN(LOG_NGAP);
  proto->Clear();

  // copy over gnbs
  hashtable_ts_to_proto<gnb_description_t, GnbDescription>(
      &state->gnbs, proto->mutable_gnbs(), gnb_to_proto, LOG_NGAP);

  // copy over amfid2associd
  hashtable_rc_t ht_rc;
  amf_ue_ngap_id_t amfid;
  sctp_assoc_id_t associd;
  void* associd_ptr = nullptr;
  auto amfid2associd = proto->mutable_amfid2associd();

  hashtable_key_array_t* keys = hashtable_ts_get_keys(&state->amfid2associd);
  if (!keys) {
    OAILOG_DEBUG(LOG_NGAP, "No keys in amfid2associd hashtable");
  } else {
    for (int i = 0; i < keys->num_keys; i++) {
      amfid = (amf_ue_ngap_id_t)keys->keys[i];
      ht_rc = hashtable_ts_get(&state->amfid2associd, (hash_key_t)amfid,
                               reinterpret_cast<void**>(&associd_ptr));
      AssertFatal(ht_rc == HASH_TABLE_OK, "amfid not in amfid2associd");

      if (associd_ptr) {
        associd = (sctp_assoc_id_t)(uintptr_t)associd_ptr;
        (*amfid2associd)[amfid] = associd;
      }
    }
    FREE_HASHTABLE_KEY_ARRAY(keys);
  }

  proto->set_num_gnbs(state->num_gnbs);
  OAILOG_FUNC_OUT(LOG_NGAP);
}

void NgapStateConverter::proto_to_state(const NgapState& proto,
                                        ngap_state_t* state) {
  OAILOG_FUNC_IN(LOG_NGAP);
  proto_to_hashtable_ts<GnbDescription, gnb_description_t>(
      proto.gnbs(), &state->gnbs, proto_to_gnb, LOG_NGAP);

  hashtable_rc_t ht_rc;
  auto amfid2associd = proto.amfid2associd();
  for (auto const& kv : amfid2associd) {
    amf_ue_ngap_id_t amfid = (amf_ue_ngap_id_t)kv.first;
    sctp_assoc_id_t associd = (sctp_assoc_id_t)kv.second;

    ht_rc = hashtable_ts_insert(&state->amfid2associd, (hash_key_t)amfid,
                                (void*)(uintptr_t)associd);
    AssertFatal(ht_rc == HASH_TABLE_OK, "failed to insert associd");
  }

  state->num_gnbs = proto.num_gnbs();
  OAILOG_FUNC_OUT(LOG_NGAP);
}

void NgapStateConverter::gnb_to_proto(gnb_description_t* gnb,
                                      GnbDescription* proto) {
  OAILOG_FUNC_IN(LOG_NGAP);
  proto->Clear();

  proto->set_gnb_id(gnb->gnb_id);
  proto->set_ng_state(gnb->ng_state);
  proto->set_gnb_name(gnb->gnb_name);
  proto->set_default_paging_drx(gnb->default_paging_drx);
  proto->set_nb_ue_associated(gnb->nb_ue_associated);
  proto->set_sctp_assoc_id(gnb->sctp_assoc_id);
  proto->set_next_sctp_stream(gnb->next_sctp_stream);
  proto->set_instreams(gnb->instreams);
  proto->set_outstreams(gnb->outstreams);

  // store ue_ids
  hashtable_uint64_ts_to_proto(&gnb->ue_id_coll, proto->mutable_ue_ids());
  supported_ta_list_to_proto(&gnb->supported_ta_list,
                             proto->mutable_supported_ta_list());
  OAILOG_FUNC_OUT(LOG_NGAP);
}

void NgapStateConverter::proto_to_gnb(const GnbDescription& proto,
                                      gnb_description_t* gnb) {
  OAILOG_FUNC_IN(LOG_NGAP);
  memset(gnb, 0, sizeof(*gnb));

  gnb->gnb_id = proto.gnb_id();
  gnb->ng_state = (amf_ng_gnb_state_s)proto.ng_state();
  strncpy(gnb->gnb_name, proto.gnb_name().c_str(), sizeof(gnb->gnb_name));
  gnb->default_paging_drx = proto.default_paging_drx();
  gnb->nb_ue_associated = proto.nb_ue_associated();
  gnb->sctp_assoc_id = proto.sctp_assoc_id();
  gnb->next_sctp_stream = proto.next_sctp_stream();
  gnb->instreams = proto.instreams();
  gnb->outstreams = proto.outstreams();

  // load ues
  hashtable_rc_t ht_rc;
  auto ht_name = bfromcstr("ngap_ue_coll");

  hashtable_uint64_ts_init(&gnb->ue_id_coll, amf_config.max_ues, nullptr,
                           ht_name);
  bdestroy(ht_name);

  auto ue_ids = proto.ue_ids();
  for (auto const& kv : ue_ids) {
    amf_ue_ngap_id_t amf_ue_ngap_id = kv.first;
    uint64_t comp_ngap_id = kv.second;

    ht_rc = hashtable_uint64_ts_insert(
        &gnb->ue_id_coll, (hash_key_t)amf_ue_ngap_id, comp_ngap_id);
    if (ht_rc != HASH_TABLE_OK) {
      OAILOG_DEBUG(LOG_NGAP,
                   "Failed to insert amf_ue_ngap_id in ue_coll_id hashtable");
    }
  }
  proto_to_supported_ta_list(&gnb->supported_ta_list,
                             proto.supported_ta_list());
  OAILOG_FUNC_OUT(LOG_NGAP);
}
void NgapStateConverter::ue_to_proto(const m5g_ue_description_t* ue,
                                     Ngap_UeDescription* proto) {
  OAILOG_FUNC_IN(LOG_NGAP);
  proto->Clear();

  proto->set_ng_ue_state(ue->ng_ue_state);
  proto->set_gnb_ue_ngap_id(ue->gnb_ue_ngap_id);
  proto->set_amf_ue_ngap_id(ue->amf_ue_ngap_id);
  proto->set_sctp_assoc_id(ue->sctp_assoc_id);
  proto->set_sctp_stream_recv(ue->sctp_stream_recv);
  proto->set_sctp_stream_send(ue->sctp_stream_send);
  proto->mutable_ngap_ue_context_rel_timer()->set_id(
      ue->ngap_ue_context_rel_timer.id);
  proto->mutable_ngap_ue_context_rel_timer()->set_msec(
      ue->ngap_ue_context_rel_timer.msec);
  OAILOG_FUNC_OUT(LOG_NGAP);
}
void NgapStateConverter::proto_to_ue(const Ngap_UeDescription& proto,
                                     m5g_ue_description_t* ue) {
  OAILOG_FUNC_IN(LOG_NGAP);
  memset(ue, 0, sizeof(*ue));

  ue->ng_ue_state = (ng_ue_state_s)proto.ng_ue_state();
  ue->gnb_ue_ngap_id = proto.gnb_ue_ngap_id();
  ue->amf_ue_ngap_id = proto.amf_ue_ngap_id();
  ue->sctp_assoc_id = proto.sctp_assoc_id();
  ue->sctp_stream_recv = proto.sctp_stream_recv();
  ue->sctp_stream_send = proto.sctp_stream_send();
  ue->ngap_ue_context_rel_timer.id = proto.ngap_ue_context_rel_timer().id();
  ue->ngap_ue_context_rel_timer.msec = proto.ngap_ue_context_rel_timer().msec();

  ue->comp_ngap_id =
      ngap_get_comp_ngap_id(ue->sctp_assoc_id, ue->gnb_ue_ngap_id);
  OAILOG_FUNC_OUT(LOG_NGAP);
}

void NgapStateConverter::ngap_imsi_map_to_proto(
    const ngap_imsi_map_t* ngap_imsi_map, NgapImsiMap* ngap_imsi_proto) {
  OAILOG_FUNC_IN(LOG_NGAP);
  hashtable_uint64_ts_to_proto(ngap_imsi_map->amf_ue_id_imsi_htbl,
                               ngap_imsi_proto->mutable_amf_ue_id_imsi_map());
  OAILOG_FUNC_OUT(LOG_NGAP);
}
void NgapStateConverter::proto_to_ngap_imsi_map(
    const NgapImsiMap& ngap_imsi_proto, ngap_imsi_map_t* ngap_imsi_map) {
  OAILOG_FUNC_IN(LOG_NGAP);
  proto_to_hashtable_uint64_ts(ngap_imsi_proto.amf_ue_id_imsi_map(),
                               ngap_imsi_map->amf_ue_id_imsi_htbl);
  OAILOG_FUNC_OUT(LOG_NGAP);
}

void NgapStateConverter::supported_ta_list_to_proto(
    const m5g_supported_ta_list_t* supported_ta_list,
    oai::Ngap_SupportedTaList* supported_ta_list_proto) {
  OAILOG_FUNC_IN(LOG_NGAP);
  supported_ta_list_proto->set_list_count(supported_ta_list->list_count);
  for (int idx = 0; idx < supported_ta_list->list_count; idx++) {
    OAILOG_DEBUG(LOG_NGAP, "Writing Ngap_Supported TAI list at index %d", idx);
    oai::Ngap_SupportedTaiItems* supported_tai_item =
        supported_ta_list_proto->add_supported_tai_items();
    supported_tai_item_to_proto(&supported_ta_list->supported_tai_items[idx],
                                supported_tai_item);
  }
  OAILOG_FUNC_OUT(LOG_NGAP);
}

void NgapStateConverter::proto_to_supported_ta_list(
    m5g_supported_ta_list_t* supported_ta_list_state,
    const oai::Ngap_SupportedTaList& supported_ta_list_proto) {
  supported_ta_list_state->list_count = supported_ta_list_proto.list_count();
  OAILOG_FUNC_IN(LOG_NGAP);
  for (int idx = 0; idx < supported_ta_list_state->list_count; idx++) {
    OAILOG_DEBUG(LOG_MME_APP, "reading supported ta list at index %d", idx);
    proto_to_supported_tai_items(
        &supported_ta_list_state->supported_tai_items[idx],
        supported_ta_list_proto.supported_tai_items(idx));
  }
  OAILOG_FUNC_OUT(LOG_NGAP);
}

void NgapStateConverter::supported_tai_item_to_proto(
    const m5g_supported_tai_items_t* state_supported_tai_item,
    oai::Ngap_SupportedTaiItems* supported_tai_item_proto) {
  supported_tai_item_proto->set_tac(state_supported_tai_item->tac);
  supported_tai_item_proto->set_bplmnlist_count(
      state_supported_tai_item->bplmnlist_count);
  char plmn_array[PLMN_BYTES] = {0};
  OAILOG_FUNC_IN(LOG_NGAP);
  for (int idx = 0; idx < state_supported_tai_item->bplmnlist_count; idx++) {
    plmn_array[0] = static_cast<char>(
        state_supported_tai_item->bplmn_list[idx].plmn_id.mcc_digit1 +
        ASCII_ZERO);
    plmn_array[1] = static_cast<char>(
        state_supported_tai_item->bplmn_list[idx].plmn_id.mcc_digit2 +
        ASCII_ZERO);
    plmn_array[2] = static_cast<char>(
        state_supported_tai_item->bplmn_list[idx].plmn_id.mcc_digit3 +
        ASCII_ZERO);
    plmn_array[3] = static_cast<char>(
        state_supported_tai_item->bplmn_list[idx].plmn_id.mnc_digit1 +
        ASCII_ZERO);
    plmn_array[4] = static_cast<char>(
        state_supported_tai_item->bplmn_list[idx].plmn_id.mnc_digit2 +
        ASCII_ZERO);
    plmn_array[5] = static_cast<char>(
        state_supported_tai_item->bplmn_list[idx].plmn_id.mnc_digit3 +
        ASCII_ZERO);
    supported_tai_item_proto->add_bplmns(plmn_array);
  }
  OAILOG_FUNC_OUT(LOG_NGAP);
}

void NgapStateConverter::proto_to_supported_tai_items(
    m5g_supported_tai_items_t* supported_tai_item_state,
    const oai::Ngap_SupportedTaiItems& supported_tai_item_proto) {
  supported_tai_item_state->tac = supported_tai_item_proto.tac();
  supported_tai_item_state->bplmnlist_count =
      supported_tai_item_proto.bplmnlist_count();
  OAILOG_FUNC_IN(LOG_NGAP);
  for (int idx = 0; idx < supported_tai_item_state->bplmnlist_count; idx++) {
    supported_tai_item_state->bplmn_list[idx].plmn_id.mcc_digit1 =
        static_cast<int>(supported_tai_item_proto.bplmns(idx)[0]) - ASCII_ZERO;
    supported_tai_item_state->bplmn_list[idx].plmn_id.mcc_digit2 =
        static_cast<int>(supported_tai_item_proto.bplmns(idx)[1]) - ASCII_ZERO;
    supported_tai_item_state->bplmn_list[idx].plmn_id.mcc_digit3 =
        static_cast<int>(supported_tai_item_proto.bplmns(idx)[2]) - ASCII_ZERO;
    supported_tai_item_state->bplmn_list[idx].plmn_id.mnc_digit1 =
        static_cast<int>(supported_tai_item_proto.bplmns(idx)[3]) - ASCII_ZERO;
    supported_tai_item_state->bplmn_list[idx].plmn_id.mnc_digit2 =
        static_cast<int>(supported_tai_item_proto.bplmns(idx)[4]) - ASCII_ZERO;
    supported_tai_item_state->bplmn_list[idx].plmn_id.mnc_digit3 =
        static_cast<int>(supported_tai_item_proto.bplmns(idx)[5]) - ASCII_ZERO;
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

}  // namespace magma5g
