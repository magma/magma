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
#include "s1ap_state_converter.h"

using magma::lte::oai::EnbDescription;
using magma::lte::oai::S1apState;
using magma::lte::oai::UeDescription;

namespace magma {
namespace lte {

S1apStateConverter::~S1apStateConverter() = default;
S1apStateConverter::S1apStateConverter()  = default;

void S1apStateConverter::state_to_proto(s1ap_state_t* state, S1apState* proto) {
  proto->Clear();

  // copy over enbs
  hashtable_ts_to_proto<enb_description_t, EnbDescription>(
      &state->enbs, proto->mutable_enbs(), enb_to_proto, LOG_S1AP);

  // copy over mmeid2associd
  hashtable_rc_t ht_rc;
  mme_ue_s1ap_id_t mmeid;
  // Helper ptr so sctp_assoc_id can be casted from double ptr on
  // hashtable_ts_get
  void* sctp_id_ptr  = nullptr;
  auto mmeid2associd = proto->mutable_mmeid2associd();

  hashtable_key_array_t* keys = hashtable_ts_get_keys(&state->mmeid2associd);
  if (!keys) {
    OAILOG_DEBUG(LOG_S1AP, "No keys in mmeid2associd hashtable");
  } else {
    for (int i = 0; i < keys->num_keys; i++) {
      mmeid = (mme_ue_s1ap_id_t) keys->keys[i];
      ht_rc = hashtable_ts_get(
          &state->mmeid2associd, (hash_key_t) mmeid, (void**) &sctp_id_ptr);
      AssertFatal(ht_rc == HASH_TABLE_OK, "mmeid not in mmeid2associd");
      if (sctp_id_ptr) {
        sctp_assoc_id_t sctp_assoc_id =
            (sctp_assoc_id_t)(uintptr_t) sctp_id_ptr;
        (*mmeid2associd)[mmeid] = sctp_assoc_id;
      }
    }
    FREE_HASHTABLE_KEY_ARRAY(keys);
  }

  proto->set_num_enbs(state->num_enbs);
}

void S1apStateConverter::proto_to_state(
    const S1apState& proto, s1ap_state_t* state) {
  proto_to_hashtable_ts<EnbDescription, enb_description_t>(
      proto.enbs(), &state->enbs, proto_to_enb, LOG_S1AP);

  hashtable_rc_t ht_rc;
  auto mmeid2associd = proto.mmeid2associd();
  for (auto const& kv : mmeid2associd) {
    mme_ue_s1ap_id_t mmeid  = (mme_ue_s1ap_id_t) kv.first;
    sctp_assoc_id_t associd = (sctp_assoc_id_t) kv.second;

    ht_rc = hashtable_ts_insert(
        &state->mmeid2associd, (hash_key_t) mmeid, (void*) (uintptr_t) associd);
    AssertFatal(ht_rc == HASH_TABLE_OK, "failed to insert associd");
  }

  state->num_enbs = proto.num_enbs();
}

void S1apStateConverter::enb_to_proto(
    enb_description_t* enb, oai::EnbDescription* proto) {
  proto->Clear();

  proto->set_enb_id(enb->enb_id);
  proto->set_s1_state(enb->s1_state);
  proto->set_enb_name(enb->enb_name);
  proto->set_default_paging_drx(enb->default_paging_drx);
  proto->set_nb_ue_associated(enb->nb_ue_associated);
  proto->set_sctp_assoc_id(enb->sctp_assoc_id);
  proto->set_next_sctp_stream(enb->next_sctp_stream);
  proto->set_instreams(enb->instreams);
  proto->set_outstreams(enb->outstreams);
  proto->set_ran_cp_ipaddr(enb->ran_cp_ipaddr);
  proto->set_ran_cp_ipaddr_sz(enb->ran_cp_ipaddr_sz);

  // store ue_ids
  hashtable_uint64_ts_to_proto(&enb->ue_id_coll, proto->mutable_ue_ids());
  supported_ta_list_to_proto(
      &enb->supported_ta_list, proto->mutable_supported_ta_list());
}

void S1apStateConverter::proto_to_enb(
    const oai::EnbDescription& proto, enb_description_t* enb) {
  memset(enb, 0, sizeof(*enb));

  enb->enb_id   = proto.enb_id();
  enb->s1_state = (mme_s1_enb_state_s) proto.s1_state();
  strncpy(enb->enb_name, proto.enb_name().c_str(), sizeof(enb->enb_name));
  enb->default_paging_drx = proto.default_paging_drx();
  enb->nb_ue_associated   = proto.nb_ue_associated();
  enb->sctp_assoc_id      = proto.sctp_assoc_id();
  enb->next_sctp_stream   = proto.next_sctp_stream();
  enb->instreams          = proto.instreams();
  enb->outstreams         = proto.outstreams();
  strncpy(
      enb->ran_cp_ipaddr, proto.ran_cp_ipaddr().c_str(),
      sizeof(enb->ran_cp_ipaddr));
  enb->ran_cp_ipaddr_sz = proto.ran_cp_ipaddr_sz();

  // load ues
  hashtable_rc_t ht_rc;
  auto ht_name = bfromcstr("s1ap_ue_coll");

  hashtable_uint64_ts_init(
      &enb->ue_id_coll, mme_config.max_ues, nullptr, ht_name);
  bdestroy(ht_name);

  auto ue_ids = proto.ue_ids();
  for (auto const& kv : ue_ids) {
    mme_ue_s1ap_id_t mme_ue_s1ap_id = kv.first;
    uint64_t comp_s1ap_id           = kv.second;

    ht_rc = hashtable_uint64_ts_insert(
        &enb->ue_id_coll, (hash_key_t) mme_ue_s1ap_id, comp_s1ap_id);
    if (ht_rc != HASH_TABLE_OK) {
      OAILOG_DEBUG(
          LOG_S1AP, "Failed to insert mme_ue_s1ap_id in ue_coll_id hashtable");
    }
  }
  proto_to_supported_ta_list(
      &enb->supported_ta_list, proto.supported_ta_list());
}
void S1apStateConverter::ue_to_proto(
    const ue_description_t* ue, oai::UeDescription* proto) {
  proto->Clear();

  proto->set_s1_ue_state(ue->s1_ue_state);
  proto->set_enb_ue_s1ap_id(ue->enb_ue_s1ap_id);
  proto->set_mme_ue_s1ap_id(ue->mme_ue_s1ap_id);
  proto->set_sctp_assoc_id(ue->sctp_assoc_id);
  proto->set_sctp_stream_recv(ue->sctp_stream_recv);
  proto->set_sctp_stream_send(ue->sctp_stream_send);
  proto->mutable_s1ap_ue_context_rel_timer()->set_id(
      ue->s1ap_ue_context_rel_timer.id);
  proto->mutable_s1ap_ue_context_rel_timer()->set_sec(
      ue->s1ap_ue_context_rel_timer.sec);
  proto->mutable_s1ap_handover_state()->set_mme_ue_s1ap_id(
      ue->s1ap_handover_state.mme_ue_s1ap_id);
  proto->mutable_s1ap_handover_state()->set_source_enb_id(
      ue->s1ap_handover_state.source_enb_id);
  proto->mutable_s1ap_handover_state()->set_target_enb_id(
      ue->s1ap_handover_state.target_enb_id);
  proto->mutable_s1ap_handover_state()->set_target_enb_ue_s1ap_id(
      ue->s1ap_handover_state.target_enb_ue_s1ap_id);
  proto->mutable_s1ap_handover_state()->set_target_sctp_stream_recv(
      ue->s1ap_handover_state.target_sctp_stream_recv);
  proto->mutable_s1ap_handover_state()->set_target_sctp_stream_send(
      ue->s1ap_handover_state.target_sctp_stream_send);
}
void S1apStateConverter::proto_to_ue(
    const oai::UeDescription& proto, ue_description_t* ue) {
  memset(ue, 0, sizeof(*ue));

  ue->s1_ue_state                   = (s1_ue_state_s) proto.s1_ue_state();
  ue->enb_ue_s1ap_id                = proto.enb_ue_s1ap_id();
  ue->mme_ue_s1ap_id                = proto.mme_ue_s1ap_id();
  ue->sctp_assoc_id                 = proto.sctp_assoc_id();
  ue->sctp_stream_recv              = proto.sctp_stream_recv();
  ue->sctp_stream_send              = proto.sctp_stream_send();
  ue->s1ap_ue_context_rel_timer.id  = proto.s1ap_ue_context_rel_timer().id();
  ue->s1ap_ue_context_rel_timer.sec = proto.s1ap_ue_context_rel_timer().sec();
  ue->s1ap_handover_state.mme_ue_s1ap_id =
      proto.s1ap_handover_state().mme_ue_s1ap_id();
  ue->s1ap_handover_state.source_enb_id =
      proto.s1ap_handover_state().source_enb_id();
  ue->s1ap_handover_state.target_enb_id =
      proto.s1ap_handover_state().target_enb_id();
  ue->s1ap_handover_state.target_enb_ue_s1ap_id =
      proto.s1ap_handover_state().target_enb_ue_s1ap_id();
  ue->s1ap_handover_state.target_sctp_stream_recv =
      proto.s1ap_handover_state().target_sctp_stream_recv();
  ue->s1ap_handover_state.target_sctp_stream_send =
      proto.s1ap_handover_state().target_sctp_stream_send();

  ue->comp_s1ap_id =
      S1AP_GENERATE_COMP_S1AP_ID(ue->sctp_assoc_id, ue->enb_ue_s1ap_id);
}

void S1apStateConverter::s1ap_imsi_map_to_proto(
    const s1ap_imsi_map_t* s1ap_imsi_map, oai::S1apImsiMap* s1ap_imsi_proto) {
  hashtable_uint64_ts_to_proto(
      s1ap_imsi_map->mme_ue_id_imsi_htbl,
      s1ap_imsi_proto->mutable_mme_ue_id_imsi_map());
}
void S1apStateConverter::proto_to_s1ap_imsi_map(
    const oai::S1apImsiMap& s1ap_imsi_proto, s1ap_imsi_map_t* s1ap_imsi_map) {
  proto_to_hashtable_uint64_ts(
      s1ap_imsi_proto.mme_ue_id_imsi_map(), s1ap_imsi_map->mme_ue_id_imsi_htbl);
}

void S1apStateConverter::supported_ta_list_to_proto(
    const supported_ta_list_t* supported_ta_list,
    oai::SupportedTaList* supported_ta_list_proto) {
  supported_ta_list_proto->set_list_count(supported_ta_list->list_count);
  for (int idx = 0; idx < supported_ta_list->list_count; idx++) {
    OAILOG_DEBUG(LOG_S1AP, "Writing Supported TAI list at index %d", idx);
    oai::SupportedTaiItems* supported_tai_item =
        supported_ta_list_proto->add_supported_tai_items();
    supported_tai_item_to_proto(
        &supported_ta_list->supported_tai_items[idx], supported_tai_item);
  }
}

void S1apStateConverter::proto_to_supported_ta_list(
    supported_ta_list_t* supported_ta_list_state,
    const oai::SupportedTaList& supported_ta_list_proto) {
  supported_ta_list_state->list_count = supported_ta_list_proto.list_count();
  for (int idx = 0; idx < supported_ta_list_state->list_count; idx++) {
    OAILOG_DEBUG(LOG_MME_APP, "reading bearer context at index %d", idx);
    proto_to_supported_tai_items(
        &supported_ta_list_state->supported_tai_items[idx],
        supported_ta_list_proto.supported_tai_items(idx));
  }
}

void S1apStateConverter::supported_tai_item_to_proto(
    const supported_tai_items_t* state_supported_tai_item,
    oai::SupportedTaiItems* supported_tai_item_proto) {
  supported_tai_item_proto->set_tac(state_supported_tai_item->tac);
  supported_tai_item_proto->set_bplmnlist_count(
      state_supported_tai_item->bplmnlist_count);
  char plmn_array[PLMN_BYTES] = {0};
  for (int idx = 0; idx < state_supported_tai_item->bplmnlist_count; idx++) {
    plmn_array[0] =
        (char) (state_supported_tai_item->bplmns[idx].mcc_digit1 + ASCII_ZERO);
    plmn_array[1] =
        (char) (state_supported_tai_item->bplmns[idx].mcc_digit2 + ASCII_ZERO);
    plmn_array[2] =
        (char) (state_supported_tai_item->bplmns[idx].mcc_digit3 + ASCII_ZERO);
    plmn_array[3] =
        (char) (state_supported_tai_item->bplmns[idx].mnc_digit1 + ASCII_ZERO);
    plmn_array[4] =
        (char) (state_supported_tai_item->bplmns[idx].mnc_digit2 + ASCII_ZERO);
    plmn_array[5] =
        (char) (state_supported_tai_item->bplmns[idx].mnc_digit3 + ASCII_ZERO);
    supported_tai_item_proto->add_bplmns(plmn_array);
  }
}

void S1apStateConverter::proto_to_supported_tai_items(
    supported_tai_items_t* supported_tai_item_state,
    const oai::SupportedTaiItems& supported_tai_item_proto) {
  supported_tai_item_state->tac = supported_tai_item_proto.tac();
  supported_tai_item_state->bplmnlist_count =
      supported_tai_item_proto.bplmnlist_count();
  for (int idx = 0; idx < supported_tai_item_state->bplmnlist_count; idx++) {
    supported_tai_item_state->bplmns[idx].mcc_digit1 =
        (int) (supported_tai_item_proto.bplmns(idx)[0]) - ASCII_ZERO;
    supported_tai_item_state->bplmns[idx].mcc_digit2 =
        (int) (supported_tai_item_proto.bplmns(idx)[1]) - ASCII_ZERO;
    supported_tai_item_state->bplmns[idx].mcc_digit3 =
        (int) (supported_tai_item_proto.bplmns(idx)[2]) - ASCII_ZERO;
    supported_tai_item_state->bplmns[idx].mnc_digit1 =
        (int) (supported_tai_item_proto.bplmns(idx)[3]) - ASCII_ZERO;
    supported_tai_item_state->bplmns[idx].mnc_digit2 =
        (int) (supported_tai_item_proto.bplmns(idx)[4]) - ASCII_ZERO;
    supported_tai_item_state->bplmns[idx].mnc_digit3 =
        (int) (supported_tai_item_proto.bplmns(idx)[5]) - ASCII_ZERO;
  }
}
}  // namespace lte
}  // namespace magma
