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
 *-----------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

extern "C" {
#include "assertions.h"
#include "dynamic_memory_check.h"
#include "log.h"
#include "timer.h"
}

#include "mme_app_state_converter.h"

namespace magma {
namespace lte {

MmeNasStateConverter::MmeNasStateConverter() = default;
MmeNasStateConverter::~MmeNasStateConverter() = default;
/**********************************************************
*                 Hashtable <-> Proto                     *
* Functions to serialize/deserialize in-memory hashtables *
* for MME task. Only MME task inserts/removes elements in *
* the hashtables, so these calls are thread-safe.         *
* We only need to lock the UE context structure as it can *
* also be accessed by the NAS task. If hashtable is empty *
* the proto field is also empty                           *
***********************************************************/

void MmeNasStateConverter::hashtable_ts_to_proto(
  hash_table_ts_t* htbl,
  google::protobuf::Map<unsigned long, UeContext>* proto_map)
{
  hashtable_key_array_t* keys = hashtable_ts_get_keys(htbl);
  if (keys == nullptr) {
    return;
  }

  for (auto i = 0; i < keys->num_keys; i++) {
    UeContext ue_ctxt_proto;
    ue_mm_context_t* ue_context_p = NULL;
    hashtable_rc_t ht_rc =
      hashtable_ts_get(htbl, keys->keys[i], (void**) &ue_context_p);
    if (ht_rc == HASH_TABLE_OK) {
      ue_context_to_proto(ue_context_p, &ue_ctxt_proto);
      (*proto_map)[(uint32_t) keys->keys[i]] = ue_ctxt_proto;
    } else {
      OAILOG_ERROR(
        LOG_MME_APP,
        "Key %u not in mme_ue_s1ap_id_ue_context_htbl",
        keys->keys[i]);
    }
  }
}

void MmeNasStateConverter::proto_to_hashtable_ts(
  const google::protobuf::Map<unsigned long, UeContext>& proto_map,
  hash_table_ts_t* state_htbl)
{
  mme_ue_s1ap_id_t mme_ue_id;

  ue_mm_context_t* ue_context_p = nullptr;

  for (auto const& kv : proto_map) {
    mme_ue_id = (mme_ue_s1ap_id_t) kv.first;
    proto_to_ue_mm_context(&kv.second, ue_context_p);
    hashtable_rc_t ht_rc = hashtable_ts_insert(
      state_htbl, (const hash_key_t) mme_ue_id, (void*) ue_context_p);
    if (ht_rc != HASH_TABLE_OK) {
      OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to insert ue_context for mme_ue_s1ap_id %u in table %s: error:"
        "%s\n",
        mme_ue_id,
        state_htbl->name->data,
        hashtable_rc_code2string(ht_rc));
    }
  }
}

void MmeNasStateConverter::hashtable_uint64_ts_to_proto(
  hash_table_uint64_ts_t* htbl,
  google::protobuf::Map<unsigned long, unsigned long>* proto_map,
  const std::string& table_name)
{
  hashtable_key_array_t* keys = hashtable_uint64_ts_get_keys(htbl);
  if (keys == nullptr) {
    return;
  }

  for (auto i = 0; i < keys->num_keys; i++) {
    uint64_t mme_ue_id;
    hashtable_rc_t ht_rc =
      hashtable_uint64_ts_get(htbl, keys->keys[i], &mme_ue_id);
    if (ht_rc == HASH_TABLE_OK) {
      (*proto_map)[keys->keys[i]] = mme_ue_id;
    } else {
      OAILOG_ERROR(LOG_MME_APP, "Key %lu not in %s", keys->keys[i], table_name);
    }
  }

  FREE_HASHTABLE_KEY_ARRAY(keys);
}

void MmeNasStateConverter::proto_to_hashtable_uint64_ts(
  const google::protobuf::Map<unsigned long, unsigned long>& proto_map,
  hash_table_uint64_ts_t* state_htbl,
  const std::string& table_name)
{
  for (auto const& kv : proto_map) {
    uint64_t id = kv.first;
    mme_ue_s1ap_id_t mme_ue_id = kv.second;

    hashtable_rc_t ht_rc =
      hashtable_uint64_ts_insert(state_htbl, (const hash_key_t) id, mme_ue_id);
    if (ht_rc != HASH_TABLE_OK) {
      OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to insert mme_ue_s1ap_id %u in table %s: error: %s\n",
        mme_ue_id,
        table_name,
        hashtable_rc_code2string(ht_rc));
    }
  }
}

