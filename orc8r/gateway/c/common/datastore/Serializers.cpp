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

#include "includes/Serializers.h"
#include <google/protobuf/message.h>  // for Message
#include <orc8r/protos/redis.pb.h>    // for RedisState

using google::protobuf::Message;
using magma::orc8r::RedisState;
namespace magma {

std::function<bool(const Message&, std::string&, uint64_t&)>
get_proto_serializer() {
  return [](const Message& message, std::string& str_out,
            uint64_t& version) -> bool {
    auto can_parse = message.SerializeToString(&str_out);
    if (!can_parse) {
      return false;
    }
    auto redis_state = RedisState();
    redis_state.set_version(version);
    redis_state.set_serialized_msg(str_out);
    return redis_state.SerializeToString(&str_out);
  };
}

std::function<bool(const std::string&, Message&)> get_proto_deserializer() {
  return [](const std::string& str, Message& msg_out) -> bool {
    RedisState redis_state;
    auto can_parse = redis_state.ParseFromString(str);
    if (!can_parse) {
      return false;
    }
    return msg_out.ParseFromString(redis_state.serialized_msg());
  };
}
}  // namespace magma
