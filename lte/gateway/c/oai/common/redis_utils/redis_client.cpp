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

#include "redis_client.h"

#ifdef __cplusplus
extern "C" {
#endif

#include <common_defs.h>

#ifdef __cplusplus
}
#endif

#include "ServiceConfigLoader.h"

using google::protobuf::Message;

namespace magma {
namespace lte {

RedisClient::RedisClient():
  db_client_(std::make_unique<cpp_redis::client>()),
  is_connected_(false)
{
  init_db_connection(LOCALHOST);
}

void RedisClient::init_db_connection(const std::string& addr)
{
  magma::ServiceConfigLoader loader;

  auto config = loader.load_service_config("redis");
  auto port = config["port"].as<uint32_t>();

  // Make connection to db
  db_client_->connect(addr, port, nullptr);

  is_connected_ = true;
}

int RedisClient::serialize(
  const Message& proto_msg,
  std::string& str_to_serialize)
{
  if (!proto_msg.SerializeToString(&str_to_serialize)) {
    return RETURNerror;
  }
  return RETURNok;
}

int RedisClient::deserialize(
  Message& proto_msg,
  const std::string& str_to_deserialize)
{
  if (!proto_msg.ParseFromString(str_to_deserialize)) {
    return RETURNerror;
  }
  return RETURNok;
}

int RedisClient::write(const std::string& key, const std::string& value)
{
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

std::string RedisClient::read(const std::string& key)
{
  auto db_read_fut = db_client_->get(key);
  db_client_->sync_commit();
  auto db_read_reply = db_read_fut.get();

  if (
    db_read_reply.is_null() || db_read_reply.is_error() ||
    !db_read_reply.is_string()) {
    throw std::runtime_error("Could not read from redis");
  }

  return db_read_reply.as_string();
}

int RedisClient::write_proto(const std::string& key, const Message& proto_msg)
{
  std::string str_value;
  if (serialize(proto_msg, str_value) != RETURNok) {
    return RETURNerror;
  }

  if (write(key, str_value) != RETURNok) {
    return RETURNerror;
  }
  return RETURNok;
}

int RedisClient::read_proto(const std::string& key, Message& proto_msg)
{
  try {
    std::string str_value = read(key);

    if (deserialize(proto_msg, str_value) != RETURNok) {
      return RETURNerror;
    }

    return RETURNok;
  } catch (const std::runtime_error& e) {
    return RETURNerror;
  }
}

} // namespace lte
} // namespace magma
