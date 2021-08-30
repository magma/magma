/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#include <thread>
#include "includes/MagmaService.h"
#include "S6aGatewayImpl.h"

extern "C" {
#include "log.h"
#include "s6a_service_handler.h"
}

namespace magma {
#define S6A_ASYNC_PROXY_SERVICE "s6a_async_service"
#define S6A_ASYNC_PROXY_VERSION "1.0"

magma::S6aProxyResponderAsyncService s6a_async_service(nullptr, nullptr);
static magma::service303::MagmaService server(
    S6A_ASYNC_PROXY_SERVICE, S6A_ASYNC_PROXY_VERSION);

void stop_async_s6a_grpc_server(void) {
  s6a_async_service.stop();  // stop queue after server shuts down
}

void init_async_s6a_grpc_server(void) {
  auto async_service_handler =
      std::make_shared<magma::S6aProxyAsyncResponderHandler>();
  s6a_async_service.set_callback(
      server.GetNewCompletionQueue(), async_service_handler);
  server.AddServiceToServer(&s6a_async_service);

  server.Start();
  OAILOG_INFO(LOG_S6A, "Started async grpc server for s6a interface \n");

  std::thread proxy_thread([&]() {
    s6a_async_service.wait_for_requests();  // block here instead of on server
  });
  proxy_thread.join();
  return;
}

AsyncService::AsyncService(std::unique_ptr<ServerCompletionQueue> cq)
    : cq_(std::move(cq)) {}

void AsyncService::wait_for_requests() {
  init_call_data();
  void* tag;
  bool ok;
  running_ = true;
  while (running_) {
    if (cq_ == nullptr) {
      OAILOG_CRITICAL(LOG_S6A, "Completion queue on s6a interface is null \n");
      return;
    }
    if (!(cq_->Next(&tag, &ok))) {
      OAILOG_ERROR(
          LOG_S6A, "Completion queue on s6a interface is shutting down \n");
      return;
    }
    if (!ok) {
      continue;
    }
    static_cast<CallData*>(tag)->proceed();
  }
}

void AsyncService::stop() {
  running_ = false;
  // Pop all items in the queue until it is empty
  void* tag;
  bool ok;
  while (cq_->Next(&tag, &ok)) {
  }
  server.Stop();
  cq_->Shutdown();
}

S6aProxyResponderAsyncService::S6aProxyResponderAsyncService(
    std::unique_ptr<ServerCompletionQueue> cq,
    std::shared_ptr<S6aProxyAsyncResponderHandler> handler)
    : AsyncService(std::move(cq)), handler_(handler) {}

void S6aProxyResponderAsyncService::init_call_data(void) {
  new CancelLocationCallData(cq_.get(), *this, *handler_);
}

void S6aProxyResponderAsyncService::set_callback(
    std::unique_ptr<ServerCompletionQueue> cq,
    std::shared_ptr<S6aProxyAsyncResponderHandler> handler) {
  cq_      = std::move(cq);
  handler_ = handler;
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
std::function<void(grpc::Status, ResponseType)> AsyncGRPCRequest<
    GRPCService, RequestType, ResponseType>::get_finish_callback() {
  return [this](grpc::Status status, ResponseType response) {
    responder_.Finish(response, status, (void*) this);
  };
}

void S6aProxyAsyncResponderHandler::CancelLocation(
    ServerContext* context, const CancelLocationRequest* request,
    std::function<void(grpc::Status, CancelLocationAnswer)> response_callback) {
  auto& request_cpy = *request;
  CancelLocationAnswer ans;
  Status status;
  status = Status::OK;
  ans.set_error_code(ErrorCode::SUCCESS);
  response_callback(status, ans);
  auto imsi              = request->user_name();
  auto cancellation_type = request->cancellation_type();
  OAILOG_INFO(
      LOG_MME_APP, "Received CLR for %s of type %d\n ", imsi.c_str(),
      cancellation_type);
  if (cancellation_type == CancelLocationRequest::SUBSCRIPTION_WITHDRAWAL) {
    auto imsi_len = imsi.length();
    delete_subscriber_request(imsi.c_str(), imsi_len);
  }
  return;
}

}  // namespace magma
