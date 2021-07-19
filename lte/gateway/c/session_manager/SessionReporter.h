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

#include <folly/io/async/EventBase.h>
#include <grpc++/grpc++.h>
#include <lte/protos/session_manager.grpc.pb.h>

#include <functional>
#include <memory>

#include "includes/GRPCReceiver.h"

namespace magma {
using namespace lte;

/**
 * ReporterCallbackFn is a function type alias for the callback called when a
   response is received.
 */
template<typename ResponseType>
using ReporterCallbackFn = std::function<void(grpc::Status, ResponseType)>;

/**
 * AsyncEvbResponse is used to call a callback in a particular event loop when
 * a response is received. This is defined here to limit the dependency on folly
 * in the common library.
 */
template<typename ResponseType>
class AsyncEvbResponse : public AsyncGRPCResponse<ResponseType> {
 public:
  AsyncEvbResponse(
      folly::EventBase* base, ReporterCallbackFn<ResponseType> callback,
      uint32_t timeout_sec);

  void handle_response() override;

 private:
  folly::EventBase* base_;
};

class SessionReporter : public GRPCReceiver {
 public:
  virtual ~SessionReporter() = default;

  /**
   * Either proxy an UpdateSessionRequest gRPC call to the cloud
   * or send the request to the local PCRF/OCS on the gateway
   */
  virtual void report_updates(
      const UpdateSessionRequest& request,
      ReporterCallbackFn<UpdateSessionResponse> callback) = 0;

  /**
   * Either proxy a CreateSessionRequest gRPC call to the cloud
   * or send the request to the local PCRF/OCS on the gateway
   */
  virtual void report_create_session(
      const CreateSessionRequest& request,
      ReporterCallbackFn<CreateSessionResponse> callback) = 0;

  /**
   * Either proxy a SessionTerminateRequest gRPC call to the cloud
   * or send the request to the local PCRF/OCS on the gateway
   */
  virtual void report_terminate_session(
      const SessionTerminateRequest& request,
      ReporterCallbackFn<SessionTerminateResponse> callback) = 0;

  /**
   * get_terminate_logging_cb returns a callback function for report_terminate
   * that logs whether or not there was an error.
   */
  static ReporterCallbackFn<SessionTerminateResponse> get_terminate_logging_cb(
      const SessionTerminateRequest& request);
};

class SessionReporterImpl : public SessionReporter {
 public:
  SessionReporterImpl(
      folly::EventBase* base, std::shared_ptr<grpc::Channel> channel);

  void report_updates(
      const UpdateSessionRequest& request,
      std::function<void(grpc::Status, UpdateSessionResponse)> callback);

  void report_create_session(
      const CreateSessionRequest& request,
      std::function<void(grpc::Status, CreateSessionResponse)> callback);

  void report_terminate_session(
      const SessionTerminateRequest& request,
      std::function<void(grpc::Status, SessionTerminateResponse)> callback);

 private:
  folly::EventBase* base_;
  std::unique_ptr<CentralSessionController::Stub> stub_;
  static const uint32_t RESPONSE_TIMEOUT = 6;  // seconds
};

}  // namespace magma
