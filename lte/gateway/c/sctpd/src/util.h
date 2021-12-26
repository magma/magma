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

#include <arpa/inet.h>
#include <errno.h>
#include <stdint.h>

#include <string>

#include <lte/protos/sctpd.grpc.pb.h>

#include "orc8r/gateway/c/common/logging/magma_logging.h"

namespace magma {
namespace sctpd {

#define MLOG_perror(fname)                                       \
  do {                                                           \
    MLOG(MERROR) << fname << " error (" << std::to_string(errno) \
                 << "): " << strerror(errno);                    \
  } while (0)
#define MLOG_grpcerr(status)                                              \
  do {                                                                    \
    MLOG(MERROR) << "grpc error (" << std::to_string(status.error_code()) \
                 << "): " << status.error_message();                      \
  } while (0)

int create_sctp_sock(const InitReq& req);
int pull_peer_ipaddr(const int sd, const uint32_t assoc_id,
                     std::string& ran_cp_ipaddr);

}  // namespace sctpd
}  // namespace magma
