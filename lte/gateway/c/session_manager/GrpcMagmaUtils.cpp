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

#include "lte/gateway/c/session_manager/GrpcMagmaUtils.hpp"

#include <glog/logging.h>
#include <google/protobuf/descriptor.h>
#include <google/protobuf/message.h>
#include <stdlib.h>
#include <sstream>
#include <string>

#include "orc8r/gateway/c/common/logging/magma_logging.hpp"

#define MAGMA_PRINT_GRPC_PAYLOAD "MAGMA_PRINT_GRPC_PAYLOAD"

bool grpcLoggingEnabled = false;

// set_grpc_logging_level will only change the level in case
// MAGMA_PRINT_GRPC_PAYLOAD envar is not set
void set_grpc_logging_level(bool enable) {
  std::string val = get_env_var(MAGMA_PRINT_GRPC_PAYLOAD);
  if (val == "") {
    grpcLoggingEnabled = enable;
  } else if (val == "1") {
    grpcLoggingEnabled = true;
  } else {
    grpcLoggingEnabled = false;
  }
  MLOG(MINFO) << "print_grpc_payload set at: " << grpcLoggingEnabled;
}

std::string get_env_var(std::string const& key) {
  MLOG(MINFO) << "Checking env var " << key;
  char* val;
  val = getenv(key.c_str());
  std::string retval = "";
  if (val != NULL) {
    retval = val;
  }
  return std::string(retval);
}

void PrintGrpcMessage(const google::protobuf::Message& msg) {
  if (grpcLoggingEnabled) {
    // Lazy log strategy
    const google::protobuf::Descriptor* desc = msg.GetDescriptor();
    MLOG(MINFO) << "\n"
                << "  " << desc->full_name().c_str() << " {\n"
                << indentText(msg.DebugString(), 6) << "  }";
  }
}

std::string indentText(std::string basicString, int indent) {
  std::stringstream iss(basicString);
  std::string blanks(indent, ' ');
  std::string result = "";
  while (iss.good()) {
    std::string SingleLine;
    getline(iss, SingleLine, '\n');
    // skip empty lines
    if (SingleLine == "") {
      continue;
    }
    result += blanks;
    result += SingleLine;
    // do not add \n on the last line
    result += "\n";
  }
  return result;
}
