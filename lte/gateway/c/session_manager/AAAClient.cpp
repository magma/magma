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

#include <memory>
#include <string>
#include <utility>

#include "AAAClient.h"
#include "magma_logging.h"
#include "includes/ServiceRegistrySingleton.h"
#include "SessionState.h"

using grpc::Status;

namespace {  // anonymous

aaa::terminate_session_request create_deactivate_req(
    const std::string& radius_session_id, const std::string& imsi) {
  aaa::terminate_session_request req;
  req.set_radius_session_id(radius_session_id);
  req.set_imsi(imsi);
  return req;
}

aaa::add_sessions_request create_add_sessions_req(
    const magma::lte::SessionMap& session_map) {
  aaa::add_sessions_request req;
  for (auto it = session_map.begin(); it != session_map.end(); it++) {
    for (const auto& session : it->second) {
      aaa::context ctx;
      if (!session->is_radius_cwf_session()) {
        continue;
      }
      const auto& config = session->get_config();
      if (!config.rat_specific_context.has_wlan_context()) {
        MLOG(MWARNING) << "Session " << session->get_session_id() << " does not"
                       << " have WLAN specific session context";
        continue;
      }
      const auto& wlan_context = config.rat_specific_context.wlan_context();
      ctx.set_imsi(session->get_imsi());
      ctx.set_session_id(wlan_context.radius_session_id());
      ctx.set_acct_session_id(session->get_session_id());
      ctx.set_mac_addr(wlan_context.mac_addr());
      ctx.set_msisdn(config.common_context.msisdn());
      ctx.set_apn(config.common_context.apn());
      req.mutable_sessions()->Add()->CopyFrom(ctx);
    }
  }
  return req;
}

}  // namespace

namespace aaa {

AsyncAAAClient::AsyncAAAClient(std::shared_ptr<grpc::Channel> channel)
    : stub_(accounting::NewStub(channel)) {}

AsyncAAAClient::AsyncAAAClient()
    : AsyncAAAClient(
          magma::ServiceRegistrySingleton::Instance()->GetGrpcChannel(
              "aaa_server", magma::ServiceRegistrySingleton::LOCAL)) {}

bool AsyncAAAClient::terminate_session(
    const std::string& radius_session_id, const std::string& imsi) {
  auto req = create_deactivate_req(radius_session_id, imsi);
  terminate_session_rpc(
      req, [radius_session_id, imsi](Status status, acct_resp resp) {
        if (status.ok()) {
          MLOG(MDEBUG) << "Terminated session for Radius ID:"
                       << radius_session_id << ", IMSI: " << imsi;
        } else {
          MLOG(MERROR) << "Could not add terminate session. Radius ID:"
                       << radius_session_id << ", IMSI: " << imsi
                       << ", Error: " << status.error_message();
        }
      });
  return true;
}

bool AsyncAAAClient::add_sessions(const magma::lte::SessionMap& session_map) {
  auto req = create_add_sessions_req(session_map);
  if (req.sessions().size() == 0) {
    MLOG(MINFO) << "Not sending add_sessions request to AAA server. No AAA "
                << "sessions found";
    return true;
  }
  add_sessions_rpc(req, [this](Status status, acct_resp resp) {
    if (status.ok()) {
      MLOG(MINFO) << "Successfully added all existing sessions to AAA server";
    } else {
      MLOG(MERROR) << "Could not add existing sessions to AAA server,"
                   << " Error: " << status.error_message();
    }
  });
  return true;
}

void AsyncAAAClient::add_sessions_rpc(
    const add_sessions_request& request,
    std::function<void(Status, acct_resp)> callback) {
  auto local_resp = new magma::AsyncLocalResponse<acct_resp>(
      std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(std::move(
      stub_->Asyncadd_sessions(local_resp->get_context(), request, &queue_)));
}

void AsyncAAAClient::terminate_session_rpc(
    const terminate_session_request& request,
    std::function<void(Status, acct_resp)> callback) {
  auto local_resp = new magma::AsyncLocalResponse<acct_resp>(
      std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(std::move(stub_->Asyncterminate_session(
      local_resp->get_context(), request, &queue_)));
}

}  // namespace aaa
