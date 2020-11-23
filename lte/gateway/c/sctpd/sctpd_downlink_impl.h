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

#include <memory>
#include <thread>

#include <grpc/grpc.h>
#include <grpcpp/server_context.h>

#include <lte/protos/sctpd.grpc.pb.h>

#include "sctp_connection.h"

#define S1AP 18
#define NGAP 60

namespace magma {
namespace sctpd {

using grpc::ServerContext;
using grpc::Status;

// Implements the sctpd downlink server
class SctpdDownlinkImpl final : public SctpdDownlink::Service {
 public:
  // Construct a new SctpdDownlinkImpl service
  SctpdDownlinkImpl(SctpEventHandler &uplink_handler);

  // Implementation of SctpdDownlink.Init method (see sctpd.proto for more info)
  Status Init(ServerContext *context, const InitReq *request, InitRes *response)
    override;

  // Implementation of SctpdDownlink.SendDl method (see sctpd.proto for more info)
  Status SendDl( ServerContext *context, const SendDlReq *request, SendDlRes *response)
    override;
  Status create_sctp_connection(std::unique_ptr<SctpConnection>& sctp_connection,
    const InitReq *request, InitRes *response);

  // Close SCTP connection for this SctpdDownlink.
  void stop();

 private:
  SctpEventHandler &_uplink_handler;
  std::unique_ptr<SctpConnection> _sctp_4G_connection;
  std::unique_ptr<SctpConnection> _sctp_5G_connection;
};

} // namespace sctpd
} // namespace magma
