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
 *------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include "spgw_state_manager.h"

namespace magma {
namespace lte {

SpgwStateManager::SpgwStateManager() :
    persist_state_(false),
    is_initialized_(false),
    spgw_state_cache_p_(nullptr),
    state_dirty_(false),
    config_(nullptr) {}

SpgwStateManager& SpgwStateManager::getInstance() {
  static SpgwStateManager instance;
  return instance;
}

void SpgwStateManager::init(bool persist_state, const spgw_config_t* config) {
  is_initialized_ = true;
  persist_state_ = persist_state;
  config_ = config;
  spgw_state_cache_p_ = create_spgw_state();
}

spgw_state_t* SpgwStateManager::get_spgw_state(bool read_from_db) {
  AssertFatal(
      is_initialized_,
      "SpgwStateManager init() function should be called to initialize state.");

  this->state_dirty_ = true;

  AssertFatal(spgw_state_cache_p_ != nullptr, "SPGW state cache is NULL");

  if (persist_state_ && read_from_db) {
    free_spgw_state();
    spgw_state_cache_p_ = create_spgw_state();
    read_state_from_db();
  }

  return spgw_state_cache_p_;
}

spgw_state_t* SpgwStateManager::create_spgw_state() {
  AssertFatal(
      is_initialized_,
      "SpgwStateManager init() function should be called to initialize state.");

  // Allocating spgw_state_p
  spgw_state_t* state_p;
  state_p = (spgw_state_t*)calloc(1, sizeof(spgw_state_t));

  bstring b = bfromcstr(SGW_S11_TEID_MME_HT_NAME);
  state_p->sgw_state.s11teid2mme =
      hashtable_ts_create(SGW_STATE_CONTEXT_HT_MAX_SIZE, nullptr, nullptr, b);
  btrunc(b, 0);

  bassigncstr(b, S11_BEARER_CONTEXT_INFO_HT_NAME);
  state_p->sgw_state.s11_bearer_context_information = hashtable_ts_create(
      SGW_STATE_CONTEXT_HT_MAX_SIZE, nullptr,
      (void (*)(void**))sgw_free_s11_bearer_context_information, b);
  bdestroy_wrapper(&b);

  state_p->sgw_state.sgw_ip_address_S1u_S12_S4_up.s_addr =
      config_->sgw_config.ipv4.S1u_S12_S4_up.s_addr;

  // TODO: Refactor GTPv1u_data state
  state_p->sgw_state.gtpv1u_data.sgw_ip_address_for_S1u_S12_S4_up =
      state_p->sgw_state.sgw_ip_address_S1u_S12_S4_up;

  // Creating PGW related state structs
  state_p->pgw_state.deactivated_predefined_pcc_rules = hashtable_ts_create(
      MAX_PREDEFINED_PCC_RULES_HT_SIZE, nullptr, pgw_free_pcc_rule, nullptr);

  state_p->pgw_state.predefined_pcc_rules = hashtable_ts_create(
      MAX_PREDEFINED_PCC_RULES_HT_SIZE, nullptr, pgw_free_pcc_rule, nullptr);

  // TO DO: RANDOM
  state_p->sgw_state.tunnel_id = 0;

  state_p->sgw_state.gtpv1u_teid = 0;

  return state_p;
}

void SpgwStateManager::free_spgw_state() {
  AssertFatal(
      is_initialized_,
      "SpgwStateManager init() function should be called to initialize state.");

  if (spgw_state_cache_p_ == nullptr) {
    return;
  }

  if (hashtable_ts_destroy(spgw_state_cache_p_->sgw_state.s11teid2mme) !=
      HASH_TABLE_OK) {
    OAI_FPRINTF_ERR(
        "An error occurred while destroying SGW s11teid2mme hashtable");
  }

  if (hashtable_ts_destroy(
          spgw_state_cache_p_->sgw_state.s11_bearer_context_information) !=
      HASH_TABLE_OK) {
    OAI_FPRINTF_ERR(
        "An error occurred while destroying SGW s11_bearer_context_information "
        "hashtable");
  }

  if (spgw_state_cache_p_->pgw_state.deactivated_predefined_pcc_rules) {
    hashtable_ts_destroy(
        spgw_state_cache_p_->pgw_state.deactivated_predefined_pcc_rules);
  }

  if (spgw_state_cache_p_->pgw_state.predefined_pcc_rules) {
    hashtable_ts_destroy(spgw_state_cache_p_->pgw_state.predefined_pcc_rules);
  }
  free(spgw_state_cache_p_);
}

int SpgwStateManager::read_state_from_db() {
  AssertFatal(
      is_initialized_,
      "SpgwStateManager init() function should be called to initialize state.");

  if (persist_state_) {
    auto db_read_fut = db_client_->get(SPGW_STATE_TABLE_NAME);
    db_client_->sync_commit();
    auto db_read_reply = db_read_fut.get();

    if (db_read_reply.is_null() || db_read_reply.is_error() ||
        !db_read_reply.is_string()) {
      OAILOG_ERROR(LOG_SPGW_APP, "Failed to read state from db");
      return RETURNok;
    } else {
      gateway::spgw::SpgwState state_proto = gateway::spgw::SpgwState();
      if (!state_proto.ParseFromString(db_read_reply.as_string())) {
        OAILOG_ERROR(LOG_SPGW_APP, "Failed to parse state");
        return RETURNok;
      }

      SpgwStateConverter::spgw_proto_to_state(state_proto, spgw_state_cache_p_);
      OAILOG_INFO(LOG_SPGW_APP, "Finished reading state");
    }
  }
  return RETURNok;
}

void SpgwStateManager::write_state_to_db() {
  AssertFatal(
      is_initialized_,
      "SpgwStateManager init() function should be called to initialize state.");

  if (!state_dirty_) {
    OAILOG_ERROR(LOG_SPGW_APP,
                 "Tried to put SPGW state while it was not in use");
    return;
  }

  if (persist_state_) {
    std::string serialized_state_s;
    gateway::spgw::SpgwState state_proto = gateway::spgw::SpgwState();
    SpgwStateConverter::spgw_state_to_proto(spgw_state_cache_p_, &state_proto);

    if (!state_proto.SerializeToString(&serialized_state_s)) {
      OAILOG_ERROR(LOG_SPGW_APP, "Failed to serialize state protobuf");
      return;
    }

    auto db_write_fut =
        db_client_->set(SPGW_STATE_TABLE_NAME, serialized_state_s);
    db_client_->sync_commit();
    auto db_write_reply = db_write_fut.get();

    if (db_write_reply.is_error()) {
      OAILOG_ERROR(LOG_SPGW_APP, "Failed to write SPGW state to db");
      return;
    }

    OAILOG_INFO(LOG_SPGW_APP, "Finished writing state");
  }

  this->state_dirty_ = false;
}

int SpgwStateManager::init_db_connection(const std::string& addr) {
  // Init db client service config
  magma::ServiceConfigLoader loader;

  auto config = loader.load_service_config("redis");
  auto port = config["port"].as<uint32_t>();

  db_client_ = std::make_unique<cpp_redis::client>();
  // Make connection to db
  db_client_->connect(addr, port, nullptr);

  if (!db_client_->is_connected()) {
    OAILOG_ERROR(LOG_SPGW_APP, "Failed to connect to redis");
    return RETURNerror;
  }

  OAILOG_INFO(LOG_SPGW_APP, "Connected to redis datastore on %s:%u\n",
              addr.c_str(), port);

  return RETURNok;
}

} // namespace lte
} // namespace magma
