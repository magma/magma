/*
 *
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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

using magma::lte::gateway::s1ap::EnbDescription;
using magma::lte::gateway::s1ap::S1apState;
using magma::lte::gateway::s1ap::UeDescription;

namespace magma {
namespace lte {

S1apStateConverter::~S1apStateConverter() = default;
S1apStateConverter::S1apStateConverter() = default;

void S1apStateConverter::state_to_proto(s1ap_state_t* state, S1apState* proto)
{
  proto->Clear();

  // copy over enbs
  hashtable_ts_to_proto<enb_description_t, EnbDescription>(
    &state->enbs, proto->mutable_enbs(), enb_to_proto, LOG_S1AP);

  // copy over mmeid2associd
  hashtable_rc_t ht_rc;
  mme_ue_s1ap_id_t mmeid;
  sctp_assoc_id_t associd;
  auto mmeid2associd = proto->mutable_mmeid2associd();

  hashtable_key_array_t* keys = hashtable_ts_get_keys(&state->mmeid2associd);
  if (!keys) {
    OAILOG_DEBUG(LOG_S1AP, "No keys in mmeid2associd hashtable");
  } else {
    for (uint32_t i = 0; i < keys->num_keys; i++) {
      mmeid = (mme_ue_s1ap_id_t) keys->keys[i];
      ht_rc = hashtable_ts_get(
        &state->mmeid2associd, (hash_key_t) mmeid, (void**) &associd);
      AssertFatal(ht_rc == HASH_TABLE_OK, "mmeid not in mmeid2associd");

      (*mmeid2associd)[mmeid] = associd;
    }
    FREE_HASHTABLE_KEY_ARRAY(keys);
  }

  proto->set_num_enbs(state->num_enbs);
}

void S1apStateConverter::proto_to_state(
  const S1apState& proto,
  s1ap_state_t* state)
{
  proto_to_hashtable_ts<EnbDescription, enb_description_t>(
    proto.enbs(), &state->enbs, proto_to_enb, LOG_S1AP);

  hashtable_rc_t ht_rc;
  auto mmeid2associd = proto.mmeid2associd();
  for (auto const& kv : mmeid2associd) {
    mme_ue_s1ap_id_t mmeid = (mme_ue_s1ap_id_t) kv.first;
    sctp_assoc_id_t associd = (sctp_assoc_id_t) kv.second;

    ht_rc = hashtable_ts_insert(
      &state->mmeid2associd, (hash_key_t) mmeid, (void*) (uintptr_t) associd);
    AssertFatal(ht_rc == HASH_TABLE_OK, "failed to insert associd");
  }

  state->num_enbs = proto.num_enbs();
}

void S1apStateConverter::enb_to_proto(
  enb_description_t* enb,
  gateway::s1ap::EnbDescription* proto)
{
  proto->Clear();

  proto->set_enb_id(enb->enb_id);
  proto->set_s1_state(enb->s1_state);
  proto->set_enb_name(enb->enb_name);
  proto->set_default_paging_drx(enb->default_paging_drx);
  proto->set_nb_ue_associated(enb->nb_ue_associated);
  proto->mutable_s1ap_enb_assoc_clean_up_timer()->set_id(
    enb->s1ap_enb_assoc_clean_up_timer.id);
  proto->mutable_s1ap_enb_assoc_clean_up_timer()->set_sec(
    enb->s1ap_enb_assoc_clean_up_timer.sec);
  proto->set_sctp_assoc_id(enb->sctp_assoc_id);
  proto->set_next_sctp_stream(enb->next_sctp_stream);
  proto->set_instreams(enb->instreams);
  proto->set_outstreams(enb->outstreams);

  // store ue_ids
  hashtable_uint64_ts_to_proto(&enb->ue_id_coll, proto->mutable_ue_ids());
}

void S1apStateConverter::proto_to_enb(
  const gateway::s1ap::EnbDescription& proto,
  enb_description_t* enb)
{
  memset(enb, 0, sizeof(*enb));

  enb->enb_id = proto.enb_id();
  enb->s1_state = (mme_s1_enb_state_s) proto.s1_state();
  strncpy(enb->enb_name, proto.enb_name().c_str(), sizeof(enb->enb_name));
  enb->default_paging_drx = proto.default_paging_drx();
  enb->nb_ue_associated = proto.nb_ue_associated();
  enb->s1ap_enb_assoc_clean_up_timer.id =
    proto.s1ap_enb_assoc_clean_up_timer().id();
  enb->s1ap_enb_assoc_clean_up_timer.sec =
    proto.s1ap_enb_assoc_clean_up_timer().sec();
  enb->sctp_assoc_id = proto.sctp_assoc_id();
  enb->next_sctp_stream = proto.next_sctp_stream();
  enb->instreams = proto.instreams();
  enb->outstreams = proto.outstreams();

  // load ues
  hashtable_rc_t ht_rc;
  auto ht_name = bfromcstr("s1ap_ue_coll");

  hashtable_uint64_ts_init(&enb->ue_id_coll, mme_config.max_ues,
      nullptr, ht_name);
  bdestroy(ht_name);

  auto ue_ids = proto.ue_ids();
  for (auto const& kv : ue_ids) {
    mme_ue_s1ap_id_t mme_ue_s1ap_id = kv.first;
    uint64_t comp_s1ap_id = kv.second;

    ht_rc = hashtable_uint64_ts_insert(&enb->ue_id_coll, (hash_key_t) mme_ue_s1ap_id, comp_s1ap_id);
    if (ht_rc != HASH_TABLE_OK) {
      OAILOG_DEBUG(LOG_S1AP, "Failed to insert mme_ue_s1ap_id in ue_coll_id"
                             "hashtable");
    }
  }
}
void S1apStateConverter::ue_to_proto(
  const ue_description_t* ue,
  gateway::s1ap::UeDescription* proto)
{
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
}
void S1apStateConverter::proto_to_ue(
  const gateway::s1ap::UeDescription& proto,
  ue_description_t* ue)
{
  memset(ue, 0, sizeof(*ue));

  ue->s1_ue_state = (s1_ue_state_s) proto.s1_ue_state();
  ue->enb_ue_s1ap_id = proto.enb_ue_s1ap_id();
  ue->mme_ue_s1ap_id = proto.mme_ue_s1ap_id();
  ue->sctp_assoc_id = proto.sctp_assoc_id();
  ue->sctp_stream_recv = proto.sctp_stream_recv();
  ue->sctp_stream_send = proto.sctp_stream_send();
  ue->s1ap_ue_context_rel_timer.id = proto.s1ap_ue_context_rel_timer().id();
  ue->s1ap_ue_context_rel_timer.sec = proto.s1ap_ue_context_rel_timer().sec();
}

void S1apStateConverter::s1ap_imsi_map_to_proto(
  const s1ap_imsi_map_t* s1ap_imsi_map,
  gateway::s1ap::S1apImsiMap* s1ap_imsi_proto)
{
  hashtable_uint64_ts_to_proto(
    s1ap_imsi_map->mme_ue_id_imsi_htbl,
    s1ap_imsi_proto->mutable_mme_ue_id_imsi_map());
}
void S1apStateConverter::proto_to_s1ap_imsi_map(
  const gateway::s1ap::S1apImsiMap& s1ap_imsi_proto,
  s1ap_imsi_map_t* s1ap_imsi_map)
{
  proto_to_hashtable_uint64_ts(
    s1ap_imsi_proto.mme_ue_id_imsi_map(),
    s1ap_imsi_map->mme_ue_id_imsi_htbl);
}

} // namespace lte
} // namespace magma