void MmeNasStateConverter::guti_table_to_proto(
  const obj_hash_table_uint64_t* guti_htbl,
  google::protobuf::Map<std::string, unsigned long>* proto_map)
{
  void*** key_array_p = (void***) calloc(1, sizeof(void**));
  unsigned int size = 0;

  hashtable_rc_t ht_rc =
    obj_hashtable_uint64_ts_get_keys(guti_htbl, key_array_p, &size);
  if ((!*key_array_p) || (ht_rc != HASH_TABLE_OK)) {
    return;
  }
  for (auto i = 0; i < size; i++) {
    uint64_t mme_ue_id;
    Guti guti_proto;
    guti_to_proto(*(guti_t*) (*key_array_p)[i], &guti_proto);
    const std::string& guti_str = guti_proto.SerializeAsString();
    OAILOG_TRACE(
      LOG_MME_APP,
      "Looking for key %p with value %u\n",
      (*key_array_p)[i],
      *((*key_list)[i]));
    hashtable_rc_t ht_rc = obj_hashtable_uint64_ts_get(
      guti_htbl, (const void*) (*key_array_p)[i], sizeof(guti_t), &mme_ue_id);
    if (ht_rc == HASH_TABLE_OK) {
      (*proto_map)[guti_str] = mme_ue_id;
    } else {
      OAILOG_ERROR(LOG_MME_APP, "Key %s not in guti_ue_context_htbl", guti_str);
    }
  }
  FREE_OBJ_HASHTABLE_KEY_ARRAY(key_array_p);
}

void MmeNasStateConverter::proto_to_guti_table(
  const google::protobuf::Map<std::string, unsigned long>& proto_map,
  obj_hash_table_uint64_t* guti_htbl)
{
  for (auto const& kv : proto_map) {
    const std::string& guti_str = kv.first;
    mme_ue_s1ap_id_t mme_ue_id = kv.second;
    guti_t* guti_p = nullptr; //TODO string_to_guti(guti_str);

    hashtable_rc_t ht_rc = obj_hashtable_uint64_ts_insert(
      guti_htbl, guti_p, sizeof(*guti_p), mme_ue_id);
    if (ht_rc != HASH_TABLE_OK) {
      OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to insert mme_ue_s1ap_id %u in GUTI table, error: %s\n",
        mme_ue_id,
        hashtable_rc_code2string(ht_rc));
    }
  }
}

/*********************************************************
*                UE Context <-> Proto                    *
* Functions to serialize/desearialize UE context         *
* The caller needs to acquire a lock on UE context       *
**********************************************************/

void MmeNasStateConverter::mme_app_timer_to_proto(
  mme_app_timer_t* state_mme_timer,
  Timer* timer_proto)
{
  timer_proto->set_id(state_mme_timer->id);
  timer_proto->set_sec(state_mme_timer->sec);
}

void MmeNasStateConverter::proto_to_mme_app_timer(
  const Timer& timer_proto,
  mme_app_timer_t* state_mme_app_timer)
{
  state_mme_app_timer->id = timer_proto.id();
  state_mme_app_timer->sec = timer_proto.sec();
}

void MmeNasStateConverter::sgs_context_to_proto(
  sgs_context_t* state_sgs_context,
  SgsContext* sgs_context_proto)
{
  // TODO
}

