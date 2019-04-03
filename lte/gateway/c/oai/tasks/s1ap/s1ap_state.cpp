/*
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
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include "s1ap_state.h"

#include <memory.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#include <cpp_redis/cpp_redis>

#include "lte/protos/s1ap_state.pb.h"

#include "ServiceConfigLoader.h"

extern "C" {
#include "bstrlib.h"

#include "assertions.h"
#include "dynamic_memory_check.h"

#include "mme_config.h"
}

using magma::lte::gateway::s1ap::EnbDescription;
using magma::lte::gateway::s1ap::S1apState;
using magma::lte::gateway::s1ap::UeDescription;

s1ap_state_t *s1ap_state_new(void);
void s1ap_state_free(s1ap_state_t *state);

s1ap_state_t *s1ap_state_from_redis(void);
void s1ap_state_to_redis(s1ap_state_t *state);

void state2proto(S1apState *proto, s1ap_state_t *state);
void proto2state(s1ap_state_t *state, S1apState *proto);

void enb2proto(EnbDescription *proto, enb_description_t *enb);
void proto2enb(enb_description_t *enb, EnbDescription *proto);

void ue2proto(UeDescription *proto, ue_description_t *ue);
void proto2ue(ue_description_t *ue, UeDescription *proto);

bool s1ap_enb_compare_by_enb_id_cb(
  const hash_key_t keyP,
  void *const elementP,
  void *parameterP,
  void **unused_res);
bool s1ap_enb_find_ue_by_mme_ue_id_cb(
  __attribute__((unused)) const hash_key_t keyP,
  void *const elementP,
  void *parameterP,
  void **resultP);

bool in_use = false;
std::shared_ptr<cpp_redis::client> client = nullptr;
s1ap_state_t *state_cache = NULL;

int s1ap_state_init(void)
{
  magma::ServiceConfigLoader loader;

  auto config = loader.load_service_config("redis");
  auto port = config["port"].as<uint32_t>();

  client = std::make_shared<cpp_redis::client>();
  client->connect("127.0.0.1", port, nullptr);

  if (!client->is_connected()) return -1;

  in_use = false;
  state_cache = s1ap_state_from_redis();

  return 0;
}

void s1ap_state_exit(void)
{
  AssertFatal(!in_use, "Exiting without committing s1ap state");

  s1ap_state_free(state_cache);

  client = nullptr;
  state_cache = NULL;
}

s1ap_state_t *s1ap_state_get(void)
{
  AssertFatal(state_cache != NULL, "s1ap state cache was NULL");
  AssertFatal(!in_use, "Tried to get s1ap_state twice without put'ing it");

  in_use = true;

  return state_cache;
}

void s1ap_state_put(s1ap_state_t *state)
{
  AssertFatal(in_use, "Tried to put s1ap_state while it was not in use");

  state_cache = state;

  s1ap_state_to_redis(state_cache);

  in_use = false;
}

enb_description_t *s1ap_state_get_enb(
  s1ap_state_t *state,
  sctp_assoc_id_t assoc_id)
{
  enb_description_t *enb = NULL;

  hashtable_ts_get(&state->enbs, (const hash_key_t) assoc_id, (void **) &enb);

  return enb;
}

ue_description_t *s1ap_state_get_ue_enbid(
  s1ap_state_t *state,
  enb_description_t *enb,
  enb_ue_s1ap_id_t enb_ue_s1ap_id)
{
  ue_description_t *ue = NULL;

  hashtable_ts_get(
    &enb->ue_coll, (const hash_key_t) enb_ue_s1ap_id, (void **) &ue);

  return ue;
}

ue_description_t *s1ap_state_get_ue_mmeid(
  s1ap_state_t *state,
  mme_ue_s1ap_id_t mme_ue_s1ap_id)
{
  ue_description_t *ue = NULL;

  hashtable_ts_apply_callback_on_elements(
    &state->enbs,
    s1ap_enb_find_ue_by_mme_ue_id_cb,
    &mme_ue_s1ap_id,
    (void **) &ue);

  return ue;
}

s1ap_state_t *s1ap_state_new(void)
{
  s1ap_state_t *state;
  hash_table_ts_t *ht;
  bstring ht_name;

  state = (s1ap_state_t *) calloc(1, sizeof(*state));
  if (state == NULL) return NULL;

  ht_name = bfromcstr("s1ap_eNB_coll");
  ht = hashtable_ts_init(
    &state->enbs, mme_config.max_enbs, NULL, free_wrapper, ht_name);
  bdestroy(ht_name);

  if (ht == NULL) {
    free(state);
    return NULL;
  }

  ht_name = bfromcstr("s1ap_mme_id2assoc_id_coll");
  ht = hashtable_ts_init(
    &state->mmeid2associd,
    mme_config.max_ues,
    NULL,
    hash_free_int_func,
    ht_name);
  bdestroy(ht_name);

  if (ht == NULL) {
    hashtable_ts_destroy(&state->enbs);
    free(state);
    return NULL;
  }

  state->num_enbs = 0;

  return state;
}

void s1ap_state_free(s1ap_state_t *state)
{
  // TODO leaks memory hard. need to free enb_desc's, ue_desc's, and ue_coll's
  if (hashtable_ts_destroy(&state->enbs) != HASH_TABLE_OK) {
    OAI_FPRINTF_ERR("An error occured while destroying s1 eNB hash table");
  }
  if (hashtable_ts_destroy(&state->mmeid2associd) != HASH_TABLE_OK) {
    OAI_FPRINTF_ERR("An error occured while destroying assoc_id hash table");
  }
  free(state);
}

s1ap_state_t *s1ap_state_from_redis(void)
{
  // todo implement me
  return s1ap_state_new();
}

void s1ap_state_to_redis(s1ap_state_t *state)
{
  // todo implement me
  return;
}

void state2proto(S1apState *proto, s1ap_state_t *state)
{
  int i;
  hashtable_rc_t ht_rc;
  hashtable_key_array_t *keys;

  mme_ue_s1ap_id_t mmeid;
  sctp_assoc_id_t associd;
  enb_description_t *enb;

  EnbDescription enb_proto;

  proto->Clear();

  // copy over enbs
  auto enbs = proto->mutable_enbs();
  keys = hashtable_ts_get_keys(&state->enbs);
  AssertFatal(keys != NULL, "failed to get keys for enbs");
  for (i = 0; i < keys->num_keys; i++) {
    associd = (sctp_assoc_id_t) keys->keys[i];
    ht_rc =
      hashtable_ts_get(&state->enbs, (hash_key_t) associd, (void **) &enb);
    AssertFatal(ht_rc == HASH_TABLE_OK, "associd not in enbs");

    enb2proto(&enb_proto, enb);
    (*enbs)[associd] = enb_proto;
  }
  free(keys->keys);
  free(keys);

  // copy over mmeid2associd
  auto mmeid2associd = proto->mutable_mmeid2associd();
  keys = hashtable_ts_get_keys(&state->mmeid2associd);
  AssertFatal(keys != NULL, "failed to get keys for mmeid2associd");
  for (i = 0; i < keys->num_keys; i++) {
    mmeid = (mme_ue_s1ap_id_t) keys->keys[i];
    ht_rc = hashtable_ts_get(
      &state->mmeid2associd, (hash_key_t) mmeid, (void **) &associd);
    AssertFatal(ht_rc == HASH_TABLE_OK, "mmeid not in mmeid2associd");

    (*mmeid2associd)[mmeid] = associd;
  }
  free(keys->keys);
  free(keys);

  proto->set_num_enbs(state->num_enbs);
}

// expects hashtables in state to be created already
void proto2state(s1ap_state_t *state, S1apState *proto)
{
  hashtable_rc_t ht_rc;
  enb_description_t *enb;

  auto enbs = proto->enbs();
  for (auto const &kv : enbs) {
    sctp_assoc_id_t associd = kv.first;
    EnbDescription enb_proto = kv.second;

    enb = (enb_description_t *) malloc(sizeof(*enb));
    AssertFatal(enb != NULL, "failed to alloc new enb_desc");

    proto2enb(enb, &enb_proto);
    ht_rc = hashtable_ts_insert(&state->enbs, (hash_key_t) associd, enb);
    AssertFatal(ht_rc == HASH_TABLE_OK, "failed to insert enb");
  }

  auto mmeid2associd = proto->mmeid2associd();
  for (auto const &kv : mmeid2associd) {
    mme_ue_s1ap_id_t mmeid = (mme_ue_s1ap_id_t) kv.first;
    sctp_assoc_id_t associd = (sctp_assoc_id_t) kv.second;

    ht_rc = hashtable_ts_insert(
      &state->mmeid2associd, (hash_key_t) mmeid, (void *) (uintptr_t) associd);
    AssertFatal(ht_rc == HASH_TABLE_OK, "failed to insert associd");
  }

  state->num_enbs = proto->num_enbs();
}

void enb2proto(EnbDescription *proto, enb_description_t *enb)
{
  int i;
  hashtable_rc_t ht_rc;
  hashtable_key_array_t *keys;

  enb_ue_s1ap_id_t enbueid;
  ue_description_t *ue;

  UeDescription ue_proto;

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

  // store ues
  auto ues = proto->mutable_ues();
  keys = hashtable_ts_get_keys(&enb->ue_coll);
  AssertFatal(keys != NULL, "failed to get keys for ue_coll");
  for (i = 0; i < keys->num_keys; i++) {
    enbueid = (mme_ue_s1ap_id_t) keys->keys[i];
    ht_rc =
      hashtable_ts_get(&enb->ue_coll, (hash_key_t) enbueid, (void **) &ue);
    AssertFatal(ht_rc == HASH_TABLE_OK, "enbueid not in ue_coll");
    AssertFatal(ue->enb == enb, "tried to commit ue assigned to wrong enb");

    ue2proto(&ue_proto, ue);
    (*ues)[enbueid] = ue_proto;
  }
  free(keys->keys);
  free(keys);
}

void proto2enb(enb_description_t *enb, EnbDescription *proto)
{
  hashtable_rc_t ht_rc;
  ue_description_t *ue;

  memset(enb, 0, sizeof(*enb));

  enb->enb_id = proto->enb_id();
  enb->s1_state = (mme_s1_enb_state_s) proto->s1_state();
  strncpy(enb->enb_name, proto->enb_name().c_str(), sizeof(enb->enb_name));
  enb->default_paging_drx = proto->default_paging_drx();
  enb->nb_ue_associated = proto->nb_ue_associated();
  enb->s1ap_enb_assoc_clean_up_timer.id =
    proto->s1ap_enb_assoc_clean_up_timer().id();
  enb->s1ap_enb_assoc_clean_up_timer.sec =
    proto->s1ap_enb_assoc_clean_up_timer().sec();
  enb->sctp_assoc_id = proto->sctp_assoc_id();
  enb->next_sctp_stream = proto->next_sctp_stream();
  enb->instreams = proto->instreams();
  enb->outstreams = proto->outstreams();

  // load ues
  auto ht_name = bfromcstr("s1ap_ue_coll");
  auto ht = hashtable_ts_init(
    &enb->ue_coll, mme_config.max_ues, NULL, free_wrapper, ht_name);
  bdestroy(ht_name);
  AssertFatal(ht != NULL, "failed to init ue_coll");

  auto ues = proto->ues();
  for (auto const &kv : ues) {
    enb_ue_s1ap_id_t enbueid = kv.first;
    UeDescription ue_proto = kv.second;

    ue = (ue_description_t *) malloc(sizeof(*ue));
    AssertFatal(ue != NULL, "failed to alloc new ue description");

    proto2ue(ue, &ue_proto);
    ue->enb = enb; // ue's are linked to parent enb

    ht_rc = hashtable_ts_insert(&enb->ue_coll, (hash_key_t) enbueid, ue);
    AssertFatal(ht_rc == HASH_TABLE_OK, "failed to insert ue");
  }
}

void ue2proto(UeDescription *proto, ue_description_t *ue)
{
  proto->Clear();

  proto->set_s1_ue_state(ue->s1_ue_state);
  proto->set_enb_ue_s1ap_id(ue->enb_ue_s1ap_id);
  proto->set_mme_ue_s1ap_id(ue->mme_ue_s1ap_id);
  proto->set_sctp_stream_recv(ue->sctp_stream_recv);
  proto->set_sctp_stream_send(ue->sctp_stream_send);
  proto->set_s11_sgw_teid(ue->s11_sgw_teid);
  proto->mutable_s1ap_ue_context_rel_timer()->set_id(
    ue->s1ap_ue_context_rel_timer.id);
  proto->mutable_s1ap_ue_context_rel_timer()->set_sec(
    ue->s1ap_ue_context_rel_timer.sec);
}

void proto2ue(ue_description_t *ue, UeDescription *proto)
{
  memset(ue, 0, sizeof(*ue));

  ue->s1_ue_state = (s1_ue_state_s) proto->s1_ue_state();
  ue->enb_ue_s1ap_id = proto->enb_ue_s1ap_id();
  ue->mme_ue_s1ap_id = proto->mme_ue_s1ap_id();
  ue->sctp_stream_recv = proto->sctp_stream_recv();
  ue->sctp_stream_send = proto->sctp_stream_send();
  ue->s11_sgw_teid = proto->s11_sgw_teid();
  ue->s1ap_ue_context_rel_timer.id = proto->s1ap_ue_context_rel_timer().id();
  ue->s1ap_ue_context_rel_timer.sec = proto->s1ap_ue_context_rel_timer().sec();
}

bool s1ap_ue_compare_by_mme_ue_id_cb(
  __attribute__((unused)) const hash_key_t keyP,
  void *const elementP,
  void *parameterP,
  void **resultP)
{
  mme_ue_s1ap_id_t *mme_ue_s1ap_id_p = (mme_ue_s1ap_id_t *) parameterP;
  ue_description_t *ue_ref = (ue_description_t *) elementP;
  if (*mme_ue_s1ap_id_p == ue_ref->mme_ue_s1ap_id) {
    *resultP = elementP;
    OAILOG_TRACE(
      LOG_S1AP,
      "Found ue_ref %p mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n",
      ue_ref,
      ue_ref->mme_ue_s1ap_id);
    return true;
  }
  return false;
}

bool s1ap_enb_find_ue_by_mme_ue_id_cb(
  __attribute__((unused)) const hash_key_t keyP,
  void *const elementP,
  void *parameterP,
  void **resultP)
{
  enb_description_t *enb_ref = (enb_description_t *) elementP;

  hashtable_ts_apply_callback_on_elements(
    (hash_table_ts_t *const) & enb_ref->ue_coll,
    s1ap_ue_compare_by_mme_ue_id_cb,
    parameterP,
    resultP);
  if (*resultP) {
    OAILOG_TRACE(
      LOG_S1AP,
      "Found ue_ref %p mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n",
      *resultP,
      ((ue_description_t *) (*resultP))->mme_ue_s1ap_id);
    return true;
  }
  return false;
}
