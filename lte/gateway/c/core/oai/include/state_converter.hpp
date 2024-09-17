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
 *------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
#include "lte/gateway/c/core/oai/lib/hashtable/hashtable.h"

#ifdef __cplusplus
}
#endif

#include <functional>

#include <google/protobuf/map.h>
#include <functional>
#include "lte/gateway/c/core/oai/include/map.h"
#include "lte/gateway/c/core/oai/include/proto_map.hpp"
#include "lte/protos/oai/common_types.pb.h"

namespace magma {
namespace lte {

#define BSTRING_TO_STRING(bstr, str_ptr)                \
  do {                                                  \
    *str_ptr = std::string(bdata(bstr), blength(bstr)); \
  } while (0) /* Convert bstring to std::string */
#define STRING_TO_BSTRING(str, bstr) \
  do {                               \
    bstr = bfromcstr(str.c_str());   \
  } while (0) /* Convert bstring to std::string */

/**
 * StateConverter is a base class for state conversion between tasks state
 * structs and protobuf objects. This class is used to support specific state
 * conversion for each task that extends from it. The class doesn't hold memory,
 * all memory is owned by caller.
 */
class StateConverter {
 protected:
  StateConverter();
  ~StateConverter();

  /**
   * Function that converts hashtable struct to protobuf Map instance, using
   * a conversion function to convert each node of the hashtable, memory
   * is owned by the caller.
   * @tparam NodeType struct type of hashmap node entry
   * @tparam ProtoMessage protobuf type for proto map value entry
   * @param state_ht hashtable_ts_t struct to convert from
   * @param proto_map protobuf Map instance to convert to
   * @param conversion_callable conversion function for each entry of hashtable
   * @param log_task_level log level for task (LOG_MME_APP, LOG_SPGW_APP)
   */
  // TODO(rsarwad): Shall be removed once all modules are converted to cpp and
  // use protobuf map instead of hash table
  template <typename NodeType, typename ProtoMessage>
  static void hashtable_ts_to_proto(
      hash_table_ts_t* state_ht,
      google::protobuf::Map<unsigned int, ProtoMessage>* proto_map,
      std::function<void(NodeType*, ProtoMessage*)> conversion_callable,
      log_proto_t log_task_level) {
    hashtable_key_array_t* ht_keys = hashtable_ts_get_keys(state_ht);
    hashtable_rc_t ht_rc;
    if (ht_keys == nullptr) {
      return;
    }

    for (int i = 0; i < ht_keys->num_keys; i++) {
      NodeType* node;
      ht_rc = hashtable_ts_get(state_ht, (hash_key_t)ht_keys->keys[i],
                               (void**)&node);
      if (ht_rc == HASH_TABLE_OK) {
        ProtoMessage proto;
        conversion_callable((NodeType*)node, &proto);
        (*proto_map)[ht_keys->keys[i]] = proto;
      } else {
        OAILOG_ERROR(log_task_level, "Key %lu not found on %s hashtable",
                     ht_keys->keys[i], state_ht->name->data);
      }
    }
    FREE_HASHTABLE_KEY_ARRAY(ht_keys);
  }

  template <typename ProtoMessage, typename NodeType>
  static void proto_to_hashtable_ts(
      const google::protobuf::Map<unsigned int, ProtoMessage>& proto_map,
      hash_table_ts_t* state_ht,
      std::function<void(const ProtoMessage&, NodeType*)> conversion_callable,
      log_proto_t log_task_level) {
    for (const auto& entry : proto_map) {
      auto proto = entry.second;
      NodeType* node_type;
      node_type = (NodeType*)calloc(1, sizeof(NodeType));
      conversion_callable(proto, node_type);
      auto ht_rc =
          hashtable_ts_insert(state_ht, (hash_key_t)entry.first, node_type);
      if (ht_rc != HASH_TABLE_OK) {
        if (ht_rc == HASH_TABLE_INSERT_OVERWRITTEN_DATA) {
          OAILOG_INFO(LOG_SPGW_APP, "Overwriting data on key: %i", entry.first);
        } else {
          OAILOG_ERROR(log_task_level, "Failed to insert node on hashtable %s",
                       state_ht->name->data);
        }
      }
    }
  }

  template <typename state_map_t, typename NodeType, typename ProtoMessage>
  static void state_map_to_proto(
      state_map_t state_map,
      google::protobuf::Map<unsigned int, ProtoMessage>* proto_map,
      std::function<void(NodeType*, ProtoMessage*)> conversion_callable,
      log_proto_t log_task_level) {
    if (!(state_map.size())) {
      return;
    }

    for (auto itr = state_map.map->begin(); itr != state_map.map->end();
         itr++) {
      NodeType* node = itr->second;
      if (node) {
        ProtoMessage proto;
        conversion_callable(reinterpret_cast<NodeType*>(node), &proto);
        (*proto_map)[itr->first] = proto;
      } else {
        OAILOG_ERROR(log_task_level, "Key %lu without value in %s map",
                     itr->first, state_map.get_name());
      }
    }
  }