void MmeNasStateConverter::proto_to_sgs_context(
  const SgsContext& sgs_context_proto,
  sgs_context_t* state_sgs_context)
{
  // TODO
}

void MmeNasStateConverter::ue_context_to_proto(
  ue_mm_context_t* ue_ctxt,
  UeContext* ue_ctxt_proto)
{
  ue_ctxt_proto->Clear();

  ue_ctxt_proto->set_imsi(ue_ctxt->emm_context._imsi64);
  ue_ctxt_proto->set_imsi_len(ue_ctxt->emm_context._imsi.length);

  char* msisdn_buffer = bstr2cstr(ue_ctxt->msisdn, (char) '?');
  if (msisdn_buffer) {
    ue_ctxt_proto->set_msisdn(msisdn_buffer);
    bcstrfree(msisdn_buffer);
  }

  ue_ctxt_proto->set_imsi_auth(ue_ctxt->imsi_auth);
  ue_ctxt_proto->set_rel_cause(ue_ctxt->ue_context_rel_cause);
  ue_ctxt_proto->set_mm_state(ue_ctxt->mm_state);
  ue_ctxt_proto->set_ecm_state(ue_ctxt->ecm_state);

  //TODO EmmContext* emm_ctx = ue_ctxt_proto->mutable_emm_context();
  //emm_context_to_proto(&ue_ctxt->emm_context, emm_ctx);
  ue_ctxt_proto->set_sctp_assoc_id_key(ue_ctxt->sctp_assoc_id_key);
  ue_ctxt_proto->set_enb_ue_s1ap_id(ue_ctxt->enb_ue_s1ap_id);
  ue_ctxt_proto->set_mme_ue_s1ap_id(ue_ctxt->mme_ue_s1ap_id);

  Tai* tai = ue_ctxt_proto->mutable_serving_cell_tai();
  char tai_digits[5];
  tai_digits[0] = ue_ctxt->serving_cell_tai.mcc_digit2;
  tai_digits[1] = ue_ctxt->serving_cell_tai.mcc_digit1;
  tai_digits[2] = ue_ctxt->serving_cell_tai.mnc_digit3;
  tai_digits[3] = ue_ctxt->serving_cell_tai.mcc_digit3;
  tai_digits[4] = ue_ctxt->serving_cell_tai.mnc_digit2;
  tai_digits[5] = ue_ctxt->serving_cell_tai.mnc_digit2;
  tai->set_mcc_mnc(tai_digits);
  tai->set_tac(ue_ctxt->serving_cell_tai.tac);

  Ecgi* ecgi = ue_ctxt_proto->mutable_e_utran_cgi();

  ue_ctxt_proto->set_cell_age((long int) ue_ctxt->cell_age);

  ApnConfigProfile* apn_cfg = ue_ctxt_proto->mutable_apn_config();
}

