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

#include "redis_client.h"

#ifdef __cplusplus
extern "C" {
#endif

#include <common_defs.h>

#ifdef __cplusplus
}
#endif

#include "includes/ServiceConfigLoader.h"
#include <yaml-cpp/yaml.h>  // IWYU pragma: keep

using google::protobuf::Message;

namespace magma {
namespace lte {

RedisClient::RedisClient(bool init_connection)
    : db_client_(std::make_unique<cpp_redis::client>()), is_connected_(false) {
  if (init_connection) {
    init_db_connection();
  }
}

void RedisClient::init_db_connection() {
  magma::ServiceConfigLoader loader;

  auto config = loader.load_service_config("redis");
  auto addr   = config["bind"].as<std::string>();
  auto port   = config["port"].as<uint32_t>();

  // Make connection to db
  db_client_->connect(addr, port, nullptr);

  is_connected_ = true;
}

status_code_e RedisClient::write(
    const std::string& key, const std::string& value) {
  if (!is_connected()) {
    return RETURNerror;
  }

  auto db_write_fut = db_client_->set(key, value);
  db_client_->sync_commit();
  auto db_write_reply = db_write_fut.get();

  if (db_write_reply.is_error()) {
    return RETURNerror;
  }

  return RETURNok;
}

std::string RedisClient::read(const std::string& key) {
  auto db_read_fut = db_client_->get(key);
  db_client_->sync_commit();
  auto db_read_reply = db_read_fut.get();

  if (db_read_reply.is_null()) {
    return "";
  }

  if (db_read_reply.is_error() || !db_read_reply.is_string()) {
    throw std::runtime_error("Could not read from redis");
  }

  return db_read_reply.as_string();
}

status_code_e RedisClient::write_proto_str(
    const std::string& key, const std::string& proto_msg, uint64_t version) {
  orc8r::RedisState wrapper_proto = orc8r::RedisState();
  wrapper_proto.set_serialized_msg(proto_msg);
  wrapper_proto.set_version(version);

  std::string str_value;
  if (serialize(wrapper_proto, str_value) != RETURNok) {
    return RETURNerror;
  }
  if (write(key, str_value) != RETURNok) {
    return RETURNerror;
  }
  return RETURNok;
}

status_code_e RedisClient::read_proto(
    const std::string& key, Message& proto_msg) {
  orc8r::RedisState wrapper_proto = orc8r::RedisState();
  if (read_redis_state(key, wrapper_proto) != RETURNok) {
    return RETURNerror;
  }

  std::string wrapped_val = wrapper_proto.serialized_msg();
  if (deserialize(proto_msg, wrapped_val) != RETURNok) {
    return RETURNerror;
  }

  return RETURNok;
}

int RedisClient::read_version(const std::string& key) {
  orc8r::RedisState wrapper_proto = orc8r::RedisState();
  if (read_redis_state(key, wrapper_proto) != RETURNok) {
    return RETURNerror;
  }

  return wrapper_proto.version();
}

status_code_e RedisClient::clear_keys(
    const std::vector<std::string>& keys_to_clear) {
  auto db_write = db_client_->del(keys_to_clear);
  db_client_->sync_commit();
  auto reply = db_write.get();

  if (reply.is_error()) {
    return RETURNerror;
  }

  return RETURNok;
}

std::vector<std::string> RedisClient::get_keys(const std::string& pattern) {
  size_t cursor = 0;
  std::vector<std::string> replies;
  do {
    auto reply_future = db_client_->scan(cursor, pattern);
    db_client_->sync_commit();
    auto db_read_reply = reply_future.get();

    if (db_read_reply.is_null()) {
      return replies;
    }

    if (db_read_reply.is_error() || !db_read_reply.is_array()) {
      throw std::runtime_error("Could not read from redis");
    }
    // First result is cursor, second result is pattern matched keys
    auto response      = db_read_reply.as_array();
    auto returned_keys = response[1];

    for (const auto& reply : returned_keys.as_array()) {
      replies.emplace_back(reply.as_string());
    }

    cursor = std::stoi(response[0].as_string());
    ;
  } while (cursor != 0);

  return replies;
}

status_code_e RedisClient::read_redis_state(
    const std::string& key, orc8r::RedisState& state_out) {
  try {
    std::string str_value = read(key);
    if (deserialize(state_out, str_value) != RETURNok) {
      return RETURNerror;
    }
    return RETURNok;
  } catch (const std::runtime_error& e) {
    return RETURNerror;
  }
}

status_code_e RedisClient::serialize(
    const Message& proto_msg, std::string& str_to_serialize) {
  if (!proto_msg.SerializeToString(&str_to_serialize)) {
    return RETURNerror;
  }
  return RETURNok;
}

status_code_e RedisClient::deserialize(
    Message& proto_msg, const std::string& str_to_deserialize) {
  if (!proto_msg.ParseFromString(str_to_deserialize)) {
    return RETURNerror;
  }
  return RETURNok;
}

}  // namespace lte
}  // namespace magma
