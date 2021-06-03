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
#include <grpc/impl/codegen/port_platform.h>
#include <stdarg.h>
#include <stdlib.h>

#include <chrono>
#include <ctime>
#include <memory>

#include "magma_logging.h"
#include "SessionManagerServer.h"

using grpc::Status;

namespace magma {
auto timenow = chrono::system_clock::to_time_t(chrono::system_clock::now());
AsyncService::AsyncService(std::unique_ptr<ServerCompletionQueue> cq)
    : cq_(std::move(cq)) {}

void AsyncService::wait_for_requests() {
  init_call_data();
  void* tag;
  bool ok;
  running_ = true;
  while (running_) {
    if (cq_ == nullptr || !cq_->Next(&tag, &ok)) {
      // queue shut down
      MLOG(MINFO) << "sessiond request completion queue shutting down";
      return;
    }
    if (!ok) {
      MLOG(MINFO)
          << "sessiond server encountered error while processing request";
      continue;
    }
    static_cast<CallData*>(tag)->proceed();
  }
}

void AsyncService::stop() {
  running_ = false;
  cq_->Shutdown();
  // Pop all items in the queue until it is empty
  // https://github.com/grpc/grpc/issues/8610
  void* tag;
  bool ok;
  while (cq_->Next(&tag, &ok)) {
  }
}

LocalSessionManagerAsyncService::LocalSessionManagerAsyncService(
    std::unique_ptr<ServerCompletionQueue> cq,
    std::unique_ptr<LocalSessionManagerHandler> handler)
    : AsyncService(std::move(cq)), handler_(std::move(handler)) {}

void LocalSessionManagerAsyncService::init_call_data() {
  new ReportRuleStatsCallData(cq_.get(), *this, *handler_);
  new CreateSessionCallData(cq_.get(), *this, *handler_);
  new EndSessionCallData(cq_.get(), *this, *handler_);
  new BindPolicy2BearerCallData(cq_.get(), *this, *handler_);
  new SetSessionRulesCallData(cq_.get(), *this, *handler_);
  new UpdateTunnelIdsCallData(cq_.get(), *this, *handler_);
}

/*Landing object invocation object call for 5G*/
AmfPduSessionSmContextAsyncService::AmfPduSessionSmContextAsyncService(
    std::unique_ptr<ServerCompletionQueue> cq,
    std::unique_ptr<SetMessageManager> handler)
    : AsyncService(std::move(cq)), handler_(std::move(handler)) {}

void AmfPduSessionSmContextAsyncService::init_call_data() {
  new SetAmfSessionContextCallData(cq_.get(), *this, *handler_);
  MLOG(MINFO) << "Initializing new call data for SetAmfSessionContext";
}

/*Landing object invocation object call for 5G*/
SetInterfaceForUserPlaneAsyncService::SetInterfaceForUserPlaneAsyncService(
    std::unique_ptr<ServerCompletionQueue> cq,
    std::unique_ptr<UpfMsgManageHandler> handler)
    : AsyncService(std::move(cq)), handler_(std::move(handler)) {}

void SetInterfaceForUserPlaneAsyncService::init_call_data() {
  MLOG(MINFO) << "Initializing new call data for SetUpfNodeStateCallData";
  new SetUPFNodeStateCallData(cq_.get(), *this, *handler_);
}

SessionProxyResponderAsyncService::SessionProxyResponderAsyncService(
    std::unique_ptr<ServerCompletionQueue> cq,
    std::shared_ptr<SessionProxyResponderHandler> handler)
    : AsyncService(std::move(cq)), handler_(handler) {}

void SessionProxyResponderAsyncService::init_call_data() {
  new ChargingReAuthCallData(cq_.get(), *this, *handler_);
  new PolicyReAuthCallData(cq_.get(), *this, *handler_);
}

AbortSessionResponderAsyncService::AbortSessionResponderAsyncService(
    std::unique_ptr<ServerCompletionQueue> cq,
    std::shared_ptr<SessionProxyResponderHandler> handler)
    : AsyncService(std::move(cq)), handler_(handler) {}

void AbortSessionResponderAsyncService::init_call_data() {
  new AbortSessionCallData(cq_.get(), *this, *handler_);
}

template<class GRPCService, class RequestType, class ResponseType>
AsyncGRPCRequest<GRPCService, RequestType, ResponseType>::AsyncGRPCRequest(
    ServerCompletionQueue* cq, GRPCService& service)
    : cq_(cq), status_(PROCESS), responder_(&ctx_), service_(service) {}

// Internal state management logic for AsyncGRPCRequest
// Once a request has started processing, create a new AsyncGRPCRequest to
// standby for new requests. After a request has finished processing, delete the
// object.
template<class GRPCService, class RequestType, class ResponseType>
void AsyncGRPCRequest<GRPCService, RequestType, ResponseType>::proceed() {
  if (status_ == PROCESS) {
    clone();  // Create another stand by CallData
    process();
    status_ = FINISH;
  } else {
    GPR_ASSERT(status_ == FINISH);
    delete this;
  }
}

template<class GRPCService, class RequestType, class ResponseType>
std::function<void(Status, ResponseType)> AsyncGRPCRequest<
    GRPCService, RequestType, ResponseType>::get_finish_callback() {
  return [this](Status status, ResponseType response) {
    responder_.Finish(response, status, (void*) this);
  };
}

}  // namespace magma