void MmeNasStateConverter::proto_to_ue_mm_context(
  const UeContext* ue_context_proto,
  ue_mm_context_t* state_ue_mm_context)
{
  state_ue_mm_context->emm_context._imsi64 = ue_context_proto->imsi();
  state_ue_mm_context->emm_context._imsi.length = ue_context_proto->imsi_len();
  state_ue_mm_context->msisdn = bfromcstr(ue_context_proto->msisdn().c_str());
  state_ue_mm_context->imsi_auth = ue_context_proto->imsi_auth();
  state_ue_mm_context->ue_context_rel_cause =
    static_cast<enum s1cause>(ue_context_proto->rel_cause());
  state_ue_mm_context->mm_state =
    static_cast<mm_state_t>(ue_context_proto->mm_state());
  state_ue_mm_context->ecm_state =
    static_cast<ecm_state_t>(ue_context_proto->ecm_state());
  // TODO: proto_to_emm_context(ue_context_proto->emm_context(),
  // state_ue_mm_context->emm_context);

  state_ue_mm_context->sctp_assoc_id_key =
    ue_context_proto->sctp_assoc_id_key();
  state_ue_mm_context->enb_ue_s1ap_id = ue_context_proto->enb_ue_s1ap_id();
  state_ue_mm_context->enb_s1ap_id_key = ue_context_proto->enb_s1ap_id_key();
  state_ue_mm_context->mme_ue_s1ap_id = ue_context_proto->mme_ue_s1ap_id();

  // TODO: all functions to be added in Nas state converter
  //proto_to_tai(ue_context_proto->serving_cell_tai(),
  //  state_ue_mm_context->serving_cell_tai);
  //proto_to_tai_list(ue_context_proto->tai_list(),
  //  state_ue_mm_context->tai_list);
  //proto_to_tai(ue_context_proto->tai_last_tau(),
  //  state_ue_mm_context->tai_last_tau);
  //proto_to_ecgi(ue_context_proto->e_utran_cgi(),
  //  state_ue_mm_context->e_utran_cgi);
  //proto_to_apn_config_profile(ue_context_proto->apn_config(),
  //  &state_ue_mm_context->apn_config);

  state_ue_mm_context->cell_age = ue_context_proto->cell_age();
  proto_to_sgs_context(
    ue_context_proto->sgs_context(), state_ue_mm_context->sgs_context);
  proto_to_mme_app_timer(
    ue_context_proto->mobile_reachability_timer(),
    &state_ue_mm_context->mobile_reachability_timer);
  proto_to_mme_app_timer(
    ue_context_proto->implicit_detach_timer(),
    &state_ue_mm_context->implicit_detach_timer);
  proto_to_mme_app_timer(
    ue_context_proto->initial_context_setup_rsp_timer(),
    &state_ue_mm_context->initial_context_setup_rsp_timer);
  proto_to_mme_app_timer(
    ue_context_proto->ue_context_modification_timer(),
    &state_ue_mm_context->ue_context_modification_timer);
  proto_to_mme_app_timer(
    ue_context_proto->paging_response_timer(),
    &state_ue_mm_context->paging_response_timer);

  return;
}

/*********************************************************
*                MME app state<-> Proto                  *
* Functions to serialize/desearialize MME app state      *
* The caller is responsible for all memory management    *
**********************************************************/
void MmeNasStateConverter::mme_nas_state_to_proto(
  mme_app_desc_t* mme_nas_state_p,
  MmeNasState* state_proto)
{
  state_proto->set_nb_enb_connected(mme_nas_state_p->nb_enb_connected);
  state_proto->set_nb_ue_attached(mme_nas_state_p->nb_ue_attached);
  state_proto->set_nb_ue_connected(mme_nas_state_p->nb_ue_connected);
  state_proto->set_nb_default_eps_bearers(
    mme_nas_state_p->nb_default_eps_bearers);
  state_proto->set_nb_s1u_bearers(mme_nas_state_p->nb_s1u_bearers);

  // copy mme_ue_contexts
  auto mme_ue_ctxts_proto = state_proto->mutable_mme_ue_contexts();
  mme_ue_ctxts_proto->set_nb_ue_managed(
    mme_nas_state_p->mme_ue_contexts.nb_ue_managed);
  mme_ue_ctxts_proto->set_nb_ue_idle(
    mme_nas_state_p->mme_ue_contexts.nb_ue_idle);
  mme_ue_ctxts_proto->set_nb_bearers_managed(
    mme_nas_state_p->mme_ue_contexts.nb_bearers_managed);
  mme_ue_ctxts_proto->set_nb_ue_since_last_stat(
    mme_nas_state_p->mme_ue_contexts.nb_ue_since_last_stat);
  mme_ue_ctxts_proto->set_nb_bearers_since_last_stat(
    mme_nas_state_p->mme_ue_contexts.nb_bearers_since_last_stat);

  hashtable_uint64_ts_to_proto(
    mme_nas_state_p->mme_ue_contexts.imsi_ue_context_htbl,
    mme_ue_ctxts_proto->mutable_imsi_ue_id_htbl(),
    "imsi_ue_context_htbl");
  hashtable_uint64_ts_to_proto(
    mme_nas_state_p->mme_ue_contexts.tun11_ue_context_htbl,
    mme_ue_ctxts_proto->mutable_tun11_ue_id_htbl(),
    "tun11_ue_context_htbl");
  hashtable_ts_to_proto(
    mme_nas_state_p->mme_ue_contexts.mme_ue_s1ap_id_ue_context_htbl,
    mme_ue_ctxts_proto->mutable_mme_ue_id_ue_ctxt_htbl());
  hashtable_uint64_ts_to_proto(
    mme_nas_state_p->mme_ue_contexts.enb_ue_s1ap_id_ue_context_htbl,
    mme_ue_ctxts_proto->mutable_enb_ue_id_ue_id_htbl(),
    "enb_ue_s1ap_id_ue_context_htbl");
  guti_table_to_proto(
    mme_nas_state_p->mme_ue_contexts.guti_ue_context_htbl,
    mme_ue_ctxts_proto->mutable_guti_ue_id_htbl());
  return;
}

