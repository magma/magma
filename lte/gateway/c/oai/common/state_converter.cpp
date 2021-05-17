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

#include "state_converter.h"

namespace magma {
namespace lte {

StateConverter::StateConverter()  = default;
StateConverter::~StateConverter() = default;

/*************************************************/
/*        Common Types -> Proto                 */
/*************************************************/

void StateConverter::plmn_to_chars(const plmn_t& state_plmn, char* plmn_array) {
  plmn_array[0] = (char) (state_plmn.mcc_digit1 + ASCII_ZERO);
  plmn_array[1] = (char) (state_plmn.mcc_digit2 + ASCII_ZERO);
  plmn_array[2] = (char) (state_plmn.mcc_digit3 + ASCII_ZERO);
  plmn_array[3] = (char) (state_plmn.mnc_digit1 + ASCII_ZERO);
  plmn_array[4] = (char) (state_plmn.mnc_digit2 + ASCII_ZERO);
  plmn_array[5] = (char) (state_plmn.mnc_digit3 + ASCII_ZERO);
}

void StateConverter::guti_to_proto(
    const guti_t& state_guti, oai::Guti* guti_proto) {
  guti_proto->Clear();

  char plmn_array[PLMN_BYTES];
  plmn_to_chars(state_guti.gummei.plmn, plmn_array);
  guti_proto->set_plmn(plmn_array);
  guti_proto->set_mme_gid(state_guti.gummei.mme_gid);
  guti_proto->set_mme_code(state_guti.gummei.mme_code);
  guti_proto->set_m_tmsi(state_guti.m_tmsi);
}

void StateConverter::ecgi_to_proto(
    const ecgi_t& state_ecgi, oai::Ecgi* ecgi_proto) {
  ecgi_proto->Clear();

  char plmn_array[PLMN_BYTES];
  plmn_to_chars(state_ecgi.plmn, plmn_array);
  ecgi_proto->set_plmn(plmn_array);
  ecgi_proto->set_enb_id(state_ecgi.cell_identity.enb_id);
  ecgi_proto->set_cell_id(state_ecgi.cell_identity.cell_id);
  ecgi_proto->set_empty(state_ecgi.cell_identity.empty);
}

void StateConverter::proto_to_ecgi(
    const oai::Ecgi& ecgi_proto, ecgi_t* state_ecgi) {
  strncpy((char*) &state_ecgi->plmn, ecgi_proto.plmn().c_str(), PLMN_BYTES);

  state_ecgi->cell_identity.enb_id  = ecgi_proto.enb_id();
  state_ecgi->cell_identity.cell_id = ecgi_proto.cell_id();
  state_ecgi->cell_identity.empty   = ecgi_proto.empty();
}

void StateConverter::eps_subscribed_qos_profile_to_proto(
    const eps_subscribed_qos_profile_t& state_eps_subscribed_qos_profile,
    oai::EpsSubscribedQosProfile* eps_subscribed_qos_profile_proto) {
  eps_subscribed_qos_profile_proto->set_qci(
      state_eps_subscribed_qos_profile.qci);
  eps_subscribed_qos_profile_proto->set_priority_level(
      state_eps_subscribed_qos_profile.allocation_retention_priority
          .priority_level);
  eps_subscribed_qos_profile_proto->set_pre_emption_vulnerability(
      state_eps_subscribed_qos_profile.allocation_retention_priority
          .pre_emp_vulnerability);
  eps_subscribed_qos_profile_proto->set_pre_emption_capability(
      state_eps_subscribed_qos_profile.allocation_retention_priority
          .pre_emp_capability);
}

void StateConverter::proto_to_eps_subscribed_qos_profile(
    const oai::EpsSubscribedQosProfile& eps_subscribed_qos_profile_proto,
    eps_subscribed_qos_profile_t* state_eps_subscribed_qos_profile) {
  state_eps_subscribed_qos_profile->qci =
      eps_subscribed_qos_profile_proto.qci();
  state_eps_subscribed_qos_profile->allocation_retention_priority
      .priority_level = eps_subscribed_qos_profile_proto.priority_level();
  state_eps_subscribed_qos_profile->allocation_retention_priority
      .pre_emp_capability =
      (pre_emption_capability_t)
          eps_subscribed_qos_profile_proto.pre_emption_capability();
  state_eps_subscribed_qos_profile->allocation_retention_priority
      .pre_emp_vulnerability =
      (pre_emption_vulnerability_t)
          eps_subscribed_qos_profile_proto.pre_emption_vulnerability();
}

void StateConverter::ambr_to_proto(
    const ambr_t& state_ambr, oai::Ambr* ambr_proto) {
  ambr_proto->set_br_ul(state_ambr.br_ul);
  ambr_proto->set_br_dl(state_ambr.br_dl);
  ambr_proto->set_br_unit(
      static_cast<magma::lte::oai::Ambr::BitrateUnitsAMBR>(state_ambr.br_unit));
}

void StateConverter::proto_to_ambr(
    const oai::Ambr& ambr_proto, ambr_t* state_ambr) {
  state_ambr->br_ul   = ambr_proto.br_ul();
  state_ambr->br_dl   = ambr_proto.br_dl();
  state_ambr->br_unit = (apn_ambr_bitrate_unit_t) ambr_proto.br_unit();
}

void StateConverter::apn_configuration_to_proto(
    const apn_configuration_t& state_apn_configuration,
    oai::ApnConfig* apn_config_proto) {
  apn_config_proto->set_context_identifier(
      state_apn_configuration.context_identifier);
  for (int i = 0; i < state_apn_configuration.nb_ip_address; i++) {
    BSTRING_TO_STRING(
        ip_address_to_bstring(&state_apn_configuration.ip_address[i]),
        apn_config_proto->add_ip_address());
  }
  apn_config_proto->set_pdn_type(state_apn_configuration.pdn_type);
  apn_config_proto->set_service_selection(
      &state_apn_configuration.service_selection,
      state_apn_configuration.service_selection_length);
  eps_subscribed_qos_profile_to_proto(
      state_apn_configuration.subscribed_qos,
      apn_config_proto->mutable_subscribed_qos());
  ambr_to_proto(state_apn_configuration.ambr, apn_config_proto->mutable_ambr());
}

void StateConverter::proto_to_apn_configuration(
    const oai::ApnConfig& apn_config_proto,
    apn_configuration_t* state_apn_configuration) {
  state_apn_configuration->context_identifier =
      apn_config_proto.context_identifier();
  state_apn_configuration->nb_ip_address = apn_config_proto.ip_address_size();
  for (int i = 0; i < state_apn_configuration->nb_ip_address; i++) {
    bstring state_ip_addr_str;
    STRING_TO_BSTRING(apn_config_proto.ip_address(i), state_ip_addr_str);
    bstring_to_ip_address(
        state_ip_addr_str, &state_apn_configuration->ip_address[i]);
  }
  state_apn_configuration->pdn_type = apn_config_proto.pdn_type();
  strcpy(
      state_apn_configuration->service_selection,
      apn_config_proto.service_selection().c_str());
  state_apn_configuration->service_selection_length =
      apn_config_proto.service_selection().length();
  proto_to_eps_subscribed_qos_profile(
      apn_config_proto.subscribed_qos(),
      &state_apn_configuration->subscribed_qos);
  proto_to_ambr(apn_config_proto.ambr(), &state_apn_configuration->ambr);
}

void StateConverter::apn_config_profile_to_proto(
    const apn_config_profile_t& state_apn_config_profile,
    oai::ApnConfigProfile* apn_config_profile_proto) {
  apn_config_profile_proto->set_context_identifier(
      state_apn_config_profile.context_identifier);
  apn_config_profile_proto->set_all_apn_conf_ind(
      state_apn_config_profile.all_apn_conf_ind);
  for (int i = 0; i < state_apn_config_profile.nb_apns; i++) {
    apn_configuration_to_proto(
        state_apn_config_profile.apn_configuration[i],
        apn_config_profile_proto->add_apn_configs());
  }
}

void StateConverter::proto_to_apn_config_profile(
    const oai::ApnConfigProfile& apn_config_profile_proto,
    apn_config_profile_t* state_apn_config_profile) {
  state_apn_config_profile->context_identifier =
      apn_config_profile_proto.context_identifier();
  state_apn_config_profile->all_apn_conf_ind =
      (all_apn_conf_ind_t) apn_config_profile_proto.all_apn_conf_ind();
  state_apn_config_profile->nb_apns =
      apn_config_profile_proto.apn_configs_size();
  for (int i = 0; i < state_apn_config_profile->nb_apns; i++) {
    proto_to_apn_configuration(
        apn_config_profile_proto.apn_configs(i),
        &state_apn_config_profile->apn_configuration[i]);
  }
}

void StateConverter::hashtable_uint64_ts_to_proto(
    hash_table_uint64_ts_t* htbl,
    google::protobuf::Map<unsigned long, unsigned long>* proto_map) {
  hashtable_key_array_t* keys = hashtable_uint64_ts_get_keys(htbl);
  if (keys == nullptr) {
    return;
  }

  for (auto i = 0; i < keys->num_keys; i++) {
    uint64_t val;
    hashtable_rc_t ht_rc = hashtable_uint64_ts_get(htbl, keys->keys[i], &val);
    if (ht_rc == HASH_TABLE_OK) {
      (*proto_map)[keys->keys[i]] = val;
    } else {
      OAILOG_ERROR(
          LOG_UTIL, "Key %lu not in %s", keys->keys[i], htbl->name->data);
    }
  }

  FREE_HASHTABLE_KEY_ARRAY(keys);
}

void StateConverter::proto_to_hashtable_uint64_ts(
    const google::protobuf::Map<unsigned long, unsigned long>& proto_map,
    hash_table_uint64_ts_t* state_htbl) {
  for (auto const& kv : proto_map) {
    uint64_t id  = kv.first;
    uint64_t val = kv.second;

    hashtable_rc_t ht_rc =
        hashtable_uint64_ts_insert(state_htbl, (const hash_key_t) id, val);
    if (ht_rc != HASH_TABLE_OK) {
      OAILOG_ERROR(
          LOG_UTIL, "Failed to insert value %lu in table %s: error: %s\n", val,
          state_htbl->name->data, hashtable_rc_code2string(ht_rc));
    }
  }
}

}  // namespace lte
}  // namespace magma
