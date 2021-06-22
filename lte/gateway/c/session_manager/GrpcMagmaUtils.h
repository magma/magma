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

#include <string>  // for string

namespace google {
namespace protobuf {
class Message;
}
}  // namespace google

void set_grpc_logging_level(bool enable);

std::string get_env_var(std::string const& key);

void PrintGrpcMessage(const google::protobuf::Message& message);

std::string indentText(std::string basicString, int indent);