void MmeNasStateConverter::mme_nas_proto_to_state(
  MmeNasState* state_proto,
  mme_app_desc_t* mme_nas_state_p)
{
  mme_nas_state_p->nb_enb_connected = state_proto->nb_enb_connected();
  mme_nas_state_p->nb_ue_attached = state_proto->nb_ue_attached();
  mme_nas_state_p->nb_ue_connected = state_proto->nb_ue_connected();
  mme_nas_state_p->nb_default_eps_bearers =
    state_proto->nb_default_eps_bearers();
  mme_nas_state_p->nb_s1u_bearers = state_proto->nb_s1u_bearers();

  // copy mme_ue_contexts
  MmeUeContext mme_ue_ctxts_proto = state_proto->mme_ue_contexts();
  mme_nas_state_p->mme_ue_contexts.nb_ue_managed =
    mme_ue_ctxts_proto.nb_ue_managed();
  mme_nas_state_p->mme_ue_contexts.nb_ue_idle = mme_ue_ctxts_proto.nb_ue_idle();
  mme_nas_state_p->mme_ue_contexts.nb_bearers_managed =
    mme_ue_ctxts_proto.nb_bearers_managed();
  mme_nas_state_p->mme_ue_contexts.nb_ue_since_last_stat =
    mme_ue_ctxts_proto.nb_ue_since_last_stat();
  mme_nas_state_p->mme_ue_contexts.nb_bearers_since_last_stat =
    mme_ue_ctxts_proto.nb_bearers_since_last_stat();

  mme_ue_context_t* mme_ue_ctxt_state = &mme_nas_state_p->mme_ue_contexts;
  // copy maps to hashtables
  proto_to_hashtable_ts(
    mme_ue_ctxts_proto.mme_ue_id_ue_ctxt_htbl(),
    mme_ue_ctxt_state->mme_ue_s1ap_id_ue_context_htbl);
  proto_to_hashtable_uint64_ts(
    mme_ue_ctxts_proto.imsi_ue_id_htbl(),
    mme_ue_ctxt_state->imsi_ue_context_htbl,
    "imsi_ue_context_htbl");
  proto_to_hashtable_uint64_ts(
    mme_ue_ctxts_proto.tun11_ue_id_htbl(),
    mme_ue_ctxt_state->imsi_ue_context_htbl,
    "tun11_ue_context_htbl");
  proto_to_hashtable_uint64_ts(
    mme_ue_ctxts_proto.enb_ue_id_ue_id_htbl(),
    mme_ue_ctxt_state->enb_ue_s1ap_id_ue_context_htbl,
    "enb_ue_s1ap_id_ue_context_htbl");
  proto_to_guti_table(
    mme_ue_ctxts_proto.guti_ue_id_htbl(),
    mme_ue_ctxt_state->guti_ue_context_htbl);
}
} // namespace lte
} // namespace magma
