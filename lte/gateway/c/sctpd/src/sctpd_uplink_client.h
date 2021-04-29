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

#include <grpcpp/grpcpp.h>

#include <lte/protos/sctpd.grpc.pb.h>

namespace magma {
namespace sctpd {

using grpc::Channel;

// Grpc uplink client to allow sctpd to signal MME
class SctpdUplinkClient {
 public:
  // Construct SctpdUplinkClient with the specified channel
  explicit SctpdUplinkClient(std::shared_ptr<Channel> channel);

  // Send an uplink packet to MME (see sctpd.proto for more info)
  virtual int sendUl(const SendUlReq& req, SendUlRes* res);
  // Notify MME of new association (see sctpd.proto for more info)
  virtual int newAssoc(const NewAssocReq& req, NewAssocRes* res);
  // Notify MME of closing/reseting association (see sctpd.proto for more info)
  virtual int closeAssoc(const CloseAssocReq& req, CloseAssocRes* res);

 private:
  // Stub used for client to communicate with server
  std::unique_ptr<SctpdUplink::Stub> _stub;
  // GRPC call timeout
  static const uint32_t RESPONSE_TIMEOUT = 2;  // seconds
};

}  // namespace sctpd
}  // namespace magma
