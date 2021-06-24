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
 *-----------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

extern "C" {
#include "bytes_to_ie.h"
#include "dynamic_memory_check.h"
#include "ie_to_bytes.h"
#include "log.h"
#include "timer.h"
}

#include "mme_app_state_converter.h"
#include "nas_state_converter.h"

namespace magma {
namespace lte {

MmeNasStateConverter::MmeNasStateConverter()  = default;
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
    google::protobuf::Map<unsigned long, oai::UeContext>* proto_map) {
  hashtable_key_array_t* keys = hashtable_ts_get_keys(htbl);
  if (keys == nullptr) {
    return;
  }

  for (auto i = 0; i < keys->num_keys; i++) {
    oai::UeContext ue_ctxt_proto;
    ue_mm_context_t* ue_context_p = NULL;
    hashtable_rc_t ht_rc =
        hashtable_ts_get(htbl, keys->keys[i], (void**) &ue_context_p);
    if (ht_rc == HASH_TABLE_OK) {
      ue_context_to_proto(ue_context_p, &ue_ctxt_proto);
      (*proto_map)[(uint32_t) keys->keys[i]] = ue_ctxt_proto;
    } else {
      OAILOG_ERROR(
          LOG_MME_APP, "Key %lu not in mme_ue_s1ap_id_ue_context_htbl",
          keys->keys[i]);
    }
  }
  FREE_HASHTABLE_KEY_ARRAY(keys);
}

void MmeNasStateConverter::proto_to_hashtable_ts(
    const google::protobuf::Map<unsigned long, oai::UeContext>& proto_map,
    hash_table_ts_t* state_htbl) {
  OAILOG_DEBUG(LOG_MME_APP, "Converting proto to hashtable_ts");
  mme_ue_s1ap_id_t mme_ue_id;

  for (auto const& kv : proto_map) {
    mme_ue_id = (mme_ue_s1ap_id_t) kv.first;
    OAILOG_DEBUG(
        LOG_MME_APP, "Reading ue_context for " MME_UE_S1AP_ID_FMT, mme_ue_id);
    ue_mm_context_t* ue_context_p = mme_create_new_ue_context();
    if (!ue_context_p) {
      OAILOG_ERROR(LOG_MME_APP, "Could not allocate new UE context");
      continue;
    }
    proto_to_ue_mm_context(kv.second, ue_context_p);
    hashtable_rc_t ht_rc = hashtable_ts_insert(
        state_htbl, (const hash_key_t) mme_ue_id, (void*) ue_context_p);
    if (ht_rc != HASH_TABLE_OK) {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Failed to insert ue_context for mme_ue_s1ap_id %u in table %s: "
          "error:"
          "%s\n",
          mme_ue_id, state_htbl->name->data, hashtable_rc_code2string(ht_rc));
    }
    OAILOG_DEBUG(LOG_MME_APP, "Written one key into hashtable_ts");
  }
}

char* MmeNasStateConverter::mme_app_convert_guti_to_string(guti_t* guti_p) {
#define GUTI_STRING_LEN 21
  char* str = (char*) calloc(1, sizeof(char) * GUTI_STRING_LEN);
  snprintf(
      str, GUTI_STRING_LEN, "%x%x%x%x%x%x%04x%02x%08x",
      guti_p->gummei.plmn.mcc_digit1, guti_p->gummei.plmn.mcc_digit2,
      guti_p->gummei.plmn.mcc_digit3, guti_p->gummei.plmn.mnc_digit1,
      guti_p->gummei.plmn.mnc_digit2, guti_p->gummei.plmn.mnc_digit3,
      guti_p->gummei.mme_gid, guti_p->gummei.mme_code, guti_p->m_tmsi);
  return (str);
}

void MmeNasStateConverter::guti_table_to_proto(
    const obj_hash_table_uint64_t* guti_htbl,
    google::protobuf::Map<std::string, unsigned long>* proto_map) {
  void*** key_array_p = (void***) calloc(1, sizeof(void**));
  unsigned int size   = 0;

  hashtable_rc_t ht_rc =
      obj_hashtable_uint64_ts_get_keys(guti_htbl, key_array_p, &size);
  if ((!*key_array_p) || (ht_rc != HASH_TABLE_OK)) {
    FREE_OBJ_HASHTABLE_KEY_ARRAY(key_array_p);
    return;
  }
  for (unsigned int i = 0; i < size; i++) {
    uint64_t mme_ue_id;

    char* str = mme_app_convert_guti_to_string((guti_t*) (*key_array_p)[i]);
    std::string guti_str(str);
    free(str);
    OAILOG_TRACE(
        LOG_MME_APP, "Looking for key %p with value %u\n", (*key_array_p)[i],
        *((*key_list)[i]));
    hashtable_rc_t ht_rc = obj_hashtable_uint64_ts_get(
        guti_htbl, (const void*) (*key_array_p)[i], sizeof(guti_t), &mme_ue_id);
    if (ht_rc == HASH_TABLE_OK) {
      (*proto_map)[guti_str] = mme_ue_id;
    } else {
      OAILOG_ERROR(
          LOG_MME_APP, "Key %s not in guti_ue_context_htbl", guti_str.c_str());
    }
  }
  FREE_OBJ_HASHTABLE_KEY_ARRAY(key_array_p);
}

