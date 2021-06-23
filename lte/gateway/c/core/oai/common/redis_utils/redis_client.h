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

#include <string>

#include <cpp_redis/cpp_redis>
#include <google/protobuf/message.h>

#include <common_defs.h>
#include "orc8r/protos/redis.pb.h"

namespace magma {
namespace lte {

class RedisClient {
 public:
  explicit RedisClient(bool init_connection);
  ~RedisClient() = default;

  /**
   * Initializes a connection to the redis datastore configured in redis.yml
   * @return response code of success / error with db connection
   */
  void init_db_connection();

  /**
   * Returns the value on redis db mapped to a key
   * @param key
   * @return string repr of value
   */
  std::string read(const std::string& key);

  /**
   * Writes a str value to redis mapped to str key
   * @param key
   * @param value
   * @return response code of operation
   */
  status_code_e write(const std::string& key, const std::string& value);

  /**
   * Writes a protobuf object to redis
   * @param key
   * @param proto_msg
   * @param version
   * @return response code of operation
   */
  status_code_e write_proto_str(
      const std::string& key, const std::string& proto_msg, uint64_t version);

  /**
   * Converts protobuf Message and parses it to string
   * @param proto_msg
   * @param str_to_serialize
   * @return response code of operation
   */
  static status_code_e serialize(
      const google::protobuf::Message& proto_msg,
      std::string& str_to_serialize);

  /**
   * Reads value from redis mapped to key and returns proto object
   * @param key
   * @return response code of operation
   */
  status_code_e read_proto(
      const std::string& key, google::protobuf::Message& proto_msg);

  int read_version(const std::string& key);

  status_code_e clear_keys(const std::vector<std::string>& keys_to_clear);

  std::vector<std::string> get_keys(const std::string& pattern);

  bool is_connected() const { return is_connected_; }

 private:
  std::unique_ptr<cpp_redis::client> db_client_;
  bool is_connected_;

  /**
   * Read the wrapper RedisState value from Redis for a key
   * @param key
   * @param state_out
   * @return response code of operation
   */
  status_code_e read_redis_state(
      const std::string& key, orc8r::RedisState& state_out);

  /**
   * Takes a string and parses it to protobuf Message
   * @param proto_msg
   * @param str_to_deserialize
   * @return response code of operation
   */
  static status_code_e deserialize(
      google::protobuf::Message& proto_msg,
      const std::string& str_to_deserialize);
};

}  // namespace lte
}  // namespace magma