  template <typename state_map_t, typename ProtoMessage, typename NodeType>
  static void proto_to_state_map(
      const google::protobuf::Map<unsigned int, ProtoMessage>& proto_map,
      state_map_t state_map,
      std::function<void(const ProtoMessage&, NodeType*)> conversion_callable,
      log_proto_t log_task_level) {
    for (const auto& entry : proto_map) {
      bool failed_to_write = true;
      auto proto = entry.second;
      NodeType* node_type;
      node_type = reinterpret_cast<NodeType*>(new NodeType());
      conversion_callable(proto, node_type);
      proto_map_rc_t rc = state_map.insert(entry.first, node_type);
      if (rc != PROTO_MAP_OK) {
        if (rc == PROTO_MAP_KEY_ALREADY_EXISTS) {
          if ((state_map.remove(entry.first)) == PROTO_MAP_OK) {
            if ((state_map.insert(entry.first, node_type)) == PROTO_MAP_OK) {
              OAILOG_INFO(LOG_SPGW_APP, "Overwriting data on key: %i",
                          entry.first);
              failed_to_write = false;
            }
          }
        }
      } else {
        failed_to_write = false;
      }
      if (failed_to_write) {
        OAILOG_ERROR(log_task_level, "Failed to insert key in map %s",
                     state_map.get_name());
      }
    }
  }

  static void hashtable_uint64_ts_to_proto(
      hash_table_uint64_ts_t* htbl,
      google::protobuf::Map<unsigned long, unsigned long>* proto_map);

  static void proto_to_hashtable_uint64_ts(
      const google::protobuf::Map<unsigned long, unsigned long>& proto_map,
      hash_table_uint64_ts_t* state_htbl);

  static void guti_to_proto(const guti_t& guti_state, oai::Guti* guti_proto);
  static void proto_to_guti(const oai::Guti& guti_proto, guti_t* state_guti);

  static void ecgi_to_proto(const ecgi_t& state_ecgi, oai::Ecgi* ecgi_proto);
  static void proto_to_ecgi(const oai::Ecgi& ecgi_proto, ecgi_t* state_ecgi);

  static void eps_subscribed_qos_profile_to_proto(
      const eps_subscribed_qos_profile_t& state_eps_subscribed_qos_profile,
      oai::EpsSubscribedQosProfile* eps_subscribed_qos_profile_proto);
  static void ambr_to_proto(const ambr_t& state_ambr, oai::Ambr* ambr_proto);
  static void apn_configuration_to_proto(
      const apn_configuration_t& state_apn_configuration,
      oai::ApnConfig* apn_config_proto);
  static void apn_config_profile_to_proto(
      const apn_config_profile_t& state_apn_config_profile,
      oai::ApnConfigProfile* apn_config_profile_proto);

  static void proto_to_eps_subscribed_qos_profile(
      const oai::EpsSubscribedQosProfile& eps_subscribed_qos_profile_proto,
      eps_subscribed_qos_profile_t* state_eps_subscribed_qos_profile);
  static void proto_to_ambr(const oai::Ambr& ambr_proto, ambr_t* state_ambr);
  static void proto_to_apn_configuration(
      const oai::ApnConfig& apn_config_proto,
      apn_configuration_t* state_apn_configuration);
  static void proto_to_apn_config_profile(
      const oai::ApnConfigProfile& apn_config_profile_proto,
      apn_config_profile_t* state_apn_config_profile);

  /***********************************************************
   *                 Map <-> Proto
   * Functions to serialize/deserialize in-memory maps
   * for AMF task. Only AMF task inserts/removes elements in
   * the maps, so these calls are thread-safe.
   * We only need to lock the UE context structure as it can
   * also be accessed by the NAS task. If map is empty
   * the proto field is also empty
   ***********************************************************/

  static void map_uint64_uint64_to_proto(
      map_uint64_uint64_t map,
      google::protobuf::Map<uint64_t, uint64_t>* proto_map);

  static void proto_to_map_uint64_uint64(
      const google::protobuf::Map<uint64_t, uint64_t>& proto_map,
      map_uint64_uint64_t* map);

 private:
  static void plmn_to_chars(const plmn_t& state_plmn, char* plmn_array);
  static void chars_to_plmn(const char* plmn_array, plmn_t* state_plmn);
};

}  // namespace lte
}  // namespace magma
