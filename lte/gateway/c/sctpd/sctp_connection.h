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

#include <atomic>
#include <memory>
#include <string>
#include <thread>

#include <lte/protos/sctpd.grpc.pb.h>

#include "sctp_desc.h"

struct sctp_assoc_change;

namespace magma {
namespace sctpd {

// Describes status of Sctp event (up or down stream)
enum class SctpStatus {
  OK,                      // Sctp event was ok
  FAILURE,                 // General failure - nonfatal
  DISCONNECT,              // Sctp assoc disconnected
  NEW_ASSOC_NOTIF_FAILED,  // GRPC call for new assoc notification failed
};

// Interface for upstream Sctp event handling
class SctpEventHandler {
 public:
  // Specification for NewAssoc handler function
  virtual int HandleNewAssoc(
      uint32_t ppid, uint32_t assoc_id, uint32_t instreams, uint32_t outstreams,
      std::string& ran_cp_ipaddr) = 0;

  // Specification for CloseAssoc handler function
  virtual void HandleCloseAssoc(
      uint32_t ppid, uint32_t assoc_id, bool reset) = 0;

  // Specification for Recv handler function
  virtual void HandleRecv(
      uint32_t ppid, uint32_t assoc_id, uint32_t stream,
      const std::string& payload) = 0;
};

// Manages Sctp connection including setup/teardown and send/recv
class SctpConnection {
 public:
  // Construct as per the InitReq and sending upstream events to handler
  SctpConnection(const InitReq& req, SctpEventHandler& handler);

  // Start SCTP connection and begin listening/relaying events to handler
  void Start();
  // Close the SCTP connection and its associations - blocking call
  void Close();

  // Send a message on the Sctp connection to (assoc_id, stream)
  void Send(uint32_t assoc_id, uint32_t stream, const std::string& msg);

 private:
  // Listener loop run in separate thread by Start
  void Listen();
  // Handle an event on a client socket
  SctpStatus HandleClientSock(int sd);
  // Handle an association change event for an association sd/change
  SctpStatus HandleAssocChange(int sd, struct sctp_assoc_change* change);
  // Handle a comup event on an association sd/change
  SctpStatus HandleComUp(int sd, struct sctp_assoc_change* change);
  // Handle a comdown event on an association keyed by assoc_id
  SctpStatus HandleComDown(uint32_t assoc_id);
  // Handle a reset event on an association keyed by assoc_id
  SctpStatus HandleReset(uint32_t assoc_id);

  // Flag is set to true when the connection is closing
  std::atomic<bool> _done;
  // Concrete handler instance to handle upstream events
  SctpEventHandler& _handler;
  // PPID to use with Sctp connections
  int _ppid;
  // Keeps track of sctp and assocation info
  SctpDesc _sctp_desc;
  // Thread for sctp listener to run on
  std::unique_ptr<std::thread> _thread;
};

}  // namespace sctpd
};  // namespace magma