void MmeNasStateConverter::mme_app_convert_string_to_guti(
    guti_t* guti_p, const std::string& guti_str) {
  int idx                   = 0;
  std::size_t chars_to_read = 1;
#define HEX_BASE_VAL 16
  guti_p->gummei.plmn.mcc_digit1 = std::stoul(
      guti_str.substr(idx++, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  guti_p->gummei.plmn.mcc_digit2 = std::stoul(
      guti_str.substr(idx++, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  guti_p->gummei.plmn.mcc_digit3 = std::stoul(
      guti_str.substr(idx++, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  guti_p->gummei.plmn.mnc_digit1 = std::stoul(
      guti_str.substr(idx++, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  guti_p->gummei.plmn.mnc_digit2 = std::stoul(
      guti_str.substr(idx++, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  guti_p->gummei.plmn.mnc_digit3 = std::stoul(
      guti_str.substr(idx++, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  chars_to_read          = 4;
  guti_p->gummei.mme_gid = std::stoul(
      guti_str.substr(idx, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  idx += chars_to_read;
  chars_to_read           = 2;
  guti_p->gummei.mme_code = std::stoul(
      guti_str.substr(idx, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  idx += chars_to_read;
  chars_to_read = 8;
  guti_p->m_tmsi =
      std::stoul(guti_str.substr(idx, chars_to_read), 0, HEX_BASE_VAL);

  OAILOG_DEBUG_GUTI(guti_p);
}

void MmeNasStateConverter::proto_to_guti_table(
    const google::protobuf::Map<std::string, unsigned long>& proto_map,
    obj_hash_table_uint64_t* guti_htbl) {
  for (auto const& kv : proto_map) {
    mme_ue_s1ap_id_t mme_ue_id = kv.second;
    guti_t* guti_p             = (guti_t*) calloc(1, sizeof(guti_t));

    mme_app_convert_string_to_guti(guti_p, kv.first);
    hashtable_rc_t ht_rc = obj_hashtable_uint64_ts_insert(
        guti_htbl, guti_p, sizeof(*guti_p), mme_ue_id);
    if (ht_rc != HASH_TABLE_OK) {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Failed to insert mme_ue_s1ap_id %u in GUTI table, error: %s\n",
          mme_ue_id, hashtable_rc_code2string(ht_rc));
    }
    free_wrapper((void**) &guti_p);
  }
}

/*********************************************************
 *                UE Context <-> Proto                    *
 * Functions to serialize/desearialize UE context         *
 * The caller needs to acquire a lock on UE context       *
 **********************************************************/

void MmeNasStateConverter::mme_app_timer_to_proto(
    const mme_app_timer_t& state_mme_timer, oai::Timer* timer_proto) {
  timer_proto->set_id(state_mme_timer.id);
  timer_proto->set_sec(state_mme_timer.sec);
}

void MmeNasStateConverter::proto_to_mme_app_timer(
    const oai::Timer& timer_proto, mme_app_timer_t* state_mme_app_timer) {
  state_mme_app_timer->id  = timer_proto.id();
  state_mme_app_timer->sec = timer_proto.sec();
}

void MmeNasStateConverter::sgs_context_to_proto(
    sgs_context_t* state_sgs_context, oai::SgsContext* sgs_context_proto) {
  // TODO
}

void MmeNasStateConverter::proto_to_sgs_context(
    const oai::SgsContext& sgs_context_proto,
    sgs_context_t* state_sgs_context) {
  // TODO
}

void MmeNasStateConverter::fteid_to_proto(
    const fteid_t& state_fteid, oai::Fteid* fteid_proto) {
  if (state_fteid.ipv4) {
    fteid_proto->set_ipv4_address(state_fteid.ipv4_address.s_addr);
  } else if (state_fteid.ipv6) {
    fteid_proto->set_ipv6_address(
        &state_fteid.ipv6_address, TRAFFIC_FLOW_TEMPLATE_IPV6_ADDR_SIZE);
  }
  fteid_proto->set_interface_type(state_fteid.interface_type);
  fteid_proto->set_teid(state_fteid.teid);
}

void MmeNasStateConverter::proto_to_fteid(
    const oai::Fteid& fteid_proto, fteid_t* state_fteid) {
  if (fteid_proto.ipv4_address()) {
    state_fteid->ipv4                = 1;
    state_fteid->ipv4_address.s_addr = fteid_proto.ipv4_address();
  } else if (fteid_proto.ipv6_address().length() > 0) {
    state_fteid->ipv6 = 1;
    memcpy(
        &state_fteid->ipv6_address, fteid_proto.ipv6_address().c_str(),
        TRAFFIC_FLOW_TEMPLATE_IPV6_ADDR_SIZE);
  }
  state_fteid->interface_type = (interface_type_t) fteid_proto.interface_type();
  state_fteid->teid           = fteid_proto.teid();
}

void MmeNasStateConverter::bearer_context_to_proto(
    const bearer_context_t& state_bearer_context,
    oai::BearerContext* bearer_context_proto) {
  bearer_context_proto->set_ebi(state_bearer_context.ebi);
  bearer_context_proto->set_transaction_identifier(
      state_bearer_context.transaction_identifier);
  fteid_to_proto(
      state_bearer_context.s_gw_fteid_s1u,
      bearer_context_proto->mutable_s_gw_fteid_s1u());
  fteid_to_proto(
      state_bearer_context.p_gw_fteid_s5_s8_up,
      bearer_context_proto->mutable_p_gw_fteid_s5_s8_up());
  bearer_context_proto->set_qci(state_bearer_context.qci);
  bearer_context_proto->set_pdn_cx_id(state_bearer_context.pdn_cx_id);
  NasStateConverter::esm_ebr_context_to_proto(
      state_bearer_context.esm_ebr_context,
      bearer_context_proto->mutable_esm_ebr_context());
  fteid_to_proto(
      state_bearer_context.enb_fteid_s1u,
      bearer_context_proto->mutable_enb_fteid_s1u());
  bearer_context_proto->set_priority_level(state_bearer_context.priority_level);
  bearer_context_proto->set_preemption_vulnerability(
      state_bearer_context.preemption_vulnerability);
  bearer_context_proto->set_preemption_capability(
      state_bearer_context.preemption_capability);

  /* TODO
   * SpgwStateConverter::traffic_flow_template_to_proto(
   * state_bearer_context.saved_tft,
   * bearer_context_proto->mutable_saved_tft());
   */

  /* TODO
   * SpgwStateConverter::eps_bearer_qos_to_proto(
   * state_bearer_context.saved_qos,
   * bearer_context_proto->mutable_saved_qos());
   */
}

void MmeNasStateConverter::proto_to_bearer_context(
    const oai::BearerContext& bearer_context_proto,
    bearer_context_t* state_bearer_context) {
  state_bearer_context->ebi = bearer_context_proto.ebi();
  state_bearer_context->transaction_identifier =
      bearer_context_proto.transaction_identifier();
  proto_to_fteid(
      bearer_context_proto.s_gw_fteid_s1u(),
      &state_bearer_context->s_gw_fteid_s1u);
  proto_to_fteid(
      bearer_context_proto.p_gw_fteid_s5_s8_up(),
      &state_bearer_context->p_gw_fteid_s5_s8_up);
  state_bearer_context->qci       = bearer_context_proto.qci();
  state_bearer_context->pdn_cx_id = bearer_context_proto.pdn_cx_id();
  NasStateConverter::proto_to_esm_ebr_context(
      bearer_context_proto.esm_ebr_context(),
      &state_bearer_context->esm_ebr_context);
  proto_to_fteid(
      bearer_context_proto.enb_fteid_s1u(),
      &state_bearer_context->enb_fteid_s1u);
  state_bearer_context->priority_level = bearer_context_proto.priority_level();
  state_bearer_context->preemption_vulnerability =
      (pre_emption_vulnerability_t)
          bearer_context_proto.preemption_vulnerability();
  state_bearer_context->preemption_capability =
      (pre_emption_capability_t) bearer_context_proto.preemption_capability();

  /* TODO
   * proto_to_gateway.spgw._traffic_flow_template(
   * bearer_context_proto.saved_tft(),
   * &state_bearer_context->saved_tft);
   */

  /* TODO proto_to_bearer_qos(bearer_context_proto.saved_qos(),
   * &state_bearer_context->saved_qos);
   */
}

void MmeNasStateConverter::bearer_context_list_to_proto(
    const ue_mm_context_t& state_ue_context, oai::UeContext* ue_context_proto) {
  for (int i = 0; i < BEARERS_PER_UE; i++) {
    oai::BearerContext* bearer_ctxt_proto =
        ue_context_proto->add_bearer_contexts();
    if (state_ue_context.bearer_contexts[i]) {
      OAILOG_DEBUG(
          LOG_MME_APP, "writing bearer context at index %d with ebi %d", i,
          state_ue_context.bearer_contexts[i]->ebi);
      bearer_ctxt_proto->set_validity(oai::BearerContext::VALID);
      bearer_context_to_proto(
          *state_ue_context.bearer_contexts[i], bearer_ctxt_proto);
    } else {
      bearer_ctxt_proto->set_validity(oai::BearerContext::INVALID);
    }
  }
}

void MmeNasStateConverter::proto_to_bearer_context_list(
    const oai::UeContext& ue_context_proto, ue_mm_context_t* state_ue_context) {
  for (int i = 0; i < BEARERS_PER_UE; i++) {
    if (ue_context_proto.bearer_contexts(i).validity() ==
        oai::BearerContext::VALID) {
      OAILOG_DEBUG(LOG_MME_APP, "reading bearer context at index %d", i);
      auto* eps_bearer_ctxt =
          (bearer_context_t*) calloc(1, sizeof(bearer_context_t));
      proto_to_bearer_context(
          ue_context_proto.bearer_contexts(i), eps_bearer_ctxt);
      state_ue_context->bearer_contexts[i] = eps_bearer_ctxt;
      if (state_ue_context->bearer_contexts[i]->esm_ebr_context.args) {
        state_ue_context->bearer_contexts[i]->esm_ebr_context.args->ctx =
            &state_ue_context->emm_context;
      }
    } else {
      state_ue_context->bearer_contexts[i] = nullptr;
    }
  }
}

void MmeNasStateConverter::esm_pdn_to_proto(
    const esm_pdn_t& state_esm_pdn, oai::EsmPdn* esm_pdn_proto) {
  esm_pdn_proto->set_pti(state_esm_pdn.pti);
  esm_pdn_proto->set_is_emergency(state_esm_pdn.is_emergency);
  esm_pdn_proto->set_ambr(state_esm_pdn.ambr);
  esm_pdn_proto->set_addr_realloc(state_esm_pdn.addr_realloc);
  esm_pdn_proto->set_n_bearers(state_esm_pdn.n_bearers);
  esm_pdn_proto->set_pt_state(state_esm_pdn.pt_state);
}

void MmeNasStateConverter::proto_to_esm_pdn(
    const oai::EsmPdn& esm_pdn_proto, esm_pdn_t* state_esm_pdn) {
  state_esm_pdn->pti          = esm_pdn_proto.pti();
  state_esm_pdn->is_emergency = esm_pdn_proto.is_emergency();
  state_esm_pdn->ambr         = esm_pdn_proto.ambr();
  state_esm_pdn->addr_realloc = esm_pdn_proto.addr_realloc();
  state_esm_pdn->n_bearers    = esm_pdn_proto.n_bearers();
  state_esm_pdn->pt_state     = (esm_pt_state_e) esm_pdn_proto.pt_state();
}

void MmeNasStateConverter::pdn_context_to_proto(
    const pdn_context_t& state_pdn_context,
    oai::PdnContext* pdn_context_proto) {
  pdn_context_proto->set_context_identifier(
      state_pdn_context.context_identifier);
  BSTRING_TO_STRING(
      state_pdn_context.apn_in_use, pdn_context_proto->mutable_apn_in_use());
  BSTRING_TO_STRING(
      state_pdn_context.apn_subscribed,
      pdn_context_proto->mutable_apn_subscribed());
  pdn_context_proto->set_pdn_type(state_pdn_context.pdn_type);
  bstring bstr_buffer = paa_to_bstring(&state_pdn_context.paa);
  BSTRING_TO_STRING(bstr_buffer, pdn_context_proto->mutable_paa());
  BSTRING_TO_STRING(
      state_pdn_context.apn_oi_replacement,
      pdn_context_proto->mutable_apn_oi_replacement());
  bdestroy(bstr_buffer);
  bstr_buffer = ip_address_to_bstring(&state_pdn_context.p_gw_address_s5_s8_cp);
  BSTRING_TO_STRING(
      bstr_buffer, pdn_context_proto->mutable_p_gw_address_s5_s8_cp());
  bdestroy(bstr_buffer);
  pdn_context_proto->set_p_gw_teid_s5_s8_cp(
      state_pdn_context.p_gw_teid_s5_s8_cp);
  eps_subscribed_qos_profile_to_proto(
      state_pdn_context.default_bearer_eps_subscribed_qos_profile,
      pdn_context_proto->mutable_default_bearer_qos_profile());
  StateConverter::ambr_to_proto(
      state_pdn_context.subscribed_apn_ambr,
      pdn_context_proto->mutable_subscribed_apn_ambr());
  StateConverter::ambr_to_proto(
      state_pdn_context.p_gw_apn_ambr,
      pdn_context_proto->mutable_p_gw_apn_ambr());
  pdn_context_proto->set_default_ebi(state_pdn_context.default_ebi);
  for (int i = 0; i < BEARERS_PER_UE; i++) {
    pdn_context_proto->add_bearer_contexts(
        state_pdn_context.bearer_contexts[i]);
  }

  bstr_buffer = ip_address_to_bstring(&state_pdn_context.s_gw_address_s11_s4);
  BSTRING_TO_STRING(
      bstr_buffer, pdn_context_proto->mutable_s_gw_address_s11_s4());
  bdestroy_wrapper(&bstr_buffer);
  pdn_context_proto->set_s_gw_teid_s11_s4(state_pdn_context.s_gw_teid_s11_s4);
  esm_pdn_to_proto(
      state_pdn_context.esm_data, pdn_context_proto->mutable_esm_data());
  pdn_context_proto->set_is_active(state_pdn_context.is_active);
  if (state_pdn_context.pco != nullptr) {
    NasStateConverter::protocol_configuration_options_to_proto(
        *state_pdn_context.pco, pdn_context_proto->mutable_pco());
  }
}

void MmeNasStateConverter::proto_to_pdn_context(
    const oai::PdnContext& pdn_context_proto,
    pdn_context_t* state_pdn_context) {
  state_pdn_context->context_identifier =
      pdn_context_proto.context_identifier();
  STRING_TO_BSTRING(
      pdn_context_proto.apn_in_use(), state_pdn_context->apn_in_use);
  STRING_TO_BSTRING(
      pdn_context_proto.apn_subscribed(), state_pdn_context->apn_subscribed);
  state_pdn_context->pdn_type = pdn_context_proto.pdn_type();
  bstring bstr_buffer;
  STRING_TO_BSTRING(pdn_context_proto.paa(), bstr_buffer);
  bstring_to_paa(bstr_buffer, &state_pdn_context->paa);
  bdestroy(bstr_buffer);
  STRING_TO_BSTRING(
      pdn_context_proto.apn_oi_replacement(),
      state_pdn_context->apn_oi_replacement);
  STRING_TO_BSTRING(pdn_context_proto.p_gw_address_s5_s8_cp(), bstr_buffer);
  bstring_to_ip_address(bstr_buffer, &state_pdn_context->p_gw_address_s5_s8_cp);
  bdestroy(bstr_buffer);
  state_pdn_context->p_gw_teid_s5_s8_cp =
      pdn_context_proto.p_gw_teid_s5_s8_cp();
  proto_to_eps_subscribed_qos_profile(
      pdn_context_proto.default_bearer_qos_profile(),
      &state_pdn_context->default_bearer_eps_subscribed_qos_profile);
  proto_to_ambr(
      pdn_context_proto.subscribed_apn_ambr(),
      &state_pdn_context->subscribed_apn_ambr);
  proto_to_ambr(
      pdn_context_proto.p_gw_apn_ambr(), &state_pdn_context->p_gw_apn_ambr);
  state_pdn_context->default_ebi = pdn_context_proto.default_ebi();
  for (int i = 0; i < BEARERS_PER_UE; i++) {
    state_pdn_context->bearer_contexts[i] =
        pdn_context_proto.bearer_contexts(i);
  }
  STRING_TO_BSTRING(pdn_context_proto.s_gw_address_s11_s4(), bstr_buffer);
  bstring_to_ip_address(bstr_buffer, &state_pdn_context->s_gw_address_s11_s4);
  bdestroy_wrapper(&bstr_buffer);
  state_pdn_context->s_gw_teid_s11_s4 = pdn_context_proto.s_gw_teid_s11_s4();
  proto_to_esm_pdn(pdn_context_proto.esm_data(), &state_pdn_context->esm_data);
  state_pdn_context->is_active = pdn_context_proto.is_active();
  if (pdn_context_proto.has_pco()) {
    state_pdn_context->pco = (protocol_configuration_options_t*) calloc(
        1, sizeof(protocol_configuration_options_t));
    NasStateConverter::proto_to_protocol_configuration_options(
        pdn_context_proto.pco(), state_pdn_context->pco);
  }
}

void MmeNasStateConverter::pdn_context_list_to_proto(
    const ue_mm_context_t& state_ue_context, oai::UeContext* ue_context_proto) {
  for (int i = 0; i < MAX_APN_PER_UE; i++) {
    if (state_ue_context.pdn_contexts[i] != nullptr) {
      OAILOG_DEBUG(LOG_MME_APP, "Writing PDN context at index %d", i);
      oai::PdnContext* pdn_ctxt_proto = ue_context_proto->add_pdn_contexts();
      pdn_context_to_proto(*state_ue_context.pdn_contexts[i], pdn_ctxt_proto);
    }
  }
}

void MmeNasStateConverter::proto_to_pdn_context_list(
    const oai::UeContext& ue_context_proto, ue_mm_context_t* state_ue_context) {
  for (int i = 0; i < ue_context_proto.pdn_contexts_size(); i++) {
    OAILOG_DEBUG(LOG_MME_APP, "Reading PDN context at index %d", i);
    auto* pdn_context_p = (pdn_context_t*) calloc(1, sizeof(pdn_context_t));
    proto_to_pdn_context(ue_context_proto.pdn_contexts(i), pdn_context_p);
    state_ue_context->pdn_contexts[i] = pdn_context_p;
  }
}

void MmeNasStateConverter::regional_subscription_to_proto(
    const ue_mm_context_t& state_ue_context, oai::UeContext* ue_context_proto) {
  for (int itr = 0; itr < state_ue_context.num_reg_sub; itr++) {
    oai::Regional_subscription* reg_sub_proto = ue_context_proto->add_reg_sub();
    reg_sub_proto->set_zone_code(
        (const char*) state_ue_context.reg_sub[itr].zone_code);
    OAILOG_DEBUG(LOG_MME_APP, "Writing regional_subscription at index %d", itr);
  }
}

void MmeNasStateConverter::proto_to_regional_subscription(
    const oai::UeContext& ue_context_proto, ue_mm_context_t* state_ue_context) {
  for (int itr = 0; itr < ue_context_proto.num_reg_sub(); itr++) {
    memcpy(
        state_ue_context->reg_sub[itr].zone_code,
        ue_context_proto.reg_sub(itr).zone_code().c_str(),
        ue_context_proto.reg_sub(itr).zone_code().length());
    OAILOG_DEBUG(LOG_MME_APP, "Reading regional_subscription at index %d", itr);
  }
}

void MmeNasStateConverter::ue_context_to_proto(
    const ue_mm_context_t* state_ue_context, oai::UeContext* ue_context_proto) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  ue_context_proto->Clear();

  char* msisdn_buffer = bstr2cstr(state_ue_context->msisdn, (char) '?');
  if (msisdn_buffer) {
    ue_context_proto->set_msisdn(msisdn_buffer);
    bcstrfree(msisdn_buffer);
  } else {
    ue_context_proto->set_msisdn("");
  }

  ue_context_proto->set_rel_cause(state_ue_context->ue_context_rel_cause);
  ue_context_proto->set_mm_state(state_ue_context->mm_state);
  ue_context_proto->set_ecm_state(state_ue_context->ecm_state);

  oai::EmmContext* emm_ctx = ue_context_proto->mutable_emm_context();
  NasStateConverter::emm_context_to_proto(
      &state_ue_context->emm_context, emm_ctx);
  ue_context_proto->set_sctp_assoc_id_key(state_ue_context->sctp_assoc_id_key);
  ue_context_proto->set_enb_ue_s1ap_id(state_ue_context->enb_ue_s1ap_id);
  ue_context_proto->set_enb_s1ap_id_key(state_ue_context->enb_s1ap_id_key);
  ue_context_proto->set_mme_ue_s1ap_id(state_ue_context->mme_ue_s1ap_id);

  ue_context_proto->set_attach_type(state_ue_context->attach_type);
  ue_context_proto->set_sgs_detach_type(state_ue_context->sgs_detach_type);
  StateConverter::ecgi_to_proto(
      state_ue_context->e_utran_cgi, ue_context_proto->mutable_e_utran_cgi());

  ue_context_proto->set_cell_age((long int) state_ue_context->cell_age);

  char lai_bytes[IE_LENGTH_LAI];
  lai_to_bytes(&state_ue_context->lai, lai_bytes);
  ue_context_proto->set_lai(lai_bytes, IE_LENGTH_LAI);
  StateConverter::apn_config_profile_to_proto(
      state_ue_context->apn_config_profile,
      ue_context_proto->mutable_apn_config());
  ue_context_proto->set_subscriber_status(state_ue_context->subscriber_status);
  ue_context_proto->set_network_access_mode(
      state_ue_context->network_access_mode);
  ue_context_proto->set_access_restriction_data(
      state_ue_context->access_restriction_data);
  BSTRING_TO_STRING(
      state_ue_context->apn_oi_replacement,
      ue_context_proto->mutable_apn_oi_replacement());
  ue_context_proto->set_mme_teid_s11(state_ue_context->mme_teid_s11);
  StateConverter::ambr_to_proto(
      state_ue_context->subscribed_ue_ambr,
      ue_context_proto->mutable_subscribed_ue_ambr());
  StateConverter::ambr_to_proto(
      state_ue_context->used_ue_ambr, ue_context_proto->mutable_used_ue_ambr());
  StateConverter::ambr_to_proto(
      state_ue_context->used_ambr, ue_context_proto->mutable_used_ambr());
  ue_context_proto->set_nb_active_pdn_contexts(
      state_ue_context->nb_active_pdn_contexts);
  pdn_context_list_to_proto(*state_ue_context, ue_context_proto);
  bearer_context_list_to_proto(*state_ue_context, ue_context_proto);
  if (state_ue_context->ue_radio_capability) {
    BSTRING_TO_STRING(
        state_ue_context->ue_radio_capability,
        ue_context_proto->mutable_ue_radio_capability());
  }
  ue_context_proto->set_send_ue_purge_request(
      state_ue_context->send_ue_purge_request);
  ue_context_proto->set_hss_initiated_detach(
      state_ue_context->hss_initiated_detach);
  ue_context_proto->set_location_info_confirmed_in_hss(
      state_ue_context->location_info_confirmed_in_hss);
  ue_context_proto->set_ppf(state_ue_context->ppf);
  ue_context_proto->set_subscription_known(
      state_ue_context->subscription_known);
  ue_context_proto->set_path_switch_req(state_ue_context->path_switch_req);
  ue_context_proto->set_granted_service(state_ue_context->granted_service);

  ue_context_proto->set_num_reg_sub(state_ue_context->num_reg_sub);
  regional_subscription_to_proto(*state_ue_context, ue_context_proto);
  ue_context_proto->set_cs_fallback_indicator(
      state_ue_context->cs_fallback_indicator);
  sgs_context_to_proto(
      state_ue_context->sgs_context, ue_context_proto->mutable_sgs_context());
  mme_app_timer_to_proto(
      state_ue_context->mobile_reachability_timer,
      ue_context_proto->mutable_mobile_reachability_timer());
  mme_app_timer_to_proto(
      state_ue_context->implicit_detach_timer,
      ue_context_proto->mutable_implicit_detach_timer());
  mme_app_timer_to_proto(
      state_ue_context->initial_context_setup_rsp_timer,
      ue_context_proto->mutable_initial_context_setup_rsp_timer());
  mme_app_timer_to_proto(
      state_ue_context->ue_context_modification_timer,
      ue_context_proto->mutable_ue_context_modification_timer());
  mme_app_timer_to_proto(
      state_ue_context->paging_response_timer,
      ue_context_proto->mutable_paging_response_timer());
  ue_context_proto->set_rau_tau_timer(state_ue_context->rau_tau_timer);
  mme_app_timer_to_proto(
      state_ue_context->ulr_response_timer,
      ue_context_proto->mutable_ulr_response_timer());
  ue_context_proto->mutable_time_mobile_reachability_timer_started()
      ->set_seconds(state_ue_context->time_mobile_reachability_timer_started);
  ue_context_proto->mutable_time_implicit_detach_timer_started()->set_seconds(
      state_ue_context->time_implicit_detach_timer_started);
  ue_context_proto->mutable_time_paging_response_timer_started()->set_seconds(
      state_ue_context->time_paging_response_timer_started);
  ue_context_proto->set_paging_retx_count(state_ue_context->paging_retx_count);
  ue_context_proto->mutable_time_ics_rsp_timer_started()->set_seconds(
      state_ue_context->time_ics_rsp_timer_started);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

void MmeNasStateConverter::proto_to_ue_mm_context(
    const oai::UeContext& ue_context_proto,
    ue_mm_context_t* state_ue_mm_context) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  state_ue_mm_context->msisdn = bfromcstr(ue_context_proto.msisdn().c_str());
  state_ue_mm_context->ue_context_rel_cause =
      static_cast<enum s1cause>(ue_context_proto.rel_cause());
  state_ue_mm_context->mm_state =
      static_cast<mm_state_t>(ue_context_proto.mm_state());
  state_ue_mm_context->ecm_state =
      static_cast<ecm_state_t>(ue_context_proto.ecm_state());
  NasStateConverter::proto_to_emm_context(
      ue_context_proto.emm_context(), &state_ue_mm_context->emm_context);

  state_ue_mm_context->sctp_assoc_id_key = ue_context_proto.sctp_assoc_id_key();
  state_ue_mm_context->enb_ue_s1ap_id    = ue_context_proto.enb_ue_s1ap_id();
  state_ue_mm_context->enb_s1ap_id_key   = ue_context_proto.enb_s1ap_id_key();
  state_ue_mm_context->mme_ue_s1ap_id    = ue_context_proto.mme_ue_s1ap_id();

  StateConverter::proto_to_ecgi(
      ue_context_proto.e_utran_cgi(), &state_ue_mm_context->e_utran_cgi);
  state_ue_mm_context->cell_age = ue_context_proto.cell_age();
  bytes_to_lai(ue_context_proto.lai().c_str(), &state_ue_mm_context->lai);

  StateConverter::proto_to_apn_config_profile(
      ue_context_proto.apn_config(), &state_ue_mm_context->apn_config_profile);

  state_ue_mm_context->subscriber_status =
      (subscriber_status_t) ue_context_proto.subscriber_status();
  state_ue_mm_context->network_access_mode =
      (network_access_mode_t) ue_context_proto.network_access_mode();
  state_ue_mm_context->access_restriction_data =
      ue_context_proto.access_restriction_data();
  if (ue_context_proto.apn_oi_replacement().length() > 0) {
    state_ue_mm_context->apn_oi_replacement = bfromcstr_with_str_len(
        ue_context_proto.apn_oi_replacement().c_str(),
        ue_context_proto.apn_oi_replacement().length());
  }
  state_ue_mm_context->mme_teid_s11 = ue_context_proto.mme_teid_s11();
  StateConverter::proto_to_ambr(
      ue_context_proto.subscribed_ue_ambr(),
      &state_ue_mm_context->subscribed_ue_ambr);
  StateConverter::proto_to_ambr(
      ue_context_proto.used_ue_ambr(), &state_ue_mm_context->used_ue_ambr);
  StateConverter::proto_to_ambr(
      ue_context_proto.used_ambr(), &state_ue_mm_context->used_ambr);
  state_ue_mm_context->nb_active_pdn_contexts =
      ue_context_proto.nb_active_pdn_contexts();
  proto_to_pdn_context_list(ue_context_proto, state_ue_mm_context);
  proto_to_bearer_context_list(ue_context_proto, state_ue_mm_context);
  state_ue_mm_context->ue_radio_capability = nullptr;
  if (ue_context_proto.ue_radio_capability().length() > 0) {
    state_ue_mm_context->ue_radio_capability = bfromcstr_with_str_len(
        ue_context_proto.ue_radio_capability().c_str(),
        ue_context_proto.ue_radio_capability().length());
  }
  state_ue_mm_context->send_ue_purge_request =
      ue_context_proto.send_ue_purge_request();
  state_ue_mm_context->hss_initiated_detach =
      ue_context_proto.hss_initiated_detach();
  state_ue_mm_context->location_info_confirmed_in_hss =
      ue_context_proto.location_info_confirmed_in_hss();
  state_ue_mm_context->ppf = ue_context_proto.ppf();
  state_ue_mm_context->subscription_known =
      ue_context_proto.subscription_known();
  state_ue_mm_context->num_reg_sub = ue_context_proto.num_reg_sub();
  proto_to_regional_subscription(ue_context_proto, state_ue_mm_context);
  state_ue_mm_context->path_switch_req = ue_context_proto.path_switch_req();
  state_ue_mm_context->granted_service =
      (granted_service_t) ue_context_proto.granted_service();
  state_ue_mm_context->cs_fallback_indicator =
      ue_context_proto.cs_fallback_indicator();

  proto_to_sgs_context(
      ue_context_proto.sgs_context(), state_ue_mm_context->sgs_context);

  proto_to_mme_app_timer(
      ue_context_proto.mobile_reachability_timer(),
      &state_ue_mm_context->mobile_reachability_timer);
  proto_to_mme_app_timer(
      ue_context_proto.implicit_detach_timer(),
      &state_ue_mm_context->implicit_detach_timer);
  proto_to_mme_app_timer(
      ue_context_proto.initial_context_setup_rsp_timer(),
      &state_ue_mm_context->initial_context_setup_rsp_timer);
  proto_to_mme_app_timer(
      ue_context_proto.ue_context_modification_timer(),
      &state_ue_mm_context->ue_context_modification_timer);
  proto_to_mme_app_timer(
      ue_context_proto.ulr_response_timer(),
      &state_ue_mm_context->ulr_response_timer);
  proto_to_mme_app_timer(
      ue_context_proto.paging_response_timer(),
      &state_ue_mm_context->paging_response_timer);
  state_ue_mm_context->time_mobile_reachability_timer_started =
      ue_context_proto.time_mobile_reachability_timer_started().seconds();
  state_ue_mm_context->time_implicit_detach_timer_started =
      ue_context_proto.time_implicit_detach_timer_started().seconds();
  state_ue_mm_context->time_paging_response_timer_started =
      ue_context_proto.time_paging_response_timer_started().seconds();
  state_ue_mm_context->paging_retx_count = ue_context_proto.paging_retx_count();
  state_ue_mm_context->time_ics_rsp_timer_started =
      ue_context_proto.time_ics_rsp_timer_started().seconds();
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

/*********************************************************
 *                MME app state<-> Proto                  *
 * Functions to serialize/desearialize MME app state      *
 * The caller is responsible for all memory management    *
 **********************************************************/
void MmeNasStateConverter::state_to_proto(
    const mme_app_desc_t* mme_nas_state_p, oai::MmeNasState* state_proto) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  state_proto->set_nb_ue_attached(mme_nas_state_p->nb_ue_attached);
  state_proto->set_nb_ue_connected(mme_nas_state_p->nb_ue_connected);
  state_proto->set_nb_default_eps_bearers(
      mme_nas_state_p->nb_default_eps_bearers);
  state_proto->set_nb_s1u_bearers(mme_nas_state_p->nb_s1u_bearers);
  state_proto->set_nb_ue_managed(mme_nas_state_p->nb_ue_managed);
  state_proto->set_nb_ue_idle(mme_nas_state_p->nb_ue_idle);
  state_proto->set_nb_bearers_managed(mme_nas_state_p->nb_bearers_managed);
  state_proto->set_nb_ue_since_last_stat(
      mme_nas_state_p->nb_ue_since_last_stat);
  state_proto->set_nb_bearers_since_last_stat(
      mme_nas_state_p->nb_bearers_since_last_stat);
  state_proto->set_mme_app_ue_s1ap_id_generator(
      mme_nas_state_p->mme_app_ue_s1ap_id_generator);

  // copy mme_ue_contexts
  auto mme_ue_ctxts_proto = state_proto->mutable_mme_ue_contexts();

  OAILOG_DEBUG(LOG_MME_APP, "IMSI table to proto");
  hashtable_uint64_ts_to_proto(
      mme_nas_state_p->mme_ue_contexts.imsi_mme_ue_id_htbl,
      mme_ue_ctxts_proto->mutable_imsi_ue_id_htbl());
  OAILOG_DEBUG(LOG_MME_APP, "Tunnel table to proto");
  hashtable_uint64_ts_to_proto(
      mme_nas_state_p->mme_ue_contexts.tun11_ue_context_htbl,
      mme_ue_ctxts_proto->mutable_tun11_ue_id_htbl());
  OAILOG_DEBUG(LOG_MME_APP, "Enb_Ue_S1ap_id table to proto");
  hashtable_uint64_ts_to_proto(
      mme_nas_state_p->mme_ue_contexts.enb_ue_s1ap_id_ue_context_htbl,
      mme_ue_ctxts_proto->mutable_enb_ue_id_ue_id_htbl());
  guti_table_to_proto(
      mme_nas_state_p->mme_ue_contexts.guti_ue_context_htbl,
      mme_ue_ctxts_proto->mutable_guti_ue_id_htbl());
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

void MmeNasStateConverter::proto_to_state(
    const oai::MmeNasState& state_proto, mme_app_desc_t* mme_nas_state_p) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  mme_nas_state_p->nb_ue_attached  = state_proto.nb_ue_attached();
  mme_nas_state_p->nb_ue_connected = state_proto.nb_ue_connected();
  mme_nas_state_p->nb_default_eps_bearers =
      state_proto.nb_default_eps_bearers();
  mme_nas_state_p->nb_s1u_bearers        = state_proto.nb_s1u_bearers();
  mme_nas_state_p->nb_ue_managed         = state_proto.nb_ue_managed();
  mme_nas_state_p->nb_ue_idle            = state_proto.nb_ue_idle();
  mme_nas_state_p->nb_bearers_managed    = state_proto.nb_bearers_managed();
  mme_nas_state_p->nb_ue_since_last_stat = state_proto.nb_ue_since_last_stat();
  mme_nas_state_p->nb_bearers_since_last_stat =
      state_proto.nb_bearers_since_last_stat();
  mme_nas_state_p->mme_app_ue_s1ap_id_generator =
      state_proto.mme_app_ue_s1ap_id_generator();
  if (mme_nas_state_p->mme_app_ue_s1ap_id_generator == 0) {  // uninitialized
    mme_nas_state_p->mme_app_ue_s1ap_id_generator = 1;
  }
  OAILOG_INFO(LOG_MME_APP, "Done reading MME statistics from data store");

  // copy mme_ue_contexts
  oai::MmeUeContext mme_ue_ctxts_proto = state_proto.mme_ue_contexts();

  mme_ue_context_t* mme_ue_ctxt_state = &mme_nas_state_p->mme_ue_contexts;
  // copy maps to hashtables
  OAILOG_INFO(LOG_MME_APP, "Hashtable MME UE ID => IMSI");
  proto_to_hashtable_uint64_ts(
      mme_ue_ctxts_proto.imsi_ue_id_htbl(),
      mme_ue_ctxt_state->imsi_mme_ue_id_htbl);
  OAILOG_INFO(LOG_MME_APP, "Hashtable TEID 11 => MME UE ID");
  proto_to_hashtable_uint64_ts(
      mme_ue_ctxts_proto.tun11_ue_id_htbl(),
      mme_ue_ctxt_state->tun11_ue_context_htbl);
  OAILOG_INFO(LOG_MME_APP, "Hashtable ENB UE S1AP ID => MME UE ID");
  proto_to_hashtable_uint64_ts(
      mme_ue_ctxts_proto.enb_ue_id_ue_id_htbl(),
      mme_ue_ctxt_state->enb_ue_s1ap_id_ue_context_htbl);
  proto_to_guti_table(
      mme_ue_ctxts_proto.guti_ue_id_htbl(),
      mme_ue_ctxt_state->guti_ue_context_htbl);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

void MmeNasStateConverter::ue_to_proto(
    const ue_mm_context_t* ue_ctxt, oai::UeContext* ue_ctxt_proto) {
  ue_context_to_proto(ue_ctxt, ue_ctxt_proto);
}

void MmeNasStateConverter::proto_to_ue(
    const oai::UeContext& ue_ctxt_proto, ue_mm_context_t* ue_ctxt) {
  proto_to_ue_mm_context(ue_ctxt_proto, ue_ctxt);
}
}  // namespace lte
}  // namespace magma
