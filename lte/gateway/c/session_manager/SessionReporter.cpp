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
#include <glog/logging.h>

#include <iostream>
#include <utility>

#include "GrpcMagmaUtils.h"
#include "magma_logging.h"
#include "includes/ServiceRegistrySingleton.h"
#include "SessionReporter.h"

namespace magma {

template<class ResponseType>
AsyncEvbResponse<ResponseType>::AsyncEvbResponse(
    folly::EventBase* base,
    std::function<void(grpc::Status, ResponseType)> callback,
    uint32_t timeout_sec)
    : AsyncGRPCResponse<ResponseType>(callback, timeout_sec), base_(base) {}

template<class ResponseType>
void AsyncEvbResponse<ResponseType>::handle_response() {
  base_->runInEventBaseThread([this]() {
    this->callback_(this->status_, this->response_);
    delete this;
  });
}

ReporterCallbackFn<SessionTerminateResponse>
SessionReporter::get_terminate_logging_cb(
    const SessionTerminateRequest& request) {
  return [request](grpc::Status status, SessionTerminateResponse response) {
    if (!status.ok()) {
      MLOG(MERROR) << "Failed to terminate session in controller for "
                   << request.session_id() << ": " << status.error_message();
    } else {
      MLOG(MDEBUG) << "Termination successful in controller for "
                   << request.session_id();
    }
  };
}

SessionReporterImpl::SessionReporterImpl(
    folly::EventBase* base, std::shared_ptr<grpc::Channel> channel)
    : base_(base), stub_(CentralSessionController::NewStub(channel)) {}

void SessionReporterImpl::report_updates(
    const UpdateSessionRequest& request,
    ReporterCallbackFn<UpdateSessionResponse> callback) {
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request));

  auto controller_response = new AsyncEvbResponse<UpdateSessionResponse>(
      base_, callback, RESPONSE_TIMEOUT);
  controller_response->set_response_reader(std::move(stub_->AsyncUpdateSession(
      controller_response->get_context(), request, &queue_)));
}

/*void SessionReporterImpl::report_update_callback(UpdateSessionRequest &request,
                                    SessionMap session_map_ptr, SessionUpdate session_update, grpc::Status status, 
                                    UpdateSessionResponse &response){
        PrintGrpcMessage(
        static_cast<const google::protobuf::Message&>(response));

        // clear all the reporting flags
        // TODO this could be done in one go with the SessionStore update below
        session_store_.set_and_save_reporting_flag(false, request, session_uc);
        auto updates_by_session = UpdateRequestsBySession(request);
        if (!status.ok()) {
          MLOG(MERROR)
              << "UpdateSession request to FeG/PolicyDB failed entirely: "
              << status.error_message();
          handle_update_failure(
              *session_map_ptr, updates_by_session, session_uc);
          report_session_update_event_failure(
              *session_map_ptr, updates_by_session, status.error_message());
          session_store_.update_sessions(session_uc);
          return;
        }
        // Success!
        update_session_credits_and_rules(
            *session_map_ptr, response, session_uc);
        report_session_update_event(*session_map_ptr, updates_by_session);
        session_store_.update_sessions(session_uc);
      });
                                      
                                      }*/

void SessionReporterImpl::report_create_session(
    const CreateSessionRequest& request,
    ReporterCallbackFn<CreateSessionResponse> callback) {
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request));
  auto controller_response = new AsyncEvbResponse<CreateSessionResponse>(
      base_, callback, RESPONSE_TIMEOUT);
  controller_response->set_response_reader(std::move(stub_->AsyncCreateSession(
      controller_response->get_context(), request, &queue_)));
}

void SessionReporterImpl::report_terminate_session(
    const SessionTerminateRequest& request,
    ReporterCallbackFn<SessionTerminateResponse> callback) {
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request));
  auto controller_response = new AsyncEvbResponse<SessionTerminateResponse>(
      base_, callback, RESPONSE_TIMEOUT);
  controller_response->set_response_reader(
      std::move(stub_->AsyncTerminateSession(
          controller_response->get_context(), request, &queue_)));
}

}  // namespace magma
