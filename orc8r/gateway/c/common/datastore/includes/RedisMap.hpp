/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#pragma once

#include "ObjectMap.h"
#include "magma_logging.h"
#include <orc8r/protos/redis.pb.h>

using magma::orc8r::RedisState;

namespace magma {

/**
 * RedisMap stores objects using the redis hash structure. This map requires a
 * serializer and deserializer to store objects as strings in the redis store
 */
template <typename ObjectType>
class RedisMap : public ObjectMap<ObjectType> {
public:
  RedisMap(
    std::shared_ptr<cpp_redis::client> client,
    const std::string& hash,
    std::function<bool(const ObjectType&, std::string&, uint64_t&)> serializer,
    std::function<bool(const std::string&, ObjectType&)> deserializer)
    : client_(client),
      hash_(hash),
      serializer_(serializer),
      deserializer_(deserializer) {}

  /**
   * set serializes the object passed into a string and stores it at the key.
   * Returns false if the operation was unsuccessful
   */
  ObjectMapResult set(
      const std::string& key,
      const ObjectType& object) override {
    uint64_t version;
    auto res = this->get_version(key, version);
    if (res != SUCCESS) {
      MLOG(MERROR) << "Unable to get version for key successfully " << key;
      return res;
    }
    std::string value;
    auto new_version = version + (uint64_t) 1u;
    if (!serializer_(object, value, new_version)) {
      MLOG(MERROR) << "Unable to serialize value for key " << key;
      return SERIALIZE_FAIL;
    }
    auto hset_future = client_->hset(hash_, key, value);
    client_->sync_commit();
    bool is_error = hset_future.get().is_error();
    if (is_error) {
      MLOG(MERROR) << "Error setting value in redis for key " << key;
      return CLIENT_ERROR;
    }
    return SUCCESS;
  }

  /**
   * get returns the object located at key. If the key was not found or the
   * operation was unsuccessful, this returns false
   */
  ObjectMapResult get(const std::string& key, ObjectType& object_out) override {
    auto hget_future = client_->hget(hash_, key);
    client_->sync_commit();
    auto reply = hget_future.get();
    if (reply.is_null()) {
      // value just doesn't exist
      return KEY_NOT_FOUND;
    } else if (reply.is_error()) {
      MLOG(MERROR) << "Unable to get value for key " << key;
      return CLIENT_ERROR;
    } else if (!reply.is_string()) {
      MLOG(MERROR) << "Value was not string for key " << key;
      return INCORRECT_VALUE_TYPE;
    }
    if (!deserializer_(reply.as_string(), object_out)) {
      MLOG(MERROR) << "Failed to deserialize key " << key
        << " with value " << reply.as_string();
      return DESERIALIZE_FAIL;
    }
    return SUCCESS;
  }

  /**
   * getall returns all values stored in the hash
   */
  ObjectMapResult getall(std::vector<ObjectType>& values_out) override {
      return getall(values_out, nullptr);
  }

  /**
   * getall is an extra overloaded function that also returns the keys of values
   * that failed to be deserialized.
   */
  ObjectMapResult getall(
    std::vector<ObjectType>& values_out,
    std::vector<std::string>* failed_keys) {
    auto hgetall_future = client_->hgetall(hash_);
    client_->sync_commit();
    auto reply = hgetall_future.get();
    if (reply.is_error()) {
      MLOG(MERROR) << "unable to perform hvals command";
      return CLIENT_ERROR;
    } else if (reply.is_null()) {
      // fine, just no values
      return SUCCESS;
    }
    auto array = reply.as_array();
    for (unsigned int i = 0; i < array.size(); i += 2) {
      auto key_reply = array[i];
      if (!key_reply.is_string()) {
        // this should essentially never happen
        MLOG(MERROR) << "Non string key found";
      }
      auto key = key_reply.as_string();
      auto value_reply = array[i+1];
      if (!value_reply.is_string()) {
        MLOG(MERROR) << "Non string value found";
        if (failed_keys != nullptr) failed_keys->push_back(key);
        continue;
      }
      ObjectType obj;
      if (!deserializer_(value_reply.as_string(), obj)) {
        MLOG(MERROR) << "Unable to deserialize value in map";
        if (failed_keys != nullptr) failed_keys->push_back(key);
        continue;
      }
      values_out.push_back(obj);
    }
    return SUCCESS;
  }

private:
  /*
   * Return the version of the value for key *key*. Returns 0 if
   * key is not in the map
   */
  ObjectMapResult get_version(const std::string& key, uint64_t& version) {
    auto redis_state = RedisState();
    auto res = this->get_redis_state(key, redis_state);
    if (res == KEY_NOT_FOUND) {
      version = (uint64_t) 0u;
      return SUCCESS;
    } else if (res != SUCCESS) {
      return res;
    }
    version = redis_state.version();
    return res;
  }

  /**
   * Gets the RedisState result that stores the actual value we want
   */
  ObjectMapResult get_redis_state(
      const std::string& key,
      RedisState& redis_state) {
    auto hget_future = client_->hget(hash_, key);
    client_->sync_commit();
    auto reply = hget_future.get();
    if (reply.is_null()) {
      // value just doesn't exist
      return KEY_NOT_FOUND;
    } else if (reply.is_error()) {
      MLOG(MERROR) << "Unable to get value for key " << key;
      return CLIENT_ERROR;
    } else if (!reply.is_string()) {
      MLOG(MERROR) << "Value was not string for key " << key;
      return INCORRECT_VALUE_TYPE;
    }
    if (!redis_state.ParseFromString(reply.as_string())) {
      MLOG(MERROR) << "Failed to deserialize key " << key
                   << " with value " << reply.as_string();
      return DESERIALIZE_FAIL;
    }
    return SUCCESS;
  }

  std::shared_ptr<cpp_redis::client> client_;
  std::string hash_;
  std::function<bool(const ObjectType&, std::string&, uint64_t&)> serializer_;
  std::function<bool(const std::string&, ObjectType&)> deserializer_;
};

}
