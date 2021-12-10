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

#include <google/protobuf/message.h>  // for Message
#include <istream>                    // for istream
#include <string>                     // for string

namespace google {
namespace protobuf {
class Message;
}  // namespace protobuf
}  // namespace google

namespace magma {
/**
 * load_service_mconfig_from_file loads an mconfig from the statically defined
 * files.
 * @param service_name - name of service to load
 * @param message - pointer to protobuf message to load file to. Note that
 *                  this should match the message defined in mconfigs.proto.
 * @returns true if message was parsed successfully, false if the file reading
 *          failed, the service configuration is missing, or the message
 *          passed is not the right type
 */
// TODO(#6151): migrate to an absl::Status return type.
bool load_service_mconfig_from_file(const std::string& service_name,
                                    google::protobuf::Message* message);

bool load_service_mconfig(const std::string& service_name,
                          std::istream* config_stream,
                          google::protobuf::Message* message);

void get_mconfig_file(std::ifstream* file);
}  // namespace magma
